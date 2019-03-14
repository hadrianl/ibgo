package ibgo

import (
	"bufio"
	"bytes"
	"fmt"
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
	clientId         int8
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
}

func (ic *IbClient) getReqId() int {
	ic.reqIdSeq++
	return ic.reqIdSeq
}

func NewIbClient(host string, port int, clientId int8) *IbClient {
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

func (ic *IbClient) Connect(host string, port int, clientId int8) error {
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

func (ic *IbClient) Disconnect() err {
	err := ic.conn.disconnect()
	return err
}

func (ic *IbClient) HandShake() error {
	var msg bytes.Buffer
	head := []byte("API\x00")
	// minVer := []byte("100")
	// maxVer := []byte("148")
	// connectOptions := []byte("")
	msg.Write(head)
	msg.Write([]byte("\x00\x00\x00\tv100..148"))
	fmt.Println(msg.Bytes())
	if _, err := ic.writer.Write(msg.Bytes()); err != nil {
		return err
	}
	// serverShake := []byte
	// if _, err := ic.reader.Read(serverShake);err != nil {
	// 	return err
	// }
	// err := ibconn.sendMsg(msg.Bytes())
	return nil
}

func (ic *IbClient) reset() {
	ic.reqIdSeq = 0
	ic.conn = &IbConnection{}
	ic.conn.reset()
	ic.reader = bufio.NewReaderSize(ic.conn, 4096)
	ic.writer = bufio.NewWriterSize(ic.conn, 4096)
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
	buf := make([]byte, 0, 4096)
	for {
		// buf := []byte
		select {
		case _, err := ic.reader.Read(buf):
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(buf)
		case m <- ic.conn.event.disconnected: // got msg from disconnected channel and break the loop
			break

		}
	}
}

func (ic *IbClient) Run() {
	fmt.Println("setup receiver")
	go ic.goReceive()
}
