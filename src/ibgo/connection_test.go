package ibgo

import (
	"fmt"
	"testing"
)

func TestConnection(t *testing.T) {
	conn := IbConnection{}
	conn.Connect("127.0.0.1", 7497)
	msg := conn.recvMsg()
	fmt.Println(string(msg))
	conn.Disconnect()
}
