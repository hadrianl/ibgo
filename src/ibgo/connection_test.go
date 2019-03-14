package ibgo

import (
	"fmt"
	"testing"
)

func TestConnection(t *testing.T) {
	conn := &IbConnection{}
	conn.connect("127.0.0.1", 7497)
	buf := make([]byte, 0)
	_, err := conn.Read(buf)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(buf))
	conn.disconnect()
}
