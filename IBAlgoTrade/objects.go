package IBAlgoTrade

import (
	"time"

	"github.com/hadrianl/ibgo/ibapi"
)

// func ContractDetails(Contract ibapi.Contract) *ibapi.ContractDetails {
// 	return &ContractDetails{Contract: *Contract()}
// }

// func ContractDescription(Contract ibapi.Contract) *ibapi.ContractDescription {
// 	return &ContractDescription{Contract: *Contract()}
// }

// type ComboLeg struct {
// 	ContractID int64
// 	Ratio      int64
// 	Action     string
// 	Exchange   string
// 	OpenClose  int64

// 	// for stock legs when doing short sale
// 	ShortSaleSlot      int64
// 	DesignatedLocation string
// 	ExemptCode         int64 `default:"-1"`
// }

// type DeltaNeutralContract struct {
// 	ibapi.DeltaNeutralContract
// }

// type OrderComboLeg struct {
// 	Price float64 `default:"UNSETFLOAT"`
// }

// type OrderState struct {
// 	Status                  string
// 	InitialMarginBefore     string
// 	InitialMarginChange     string
// 	InitialMarginAfter      string
// 	MaintenanceMarginBefore string
// 	MaintenanceMarginChange string
// 	MaintenanceMarginAfter  string
// 	EquityWithLoanBefore    string
// 	EquityWithLoanChange    string
// 	EquityWithLoanAfter     string
// 	Commission              float64 `default:"UNSETFLOAT"`
// 	MinCommission           float64 `default:"UNSETFLOAT"`
// 	MaxCommission           float64 `default:"UNSETFLOAT"`
// 	CommissionCurrency      string
// 	WarningText             string
// }

// type ScannerSubscription struct {
// 	NumberOfRows             int
// 	Instrument               string
// 	LocationCode             string
// 	ScanCode                 string
// 	AbovePrice               float64 `default:"UNSETFLOAT"`
// 	BelowPrice               float64 `default:"UNSETFLOAT"`
// 	AboveVolume              int64   `default:"UNSETINT"`
// 	MarketCapAbove           float64 `default:"UNSETFLOAT"`
// 	MarketCapBelow           float64 `default:"UNSETFLOAT"`
// 	MoodyRatingAbove         string
// 	MoodyRatingBelow         string
// 	SpRatingAbove            string
// 	SpRatingBelow            string
// 	MaturityDateAbove        string
// 	MaturityDateBelow        string
// 	CouponRateAbove          float64 `default:"UNSETFLOAT"`
// 	CouponRateBelow          float64 `default:"UNSETFLOAT"`
// 	ExcludeConvertible       bool
// 	AverageOptionVolumeAbove int64 `default:"UNSETINT"`
// 	ScannerSettingPairs      string
// 	StockTypeFilter          string
// }

// type SoftDollarTier struct {
// 	ibapi.SoftDollarTier
// }

// func Execution() *ibapi.Execution {
// 	return &ibapi.Execution{}
// }

// func CommissionReport() *ibapi.CommissionReport {
// 	return &ibapi.CommissionReport{}
// }

// func ExecutionFilter() *ibapi.ExecutionFilter {
// 	return &ibapi.ExecutionFilter{}
// }

// func BarData() *ibapi.BarData {
// 	return &ibapi.BarData{}
// }

// func RealTimeBar() *ibapi.RealTimeBar {
// 	return &ibapi.RealTimeBar{}
// }

// func TickAttrib() *ibapi.TickAttrib {
// 	return &ibapi.TickAttrib{}
// }

// func TickAttribBidAsk() *ibapi.TickAttribBidAsk {
// 	return &ibapi.TickAttribBidAsk{}
// }

// func TickAttribLast() *ibapi.TickAttribLast {
// 	return &ibapi.TickAttribLast{}
// }

// func HistogramData() *ibapi.HistogramData {
// 	return &ibapi.HistogramData{}
// }

// func NewsProvider() *ibapi.NewsProvider {
// 	return &ibapi.NewsProvider{}
// }

// func DepthMktDataDescription() *ibapi.DepthMktDataDescription {
// 	var depthDescrip *ibapi.DepthMktDataDescription
// 	ibapi.InitDefault(depthDescrip)
// 	return depthDescrip
// }

// type PnL struct {
// 	Account      string
// 	ModelCode    string
// 	DailyPnL     float64
// 	UnrealizePnL float64
// 	RealizePnL   float64
// }

// type PnLSingle struct {
// 	ContractID int64
// 	Position   int64
// 	value      float64
// 	PnL        PnL
// }

// type BarList []interface{}

// type BarDataList struct {
// 	BarList
// 	ReqID          int64
// 	Contract       Contract
// 	EndDateTime    string
// 	Duration       string
// 	BarSizeSetting string
// 	WhatToShow     string
// 	useRTH         bool
// 	FormatDate     int
// 	KeepUpToDate   bool
// 	ChartOptions   []ibapi.TagValue
// }

// type RealTimeBarList struct {
// 	BarList
// 	ReqID               int64
// 	Contract            Contract
// 	EndDateTime         string
// 	BarSize             string
// 	WhatToShow          string
// 	useRTH              bool
// 	RealTimeBarsOptions []ibapi.TagValue
// }

// type ScanDataList struct {
// 	ReqID                            int64
// 	Subscription                     ScannerSubscription
// 	ScannerSubscriptionOptions       []ibapi.TagValue
// 	ScannerSubscriptionFilterOptions []ibapi.TagValue
// }

// type AccountValue struct {
// 	Account   string
// 	TagValue  ibapi.TagValue
// 	Currency  string
// 	ModelCode string
// }

// type Tick struct {
// 	Price float64
// 	Size  int64
// }

// type TickData struct {
// 	Time     time.Time
// 	TickType string
// 	Tick     Tick
// }

// func HistoricalTick() *ibapi.HistoricalTick {
// 	return &ibapi.HistoricalTick{}
// }

// func HistoricalTickBidAsk() *ibapi.HistoricalTickBidAsk {
// 	return &ibapi.HistoricalTickBidAsk{}
// }

// func HistoricalTickLast() *ibapi.HistoricalTickLast {
// 	return &ibapi.HistoricalTickLast{}
// }

// type TickByTickAllLast struct {
// 	time              time.Time
// 	TickType          string
// 	Tick              Tick
// 	TickAttribLast    TickAttribLast
// 	Exchange          string
// 	SpecialConditions []ibapi.OrderConditioner
// }

// type TickByTickBidAsk struct {
// 	time             time.Time
// 	Bid              Tick
// 	Ask              Tick
// 	TickAttribBidAsk ibapi.TickAttribBidAsk
// }

// type TickByTickMidPoint struct {
// 	time     time.Time
// 	MidPoint float64
// }

// type MarketDepthData struct {
// 	time        time.Time
// 	Position    int64
// 	MarketMaker string
// 	Opertation  string
// 	Side        string
// 	Tick        Tick
// }

// type DOMLevel struct {
// 	Tick        Tick
// 	MarketMaker string
// }

// type BracketOrder struct {
// 	Parent     ibapi.Order
// 	TakeProfit ibapi.Order
// 	StopLoss   ibapi.Order
// }

type TradeLogEntry struct {
	Time    time.Time
	status  string
	message string
}

// type PriceIncrement struct {
// 	ibapi.PriceIncrement
// }

// type ScanData struct {
// 	ContractDetails ContractDetails
// 	Rank            int64
// 	Distance        string
// 	Benchmark       string
// 	Projection      string
// 	Legs            string
// }

type PortfolioItem struct {
	Contract      ibapi.Contract
	Position      float64
	MarketPrice   float64
	MarketValue   float64
	AverageCost   float64
	UnrealizedPNL float64
	RealizedPNL   float64
	Account       string
}

// type Position struct {
// 	Account     string
// 	Contract    Contract
// 	Position    int64
// 	AverageCost float64
// }

type Fill struct {
	Time             time.Time
	Contract         ibapi.Contract
	Execution        ibapi.Execution
	CommissionReport ibapi.CommissionReport
}

// type OptionComputation struct {
// 	ImpliedVol      float64
// 	Delta           float64
// 	Gamma           float64
// 	Vega            float64
// 	Theta           float64
// 	PvDividend      float64
// 	OptionPrice     float64
// 	UnderlyingPrice float64
// }

// type OptionChain struct {
// 	Exchange             string
// 	UnderlyingContractID int64
// 	TradingClass         string
// 	Multiplier           string
// 	Expirations          []string
// 	Strikes              []float64
// }

// type Dividends struct {
// 	Past12Months float64
// 	Next12Months float64
// 	NextDate     string
// 	NextAmount   int64
// }

// type NewsArticle struct {
// 	ArticleType string
// 	ArticleText string
// }

// type HistoricalNews struct {
// 	time         time.Time
// 	ProviderCode string
// 	ArticleID    string
// 	Headline     string
// }

// type NewsTick struct {
// 	TimeStamp    int64
// 	ProviderCode string
// 	ArticleID    string
// 	Headline     string
// 	ExtraData    string
// }

// type NewsBulletin struct {
// 	MsgID         int64
// 	MsgType       int64
// 	Message       string
// 	OriginExchang string
// }

// type FamilyCode struct {
// 	ibapi.FamilyCode
// }

// type SmartComponent struct {
// 	ibapi.SmartComponent
// }

// type ConnectionStats struct {
// 	StartTime    time.Time
// 	Duration     int64
// 	NumBytesRecv int64
// 	NumBytesSend int64
// 	NumMsgRecv   int64
// 	NumMsgSend   int64
// }

type AccountValues struct {
	TagValues map[string][3]string
	Account   string
	// Currency  string
	// ModelCode string
}

type AccountSummary struct {
	TagValues map[string][2]string
	Account   string
	// Currency  string
	// ModelCode string
}
