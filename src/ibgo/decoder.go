package ibgo

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"
)

const (
	TIME_FORMAT string = "2006-01-02 15:04:05 +0700 CST"
)

type IbDecoder struct {
	wrapper       IbWrapper
	version       Version
	msgId2process map[IN]func([][]byte)
	errChan       chan error
}

//NewIbDecoder create a decoder to decode the fileds based on version
func NewIbDecoder(wrapper IbWrapper, version Version) *IbDecoder {
	decoder := IbDecoder{}
	decoder.wrapper = wrapper
	decoder.version = version
	decoder.errChan = make(chan error, 30)
	return &decoder
}

func (d *IbDecoder) setVersion(version Version) {
	d.version = version
}

func (d *IbDecoder) interpret(fs ...[]byte) {
	if len(fs) == 0 {
		return
	}

	// if decode error ocours,handle the error
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("!!!!!!errDeocde!!!!!!->%v", err)
		}
	}()

	MsgId, _ := strconv.ParseInt(string(fs[0]), 10, 64)
	if processer, ok := d.msgId2process[IN(MsgId)]; ok {
		processer(fs[1:])
	} else {
		log.Printf("MsgId: %v -> MsgBuf: %v", MsgId, fs[1:])
	}

}

// func (d *IbDecoder) interpretWithSignature(fs [][]byte, processer interface{}) {
// 	if processer == nil {
// 		fmt.Println("No processer")
// 	}

// 	processerType := reflect.TypeOf(processer)
// 	params := make([]interface{}, processerType.NumIn())
// 	for i, f := range fs[1:] {
// 		switch processerType.In(i).Kind() {
// 		case reflect.Int:
// 			param := strconv.Atoi(string(f))
// 		case reflect.Float64:
// 			param, _ := strconv.ParseFloat(string(f), 64)
// 		default:
// 			param := string(f)
// 		}
// 		params = append(params, param)

// 	}

// 	processer(params...)
// }

func (d *IbDecoder) setmsgId2process() {
	d.msgId2process = map[IN]func([][]byte){
		TICK_PRICE:       d.processTickPriceMsg,
		TICK_SIZE:        d.wrapTickSize,
		ORDER_STATUS:     d.processOrderStatusMsg,
		ERR_MSG:          d.wrapError,
		OPEN_ORDER:       d.processOpenOrder,
		ACCT_VALUE:       d.wrapUpdateAccountValue,
		PORTFOLIO_VALUE:  d.processPortfolioValueMsg,
		ACCT_UPDATE_TIME: d.wrapUpdateAccountTime,
		NEXT_VALID_ID:    d.wrapNextValidId,
		MANAGED_ACCTS:    d.wrapManagedAccounts,

		ACCT_DOWNLOAD_END: d.wrapAccountDownloadEnd,
		CURRENT_TIME:      d.wrapCurrentTime,
	}

}

func (d *IbDecoder) wrapTickSize(f [][]byte) {
	reqId := decodeInt(f[1])
	tickType := decodeInt(f[2])
	size := decodeInt(f[3])
	d.wrapper.tickSize(reqId, tickType, size)
}

func (d *IbDecoder) wrapNextValidId(f [][]byte) {
	reqId := decodeInt(f[1])
	d.wrapper.nextValidId(reqId)

}

func (d *IbDecoder) wrapManagedAccounts(f [][]byte) {
	// accNames := strings.Split(string(f[1]), ",")
	accNameField := bytes.Split(f[1], []byte{','})

	accsList := []Account{}
	for _, acc := range accNameField {
		accsList = append(accsList, Account{Name: string(acc)})
	}
	d.wrapper.managedAccounts(accsList)

}

func (d *IbDecoder) wrapUpdateAccountValue(f [][]byte) {
	tag := decodeString(f[1])
	val := decodeString(f[2])
	currency := decodeString(f[3])
	accName := decodeString(f[4])

	d.wrapper.updateAccountValue(tag, val, currency, accName)
}

func (d *IbDecoder) wrapUpdateAccountTime(f [][]byte) {
	ts := string(f[1])
	today := time.Now()
	// time.
	t, err := time.ParseInLocation("04:05", ts, time.Local)
	if err != nil {
		panic(err)
	}
	t = t.AddDate(today.Year(), int(today.Month())-1, today.Day()-1)

	d.wrapper.updateAccountTime(t)
}

func (d *IbDecoder) wrapError(f [][]byte) {
	reqId := decodeInt(f[1])
	errorCode := decodeInt(f[2])
	errorString := decodeString(f[3])

	d.wrapper.error(reqId, errorCode, errorString)
}

func (d *IbDecoder) wrapCurrentTime(f [][]byte) {
	ts := decodeInt(f[1])
	t := time.Unix(ts, 0)

	d.wrapper.currentTime(t)
}

//--------------wrap end func ---------------------------------

func (d *IbDecoder) wrapAccountDownloadEnd(f [][]byte) {
	accName := string(f[1])

	d.wrapper.accountDownloadEnd(accName)
}

// ------------------------------------------------------------------

func (d *IbDecoder) processTickPriceMsg(f [][]byte) {
	reqId := decodeInt(f[1])
	tickType := decodeInt(f[2])
	price := decodeFloat(f[3])
	size := decodeInt(f[4])
	attrMask := decodeInt(f[5])

	attrib := TickAttrib{}
	attrib.CanAutoExecute = attrMask == 1

	if d.version >= MIN_SERVER_VER_PAST_LIMIT {
		attrib.CanAutoExecute = attrMask&0x1 != 0
		attrib.PastLimit = attrMask&0x2 != 0
		if d.version >= MIN_SERVER_VER_PRE_OPEN_BID_ASK {
			attrib.PreOpen = attrMask&0x4 != 0
		}
	}

	d.wrapper.tickPrice(reqId, tickType, price, attrib)

	var sizeTickType int64
	switch tickType {
	case BID:
		sizeTickType = BID_SIZE
	case ASK:
		sizeTickType = ASK_SIZE
	case LAST:
		sizeTickType = LAST_SIZE
	case DELAYED_BID:
		sizeTickType = DELAYED_BID_SIZE
	case DELAYED_ASK:
		sizeTickType = DELAYED_ASK_SIZE
	case DELAYED_LAST:
		sizeTickType = DELAYED_LAST_SIZE
	default:
		sizeTickType = NOT_SET
	}

	if sizeTickType != NOT_SET {
		d.wrapper.tickSize(reqId, sizeTickType, size)
	}

}

func (d *IbDecoder) processOrderStatusMsg(f [][]byte) {
	if d.version < MIN_SERVER_VER_MARKET_CAP_PRICE {
		f = f[1:]
	}
	orderId := decodeInt(f[0])
	status := decodeString(f[1])

	filled := decodeFloat(f[2])

	remaining := decodeFloat(f[3])

	avgFilledPrice := decodeFloat(f[4])

	permId := decodeInt(f[5])
	parentId := decodeInt(f[6])
	lastFillPrice := decodeFloat(f[7])
	clientId := decodeInt(f[8])
	whyHeld := decodeString(f[9])

	var mktCapPrice float64
	if d.version >= MIN_SERVER_VER_MARKET_CAP_PRICE {
		mktCapPrice = decodeFloat(f[10])
	} else {
		mktCapPrice = float64(0)
	}

	d.wrapper.orderStatus(orderId, status, filled, remaining, avgFilledPrice, permId, parentId, lastFillPrice, clientId, whyHeld, mktCapPrice)

}

func (d *IbDecoder) processOpenOrder(f [][]byte) {

	log.Println("processOpenOrders:", f)
	var version int64
	if d.version < MIN_SERVER_VER_ORDER_CONTAINER {
		version = decodeInt(f[0])
		f = f[1:]
	} else {
		version = int64(d.version)
	}

	o := &Order{}
	o.OrderId = decodeInt(f[0])

	c := &Contract{}

	c.ContractId = decodeInt(f[1])
	c.Symbol = decodeString(f[2])
	c.SecurityType = decodeString(f[3])
	if t, err := time.Parse(TIME_FORMAT, decodeString(f[4])); err == nil {
		c.Expiry = t
	}

	c.Strike = decodeFloat(f[5])
	c.Right = decodeString(f[6])

	if version >= 32 {
		c.Multiplier = decodeString(f[7])
		f = f[1:]
	}
	c.Exchange = decodeString(f[7])
	c.Currency = decodeString(f[8])
	c.LocalSymbol = decodeString(f[9])
	if version >= 32 {
		c.TradingClass = decodeString(f[10])
		f = f[1:]
	}

	o.Action = decodeString(f[10])
	o.TotalQuantity = decodeFloat(f[11])
	o.OrderType = decodeString(f[12])
	o.LmtPrice = decodeFloat(f[13]) //todo: show_unset
	o.AuxPrice = decodeFloat(f[14]) //todo: show_unset
	o.Tif = decodeString(f[15])
	o.OCAGroup = decodeString(f[16])
	o.Account = decodeString(f[17])
	o.OpenClose = decodeString(f[18])

	o.Origin = decodeInt(f[19])

	o.OrderRef = decodeString(f[20])
	o.ClientId = decodeInt(f[21])
	o.PermId = decodeInt(f[22])

	o.OutsideRTH = decodeBool(f[23])
	o.Hidden = decodeBool(f[24])
	o.DiscretionaryAmount = decodeFloat(f[25])
	o.GoodAfterTime = decodeTime(f[26], "20060102")

	_ = decodeString(f[27]) //_sharesAllocation

	o.FAGroup = decodeString(f[28])
	o.FAMethod = decodeString(f[29])
	o.FAPercentage = decodeString(f[30])
	o.FAProfile = decodeString(f[31])
	o.GoodTillDate = decodeTime(f[32], "20060102")

	o.Rule80A = decodeString(f[33])
	o.PercentOffset = decodeFloat(f[34]) //show_unset
	o.SettlingFirm = decodeString(f[35])
	o.ShortSaleSlot = decodeInt(f[36])
	o.DesignatedLocation = decodeString(f[37])

	if d.version == MIN_SERVER_VER_SSHORTX_OLD {
		f = f[1:]
	} else if version >= 23 {
		o.ExemptCode = decodeInt(f[38])
		f = f[1:]
	}

	o.AuctionStrategy = decodeInt(f[38])
	o.StartingPrice = decodeFloat(f[39])   //show_unset
	o.StockRefPrice = decodeFloat(f[40])   //show_unset
	o.Delta = decodeFloat(f[41])           //show_unset
	o.StockRangeLower = decodeFloat(f[42]) //show_unset
	o.StockRangeUpper = decodeFloat(f[43]) //show_unset
	o.DisplaySize = decodeInt(f[44])

	o.BlockOrder = decodeBool(f[45])
	o.SweepToFill = decodeBool(f[46])
	o.AllOrNone = decodeBool(f[47])
	o.MinQty = decodeInt(f[48]) //show_unset
	o.OCAType = decodeInt(f[49])
	o.ETradeOnly = decodeBool(f[50])
	o.FirmQuoteOnly = decodeBool(f[51])
	o.NBBOPriceCap = decodeFloat(f[52]) //show_unset

	o.ParentID = decodeInt(f[53])
	o.TriggerMethod = decodeInt(f[54])

	o.Volatility = decodeFloat(f[55]) //show_unset
	o.VolatilityType = decodeInt(f[56])
	o.DeltaNeutralOrderType = decodeString(f[57])
	o.DeltaNeutralAuxPrice = decodeFloat(f[58])

	if version >= 27 && o.DeltaNeutralOrderType != "" {
		o.DeltaNeutralConId = decodeInt(f[59])
		o.DeltaNeutralSettlingFirm = decodeString(f[60])
		o.DeltaNeutralClearingAccount = decodeString(f[61])
		o.DeltaNeutralClearingIntent = decodeString(f[62])
		f = f[4:]
	}

	if version >= 31 && o.DeltaNeutralOrderType != "" {
		o.DeltaNeutralOpenClose = decodeString(f[59])
		o.DeltaNeutralShortSale = decodeBool(f[60])
		o.DeltaNeutralShortSaleSlot = decodeInt(f[61])
		o.DeltaNeutralDesignatedLocation = decodeString(f[62])
		f = f[4:]
	}

	o.ContinuousUpdate = decodeBool(f[59])

	o.ReferencePriceType = decodeInt(f[60])

	o.TrailStopPrice = decodeFloat(f[61])

	if version >= 30 {
		o.TrailingPercent = decodeFloat(f[62]) //show_unset
		f = f[1:]
	}

	o.BasisPoints = decodeFloat(f[62])
	o.BasisPointsType = decodeInt(f[63])
	c.ComboLegsDescription = decodeString(f[64])

	if version >= 29 {
		c.ComboLegs = []ComboLeg{}
		fmt.Println("comboLegsCount:", f[65])
		for comboLegsCount := decodeInt(f[65]); comboLegsCount > 0 && comboLegsCount != math.MaxInt64; comboLegsCount-- {
			fmt.Println("comboLegsCount:", comboLegsCount)
			comboleg := ComboLeg{}
			comboleg.ConId = decodeInt(f[66])
			comboleg.Ratio = decodeInt(f[67])
			comboleg.Action = decodeString(f[68])
			comboleg.Exchange = decodeString(f[69])
			comboleg.OpenClose = decodeInt(f[70])
			comboleg.ShortSaleSlot = decodeInt(f[71])
			comboleg.DesignatedLocation = decodeString(f[72])
			comboleg.ExemptCode = decodeInt(f[73])
			c.ComboLegs = append(c.ComboLegs, comboleg)
			f = f[8:]
		}

		o.OrderComboLegs = []OrderComboLeg{}
		for orderComboLegsCount := decodeInt(f[74]); orderComboLegsCount > 0 && orderComboLegsCount != math.MaxInt64; orderComboLegsCount-- {
			orderComboLeg := OrderComboLeg{}
			orderComboLeg.Price = decodeFloat(f[75])
			o.OrderComboLegs = append(o.OrderComboLegs, orderComboLeg)
			f = f[1:]
		}
		f = f[2:]
	}

	if version >= 26 {
		o.SmartComboRoutingParams = []TagValue{}
		for smartComboRoutingParamsCount := decodeInt(f[65]); smartComboRoutingParamsCount > 0 && smartComboRoutingParamsCount != math.MaxInt64; smartComboRoutingParamsCount-- {
			tagValue := TagValue{}
			tagValue.Tag = decodeString(f[66])
			tagValue.Value = decodeString(f[67])
			o.SmartComboRoutingParams = append(o.SmartComboRoutingParams, tagValue)
			f = f[2:]
		}

		f = f[1:]
	}

	if version >= 20 {
		o.ScaleInitLevelSize = decodeInt(f[65]) //show_unset
		o.ScaleSubsLevelSize = decodeInt(f[66]) //show_unset
	} else {
		o.NotSuppScaleNumComponents = decodeInt(f[65])
		o.ScaleInitLevelSize = decodeInt(f[66])
	}

	o.ScalePriceIncrement = decodeFloat(f[67])

	if version >= 28 && o.ScalePriceIncrement != math.MaxFloat64 && o.ScalePriceIncrement > 0.0 {
		o.ScalePriceAdjustValue = decodeFloat(f[68])
		o.ScalePriceAdjustInterval = decodeInt(f[69])
		o.ScaleProfitOffset = decodeFloat(f[70])
		o.ScaleAutoReset = decodeBool(f[71])
		o.ScaleInitPosition = decodeInt(f[72])
		o.ScaleInitFillQty = decodeInt(f[73])
		o.ScaleRandomPercent = decodeBool(f[74])
		f = f[7:]
	}

	if version >= 24 {
		o.HedgeType = decodeString(f[68])
		if o.HedgeType != "" {
			o.HedgeParam = decodeString(f[69])
			f = f[1:]
		}
		f = f[1:]
	}

	if version >= 25 {
		o.OptOutSmartRouting = decodeBool(f[68])
		f = f[1:]
	}

	if version >= 22 {
		o.NotHeld = decodeBool(f[68])
		f = f[1:]
	}

	if version >= 20 {
		deltaNeutralContractPresent := decodeBool(f[68])
		if deltaNeutralContractPresent {
			c.DeltaNeutralContract = DeltaNeutralContract{}
			c.DeltaNeutralContract.CondId = decodeInt(f[69])
			c.DeltaNeutralContract.Delta = decodeFloat(f[70])
			c.DeltaNeutralContract.Price = decodeFloat(f[71])
			f = f[3:]
		}
		f = f[1:]
	}

	if version >= 21 {
		o.AlgoStrategy = decodeString(f[68])
		if o.AlgoStrategy != "" {
			o.AlgoParams = []TagValue{}
			for algoParamsCount := decodeInt(f[69]); algoParamsCount > 0 && algoParamsCount != math.MaxInt64; algoParamsCount-- {
				tagValue := TagValue{}
				tagValue.Tag = decodeString(f[70])
				tagValue.Value = decodeString(f[71])
				o.AlgoParams = append(o.AlgoParams, tagValue)
				f = f[2:]
			}
		}
		f = f[1:]
	}

	if version >= 33 {
		o.Solictied = decodeBool(f[68])
	}

	orderState := &OrderState{}

	o.WhatIf = decodeBool(f[68])

	orderState.Status = decodeString(f[69])

	if d.version >= MIN_SERVER_VER_WHAT_IF_EXT_FIELDS {
		orderState.InitialMarginBefore = decodeString(f[70])
		orderState.MaintenanceMarginBefore = decodeString(f[71])
		orderState.EquityWithLoanBefore = decodeString(f[72])
		orderState.InitialMarginChange = decodeString(f[73])
		orderState.MaintenanceMarginChange = decodeString(f[74])
		orderState.EquityWithLoanChange = decodeString(f[75])
		f = f[6:]
	}

	orderState.InitialMarginAfter = decodeString(f[70])
	orderState.MaintenanceMarginAfter = decodeString(f[71])
	orderState.EquityWithLoanAfter = decodeString(f[72])

	orderState.Commission = decodeFloat(f[73])
	orderState.MinCommission = decodeFloat(f[74])
	orderState.MaxCommission = decodeFloat(f[75])
	orderState.WarningText = decodeString(f[76])

	if version >= 34 {
		o.RandomizeSize = decodeBool(f[77])
		o.RandomizePrice = decodeBool(f[78])
		f = f[2:]
	}

	if d.version >= MIN_SERVER_VER_PEGGED_TO_BENCHMARK {
		if o.OrderType == "PEG BENCH" {
			o.ReferenceContractId = decodeInt(f[77])
			o.IsPeggedChangeAmountDecrease = decodeBool(f[78])
			o.PeggedChangeAmount = decodeFloat(f[79])
			o.ReferenceChangeAmount = decodeFloat(f[80])
			o.ReferenceExchangeId = decodeString(f[81])
			f = f[5:]
		}

		o.Conditions = []OrderCondition{}
		for conditionsSize := decodeInt(f[77]); conditionsSize > 0; conditionsSize-- {
			tagValue := TagValue{}
			tagValue.Tag = decodeString(f[78])
			tagValue.Value = decodeString(f[79])
			o.AlgoParams = append(o.AlgoParams, tagValue)
			f = f[2:]
		}

		o.ConditionsIgnoreRth = decodeBool(f[78])
		o.ConditionsCancelOrder = decodeBool(f[79])

		f = f[3:]
	}

	o.AdjustedOrderType = decodeString(f[77])
	o.TriggerPrice = decodeFloat(f[78])
	o.TrailStopPrice = decodeFloat(f[79])
	o.LmtPriceOffset = decodeFloat(f[80])
	o.AdjustedStopPrice = decodeFloat(f[81])
	o.AdjustedTrailingAmount = decodeFloat(f[82])
	o.AdjustableTrailingUnit = decodeInt(f[83])

	if d.version >= MIN_SERVER_VER_SOFT_DOLLAR_TIER {
		name := decodeString(f[84])
		value := decodeString(f[85])
		displayName := decodeString(f[86])
		o.SoftDollarTier = SoftDollarTier{name, value, displayName}
		f = f[3:]
	}

	if d.version >= MIN_SERVER_VER_CASH_QTY {
		o.CashQty = decodeFloat(f[84])
		f = f[1:]
	}

	if d.version >= MIN_SERVER_VER_AUTO_PRICE_FOR_HEDGE {
		o.DontUseAutoPriceForHedge = decodeBool(f[84])
		f = f[1:]
	}

	if d.version >= MIN_SERVER_VER_ORDER_CONTAINER {
		o.IsOmsContainer = decodeBool(f[84])
		f = f[1:]
	}

	if d.version >= MIN_SERVER_VER_D_PEG_ORDERS {
		o.DiscretionaryUpToLimitPrice = decodeBool(f[84])
		f = f[1:]
	}

	d.wrapper.openOrder(o.OrderId, c, o, orderState)

}

func (d *IbDecoder) processPortfolioValueMsg(f [][]byte) {
	v := decodeInt(f[0])

	c := &Contract{}
	c.ContractId = decodeInt(f[1])
	c.Symbol = decodeString(f[2])
	c.SecurityType = decodeString(f[3])
	c.Expiry = decodeDate(f[4])
	c.Strike = decodeFloat(f[5])
	c.Right = decodeString(f[6])

	if v >= 7 {
		c.Multiplier = decodeString(f[7])
		c.PrimaryExchange = decodeString(f[8])
		f = f[2:]
	}

	c.Currency = decodeString(f[7])
	c.LocalSymbol = decodeString(f[8])

	if v >= 8 {
		c.TradingClass = decodeString(f[9])
		f = f[1:]
	}

	position := decodeFloat(f[9])
	marketPrice := decodeFloat(f[10])
	marketValue := decodeFloat(f[11])
	averageCost := decodeFloat(f[12])
	unrealizedPNL := decodeFloat(f[13])
	realizedPNL := decodeFloat(f[14])
	accName := decodeString(f[15])

	if v == 6 && d.version == 39 {
		c.PrimaryExchange = decodeString(f[16])
	}

	d.wrapper.updatePortfolio(c, position, marketPrice, marketValue, averageCost, unrealizedPNL, realizedPNL, accName)

}
func (d *IbDecoder) processContractDataMsg(f [][]byte) {

}
func (d *IbDecoder) processBondContractDataMsg(f [][]byte) {

}
func (d *IbDecoder) processScannerDataMsg(f [][]byte) {

}
func (d *IbDecoder) processExecutionDataMsg(f [][]byte) {

}
func (d *IbDecoder) processHistoricalDataMsg(f [][]byte) {

}
func (d *IbDecoder) processHistoricalDataUpdateMsg(f [][]byte) {

}
func (d *IbDecoder) processRealTimeBarMsg(f [][]byte) {

}
func (d *IbDecoder) processTickOptionComputationMsg(f [][]byte) {

}

func (d *IbDecoder) processDeltaNeutralValidationMsg(f [][]byte) {

}
func (d *IbDecoder) processMarketDataTypeMsg(f [][]byte) {

}
func (d *IbDecoder) processCommissionReportMsg(f [][]byte) {

}
func (d *IbDecoder) processPositionDataMsg(f [][]byte) {

}
func (d *IbDecoder) processPositionMultiMsg(f [][]byte) {

}
func (d *IbDecoder) processSecurityDefinitionOptionParameterMsg(f [][]byte) {

}
func (d *IbDecoder) processSecurityDefinitionOptionParameterEndMsg(f [][]byte) {

}
func (d *IbDecoder) processSoftDollarTiersMsg(f [][]byte) {

}
func (d *IbDecoder) processFamilyCodesMsg(f [][]byte) {

}
func (d *IbDecoder) processSymbolSamplesMsg(f [][]byte) {

}
func (d *IbDecoder) processSmartComponents(f [][]byte) {

}
func (d *IbDecoder) processTickReqParams(f [][]byte) {

}
func (d *IbDecoder) processMktDepthExchanges(f [][]byte) {

}

func (d *IbDecoder) processHeadTimestamp(f [][]byte) {

}
func (d *IbDecoder) processTickNews(f [][]byte) {

}
func (d *IbDecoder) processNewsProviders(f [][]byte) {

}
func (d *IbDecoder) processNewsArticle(f [][]byte) {

}
func (d *IbDecoder) processHistoricalNews(f [][]byte) {

}
func (d *IbDecoder) processHistoricalNewsEnd(f [][]byte) {

}
func (d *IbDecoder) processHistogramData(f [][]byte) {

}
func (d *IbDecoder) processRerouteMktDataReq(f [][]byte) {

}
func (d *IbDecoder) processRerouteMktDepthReq(f [][]byte) {

}
func (d *IbDecoder) processMarketRuleMsg(f [][]byte) {

}
func (d *IbDecoder) processPnLMsg(f [][]byte) {

}
func (d *IbDecoder) processPnLSingleMsg(f [][]byte) {

}
func (d *IbDecoder) processHistoricalTicks(f [][]byte) {

}
func (d *IbDecoder) processHistoricalTicksBidAsk(f [][]byte) {

}
func (d *IbDecoder) processHistoricalTicksLast(f [][]byte) {

}
func (d *IbDecoder) processTickByTickMsg(f [][]byte) {

}
func (d *IbDecoder) processOrderBoundMsg(f [][]byte) {

}
func (d *IbDecoder) processMarketDepthL2Msg(f [][]byte) {

}

// ----------------------------------------------------
