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
	ConId     int64
	Ratio     int64
	Action    string
	Exchange  string
	OpenClose int64

	// for stock legs when doing short sale
	ShortSaleSlot      int64
	DesignatedLocation string
	ExemptCode         int64
}

// -----------------------------------------------------
