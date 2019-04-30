package IBAlgoTrade

import (
	"time"

	"github.com/hadrianl/ibgo/ibapi"
)

type Tick struct {
	Price float64
	Size  int64
}

type Ticker struct {
	Contract ibapi.Contract
	Time     time.Time
	Bid      Tick
	Ask      Tick
	Last     Tick
	PrevBid  Tick
	PrevAsk  Tick
	PrevLast Tick
	Volume   int64
	Open     float64
	High     float64
	Low      float64
	Close    float64
	Vwap     float64
}
