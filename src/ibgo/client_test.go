package ibgo

import (
	"log"
	"testing"
	"time"
	// "time"
)

func TestClient(t *testing.T) {

	var err error
	ic := &IbClient{}
	err = ic.Connect("127.0.0.1", 7497, 0)
	if err != nil {
		log.Panic("Connect failed:", err)
		return
	}

	err = ic.HandShake()
	if err != nil {
		log.Println("HandShake failed:", err)
		return
	}

	ic.reqCurrentTime()
	ic.reqAutoOpenOrders(true)
	ic.reqAccountUpdates(true, "")

	ic.Run()
	time.Sleep(time.Second * 10)
	ic.Disconnect()
}
