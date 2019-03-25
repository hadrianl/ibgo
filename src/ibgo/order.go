package ibgo

import "time"

type Order struct {
	OrderId                       int64
	ClientId                      int64
	PermId                        int64
	Action                        string
	TotalQuantity                 float64
	OrderType                     string
	LmtPrice                      float64
	AuxPrice                      float64
	Tif                           string
	ActiveStartTime               time.Time
	ActiveStopTime                time.Time
	OCAGroup                      string
	OCAType                       int64 // 1 = CANCEL_WITH_BLOCK, 2 = REDUCE_WITH_BLOCK, 3 = REDUCE_NON_BLOCK
	OrderRef                      string
	Transmit                      bool
	ParentID                      int64
	BlockOrder                    bool
	SweepToFill                   bool
	DisplaySize                   int64
	TriggerMethod                 int64 // 0=Default, 1=Double_Bid_Ask, 2=Last, 3=Double_Last, 4=Bid_Ask, 7=Last_or_Bid_Ask, 8=Mid-point
	OutsideRTH                    bool
	Hidden                        bool
	GoodAfterTime                 time.Time
	GoodTillDate                  time.Time
	OverridePercentageConstraints bool
	Rule80A                       string // Individual = 'I', Agency = 'A', AgentOtherMember = 'W', IndividualPTIA = 'J', AgencyPTIA = 'U', AgentOtherMemberPTIA = 'M', IndividualPT = 'K', AgencyPT = 'Y', AgentOtherMemberPT = 'N'
	AllOrNone                     bool
	MinQty                        int64
	PercentOffset                 float64
	TrailStopPrice                float64
	TrailingPercent               float64
	//---- financial advisors only -----
	FAGroup      string
	FAProfile    string
	FAMethod     string
	FAPercentage string
	// ---------------------------------
	// ---------institutional only------
	OpenClose          string // O=Open, C=Close
	Origin             int64  // 0=Customer, 1=Firm
	ShortSaleSlot      int64  // 1 if you hold the shares, 2 if they will be delivered from elsewhere.  Only for Action=SSHORT
	DesignatedLocation string // used only when shortSaleSlot=2
	ExemptCode         int64
	// ---------------------------------
	// ------- SMART routing only ------
	DiscretionaryAmount float64
	ETradeOnly          bool
	FirmQuoteOnly       bool
	NBBOPriceCap        float64
	OptOutSmartRouting  bool
	// --------------------------------
	// ---BOX exchange orders only ----
	AuctionStrategy int64
	StartingPrice   float64
	StockRefPrice   float64
	Delta           float64
	// --------------------------------
	// --pegged to stock and VOL orders only--
	StockRangeLower float64
	StockRangeUpper float64

	RandomizePrice bool
	RandomizeSize  bool

	// ---VOLATILITY ORDERS ONLY--------
	Volatility                     float64
	VolatilityType                 int64
	DeltaNeutralOrderType          string
	DeltaNeutralAuxPrice           float64
	DeltaNeutralConId              int64
	DeltaNeutralSettlingFirm       string
	DeltaNeutralClearingAccount    string
	DeltaNeutralClearingIntent     string
	DeltaNeutralOpenClose          string
	DeltaNeutralShortSale          bool
	DeltaNeutralShortSaleSlot      int64
	DeltaNeutralDesignatedLocation string
	ContinuousUpdate               bool
	ReferencePriceType             int64 // 1=Average, 2 = BidOrAsk
	// DeltaNeutral                  DeltaNeutralData `when:"DeltaNeutralOrderType" cond:"is" value:""`
	// -------------------------------------
	// ------COMBO ORDERS ONLY-----------
	BasisPoints     float64 // EFP orders only
	BasisPointsType int64   // EFP orders only
	// -----------------------------------
	//-----------SCALE ORDERS ONLY------------
	ScaleInitLevelSize        int64
	ScaleSubsLevelSize        int64
	ScalePriceIncrement       float64
	ScalePriceAdjustValue     float64
	ScalePriceAdjustInterval  int64
	ScaleProfitOffset         float64
	ScaleAutoReset            bool
	ScaleInitPosition         int64
	ScaleInitFillQty          int64
	ScaleRandomPercent        bool
	ScaleTable                string
	NotSuppScaleNumComponents int64
	//--------------------------------------
	// ---------HEDGE ORDERS--------------
	HedgeType  string
	HedgeParam string
	//--------------------------------------
	//-----------Clearing info ----------------
	Account         string
	SettlingFirm    string
	ClearingAccount string // True beneficiary of the order
	ClearingIntent  string // "" (Default), "IB", "Away", "PTA" (PostTrade)
	// ----------------------------------------
	// --------- ALGO ORDERS ONLY --------------
	AlgoStrategy string

	AlgoParams              []TagValue
	SmartComboRoutingParams []TagValue
	AlgoId                  string
	// -----------------------------------------

	// ----------what if order -------------------
	WhatIf bool

	// --------------Not Held ------------------
	NotHeld   bool
	Solictied bool
	//--------------------------------------

	// ------order combo legs -----------------
	OrderComboLegs   []OrderComboLeg
	OrderMiscOptions []TagValue
	//----------------------------------------
	//-----------VER PEG2BENCH fields----------
	ReferenceContractId          int64
	PeggedChangeAmount           float64
	IsPeggedChangeAmountDecrease bool
	ReferenceChangeAmount        float64
	ReferenceExchangeId          string
	AdjustedOrderType            string
	TriggerPrice                 float64
	AdjustedStopPrice            float64
	AdjustedStopLimitPrice       float64
	AdjustedTrailingAmount       float64
	AdjustableTrailingUnit       int64
	LmtPriceOffset               float64

	Conditions            []OrderCondition
	ConditionsCancelOrder bool
	ConditionsIgnoreRth   bool

	//------ext operator--------------
	ExtOperator string

	//-----native cash quantity --------
	CashQty float64

	//--------------------------------
	Mifid2DecisionMaker   string
	Mifid2DecisionAlgo    string
	Mifid2ExecutionTrader string
	Mifid2ExecutionAlgo   string

	//-------------
	DontUseAutoPriceForHedge bool

	IsOmsContainer bool

	DiscretionaryUpToLimitPrice bool

	SoftDollarTier SoftDollarTier
}

type OrderState struct {
	Status                  string
	InitialMarginBefore     string
	InitialMarginChange     string
	InitialMarginAfter      string
	MaintenanceMarginBefore string
	MaintenanceMarginChange string
	MaintenanceMarginAfter  string
	EquityWithLoanBefore    string
	EquityWithLoanChange    string
	EquityWithLoanAfter     string
	Commission              float64 // max
	MinCommission           float64 // max
	MaxCommission           float64 // max
	CommissionCurrency      string
	WarningText             string
}

type OrderCondition struct {
	CondType                int64
	IsConjunctionConnection bool

	// Price = 1
	// Time = 3
	// Margin = 4
	// Execution = 5
	// Volume = 6
	// PercentChange = 7
}

type SoftDollarTier struct {
	Name        string
	Value       string
	DisplayName string
}
