package ibgo

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
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
	inBuffer         *bytes.Buffer
	outBuffer        *bytes.Buffer
	connectOption    []byte
	reqIdSeq         int
	msgChan          chan interface{}
	timeChan         chan time.Time
	terminatedSignal chan int
	clientVersion    Version
	serverVersion    Version
	serverTime       time.Time
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
	err := ic.conn.connect(host, port)
	if err != nil {
		panic(err)
	}

	return ic

}

func (ic *IbClient) Connect(host string, port int, clientId int64) error {
	ic.host = host
	ic.port = port
	ic.clientId = clientId
	ic.reset()
	err := ic.conn.connect(host, port)
	// if err != nil {
	// 	panic(err)
	// }
	return err
	// 连接后开始
}

func (ic *IbClient) Disconnect() error {
	err := ic.conn.disconnect()
	return err
}

// handshake with the TWS or GateWay to ensure the version
func (ic *IbClient) HandShake() error {
	var msg bytes.Buffer
	head := []byte("API\x00")
	// minVer := []byte("100")
	// maxVer := []byte("148")
	// connectOptions := []byte("")
	clientVersion := []byte("v100..148")
	sizeofCV := make([]byte, 4)
	binary.BigEndian.PutUint32(sizeofCV, uint32(len(clientVersion)))
	msg.Write(head)
	msg.Write(sizeofCV)
	msg.Write(clientVersion)
	fmt.Println("Send API Init")
	fmt.Println(msg.Bytes())
	if _, err := ic.writer.Write(msg.Bytes()); err != nil {
		return err
	}

	if err := ic.writer.Flush(); err != nil {
		return err
	}

	fmt.Println("Recv ServerInfo Init")
	if msgBuf, err := readMsgBuf(ic.reader); err != nil {
		fmt.Println(err)
		return err
	} else {
		serverInfo := splitMsgBuf(msgBuf)
		fmt.Println("Recv ServerInfo:", serverInfo)
		ic.serverVersion = Version(bytesToInt(serverInfo[0]))
		ic.serverTime = bytesToTime(serverInfo[1])
	}

	err := ic.startAPI()

	return err
}

// send the clientId to TWS or Gateway
func (ic *IbClient) startAPI() error {
	v := 2
	fields := []string{strconv.FormatInt(int64(START_API), 10), strconv.FormatInt(int64(v), 10), strconv.FormatInt(int64(ic.clientId), 10)}
	// fields.
	fmt.Println("Start API", strings.Join(fields, "\x00"))
	_, err := ic.writer.WriteString(strings.Join(fields, "\x00"))
	return err
}

func (ic *IbClient) reset() {
	ic.reqIdSeq = 0
	ic.conn = &IbConnection{}
	ic.conn.reset()
	ic.reader = bufio.NewReader(ic.conn)
	ic.writer = bufio.NewWriter(ic.conn)
	ic.inBuffer = bytes.NewBuffer(make([]byte, 0, 4096))
	ic.outBuffer = bytes.NewBuffer(make([]byte, 0, 4096))

}

func (ic *IbClient) onSocketConnected() {

}

func (ic *IbClient) onSocketHasData(newData []byte) {

}

func (ic *IbClient) onSocketDisconnected() {

}

func (ic *IbClient) onSocketHasError(e error) {

}

func (ic *IbClient) goReceive() {
	// defer
	fmt.Println("Start receive!")
	// buf := make([]byte, 0, 4096)
	for {
		// buf := []byte
		msgBuf, err := readMsgBuf(ic.reader)
		if err != nil {
			fmt.Println("readmsgBuf Error:", err)
		}

		fields := splitMsgBuf(msgBuf)
		fmt.Println(fields)
		// select {
		// // case _, err := ic.reader.Read(buf):
		// // 	if err != nil {
		// // 		fmt.Println(err)
		// // 	}
		// // 	fmt.Println(buf)
		// case <-ic.conn.event.disconnected: // got msg from disconnected channel and break the loop
		// 	break

		// }
	}
}

func (ic *IbClient) Run() {
	fmt.Println("setup receiver")
	go ic.goReceive()
}
