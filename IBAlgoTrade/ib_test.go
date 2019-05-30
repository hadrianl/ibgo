package IBAlgoTrade_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hadrianl/ibgo/IBAlgoTrade"
	log "github.com/sirupsen/logrus"
	// "github.com/hadrianl/ibgo/ibapi"
)

func TestIB(t *testing.T) {
	log.SetLevel(log.InfoLevel)
	ib := IBAlgoTrade.NewIB("192.168.2.226", 4002, 0)
	if err := ib.Connect(); err != nil {
		log.Panicf("Connect failed: %v", err)
	}
	ib.DoSomeTest()
	fmt.Println(ib.Wrapper.AccSummary)
	fmt.Println(ib.Wrapper.AccValues)
	fmt.Println(ib.Wrapper.Portfolio)
	go func() {
		for {
			for _, pnl := range ib.Wrapper.PnLs {
				fmt.Print(*pnl)
			}
			time.Sleep(5 * time.Second)
		}

	}()
	time.Sleep(time.Second * 60 * 30)
	if err := ib.DisConnect(); err != nil {
		log.Panicf("DisConnect failed: %v", err)
	}
}
