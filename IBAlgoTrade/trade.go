package IBAlgoTrade

import (
	"time"

	"github.com/hadrianl/ibgo/ibapi"
)

type Trade struct {
	Contract   ibapi.Contract
	Order      ibapi.Order
	OrderState ibapi.OrderState
	Fills      []Fill
	Log        []TradeLogEntry
}

type Fill struct {
	Time             time.Time
	Contract         ibapi.Contract
	Execution        ibapi.Execution
	CommissionReport ibapi.CommissionReport
}

type TradeLogEntry struct {
	Time    time.Time
	status  string
	message string
}
