package ibgo

import (
	"log"
	"testing"
	"time"
)

func TestIB(t *testing.T) {
	log.SetFlags(log.Lmicroseconds)
	ib := NewIB("127.0.0.1", 7497, 2)
	if err := ib.Connect(); err != nil {
		log.Panicf("Connect failed: %v", err)
	}
	ib.DoSomeTest()
	time.Sleep(time.Second * 20)
	if err := ib.DisConnect(); err != nil {
		log.Panicf("DisConnect failed: %v", err)
	}
	// time.Sleep(time.Second * 10)
	// if err := ib.Connect(); err != nil {
	// 	log.Panicf("Connect failed: %v", err)
	// }
	// ib.DoSomeTest()
	// time.Sleep(time.Second * 10)
	// if err := ib.DisConnect(); err != nil {
	// 	log.Panicf("DisConnect failed: %v", err)
	// }
}
