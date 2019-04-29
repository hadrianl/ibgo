package IBAlgoTrade

import "github.com/hadrianl/ibgo/ibapi"

const (
	PendingSubmit = "PendingSubmit"
	PendingCancel = "PendingCancel"
	PreSubmitted  = "PreSubmitted"
	Submitted     = "Submitted"
	ApiPending    = "ApiPending"
	ApiCancelled  = "ApiCancelled"
	Cancelled     = "Cancelled"
	Filled        = "Filled"
	Inactive      = "Inactive"
	// DoneStates    = [3]string{"Filled", "Cancelled", "ApiCancelled"}
	// ActiveStates  = [4]string{"PendingSubmit", "ApiPending", "PreSubmitted", "Submitted"}
)

type Trade struct {
	Contract   ibapi.Contract
	Order      ibapi.Order
	OrderState OrderStatus
	Fills      []Fill
	Log        []TradeLogEntry
}

type OrderStatus struct {
	OrderID           int64
	Status            string
	Filled            int64
	Remaining         int64
	AverageFillPrice  float64
	PermID            int64
	ParentID          int64
	LastFillPrice     float64
	ClientID          int64
	WhyHeld           string
	MarketCappedPrice float64
	LastLiquidity     int64
}

func NewOrder() *ibapi.Order {
	var order *ibapi.Order
	ibapi.InitDefault(order)

	return order
}
