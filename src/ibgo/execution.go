package ibgo

import "time"

type Execution struct {
	ExecID        string
	Time          time.Time
	AccountCode   string
	Exchange      string
	Side          string
	Shares        float64
	Price         float64
	PermID        int64
	ClientID      int64
	OrderID       int64
	Liquidation   int64
	CumQty        float64
	AveragePrice  float64
	OrderRef      string
	EVRule        string
	EVMultiplier  float64
	ModelCode     string
	LastLiquidity int64
}
