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
	CouponRateBelow          float64
	ExcludeConvertible       bool
	AverageOptionVolumeAbove int64
	ScannerSettingPairs      string
	StockTypeFilter          string
}

func NewScanData(contractDetails ContractDetails, rank int64, distance string, benchmark string, projection string, legsStr string) *ScanData {
	scanData := &ScanData{contractDetails, rank, distance, benchmark, projection, legsStr}
	return scanData
}

func NewScannerSubscription() *ScannerSubscription {
	scannerSubscription := &ScannerSubscription{}

	scannerSubscription.NumberOfRows = -1
	scannerSubscription.AbovePrice = UNSETFLOAT
	scannerSubscription.BelowPrice = UNSETFLOAT
	scannerSubscription.AboveVolume = UNSETINT
	scannerSubscription.MarketCapAbove = UNSETFLOAT
	scannerSubscription.MarketCapBelow = UNSETFLOAT

	scannerSubscription.CouponRateAbove = UNSETFLOAT
	scannerSubscription.CouponRateBelow = UNSETFLOAT
	scannerSubscription.AverageOptionVolumeAbove = UNSETINT

	return scannerSubscription
}
