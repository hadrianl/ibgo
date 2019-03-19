package ibgo

import (
	"fmt"
	"testing"
	"time"
	// "time"
)

func TestClient(t *testing.T) {
	var err error
	ic := &IbClient{}
	err = ic.Connect("127.0.0.1", 7497, 10)
	if err != nil {
		fmt.Println("Connect failed:", err)
		return
	}

	err = ic.HandShake()
	if err != nil {
		fmt.Println("HandShake failed:", err)
		return
	}

	ic.Run()
	time.Sleep(time.Second * 10)
	ic.Disconnect()
}
