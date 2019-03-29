package ibgo

type Account struct {
	Name string
}

type TickAttrib struct {
	CanAutoExecute bool
	PastLimit      bool
	PreOpen        bool
}

type AlgoParams struct {
}

type TagValue struct {
	Tag   string
	Value string
}

type OrderComboLeg struct {
	Price float64
}

// ------------ComboLeg--------------------
type ComboLegOpenClose int64
type ComboLegShortSaleSlot int64

const (
	SAME_POS       ComboLegOpenClose     = 0
	OPEN_POS                             = 1
	CLOSE_POS                            = 2
	UNKNOWN_POS                          = 3
	ClearingBroker ComboLegShortSaleSlot = 1
	ThirdParty                           = 2
)

type ComboLeg struct {
	ContractID int64
	Ratio      int64
	Action     string
	Exchange   string
	OpenClose  int64

	// for stock legs when doing short sale
	ShortSaleSlot      int64
	DesignatedLocation string
	ExemptCode         int64
}

// -----------------------------------------------------

type ExecutionFilter struct {
	ClientID     int64
	AccountCode  string
	Time         string
	Symbol       string
	SecurityType string
	Exchange     string
	Side         string
}

type BarData struct {
	Date     string
	Open     float64
	High     float64
	Low      float64
	Close    float64
	Volume   float64
	BarCount int64
	Average  float64
}

type RealTimeBar struct {
	Time    int64
	endTime int64
	Open    float64
	High    float64
	Low     float64
	Close   float64
	Volume  float64
	Wap     float64
	Count   int64
}
