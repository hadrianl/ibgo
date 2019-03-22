package ibgo

import (
	"bytes"
	"strconv"
	"time"
)

const (
	TIME_FORMAT string = "2019-03-21 17:18:00 +0800 CST"
)

type IbDecoder struct {
	wrapper       IbWrapper
	version       Version
	msgId2process map[IN]func([][]byte)
}

//NewIbDecoder create a decoder to decode the fileds based on version
func NewIbDecoder(wrapper IbWrapper, version Version) *IbDecoder {
	decoder := IbDecoder{}
	decoder.wrapper = wrapper
	decoder.version = version
	return &decoder
}

func (d *IbDecoder) setVersion(version Version) {
	d.version = version
}

func (d *IbDecoder) interpret(fs ...[]byte) {
	if len(fs) == 0 {
		return
	}

	MsgId, _ := strconv.ParseInt(string(fs[0]), 10, 64)
	processer := d.msgId2process[IN(MsgId)]
	processer(fs[1:])
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
		TICK_PRICE:    d.processTickPriceMsg,
		TICK_SIZE:     d.wrapTickSize,
		ORDER_STATUS:  d.processOrderStatusMsg,
		ERR_MSG:       d.wrapError,
		OPEN_ORDER:    d.processOpenOrder,
		NEXT_VALID_ID: d.wrapNextValidId,
		MANAGED_ACCTS: d.wrapManagedAccounts,

		CURRENT_TIME: d.wrapCurrentTime,
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
	o.LmtPrice = decodeFloat(f[13])

}
func (d *IbDecoder) processPortfolioValueMsg(f [][]byte) {

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
