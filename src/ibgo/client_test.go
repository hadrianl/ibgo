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
	err = ic.Connect("127.0.0.1", 9045, 10)
	if err != nil {
		fmt.Println("Connect failed:", err)
		return
	}

	fmt.Println("beforeHandShake")
	err = ic.HandShake()
	if err != nil {
		fmt.Println("HandShake failed:", err)
		return
	}
	fmt.Println("afterhandShake")

	ic.Run()
	time.Sleep(time.Second * 10)
	ic.Disconnect()
}
