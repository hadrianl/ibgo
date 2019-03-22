/* connection handle the */

package ibgo

import (
	"fmt"
	"log"
	"net"
	"strconv"
)

// ConnectionState
const (
	DISCONNECTED = iota
	CONNECTING
	CONNECTED
	REDIRECT
)

type IbConnection struct {
	host         string
	port         int
	clientId     int64
	conn         net.Conn
	state        int
	numBytesSent int
	numMsgSent   int
	numBytesRecv int
	numMsgRecv   int
	event        socketEvent
	em           extraMethods
}

type socketEvent struct {
	connected    chan int
	disconnected chan int
	hasError     chan error
	hasData      chan []byte
}

type extraMethods interface {
	priceSizeTick()
	tcpDataArrived()
	tcpDataProcessed()
}

func (ibconn *IbConnection) Write(msg []byte) (int, error) {
	n, err := ibconn.conn.Write(msg)
	ibconn.numBytesSent += n
	ibconn.numMsgSent++
	return n, err
}

func (ibconn *IbConnection) Read(b []byte) (int, error) {
	n, err := ibconn.conn.Read(b)
	ibconn.numBytesRecv += n
	ibconn.numMsgRecv++
	// if err != nil {
	// 	ibconn.event.hasError <- err
	// 	// ibconn.reset()
	// } else {
	// 	ibconn.event.hasData <- b
	// }

	return n, err
}

// func (ibconn *IbConnection) Receive() {
// 	buf := make([]byte, 0, 4096)
// 	ibconn.Read(buf)
// 	return buf
// }

func (ibconn *IbConnection) setState(state int) {
	ibconn.state = state
}

func (ibconn *IbConnection) reset() {
	ibconn.numBytesSent = 0
	ibconn.numBytesRecv = 0
	ibconn.numMsgSent = 0
	ibconn.numMsgRecv = 0
	ibconn.setState(DISCONNECTED)
	// ibconn.event.connected = make(chan int, 10)
	// ibconn.event.disconnected = make(chan int, 10)
	// ibconn.event.hasError = make(chan error, 100)
	// ibconn.event.hasData = make(chan []byte, 100)
}

func (ibconn *IbConnection) disconnect() error {
	// ibconn.event.disconnected <- 1
	return ibconn.conn.Close()
}

func (ibconn *IbConnection) connect(host string, port int) error {
	var err error
	var addr *net.TCPAddr
	ibconn.host = host
	ibconn.port = port
	ibconn.reset()
	server := ibconn.host + ":" + strconv.Itoa(port)
	addr, err = net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		log.Printf("ResolveTCPAddr Error: %v", err)
		return err
	}

	ibconn.conn, err = net.DialTCP("tcp4", nil, addr)
	if err != nil {
		log.Printf("DialTCP Error: %v", err)
		return err
	}

	fmt.Println("TCP Connected to:", ibconn.conn.RemoteAddr())
	// ibconn.event.connected <- 1

	return err
}
