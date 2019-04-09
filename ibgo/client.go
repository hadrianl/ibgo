package ibgo

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

const (
	MaxRequests      = 95
	RequestInternal  = 2
	MaxClientVersion = 148
)

// IbClient is the key component which is used to send request to TWS ro Gateway , such subscribe market data or place order
type IbClient struct {
	host             string
	port             int
	clientID         int64
	conn             *IbConnection
	reader           *bufio.Reader
	writer           *bufio.Writer
	wrapper          IbWrapper
	decoder          *ibDecoder
	inBuffer         *bytes.Buffer
	outBuffer        *bytes.Buffer
	connectOption    []byte
	reqIDSeq         int64
	reqChan          chan []byte
	errChan          chan error
	msgChan          chan [][]byte
	timeChan         chan time.Time
	terminatedSignal chan int
	clientVersion    Version
	serverVersion    Version
	serverTime       time.Time
	wg               sync.WaitGroup
}

func NewIbClient(wrapper IbWrapper) *IbClient {
	ic := &IbClient{}
	ic.SetWrapper(wrapper)
	ic.reset()

	return ic
}

func (ic *IbClient) GetReqID() int64 {
	ic.reqIDSeq++
	return ic.reqIDSeq
}

func (ic *IbClient) SetWrapper(wrapper IbWrapper) {
	ic.wrapper = wrapper
	ic.decoder = &ibDecoder{wrapper: ic.wrapper}
}

func (ic *IbClient) Connect(host string, port int, clientID int64) error {

	ic.host = host
	ic.port = port
	ic.clientID = clientID
	if err := ic.conn.connect(host, port); err != nil {
		return err
	}

	ic.conn.setState(CONNECTING)
	return nil
	// 连接后开始
}

func (ic *IbClient) Disconnect() error {

	ic.terminatedSignal <- 1
	ic.terminatedSignal <- 1
	ic.terminatedSignal <- 1
	if err := ic.conn.disconnect(); err != nil {
		return err
	}

	defer log.Println("Disconnected!")
	ic.wg.Wait()
	ic.conn.setState(DISCONNECTED)
	return nil
}

// handshake with the TWS or GateWay to ensure the version
func (ic *IbClient) HandShake() error {
	log.Println("Try to handShake with TWS or GateWay...")
	var msg bytes.Buffer
	head := []byte("API\x00")
	minVer := []byte(strconv.FormatInt(int64(MIN_CLIENT_VER), 10))
	maxVer := []byte(strconv.FormatInt(int64(MAX_CLIENT_VER), 10))
	connectOptions := []byte("")
	clientVersion := bytes.Join([][]byte{[]byte("v"), minVer, []byte(".."), maxVer, connectOptions}, []byte(""))
	sizeofCV := make([]byte, 4)
	binary.BigEndian.PutUint32(sizeofCV, uint32(len(clientVersion)))
	msg.Write(head)
	msg.Write(sizeofCV)
	msg.Write(clientVersion)
	log.Println("HandShake Init...")
	if _, err := ic.writer.Write(msg.Bytes()); err != nil {
		return err
	}

	if err := ic.writer.Flush(); err != nil {
		return err
	}

	log.Println("Recv ServerInfo...")
	if msgBuf, err := readMsgBuf(ic.reader); err != nil {
		return err
	} else {
		serverInfo := splitMsgBuf(msgBuf)
		v, _ := strconv.Atoi(string(serverInfo[0]))
		ic.serverVersion = Version(v)
		ic.serverTime = bytesToTime(serverInfo[1])
		ic.decoder.setVersion(ic.serverVersion) // Init Decoder
		ic.decoder.setmsgID2process()
		log.Println("ServerVersion:", ic.serverVersion)
		log.Println("ServerTime:", ic.serverTime)
	}

	if err := ic.startAPI(); err != nil {
		return err
	}

	ic.conn.setState(CONNECTED)
	ic.wrapper.connectAck()

	return nil
}

// send the clientId to TWS or Gateway
func (ic *IbClient) startAPI() error {
	var startAPI []byte
	v := 2
	if ic.serverVersion >= MIN_SERVER_VER_OPTIONAL_CAPABILITIES {
		startAPI = makeMsgBuf(int64(START_API), int64(v), ic.clientID, "")
	} else {
		startAPI = makeMsgBuf(int64(START_API), int64(v), ic.clientID)
	}

	log.Println("Start API:", startAPI)
	if _, err := ic.writer.Write(startAPI); err != nil {
		return err
	}

	err := ic.writer.Flush()

	return err
}

func (ic *IbClient) reset() {
	ic.reqIDSeq = 0
	ic.conn = &IbConnection{}
	ic.conn.reset()
	ic.reader = bufio.NewReader(ic.conn)
	ic.writer = bufio.NewWriter(ic.conn)
	ic.inBuffer = bytes.NewBuffer(make([]byte, 0, 4096))
	ic.outBuffer = bytes.NewBuffer(make([]byte, 0, 4096))
	ic.reqChan = make(chan []byte, 10)
	ic.errChan = make(chan error, 10)
	ic.msgChan = make(chan [][]byte, 100)
	ic.terminatedSignal = make(chan int, 3)
	ic.wg = sync.WaitGroup{}

}

// ---------------req func ----------------------------------------------

func (ic *IbClient) ReqCurrentTime() {
	v := 1
	msg := makeMsgBuf(REQ_CURRENT_TIME, v)

	ic.reqChan <- msg
}

// reqAutoOpenOrders will make the client access to the TWS Orders (only if clientId=0)
func (ic *IbClient) ReqAutoOpenOrders(autoBind bool) {
	v := 1
	msg := makeMsgBuf(REQ_AUTO_OPEN_ORDERS, v, autoBind)

	ic.reqChan <- msg
}

func (ic *IbClient) ReqAccountUpdates(subscribe bool, accName string) {
	v := 2
	msg := makeMsgBuf(REQ_ACCT_DATA, v, subscribe, accName)

	ic.reqChan <- msg

}

func (ic *IbClient) ReqExecutions(reqID int64, execFilter ExecutionFilter) {
	v := 3
	fields := make([]interface{}, 0)
	fields = append(fields, REQ_EXECUTIONS, v)

	if ic.serverVersion >= MIN_SERVER_VER_EXECUTION_DATA_CHAIN {
		fields = append(fields, reqID)
	}

	fields = append(fields,
		execFilter.ClientID,
		execFilter.AccountCode,
		execFilter.Time,
		execFilter.Symbol,
		execFilter.SecurityType,
		execFilter.Exchange,
		execFilter.Side)
	msg := makeMsgBuf(fields...)

	ic.reqChan <- msg
}

func (ic *IbClient) ReqHistoricalData(reqID int64, contract Contract, endDateTime string, duration string, barSize string, whatToShow string, useRTH bool, formatDate int, keepUpToDate bool, chartOptions []TagValue) {
	if ic.serverVersion < MIN_SERVER_VER_TRADING_CLASS {
		if contract.TradingClass != "" || contract.ContractID > 0 {
			ic.wrapper.error(reqID, UPDATE_TWS.Code, UPDATE_TWS.Msg)
		}
	}

	v := 6

	fields := make([]interface{}, 0)
	fields = append(fields, REQ_HISTORICAL_DATA)
	if ic.serverVersion <= MIN_SERVER_VER_SYNT_REALTIME_BARS {
		fields = append(fields, v)
	}

	fields = append(fields, reqID)

	if ic.serverVersion >= MIN_SERVER_VER_TRADING_CLASS {
		fields = append(fields, contract.ContractID)
	}

	fields = append(fields,
		contract.Symbol,
		contract.SecurityType,
		contract.Expiry,
		contract.Strike,
		contract.Right,
		contract.Multiplier,
		contract.Exchange,
		contract.PrimaryExchange,
		contract.Currency,
		contract.LocalSymbol,
	)

	if ic.serverVersion >= MIN_SERVER_VER_TRADING_CLASS {
		fields = append(fields, contract.TradingClass)
	}
	fields = append(fields,
		contract.IncludeExpired,
		endDateTime,
		barSize,
		duration,
		useRTH,
		whatToShow,
		formatDate,
	)

	if contract.SecurityType == "BAG" {
		fields = append(fields, len(contract.ComboLegs))
		for _, comboLeg := range contract.ComboLegs {
			fields = append(fields,
				comboLeg.ContractID,
				comboLeg.Ratio,
				comboLeg.Action,
				comboLeg.Exchange,
			)
		}
	}

	if ic.serverVersion >= MIN_SERVER_VER_SYNT_REALTIME_BARS {
		fields = append(fields, keepUpToDate)
	}

	if ic.serverVersion >= MIN_SERVER_VER_LINKING {
		chartOptionsStr := ""
		for _, tagValue := range chartOptions {
			chartOptionsStr += tagValue.Value
		}
		fields = append(fields, chartOptionsStr)
	}

	msg := makeMsgBuf(fields...)
	// fmt.Println(msg)

	ic.reqChan <- msg
}

//--------------------------three major goroutine -----------------------------------------------------
//goRequest will get the req from reqChan and send it to TWS
func (ic *IbClient) goRequest() {
	log.Println("Start goRequest!")
	defer log.Println("End goRequest!")
	defer ic.wg.Done()
requestLoop:
	for {
		select {
		case req := <-ic.reqChan:
			_, err := ic.writer.Write(req)
			if err != nil {
				ic.writer.Reset(ic.conn)
				ic.errChan <- err
			}
			ic.writer.Flush()
		case <-ic.terminatedSignal:
			// fmt.Println("goRequest terminate")
			break requestLoop
		}
	}

}

//goReceive receive the msg from the socket, get the fields and put them into msgChan
func (ic *IbClient) goReceive() {
	// defer
	log.Println("Start goReceive!")
	// buf := make([]byte, 0, 4096)
	defer log.Println("End goReceive!")
	defer ic.wg.Done()
	for {
		// buf := []byte
		msgBuf, err := readMsgBuf(ic.reader)
		// fmt.Printf("msgBuf: %v err: %v", msgBuf, err)
		if err, ok := err.(*net.OpError); ok {
			if !err.Temporary() {
				break
			}
			log.Println(err)
		} else if err != nil {
			ic.errChan <- err
			ic.reader.Reset(ic.conn)
		}

		if msgBuf != nil {
			fields := splitMsgBuf(msgBuf)
			ic.msgChan <- fields
		}

	}
}

//goDecode decode the fields received from the msgChan
func (ic *IbClient) goDecode() {
	log.Println("Start goDecode!")
	defer log.Println("End goDecode!")
	defer ic.wg.Done()

decodeLoop:
	for {
		// buf := []byte
		select {
		case f := <-ic.msgChan:
			ic.decoder.interpret(f...)
			// log.Println(f)
		case e := <-ic.errChan:
			fmt.Println(e)
		case <-ic.terminatedSignal:
			// fmt.Println("goDecode terminate")
			break decodeLoop
		}
	}

}

// ---------------------------------------------------------------------------------------

// Run make the event loop run, all make sense after run!
func (ic *IbClient) Run() error {
	if ic.conn.state == DISCONNECTED {
		return errors.New("IbClient is DISCONNECTED!!!")
	}
	log.Println("RUN Client")
	ic.wg.Add(3)
	go ic.goRequest()
	go ic.goReceive()
	go ic.goDecode()

	return nil
}
