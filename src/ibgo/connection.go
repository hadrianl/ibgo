package ibgo

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
)

type IbConnection struct {
	host         string
	port         int
	clientId     int8
	conn         net.Conn
	numBytesSent int
	numMsgSent   int
	numBytesRecv int
	numMsgRecv   int
	event        socketEvent
	em           extraMethods
}

type socketEvent interface {
	connected()
	disconnected()
	hasError(err error)
	hasData()
}

type extraMethods interface {
	priceSizeTick()
	tcpDataArrived()
	tcpDataProcessed()
}

func (ibconn *IbConnection) sendMsg(msg []byte) error {
	n, err := ibconn.conn.Write(msg)
	ibconn.numBytesSent += n
	ibconn.numMsgSent++
	return err
}

func (ibconn *IbConnection) recvMsg() []byte {
	result, err := ioutil.ReadAll(ibconn.conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
	}

	return result
}

func (ibconn *IbConnection) reset() {
	ibconn.numBytesSent = 0
	ibconn.numBytesRecv = 0
	ibconn.numMsgSent = 0
	ibconn.numMsgRecv = 0
}

func (ibconn *IbConnection) disconnect() error {
	return ibconn.conn.Close()
}

func (ibconn *IbConnection) connect(host string, port int) error {
	var err error
	var addr *net.TCPAddr
	ibconn.host = host
	ibconn.port = port
	ibconn.reset()
	server := ibconn.host + ":" + strconv.Itoa(port)
	fmt.Println(server)
	addr, err = net.ResolveTCPAddr("tcp4", server)
	panicError(err)

	ibconn.conn, err = net.DialTCP("tcp4", nil, addr)
	panicError(err)

	err = ibconn.handShake()
	panicError(err)

	fmt.Println("connect success!")

	return err
}

func (ibconn *IbConnection) handShake() error {
	var msg bytes.Buffer
	head := []byte("API\x00")
	// minVer := []byte("100")
	// maxVer := []byte("148")
	// connectOptions := []byte("")
	msg.Write(head)
	msg.Write([]byte("\x00\x00\x00\tv100..148"))
	fmt.Println(msg.Bytes())
	err := ibconn.sendMsg(msg.Bytes())
	return err
}

// func (ic *IbConnection) onSocketHasData(data []byte) {
// 	ic.em.tcpDataArrived()
// 	n, err := io.Copy(ec.buffer, ec.conn)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
// 	}
// 	ec.numBytesRecv += n

// 	for ec.buffer.len() <= 4 {
// 		msgLen := 4
// 	}

// }
