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

// ibDecoder help to decode the msg buf received from TWS or Gateway
type ibDecoder struct {
	wrapper       IbWrapper
	version       Version
	msgID2process map[IN]func([][]byte)
	errChan       chan error
}

func (d *ibDecoder) setVersion(version Version) {
	d.version = version
}

func (d *ibDecoder) interpret(fs ...[]byte) {
	if len(fs) == 0 {
		return
	}

	// if decode error ocours,handle the error
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("!!!!!!errDeocde!!!!!!->%v", err) //TODO: handle error
		}
	}()

	MsgID, _ := strconv.ParseInt(string(fs[0]), 10, 64)
	if processer, ok := d.msgID2process[IN(MsgID)]; ok {
		processer(fs[1:])
	} else {
		log.Printf("MsgId: %v -> MsgBuf: %v", MsgID, fs[1:])
	}

}

// func (d *ibDecoder) interpretWithSignature(fs [][]byte, processer interface{}) {
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

func (d *ibDecoder) setmsgID2process() {
	d.msgID2process = map[IN]func([][]byte){
		TICK_PRICE:                  d.processTickPriceMsg,
		TICK_SIZE:                   d.wrapTickSize,
		ORDER_STATUS:                d.processOrderStatusMsg,
		ERR_MSG:                     d.wrapError,
		OPEN_ORDER:                  d.processOpenOrder,
		ACCT_VALUE:                  d.wrapUpdateAccountValue,
		PORTFOLIO_VALUE:             d.processPortfolioValueMsg,
		ACCT_UPDATE_TIME:            d.wrapUpdateAccountTime,
		NEXT_VALID_ID:               d.wrapNextValidID,
		CONTRACT_DATA:               d.processContractDataMsg,
		EXECUTION_DATA:              d.processExecutionDataMsg,
		MARKET_DEPTH:                d.wrapUpdateMktDepth,
		MARKET_DEPTH_L2:             d.wrapUpdateMktDepthL2,
		NEWS_BULLETINS:              d.wrapUpdateNewsBulletin,
		MANAGED_ACCTS:               d.wrapManagedAccounts,
		RECEIVE_FA:                  d.wrapReceiveFA,
		HISTORICAL_DATA:             d.processHistoricalDataMsg,
		HISTORICAL_DATA_UPDATE:      d.processHistoricalDataUpdateMsg,
		BOND_CONTRACT_DATA:          d.processBondContractDataMsg,
		SCANNER_PARAMETERS:          d.wrapScannerParameters,
		SCANNER_DATA:                d.processScannerDataMsg,
		TICK_OPTION_COMPUTATION:     d.processTickOptionComputationMsg,
		TICK_GENERIC:                d.wrapTickGeneric,
		TICK_STRING:                 d.wrapTickString,
		TICK_EFP:                    d.wrapTickEFP,
		CURRENT_TIME:                d.wrapCurrentTime,
		REAL_TIME_BARS:              d.processRealTimeBarMsg,
		ACCT_DOWNLOAD_END:           d.wrapAccountDownloadEnd,
		OPEN_ORDER_END:              d.wrapOpenOrderEnd,
		EXECUTION_DATA_END:          d.wrapExecDetailsEnd,
		DELTA_NEUTRAL_VALIDATION:    d.processDeltaNeutralValidationMsg,
		TICK_SNAPSHOT_END:           d.wrapTickSnapshotEnd,
		MARKET_DATA_TYPE:            d.wrapMarketDataType,
		COMMISSION_REPORT:           d.processCommissionReportMsg,
		POSITION_DATA:               d.processPositionDataMsg,
		POSITION_END:                d.wrapPositionEnd,
		ACCOUNT_SUMMARY:             d.wrapAccountSummary,
		ACCOUNT_SUMMARY_END:         d.wrapAccountSummaryEnd,
		VERIFY_MESSAGE_API:          d.wrapVerifyMessageAPI,
		VERIFY_COMPLETED:            d.wrapVerifyCompleted,
		DISPLAY_GROUP_LIST:          d.wrapDisplayGroupList,
		DISPLAY_GROUP_UPDATED:       d.wrapDisplayGroupUpdated,
		VERIFY_AND_AUTH_MESSAGE_API: d.wrapVerifyAndAuthMessageAPI,
	}

}

func (d *ibDecoder) wrapTickSize(f [][]byte) {
	reqID := decodeInt(f[1])
	tickType := decodeInt(f[2])
	size := decodeInt(f[3])
	d.wrapper.tickSize(reqID, tickType, size)
}

func (d *ibDecoder) wrapNextValidID(f [][]byte) {
	reqID := decodeInt(f[1])
	d.wrapper.nextValidID(reqID)

}

func (d *ibDecoder) wrapManagedAccounts(f [][]byte) {
	// accNames := strings.Split(string(f[1]), ",")
	accNameField := bytes.Split(f[1], []byte{','})

	accsList := []Account{}
	for _, acc := range accNameField {
		accsList = append(accsList, Account{Name: string(acc)})
	}
	d.wrapper.managedAccounts(accsList)

}

func (d *ibDecoder) wrapUpdateAccountValue(f [][]byte) {
	tag := decodeString(f[1])
	val := decodeString(f[2])
	currency := decodeString(f[3])
	accName := decodeString(f[4])

	d.wrapper.updateAccountValue(tag, val, currency, accName)
}

func (d *ibDecoder) wrapUpdateAccountTime(f [][]byte) {
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

func (d *ibDecoder) wrapError(f [][]byte) {
	reqID := decodeInt(f[1])
	errorCode := decodeInt(f[2])
	errorString := decodeString(f[3])

	d.wrapper.error(reqID, errorCode, errorString)
}

func (d *ibDecoder) wrapCurrentTime(f [][]byte) {
	ts := decodeInt(f[1])
	t := time.Unix(ts, 0)

	d.wrapper.currentTime(t)
}

func (d *ibDecoder) wrapUpdateMktDepth(f [][]byte) {
	reqID := decodeInt(f[1])
	position := decodeInt(f[2])
	operation := decodeInt(f[3])
	side := decodeInt(f[4])
	price := decodeFloat(f[5])
	size := decodeInt(f[6])

	d.wrapper.updateMktDepth(reqID, position, operation, side, price, size)

}

func (d *ibDecoder) wrapUpdateMktDepthL2(f [][]byte) {
	reqID := decodeInt(f[1])
	position := decodeInt(f[2])
	marketMaker := decodeString(f[3])
	operation := decodeInt(f[4])
	side := decodeInt(f[5])
	price := decodeFloat(f[6])
	size := decodeInt(f[7])
	isSmartDepth := decodeBool(f[8])

	d.wrapper.updateMktDepthL2(reqID, position, marketMaker, operation, side, price, size, isSmartDepth)

}

func (d *ibDecoder) wrapUpdateNewsBulletin(f [][]byte) {
	msgID := decodeInt(f[1])
	msgType := decodeInt(f[2])
	newsMessage := decodeString(f[3])
	originExch := decodeString(f[4])

	d.wrapper.updateNewsBulletin(msgID, msgType, newsMessage, originExch)
}

func (d *ibDecoder) wrapReceiveFA(f [][]byte) {
	faData := decodeInt(f[1])
	cxml := decodeString(f[2])

	d.wrapper.receiveFA(faData, cxml)
}

func (d *ibDecoder) wrapScannerParameters(f [][]byte) {
	xml := decodeString(f[1])

	d.wrapper.scannerParameters(xml)
}

func (d *ibDecoder) wrapTickGeneric(f [][]byte) {
	reqID := decodeInt(f[1])
	tickType := decodeInt(f[2])
	value := decodeFloat(f[3])

	d.wrapper.tickGeneric(reqID, tickType, value)

}

func (d *ibDecoder) wrapTickString(f [][]byte) {
	reqID := decodeInt(f[1])
	tickType := decodeInt(f[2])
	value := decodeString(f[3])

	d.wrapper.tickString(reqID, tickType, value)

}

func (d *ibDecoder) wrapTickEFP(f [][]byte) {
	reqID := decodeInt(f[1])
	tickType := decodeInt(f[2])
	basisPoints := decodeFloat(f[3])
	formattedBasisPoints := decodeString(f[4])
	totalDividends := decodeFloat(f[5])
	holdDays := decodeInt(f[6])
	futureLastTradeDate := decodeString(f[7])
	dividendImpact := decodeFloat(f[8])
	dividendsToLastTradeDate := decodeFloat(f[9])

	d.wrapper.tickEFP(reqID, tickType, basisPoints, formattedBasisPoints, totalDividends, holdDays, futureLastTradeDate, dividendImpact, dividendsToLastTradeDate)

}

func (d *ibDecoder) wrapMarketDataType(f [][]byte) {
	reqID := decodeInt(f[1])
	marketDataType := decodeInt(f[2])

	d.wrapper.marketDataType(reqID, marketDataType)
}

func (d *ibDecoder) wrapAccountSummary(f [][]byte) {
	reqID := decodeInt(f[1])
	account := decodeString(f[2])
	tag := decodeString(f[3])
	value := decodeString(f[4])
	currency := decodeString(f[5])

	d.wrapper.accountSummary(reqID, account, tag, value, currency)
}

func (d *ibDecoder) wrapVerifyMessageAPI(f [][]byte) {
	// Deprecated Function: keep it temporarily, not know how it works
	apiData := decodeString(f[1])

	d.wrapper.verifyMessageAPI(apiData)
}

func (d *ibDecoder) wrapVerifyCompleted(f [][]byte) {
	isSuccessful := decodeBool(f[1])
	err := decodeString(f[1])

	d.wrapper.verifyCompleted(isSuccessful, err)
}

func (d *ibDecoder) wrapDisplayGroupList(f [][]byte) {
	reqID := decodeInt(f[1])
	groups := decodeString(f[2])

	d.wrapper.displayGroupList(reqID, groups)
}

func (d *ibDecoder) wrapDisplayGroupUpdated(f [][]byte) {
	reqID := decodeInt(f[1])
	contractInfo := decodeString(f[2])

	d.wrapper.displayGroupUpdated(reqID, contractInfo)
}

func (d *ibDecoder) wrapVerifyAndAuthMessageAPI(f [][]byte) {
	apiData := decodeString(f[1])
	xyzChallange := decodeString(f[2])

	d.wrapper.verifyAndAuthMessageAPI(apiData, xyzChallange)
}

//--------------wrap end func ---------------------------------

func (d *ibDecoder) wrapAccountDownloadEnd(f [][]byte) {
	accName := string(f[1])

	d.wrapper.accountDownloadEnd(accName)
}

func (d *ibDecoder) wrapOpenOrderEnd(f [][]byte) {

	d.wrapper.openOrderEnd()
}

func (d *ibDecoder) wrapExecDetailsEnd(f [][]byte) {
	reqID := decodeInt(f[1])

	d.wrapper.execDetailsEnd(reqID)
}

func (d *ibDecoder) wrapTickSnapshotEnd(f [][]byte) {
	reqID := decodeInt(f[1])

	d.wrapper.tickSnapshotEnd(reqID)
}

func (d *ibDecoder) wrapPositionEnd(f [][]byte) {
	// v := decodeInt(f[0])

	d.wrapper.positionEnd()
}

func (d *ibDecoder) wrapAccountSummaryEnd(f [][]byte) {
	reqID := decodeInt(f[1])

	d.wrapper.accountSummaryEnd(reqID)
}

// ------------------------------------------------------------------

func (d *ibDecoder) processTickPriceMsg(f [][]byte) {
	reqID := decodeInt(f[1])
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

	d.wrapper.tickPrice(reqID, tickType, price, attrib)

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
		d.wrapper.tickSize(reqID, sizeTickType, size)
	}

}

func (d *ibDecoder) processOrderStatusMsg(f [][]byte) {
	if d.version < MIN_SERVER_VER_MARKET_CAP_PRICE {
		f = f[1:]
	}
	orderID := decodeInt(f[0])
	status := decodeString(f[1])

	filled := decodeFloat(f[2])

	remaining := decodeFloat(f[3])

	avgFilledPrice := decodeFloat(f[4])

	permID := decodeInt(f[5])
	parentID := decodeInt(f[6])
	lastFillPrice := decodeFloat(f[7])
	clientID := decodeInt(f[8])
	whyHeld := decodeString(f[9])

	var mktCapPrice float64
	if d.version >= MIN_SERVER_VER_MARKET_CAP_PRICE {
		mktCapPrice = decodeFloat(f[10])
	} else {
		mktCapPrice = float64(0)
	}

	d.wrapper.orderStatus(orderID, status, filled, remaining, avgFilledPrice, permID, parentID, lastFillPrice, clientID, whyHeld, mktCapPrice)

}

func (d *ibDecoder) processOpenOrder(f [][]byte) {

	var version int64
	if d.version < MIN_SERVER_VER_ORDER_CONTAINER {
		version = decodeInt(f[0])
		f = f[1:]
	} else {
		version = int64(d.version)
	}

	o := &Order{}
	o.OrderID = decodeInt(f[0])

	c := &Contract{}

	c.ContractID = decodeInt(f[1])
	c.Symbol = decodeString(f[2])
	c.SecurityType = decodeString(f[3])
	c.Expiry = decodeString(f[4])

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
	if d.version >= MIN_SERVER_VER_FRACTIONAL_POSITIONS {
		o.TotalQuantity = decodeFloat(f[11])
	} else {
		o.TotalQuantity = float64(decodeInt(f[11]))
	}

	o.OrderType = decodeString(f[12])
	o.LimitPrice = decodeFloat(f[13]) //todo: show_unset
	o.AuxPrice = decodeFloat(f[14])   //todo: show_unset
	o.TIF = decodeString(f[15])
	o.OCAGroup = decodeString(f[16])
	o.Account = decodeString(f[17])
	o.OpenClose = decodeString(f[18])

	o.Origin = decodeInt(f[19])

	o.OrderRef = decodeString(f[20])
	o.ClientID = decodeInt(f[21])
	o.PermID = decodeInt(f[22])

	o.OutsideRTH = decodeBool(f[23])
	o.Hidden = decodeBool(f[24])
	o.DiscretionaryAmount = decodeFloat(f[25])
	o.GoodAfterTime = decodeString(f[26])

	_ = decodeString(f[27]) //_sharesAllocation

	o.FAGroup = decodeString(f[28])
	o.FAMethod = decodeString(f[29])
	o.FAPercentage = decodeString(f[30])
	o.FAProfile = decodeString(f[31])

	if d.version >= MIN_SERVER_VER_MODELS_SUPPORT {
		o.ModelCode = decodeString(f[32])
		f = f[1:]
	}

	o.GoodTillDate = decodeString(f[32])

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
		o.DeltaNeutralContractID = decodeInt(f[59])
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
		for comboLegsCount := decodeInt(f[65]); comboLegsCount > 0 && comboLegsCount != math.MaxInt64; comboLegsCount-- {
			fmt.Println("comboLegsCount:", comboLegsCount)
			comboleg := ComboLeg{}
			comboleg.ContractID = decodeInt(f[66])
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
		f = f[1:]

		o.OrderComboLegs = []OrderComboLeg{}
		for orderComboLegsCount := decodeInt(f[65]); orderComboLegsCount > 0 && orderComboLegsCount != math.MaxInt64; orderComboLegsCount-- {
			orderComboLeg := OrderComboLeg{}
			orderComboLeg.Price = decodeFloat(f[66])
			o.OrderComboLegs = append(o.OrderComboLegs, orderComboLeg)
			f = f[1:]
		}
		f = f[1:]
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

	o.ClearingAccount = decodeString(f[68])
	o.ClearingIntent = decodeString(f[69])

	if version >= 22 {
		o.NotHeld = decodeBool(f[70])
		f = f[1:]
	}

	if version >= 20 {
		deltaNeutralContractPresent := decodeBool(f[70])
		if deltaNeutralContractPresent {
			c.DeltaNeutralContract = DeltaNeutralContract{}
			c.DeltaNeutralContract.ContractID = decodeInt(f[71])
			c.DeltaNeutralContract.Delta = decodeFloat(f[72])
			c.DeltaNeutralContract.Price = decodeFloat(f[73])
			f = f[3:]
		}
		f = f[1:]
	}

	if version >= 21 {
		o.AlgoStrategy = decodeString(f[70])
		if o.AlgoStrategy != "" {
			o.AlgoParams = []TagValue{}
			for algoParamsCount := decodeInt(f[71]); algoParamsCount > 0 && algoParamsCount != math.MaxInt64; algoParamsCount-- {
				tagValue := TagValue{}
				tagValue.Tag = decodeString(f[72])
				tagValue.Value = decodeString(f[73])
				o.AlgoParams = append(o.AlgoParams, tagValue)
				f = f[2:]
			}
		}
		f = f[1:]
	}

	if version >= 33 {
		o.Solictied = decodeBool(f[70])
		f = f[1:]
	}

	orderState := &OrderState{}

	o.WhatIf = decodeBool(f[70])

	orderState.Status = decodeString(f[71])

	if d.version >= MIN_SERVER_VER_WHAT_IF_EXT_FIELDS {
		orderState.InitialMarginBefore = decodeString(f[72])
		orderState.MaintenanceMarginBefore = decodeString(f[73])
		orderState.EquityWithLoanBefore = decodeString(f[74])
		orderState.InitialMarginChange = decodeString(f[75])
		orderState.MaintenanceMarginChange = decodeString(f[76])
		orderState.EquityWithLoanChange = decodeString(f[77])
		f = f[6:]
	}

	orderState.InitialMarginAfter = decodeString(f[72])
	orderState.MaintenanceMarginAfter = decodeString(f[73])
	orderState.EquityWithLoanAfter = decodeString(f[74])

	orderState.Commission = decodeFloat(f[75])
	orderState.MinCommission = decodeFloat(f[76])
	orderState.MaxCommission = decodeFloat(f[77])
	orderState.CommissionCurrency = decodeString(f[78])
	orderState.WarningText = decodeString(f[79])

	if version >= 34 {
		o.RandomizeSize = decodeBool(f[80])
		o.RandomizePrice = decodeBool(f[81])
		f = f[2:]
	}

	if d.version >= MIN_SERVER_VER_PEGGED_TO_BENCHMARK {
		if o.OrderType == "PEG BENCH" {
			o.ReferenceContractID = decodeInt(f[80])
			o.IsPeggedChangeAmountDecrease = decodeBool(f[81])
			o.PeggedChangeAmount = decodeFloat(f[82])
			o.ReferenceChangeAmount = decodeFloat(f[83])
			o.ReferenceExchangeID = decodeString(f[84])
			f = f[5:]
		}

		o.Conditions = []OrderConditioner{}
		if conditionsSize := decodeInt(f[80]); conditionsSize > 0 && conditionsSize != math.MaxInt64 {
			for ; conditionsSize > 0; conditionsSize-- {
				conditionType := decodeInt(f[81])
				cond, condSize := InitOrderCondition(conditionType)
				cond.decode(f[82 : 82+condSize])

				o.Conditions = append(o.Conditions, cond)
				f = f[condSize+1:]
			}
			o.ConditionsIgnoreRth = decodeBool(f[81])
			o.ConditionsCancelOrder = decodeBool(f[82])
			f = f[2:]
		}

		o.AdjustedOrderType = decodeString(f[81])
		o.TriggerPrice = decodeFloat(f[82])
		o.TrailStopPrice = decodeFloat(f[83])
		o.LmtPriceOffset = decodeFloat(f[84])
		o.AdjustedStopPrice = decodeFloat(f[85])
		o.AdjustedStopLimitPrice = decodeFloat(f[86])
		o.AdjustedTrailingAmount = decodeFloat(f[87])
		o.AdjustableTrailingUnit = decodeInt(f[88])
		f = f[9:]
	}

	if d.version >= MIN_SERVER_VER_SOFT_DOLLAR_TIER {
		name := decodeString(f[80])
		value := decodeString(f[81])
		displayName := decodeString(f[82])
		o.SoftDollarTier = SoftDollarTier{name, value, displayName}
		f = f[3:]
	}

	if d.version >= MIN_SERVER_VER_CASH_QTY {
		o.CashQty = decodeFloat(f[80])
		f = f[1:]
	}

	if d.version >= MIN_SERVER_VER_AUTO_PRICE_FOR_HEDGE {
		o.DontUseAutoPriceForHedge = decodeBool(f[80])
		f = f[1:]
	}

	if d.version >= MIN_SERVER_VER_ORDER_CONTAINER {
		o.IsOmsContainer = decodeBool(f[80])
		f = f[1:]
	}

	if d.version >= MIN_SERVER_VER_D_PEG_ORDERS {
		o.DiscretionaryUpToLimitPrice = decodeBool(f[80])
		f = f[1:]
	}

	d.wrapper.openOrder(o.OrderID, c, o, orderState)

}

func (d *ibDecoder) processPortfolioValueMsg(f [][]byte) {
	v := decodeInt(f[0])

	c := &Contract{}
	c.ContractID = decodeInt(f[1])
	c.Symbol = decodeString(f[2])
	c.SecurityType = decodeString(f[3])
	c.Expiry = decodeString(f[4])
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
	var position float64
	if d.version >= MIN_SERVER_VER_FRACTIONAL_POSITIONS {
		position = decodeFloat(f[9])
	} else {
		position = float64(decodeInt(f[9]))
	}

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
func (d *ibDecoder) processContractDataMsg(f [][]byte) {
	v := decodeInt(f[1])
	var reqID int64 = 1
	if v >= 3 {
		reqID = decodeInt(f[2])
		f = f[1:]
	}

	cd := ContractDetails{}
	cd.Contract = Contract{}
	cd.Contract.Symbol = decodeString(f[2])
	cd.Contract.SecurityType = decodeString(f[3])

	lastTradeDateOrContractMonth := f[4]
	if !bytes.Equal(lastTradeDateOrContractMonth, []byte{}) {
		split := bytes.Split(lastTradeDateOrContractMonth, []byte{32})
		if len(split) > 0 {
			cd.Contract.Expiry = decodeString(split[0])
		}

		if len(split) > 1 {
			cd.LastTradeTime = decodeString(split[1])
		}
	}

	cd.Contract.Strike = decodeFloat(f[5])
	cd.Contract.Right = decodeString(f[6])
	cd.Contract.Exchange = decodeString(f[7])
	cd.Contract.Currency = decodeString(f[8])
	cd.Contract.LocalSymbol = decodeString(f[9])
	cd.MarketName = decodeString(f[10])
	cd.Contract.TradingClass = decodeString(f[11])
	cd.Contract.ContractID = decodeInt(f[12])
	cd.MinTick = decodeFloat(f[13])
	if d.version >= MIN_SERVER_VER_MD_SIZE_MULTIPLIER {
		cd.MdSizeMultiplier = decodeInt(f[14])
		f = f[1:]
	}

	cd.Contract.Multiplier = decodeString(f[14])
	cd.OrderTypes = decodeString(f[15])
	cd.ValidExchanges = decodeString(f[16])
	cd.PriceMagnifier = decodeInt(f[17])

	if v >= 4 {
		cd.UnderContractID = decodeInt(f[18])
		f = f[1:]
	}

	if v >= 5 {
		cd.LongName = decodeString(f[18])
		cd.Contract.PrimaryExchange = decodeString(f[19])
		f = f[2:]
	}

	if v >= 6 {
		cd.ContractMonth = decodeString(f[18])
		cd.Industry = decodeString(f[19])
		cd.Category = decodeString(f[20])
		cd.Subcategory = decodeString(f[21])
		cd.TimezoneID = decodeString(f[22])
		cd.TradingHours = decodeString(f[23])
		cd.LiquidHours = decodeString(f[24])
		f = f[7:]
	}

	if v >= 8 {
		cd.EVRule = decodeString(f[18])
		cd.EVMultiplier = decodeInt(f[19])
		f = f[2:]
	}

	if v >= 7 {
		cd.SecurityIDList = []TagValue{}
		for secIDListCount := decodeInt(f[18]); secIDListCount > 0 && secIDListCount != math.MaxInt64; secIDListCount-- {
			tagValue := TagValue{}
			tagValue.Tag = decodeString(f[19])
			tagValue.Value = decodeString(f[20])
			f = f[2:]
		}
		f = f[1:]
	}

	if d.version >= MIN_SERVER_VER_AGG_GROUP {
		cd.AggGroup = decodeInt(f[18])
		f = f[1:]
	}

	if d.version >= MIN_SERVER_VER_UNDERLYING_INFO {
		cd.UnderSymbol = decodeString(f[18])
		cd.UnderSecurityType = decodeString(f[19])
		f = f[2:]
	}

	if d.version >= MIN_SERVER_VER_MARKET_RULES {
		cd.MarketRuleIDs = decodeString(f[18])
		f = f[1:]
	}

	if d.version >= MIN_SERVER_VER_REAL_EXPIRATION_DATE {
		cd.RealExpirationDate = decodeString(f[18])
	}

	d.wrapper.contractDetails(reqID, &cd)

}
func (d *ibDecoder) processBondContractDataMsg(f [][]byte) {
	v := decodeInt(f[0])

	var reqID int64 = -1

	if v >= 3 {
		reqID = decodeInt(f[1])
		f = f[1:]
	}

	c := &ContractDetails{}
	c.Contract.Symbol = decodeString(f[1])
	c.Contract.SecurityType = decodeString(f[2])
	c.Cusip = decodeString(f[3])
	c.Coupon = decodeInt(f[4])

	splittedExpiry := bytes.Split(f[5], []byte{32})
	switch s := len(splittedExpiry); {
	case s > 0:
		c.Maturity = decodeString(splittedExpiry[0])
	case s > 1:
		c.LastTradeTime = decodeString(splittedExpiry[1])
	case s > 2:
		c.TimezoneID = decodeString(splittedExpiry[2])
	}

	c.IssueDate = decodeString(f[6])
	c.Ratings = decodeString(f[7])
	c.BondType = decodeString(f[8])
	c.CouponType = decodeString(f[9])
	c.Convertible = decodeBool(f[10])
	c.Callable = decodeBool(f[11])
	c.Putable = decodeBool(f[12])
	c.DescAppend = decodeString(f[13])
	c.Contract.Exchange = decodeString(f[14])
	c.Contract.Currency = decodeString(f[15])
	c.MarketName = decodeString(f[16])
	c.Contract.TradingClass = decodeString(f[17])
	c.Contract.ContractID = decodeInt(f[18])
	c.MinTick = decodeFloat(f[19])

	if d.version >= MIN_SERVER_VER_MD_SIZE_MULTIPLIER {
		c.MdSizeMultiplier = decodeInt(f[20])
		f = f[1:]
	}

	c.OrderTypes = decodeString(f[20])
	c.ValidExchanges = decodeString(f[21])
	c.NextOptionDate = decodeString(f[22])
	c.NextOptionType = decodeString(f[23])
	c.NextOptionPartial = decodeBool(f[24])
	c.Notes = decodeString(f[25])

	if v >= 4 {
		c.LongName = decodeString(f[26])
		f = f[1:]
	}

	if v >= 6 {
		c.EVRule = decodeString(f[26])
		c.EVMultiplier = decodeInt(f[27])
		f = f[2:]
	}

	if v >= 5 {
		c.SecurityIDList = []TagValue{}
		for secIDListCount := decodeInt(f[26]); secIDListCount > 0; secIDListCount-- {
			tagValue := TagValue{}
			tagValue.Tag = decodeString(f[27])
			tagValue.Value = decodeString(f[28])
			c.SecurityIDList = append(c.SecurityIDList, tagValue)
			f = f[2:]
		}
		f = f[1:]
	}

	if d.version >= MIN_SERVER_VER_AGG_GROUP {
		c.AggGroup = decodeInt(f[26])
		f = f[1:]
	}

	if d.version >= MIN_SERVER_VER_MARKET_RULES {
		c.MarketRuleIDs = decodeString(f[26])
		f = f[1:]
	}

	d.wrapper.bondContractDetails(reqID, c)

}
func (d *ibDecoder) processScannerDataMsg(f [][]byte) {
	f = f[1:]
	reqID := decodeInt(f[0])
	for numofElements := decodeInt(f[1]); numofElements > 0; numofElements-- {
		sd := ScanData{}
		sd.ContractDetails = ContractDetails{}
		sd.Rank = decodeInt(f[2])
		sd.ContractDetails.Contract.ContractID = decodeInt(f[3])
		sd.ContractDetails.Contract.Symbol = decodeString(f[4])
		sd.ContractDetails.Contract.SecurityType = decodeString(f[5])
		sd.ContractDetails.Contract.Expiry = decodeString(f[6])
		sd.ContractDetails.Contract.Strike = decodeFloat(f[7])
		sd.ContractDetails.Contract.Right = decodeString(f[8])
		sd.ContractDetails.Contract.Exchange = decodeString(f[9])
		sd.ContractDetails.Contract.Currency = decodeString(f[10])
		sd.ContractDetails.Contract.LocalSymbol = decodeString(f[11])
		sd.ContractDetails.MarketName = decodeString(f[12])
		sd.ContractDetails.Contract.TradingClass = decodeString(f[13])
		sd.Distance = decodeString(f[14])
		sd.Benchmark = decodeString(f[15])
		sd.Projection = decodeString(f[16])
		sd.Legs = decodeString(f[17])

		d.wrapper.scannerData(reqID, sd.Rank, &(sd.ContractDetails), sd.Distance, sd.Benchmark, sd.Projection, sd.Legs)
		f = f[16:]

	}

	d.wrapper.scannerDataEnd(reqID)

}
func (d *ibDecoder) processExecutionDataMsg(f [][]byte) {
	var v int64
	if d.version < MIN_SERVER_VER_LAST_LIQUIDITY {
		v = decodeInt(f[0])
		f = f[1:]
	} else {
		v = int64(d.version)
	}

	var reqID int64 = -1
	if v >= 7 {
		reqID = decodeInt(f[0])
		f = f[1:]
	}

	orderID := decodeInt(f[0])

	c := Contract{}
	c.ContractID = decodeInt(f[1])
	c.Symbol = decodeString(f[2])
	c.SecurityType = decodeString(f[3])
	c.Expiry = decodeString(f[4])
	c.Strike = decodeFloat(f[5])
	c.Right = decodeString(f[6])

	if v >= 9 {
		c.Multiplier = decodeString(f[7])
		f = f[1:]
	}

	c.Exchange = decodeString(f[7])
	c.Currency = decodeString(f[8])
	c.LocalSymbol = decodeString(f[9])

	if v >= 10 {
		c.TradingClass = decodeString(f[10])
		f = f[1:]
	}

	e := Execution{}
	e.OrderID = orderID
	e.ExecID = decodeString(f[10])
	e.Time = decodeString(f[11])
	e.AccountCode = decodeString(f[12])
	e.Exchange = decodeString(f[13])
	e.Side = decodeString(f[14])
	e.Shares = decodeFloat(f[15])
	e.Price = decodeFloat(f[16])
	e.PermID = decodeInt(f[17])
	e.ClientID = decodeInt(f[18])
	e.Liquidation = decodeInt(f[19])

	if v >= 6 {
		e.CumQty = decodeFloat(f[20])
		e.AveragePrice = decodeFloat(f[21])
		f = f[2:]
	}

	if v >= 8 {
		e.OrderRef = decodeString(f[20])
		f = f[1:]
	}

	if v >= 9 {
		e.EVRule = decodeString(f[20])
		e.EVMultiplier = decodeFloat(f[21])
		f = f[2:]
	}

	if d.version >= MIN_SERVER_VER_MODELS_SUPPORT {
		e.ModelCode = decodeString(f[20])
		f = f[1:]
	}
	if d.version >= MIN_SERVER_VER_LAST_LIQUIDITY {
		e.LastLiquidity = decodeInt(f[20])
	}

	d.wrapper.execDetails(reqID, &c, &e)

}
func (d *ibDecoder) processHistoricalDataMsg(f [][]byte) {
	if d.version < MIN_SERVER_VER_SYNT_REALTIME_BARS {
		f = f[1:]
	}

	reqID := decodeInt(f[0])
	startDatestr := decodeString(f[1])
	endDateStr := decodeString(f[2])

	for itemCount := decodeInt(f[3]); itemCount > 0; itemCount-- {
		bar := &BarData{}
		bar.Date = decodeString(f[4])
		bar.Open = decodeFloat(f[5])
		bar.High = decodeFloat(f[6])
		bar.Low = decodeFloat(f[7])
		bar.Close = decodeFloat(f[8])
		bar.Volume = decodeFloat(f[9])
		bar.Average = decodeFloat(f[10])

		if d.version < MIN_SERVER_VER_SYNT_REALTIME_BARS {
			f = f[1:]
		}
		bar.BarCount = decodeInt(f[11])
		f = f[8:]
		d.wrapper.historicalData(reqID, bar)
	}
	f = f[1:]

	d.wrapper.historicalDataEnd(reqID, startDatestr, endDateStr)

}
func (d *ibDecoder) processHistoricalDataUpdateMsg(f [][]byte) {
	reqID := decodeInt(f[0])
	bar := &BarData{}
	bar.BarCount = decodeInt(f[1])
	bar.Date = decodeString(f[2])
	bar.Open = decodeFloat(f[3])
	bar.Close = decodeFloat(f[4])
	bar.High = decodeFloat(f[5])
	bar.Low = decodeFloat(f[6])
	bar.Volume = decodeFloat(f[7])

	d.wrapper.historicalDataUpdate(reqID, bar)

}
func (d *ibDecoder) processRealTimeBarMsg(f [][]byte) {
	// f = f[1:]
	// reqID := decodeInt(f[0])

	// rtb := &RealTimeBar{}
	// rtb.Time = decodeInt(f[1])
	// rtb.Open = decodeFloat(f[2])
	// rtb.High = decodeFloat(f[3])
	// rtb.Low = decodeFloat(f[4])
	// rtb.Close = decodeFloat(f[5])
	// rtb.Volume = decodeFloat(f[6])
	// rtb.Wap = decodeFloat(f[7])
	// rtb.Count = decodeInt(f[8])

}
func (d *ibDecoder) processTickOptionComputationMsg(f [][]byte) {
	optPrice := math.MaxFloat64
	pvDividend := math.MaxFloat64
	gamma := math.MaxFloat64
	vega := math.MaxFloat64
	theta := math.MaxFloat64
	undPrice := math.MaxFloat64

	v := decodeInt(f[0])
	reqID := decodeInt(f[1])
	tickType := decodeInt(f[2])

	impliedVol := decodeFloat(f[3])
	delta := decodeFloat(f[4])

	if v >= 6 || tickType == MODEL_OPTION || tickType == DELAYED_MODEL_OPTION {
		optPrice = decodeFloat(f[5])
		pvDividend = decodeFloat(f[6])
		f = f[2:]

	}

	if v >= 6 {
		gamma = decodeFloat(f[5])
		vega = decodeFloat(f[6])
		theta = decodeFloat(f[7])
		undPrice = decodeFloat(f[8])

	}

	switch {
	case impliedVol < 0:
		impliedVol = math.MaxFloat64
		fallthrough
	case delta == -2:
		delta = math.MaxFloat64
		fallthrough
	case optPrice == -1:
		optPrice = math.MaxFloat64
		fallthrough
	case pvDividend == -1:
		pvDividend = math.MaxFloat64
		fallthrough
	case gamma == -2:
		gamma = math.MaxFloat64
		fallthrough
	case vega == -2:
		vega = math.MaxFloat64
		fallthrough
	case theta == -2:
		theta = math.MaxFloat64
		fallthrough
	case undPrice == -1:
		undPrice = math.MaxFloat64
	}

	d.wrapper.tickOptionComputation(reqID, tickType, impliedVol, delta, optPrice, pvDividend, gamma, vega, theta, undPrice)

}

func (d *ibDecoder) processDeltaNeutralValidationMsg(f [][]byte) {
	_ = decodeInt(f[0])
	reqID := decodeInt(f[1])
	deltaNeutralContract := DeltaNeutralContract{}

	deltaNeutralContract.ContractID = decodeInt(f[2])
	deltaNeutralContract.Delta = decodeFloat(f[3])
	deltaNeutralContract.Price = decodeFloat(f[4])

	d.wrapper.deltaNeutralValidation(reqID, deltaNeutralContract)

}
func (d *ibDecoder) processMarketDataTypeMsg(f [][]byte) {

}
func (d *ibDecoder) processCommissionReportMsg(f [][]byte) {
	_ = decodeInt(f[0])
	cr := CommissionReport{}
	cr.ExecId = decodeString(f[1])
	cr.Commission = decodeFloat(f[2])
	cr.Currency = decodeString(f[3])
	cr.RealizedPNL = decodeFloat(f[4])
	cr.Yield = decodeFloat(f[5])
	cr.YieldRedemptionDate = decodeInt(f[6])

	d.wrapper.commissionReport(cr)

}
func (d *ibDecoder) processPositionDataMsg(f [][]byte) {
	v := decodeInt(f[0])
	acc := decodeString(f[1])

	c := new(Contract)
	c.ContractID = decodeInt(f[2])
	c.Symbol = decodeString(f[3])
	c.SecurityType = decodeString(f[4])
	c.Expiry = decodeString(f[5])
	c.Strike = decodeFloat(f[6])
	c.Right = decodeString(f[7])
	c.Multiplier = decodeString(f[8])
	c.Currency = decodeString(f[9])
	c.LocalSymbol = decodeString(f[10])

	if v >= 2 {
		c.TradingClass = decodeString(f[11])
		f = f[1:]
	}

	var p float64
	if d.version >= MIN_SERVER_VER_FRACTIONAL_POSITIONS {
		p = decodeFloat(f[11])
	} else {
		p = float64(decodeInt(f[11]))
	}

	var avgCost float64
	if v >= 3 {
		avgCost = decodeFloat(f[12])
	}

	d.wrapper.position(acc, c, p, avgCost)

}
func (d *ibDecoder) processPositionMultiMsg(f [][]byte) {

}
func (d *ibDecoder) processSecurityDefinitionOptionParameterMsg(f [][]byte) {

}
func (d *ibDecoder) processSecurityDefinitionOptionParameterEndMsg(f [][]byte) {

}
func (d *ibDecoder) processSoftDollarTiersMsg(f [][]byte) {

}
func (d *ibDecoder) processFamilyCodesMsg(f [][]byte) {

}
func (d *ibDecoder) processSymbolSamplesMsg(f [][]byte) {

}
func (d *ibDecoder) processSmartComponents(f [][]byte) {

}
func (d *ibDecoder) processTickReqParams(f [][]byte) {

}
func (d *ibDecoder) processMktDepthExchanges(f [][]byte) {

}

func (d *ibDecoder) processHeadTimestamp(f [][]byte) {

}
func (d *ibDecoder) processTickNews(f [][]byte) {

}
func (d *ibDecoder) processNewsProviders(f [][]byte) {

}
func (d *ibDecoder) processNewsArticle(f [][]byte) {

}
func (d *ibDecoder) processHistoricalNews(f [][]byte) {

}
func (d *ibDecoder) processHistoricalNewsEnd(f [][]byte) {

}
func (d *ibDecoder) processHistogramData(f [][]byte) {

}
func (d *ibDecoder) processRerouteMktDataReq(f [][]byte) {

}
func (d *ibDecoder) processRerouteMktDepthReq(f [][]byte) {

}
func (d *ibDecoder) processMarketRuleMsg(f [][]byte) {

}
func (d *ibDecoder) processPnLMsg(f [][]byte) {

}
func (d *ibDecoder) processPnLSingleMsg(f [][]byte) {

}
func (d *ibDecoder) processHistoricalTicks(f [][]byte) {

}
func (d *ibDecoder) processHistoricalTicksBidAsk(f [][]byte) {

}
func (d *ibDecoder) processHistoricalTicksLast(f [][]byte) {

}
func (d *ibDecoder) processTickByTickMsg(f [][]byte) {

}
func (d *ibDecoder) processOrderBoundMsg(f [][]byte) {

}
func (d *ibDecoder) processMarketDepthL2Msg(f [][]byte) {

}

// ----------------------------------------------------
