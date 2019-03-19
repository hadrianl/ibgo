package ibgo

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
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
	wrapper          *IbWrapper
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
	ic.wg.Wait()
	fmt.Println("Disconnected!")

	return nil
}

// handshake with the TWS or GateWay to ensure the version
func (ic *IbClient) HandShake() error {
	fmt.Println("Try to handShake with TWS or GateWay...")
	var msg bytes.Buffer
	head := []byte("API\x00")
	minVer := []byte("100")
	maxVer := []byte("148")
	connectOptions := []byte("")
	clientVersion := bytes.Join([][]byte{[]byte("v"), minVer, []byte(".."), maxVer, connectOptions}, []byte(""))
	sizeofCV := make([]byte, 4)
	binary.BigEndian.PutUint32(sizeofCV, uint32(len(clientVersion)))
	msg.Write(head)
	msg.Write(sizeofCV)
	msg.Write(clientVersion)
	fmt.Println("HandShake Init...")
	if _, err := ic.writer.Write(msg.Bytes()); err != nil {
		return err
	}

	if err := ic.writer.Flush(); err != nil {
		return err
	}

	fmt.Println("Recv ServerInfo...")
	if msgBuf, err := readMsgBuf(ic.reader); err != nil {
		return err
	} else {
		serverInfo := splitMsgBuf(msgBuf)
		v, _ := strconv.Atoi(string(serverInfo[0]))
		ic.serverVersion = Version(v)
		ic.serverTime = bytesToTime(serverInfo[1])
		ic.decoder.setVersion(ic.serverVersion) // Init Decoder
		ic.decoder.setmsgId2process()
		fmt.Println("ServerVersion:", ic.serverVersion)
		fmt.Println("ServerTime:", ic.serverTime)
	}

	if err := ic.startAPI(); err != nil {
		return err
	}

	ic.conn.setState(CONNECTED)

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

	fmt.Println("Start API:", start_api)
	if _, err := ic.writer.Write(start_api); err != nil {
		return err
	}

	err := ic.writer.Flush()

	return err
}

func (ic *IbClient) reset() {
	ic.reqIdSeq = 0
	ic.conn = &IbConnection{}
	ic.wrapper = &IbWrapper{}
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

//goRequest will get the req from reqChan and send it to TWS
func (ic *IbClient) goRequest() {
	fmt.Println("Start Request!")
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
			break
		}
		ic.wg.Done()

	}
}

//goReceive receive the msg from the socket, get the fields and put them into msgChan
func (ic *IbClient) goReceive() {
	// defer
	fmt.Println("Start goReceive!")
	// buf := make([]byte, 0, 4096)
	defer ic.reader.Reset(ic.conn)
	for {
		// buf := []byte
		msgBuf, err := readMsgBuf(ic.reader)
		if err, ok := err.(*net.OpError); ok {
			if !err.Temporary() {
				break
			}
			fmt.Println(err)
		} else if err != nil {
			ic.errChan <- err
			fmt.Println("readmsgBuf Error:", err)
			ic.reader.Reset(ic.conn)
		}

		fields := splitMsgBuf(msgBuf)
		ic.msgChan <- fields

	}
	ic.wg.Done()
}

//goDecode decode the fields received from the msgChan
func (ic *IbClient) goDecode() {
	fmt.Println("Start goDecode!")
	for {
		// buf := []byte
		select {
		case f := <-ic.msgChan:
			ic.decoder.interpret(f...)
			fmt.Println(f)
		case <-ic.terminatedSignal:
			break
		}
	}
	ic.wg.Done()

}

func (ic *IbClient) Run() {
	fmt.Println("setup receiver")
	ic.wg.Add(3)
	go ic.goRequest()
	go ic.goReceive()
	go ic.goDecode()
}
