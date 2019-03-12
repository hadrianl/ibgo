package ibgo

import (
	"bufio"
	"bytes"
	"time"
)

const (
	MaxRequests      = 95
	RequestInternal  = 2
	MaxClientVersion = 148
)

type IbClient struct {
	host          string
	port          int
	clientId      int8
	conn          *IbConnection
	reader        *bufio.Reader
	wrapper       *IbWrapper
	buffer        *bytes.Buffer
	connectOption []byte
	reqIdSeq      int
	msgChan       chan interface{}
	timeChan      chan time.Time
}

func (ic *IbClient) getReqId() int {
	ic.reqIdSeq++
	return ic.reqIdSeq
}

func (ic *IbClient) reset() {
	ic.reqIdSeq = 0
	ic.conn.reset()

}

func NewIbClient(host string, port int, clientId int8) *IbClient {
	ic := &IbClient{
		host:     host,
		port:     port,
		clientId: clientId,
	}
	ic.conn = &IbConnection{}
	ic.reader = bufio.NewReader(ic.conn.conn)
	err := ic.conn.Connect(host, port)
	if err != nil {
		panic(err)
	}

	return ic

}

func (ic *IbClient) Connect(host string, port int, clientId int8) error {
	ic.host = host
	ic.port = port
	ic.clientId = clientId
	ic.conn = &IbConnection{}
	err := ic.conn.connect(host, port)
	return err
	// 连接后开始
}
