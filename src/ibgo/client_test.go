package ibgo

import (
	"testing"
	"time"
	// "time"
)

func TestClient(t *testing.T) {
	var err error
	ic := &IbClient{}
	err = ic.Connect("127.0.0.1", 7497, 10)
	if err != nil {
		panic(err)
	}
	t.Log("afterConnect")
	err = ic.HandShake()
	if err != nil {
		panic(err)
	}

	t.Log("afterhandShake")

	ic.Run()
	time.Sleep(time.Second * 10)
	ic.Disconnect()
}
