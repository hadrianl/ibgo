package ibgo

import (
	"fmt"
	"net"
	"strconv"
)

type IbConnection struct {
	host         string
	port         int
	clientId     int64
	conn         net.Conn
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
	if err != nil {
		ibconn.event.hasError <- err
		// ibconn.reset()
	} else {
		ibconn.event.hasData <- b
	}

	return n, err
}

// func (ibconn *IbConnection) Receive() {
// 	buf := make([]byte, 0, 4096)
// 	ibconn.Read(buf)
// 	return buf
// }

func (ibconn *IbConnection) reset() {
	ibconn.numBytesSent = 0
	ibconn.numBytesRecv = 0
	ibconn.numMsgSent = 0
	ibconn.numMsgRecv = 0
	ibconn.event.connected = make(chan int, 10)
	ibconn.event.disconnected = make(chan int, 10)
	ibconn.event.hasError = make(chan error, 100)
	ibconn.event.hasData = make(chan []byte, 100)
}

func (ibconn *IbConnection) disconnect() error {
	ibconn.event.disconnected <- 1
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
	if err != nil {
		fmt.Println("ResolveTCPAddr Error:", err)
		return err
	}

	ibconn.conn, err = net.DialTCP("tcp4", nil, addr)
	if err != nil {
		fmt.Println("DialTCP Error:", err)
		return err
	}

	fmt.Println("connect success!", ibconn.conn.RemoteAddr())
	ibconn.event.connected <- 1

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
