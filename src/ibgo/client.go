package ibgo

import (
	"bufio"
	"bytes"
	"encoding/binary"
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

type IbClient struct {
	host             string
	port             int
	clientId         int64
	conn             *IbConnection
	reader           *bufio.Reader
	writer           *bufio.Writer
	wrapper          IbWrapper
	decoder          *IbDecoder
	inBuffer         *bytes.Buffer
	outBuffer        *bytes.Buffer
	connectOption    []byte
	reqIdSeq         int
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

func (ic *IbClient) getReqId() int {
	ic.reqIdSeq++
	return ic.reqIdSeq
}

func NewIbClient(host string, port int, clientId int64) *IbClient {
	ic := &IbClient{
		host:     host,
		port:     port,
		clientId: clientId,
	}
	ic.reset()
	if err := ic.conn.connect(host, port); err != nil {
		panic(err)
	}

	return ic

}

func (ic *IbClient) Connect(host string, port int, clientId int64) error {

	ic.host = host
	ic.port = port
	ic.clientId = clientId
	ic.reset()
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
	ic.conn.setState(DISCONNECTED)
	defer log.Println("Disconnected!")
	ic.wg.Wait()

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
		ic.decoder.setmsgId2process()
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
	var start_api []byte
	v := 2
	if ic.serverVersion >= MIN_SERVER_VER_OPTIONAL_CAPABILITIES {
		start_api = makeMsg(int64(START_API), int64(v), ic.clientId, "")
	} else {
		start_api = makeMsg(int64(START_API), int64(v), ic.clientId)
	}

	log.Println("Start API:", start_api)
	if _, err := ic.writer.Write(start_api); err != nil {
		return err
	}

	err := ic.writer.Flush()

	return err
}

func (ic *IbClient) reset() {
	ic.reqIdSeq = 0
	ic.conn = &IbConnection{}
	ic.wrapper = Wrapper{ic: ic}
	ic.decoder = &IbDecoder{wrapper: ic.wrapper}
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

func (ic *IbClient) reqCurrentTime() {
	v := 1
	msg := makeMsg(REQ_CURRENT_TIME, v)

	ic.reqChan <- msg
}

// reqAutoOpenOrders will make the client access to the TWS Orders (only if clientId=0)
func (ic *IbClient) reqAutoOpenOrders(autoBind bool) {
	v := 1
	msg := makeMsg(REQ_AUTO_OPEN_ORDERS, v, autoBind)

	ic.reqChan <- msg
}

func (ic *IbClient) reqAccountUpdates(subscribe bool, accName string) {
	v := 2
	msg := makeMsg(REQ_ACCT_DATA, v, subscribe, accName)

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
func (ic *IbClient) Run() {
	log.Println("setup receiver")
	ic.wg.Add(3)
	go ic.goRequest()
	go ic.goReceive()
	go ic.goDecode()
}
