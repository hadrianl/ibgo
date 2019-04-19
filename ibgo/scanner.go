package ibgo

type ScanData struct {
	ContractDetails ContractDetails
	Rank            int64
	Distance        string
	Benchmark       string
	Projection      string
	Legs            string
}

type ScannerSubscription struct {
	NumberOfRows             int
	Instrument               string
	LocationCode             string
	ScanCode                 string
	AbovePrice               float64
	BelowPrice               float64
	AboveVolume              int64
	MarketCapAbove           float64
	MarketCapBelow           float64
	MoodyRatingAbove         string
	MoodyRatingBelow         string
	SpRatingAbove            string
	SpRatingBelow            string
	MaturityDateAbove        string
	MaturityDateBelow        string
	CouponRateAbove          float64
	CouponRateBelow          float32
	ExcludeConvertible       bool
	AverageOptionVolumeAbove int64
	ScannerSettingPairs      string
	StockTypeFilter          string
}
