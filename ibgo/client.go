package ibgo

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

const (
	MaxRequests      = 95
	RequestInternal  = 2
	MaxClientVersion = 148
)

// IbClient is the key component which is used to send request to TWS ro Gateway , such subscribe market data or place order
type IbClient struct {
	host             string
	port             int
	clientID         int64
	conn             *IbConnection
	reader           *bufio.Reader
	writer           *bufio.Writer
	wrapper          IbWrapper
	decoder          *ibDecoder
	inBuffer         *bytes.Buffer
	outBuffer        *bytes.Buffer
	connectOption    []byte
	reqIDSeq         int64
	reqChan          chan []byte
	errChan          chan error
	msgChan          chan [][]byte
	timeChan         chan time.Time
	terminatedSignal chan int
	clientVersion    Version
	serverVersion    Version
	serverTime       time.Time
	wg               sync.WaitGroup
}

func NewIbClient(wrapper IbWrapper) *IbClient {
	ic := &IbClient{}
	ic.SetWrapper(wrapper)
	ic.reset()

	return ic
}

func (ic *IbClient) setConnState(connState int) {
	OldConnState := ic.conn.state
	ic.conn.state = connState
	log.Printf("connState: %v -> %v", OldConnState, connState)
}

func (ic *IbClient) GetReqID() int64 {
	ic.reqIDSeq++
	return ic.reqIDSeq
}

//SetWrapper
func (ic *IbClient) SetWrapper(wrapper IbWrapper) {
	ic.wrapper = wrapper
	ic.decoder = &ibDecoder{wrapper: ic.wrapper}
}

//Connect
func (ic *IbClient) Connect(host string, port int, clientID int64) error {

	ic.host = host
	ic.port = port
	ic.clientID = clientID
	if err := ic.conn.connect(host, port); err != nil {
		return err
	}

	ic.setConnState(CONNECTING)
	return nil
	// 连接后开始
}

//Disconnect
func (ic *IbClient) Disconnect() error {

	ic.terminatedSignal <- 1
	ic.terminatedSignal <- 1
	ic.terminatedSignal <- 1
	if err := ic.conn.disconnect(); err != nil {
		return err
	}

	defer log.Println("Disconnected!")
	ic.wg.Wait()
	ic.setConnState(DISCONNECTED)
	return nil
}

// IsConnected check if there is a connection to TWS or GateWay
func (ic *IbClient) IsConnected() bool {
	return ic.conn.state == CONNECTED
}

// send the clientId to TWS or Gateway
func (ic *IbClient) startAPI() error {
	var startAPI []byte
	v := 2
	if ic.serverVersion >= MIN_SERVER_VER_OPTIONAL_CAPABILITIES {
		startAPI = makeMsgBuf(int64(START_API), int64(v), ic.clientID, "")
	} else {
		startAPI = makeMsgBuf(int64(START_API), int64(v), ic.clientID)
	}

	log.Println("Start API:", startAPI)
	if _, err := ic.writer.Write(startAPI); err != nil {
		return err
	}

	err := ic.writer.Flush()

	return err
}

// handshake with the TWS or GateWay to ensure the version
func (ic *IbClient) HandShake() error {
	log.Println("Try to handShake with TWS or GateWay...")
	var msg bytes.Buffer
	head := []byte("API\x00")
	minVer := []byte(strconv.FormatInt(int64(MIN_CLIENT_VER), 10))
	maxVer := []byte(strconv.FormatInt(int64(MAX_CLIENT_VER), 10))
	connectOptions := []byte("")
	clientVersion := bytes.Join([][]byte{[]byte("v"), minVer, []byte(".."), maxVer, connectOptions}, []byte(""))
	sizeofCV := make([]byte, 4)
	binary.BigEndian.PutUint32(sizeofCV, uint32(len(clientVersion)))
	msg.Write(head)
	msg.Write(sizeofCV)
	msg.Write(clientVersion)
	log.Println("HandShake Init...")
	if _, err := ic.writer.Write(msg.Bytes()); err != nil {
		return err
	}

	if err := ic.writer.Flush(); err != nil {
		return err
	}

	log.Println("Recv ServerInfo...")
	if msgBuf, err := readMsgBuf(ic.reader); err != nil {
		return err
	} else {
		serverInfo := splitMsgBuf(msgBuf)
		v, _ := strconv.Atoi(string(serverInfo[0]))
		ic.serverVersion = Version(v)
		ic.serverTime = bytesToTime(serverInfo[1])
		ic.decoder.setVersion(ic.serverVersion) // Init Decoder
		ic.decoder.setmsgID2process()
		log.Println("ServerVersion:", ic.serverVersion)
		log.Println("ServerTime:", ic.serverTime)
	}

	if err := ic.startAPI(); err != nil {
		return err
	}

	ic.setConnState(CONNECTED)
	ic.wrapper.connectAck()

	return nil
}

func (ic *IbClient) reset() {
	ic.reqIDSeq = 0
	ic.conn = &IbConnection{}
	ic.conn.reset()
	ic.reader = bufio.NewReader(ic.conn)
	ic.writer = bufio.NewWriter(ic.conn)
	ic.inBuffer = bytes.NewBuffer(make([]byte, 0, 4096))
	ic.outBuffer = bytes.NewBuffer(make([]byte, 0, 4096))
	ic.reqChan = make(chan []byte, 10)
	ic.errChan = make(chan error, 10)
	ic.msgChan = make(chan [][]byte, 100)
	ic.terminatedSignal = make(chan int, 3)
	ic.wg = sync.WaitGroup{}

}

// ---------------req func ----------------------------------------------

/*
Market Data
*/

/* ReqMktData
Call this function to request market data. The market data
        will be returned by the tickPrice and tickSize events.

        reqId: TickerId - The ticker id. Must be a unique value. When the
            market data returns, it will be identified by this tag. This is
            also used when canceling the market data.
        contract:Contract - This structure contains a description of the
            Contractt for which market data is being requested.
        genericTickList:str - A commma delimited list of generic tick types.
            Tick types can be found in the Generic Tick Types page.
            Prefixing w/ 'mdoff' indicates that top mkt data shouldn't tick.
            You can specify the news source by postfixing w/ ':<source>.
            Example: "mdoff,292:FLY+BRF"
        snapshot:bool - Check to return a single snapshot of Market data and
            have the market data subscription cancel. Do not enter any
            genericTicklist values if you use snapshots.
        regulatorySnapshot: bool - With the US Value Snapshot Bundle for stocks,
            regulatory snapshots are available for 0.01 USD each.
        mktDataOptions:TagValueList - For internal use only.
            Use default value XYZ.
*/
func (ic *IbClient) ReqMktData(reqID int64, contract Contract, genericTickList string, snapshot bool, regulatorySnapshot bool, mktDataOptions []TagValue) {
	switch {
	case ic.serverVersion < MIN_SERVER_VER_DELTA_NEUTRAL && contract.DeltaNeutralContract != nil:
		ic.wrapper.error(reqID, UPDATE_TWS.code, UPDATE_TWS.msg+"  It does not support delta-neutral orders.")
		return
	case ic.serverVersion < MIN_SERVER_VER_REQ_MKT_DATA_CONID && contract.ContractID > 0:
		ic.wrapper.error(reqID, UPDATE_TWS.code, UPDATE_TWS.msg+"  It does not support conId parameter.")
		return
	case ic.serverVersion < MIN_SERVER_VER_TRADING_CLASS && contract.TradingClass != "":
		ic.wrapper.error(reqID, UPDATE_TWS.code, UPDATE_TWS.msg+"  It does not support tradingClass parameter in reqMktData.")
		return
	}

	v := 11
	fields := make([]interface{}, 0)
	fields = append(fields,
		REQ_MKT_DATA,
		v,
		reqID,
	)

	if ic.serverVersion >= MIN_SERVER_VER_REQ_MKT_DATA_CONID {
		fields = append(fields, contract.ContractID)
	}

	fields = append(fields,
		contract.Symbol,
		contract.SecurityType,
		contract.Expiry,
		contract.Strike,
		contract.Right,
		contract.Multiplier,
		contract.Exchange,
		contract.PrimaryExchange,
		contract.Currency,
		contract.LocalSymbol)

	if ic.serverVersion >= MIN_SERVER_VER_TRADING_CLASS {
		fields = append(fields, contract.TradingClass)
	}

	if contract.SecurityType == "BAG" {
		comboLegsCount := len(contract.ComboLegs)
		fields = append(fields, comboLegsCount)
		for _, comboLeg := range contract.ComboLegs {
			fields = append(fields,
				comboLeg.ContractID,
				comboLeg.Ratio,
				comboLeg.Action,
				comboLeg.Exchange)
		}
	}

	if ic.serverVersion >= MIN_SERVER_VER_DELTA_NEUTRAL {
		if contract.DeltaNeutralContract != nil {
			fields = append(fields,
				true,
				contract.DeltaNeutralContract.ContractID,
				contract.DeltaNeutralContract.Delta,
				contract.DeltaNeutralContract.Price)
		} else {
			fields = append(fields, false)
		}
	}

	fields = append(fields,
		genericTickList,
		snapshot)

	if ic.serverVersion >= MIN_SERVER_VER_REQ_SMART_COMPONENTS {
		fields = append(fields, regulatorySnapshot)
	}

	if ic.serverVersion >= MIN_SERVER_VER_LINKING {
		if mktDataOptions != nil {
			panic("not supported")
		}
		fields = append(fields, "")
	}

	msg := makeMsgBuf(fields...)
	ic.reqChan <- msg
}

//CancelMktData
func (ic *IbClient) CancelMktData(reqID int64) {
	v := 2
	fields := make([]interface{}, 0)
	fields = append(fields,
		CANCEL_MKT_DATA,
		v,
		reqID,
	)

	msg := makeMsgBuf(fields...)

	ic.reqChan <- msg
}

/*ReqMarketDataType
The API can receive frozen market data from Trader
        Workstation. Frozen market data is the last data recorded in our system.
        During normal trading hours, the API receives real-time market data. If
        you use this function, you are telling TWS to automatically switch to
        frozen market data after the close. Then, before the opening of the next
        trading day, market data will automatically switch back to real-time
        market data.

        marketDataType:int - 1 for real-time streaming market data or 2 for
            frozen market data
*/
func (ic *IbClient) ReqMarketDataType(marketDataType int64) {
	if ic.serverVersion < MIN_SERVER_VER_REQ_MARKET_DATA_TYPE {
		ic.wrapper.error(NO_VALID_ID, UPDATE_TWS.code, UPDATE_TWS.msg+"  It does not support market data type requests.")
		return
	}

	v := 1
	fields := make([]interface{}, 0)
	fields = append(fields, REQ_MARKET_DATA_TYPE, v, marketDataType)

	msg := makeMsgBuf(fields...)

	ic.reqChan <- msg
}

func (ic *IbClient) ReqSmartComponents(reqID int64, bboExchange string) {
	if ic.serverVersion < MIN_SERVER_VER_REQ_SMART_COMPONENTS {
		ic.wrapper.error(NO_VALID_ID, UPDATE_TWS.code, UPDATE_TWS.msg+"  It does not support smart components request.")
		return
	}

	msg := makeMsgBuf(REQ_SMART_COMPONENTS, reqID, bboExchange)

	ic.reqChan <- msg
}

func (ic *IbClient) ReqMarketRule(marketRuleID int64) {
	if ic.serverVersion < MIN_SERVER_VER_MARKET_RULES {
		ic.wrapper.error(NO_VALID_ID, UPDATE_TWS.code, UPDATE_TWS.msg+" It does not support market rule requests.")
		return
	}

	msg := makeMsgBuf(REQ_MARKET_RULE, marketRuleID)

	ic.reqChan <- msg
}

func (ic *IbClient) ReqTickByTickData(reqID int64, contract *Contract, tickType string, numberOfTicks int64, ignoreSize bool) {
	if ic.serverVersion < MIN_SERVER_VER_TICK_BY_TICK {
		ic.wrapper.error(NO_VALID_ID, UPDATE_TWS.code, UPDATE_TWS.msg+" It does not support tick-by-tick data requests.")
		return
	}

	if ic.serverVersion < MIN_SERVER_VER_TICK_BY_TICK_IGNORE_SIZE {
		ic.wrapper.error(NO_VALID_ID, UPDATE_TWS.code, UPDATE_TWS.msg+" It does not support ignoreSize and numberOfTicks parameters in tick-by-tick data requests.")
		return
	}

	fields := make([]interface{}, 0)
	fields = append(fields, REQ_TICK_BY_TICK_DATA,
		reqID,
		contract.ContractID,
		contract.Symbol,
		contract.SecurityType,
		contract.Expiry,
		contract.Strike,
		contract.Right,
		contract.Multiplier,
		contract.Exchange,
		contract.PrimaryExchange,
		contract.Currency,
		contract.LocalSymbol,
		contract.TradingClass,
		tickType)

	if ic.serverVersion >= MIN_SERVER_VER_TICK_BY_TICK_IGNORE_SIZE {
		fields = append(fields, numberOfTicks, ignoreSize)
	}

	msg := makeMsgBuf(fields)

	ic.reqChan <- msg
}

func (ic *IbClient) CancelTickByTickData(reqID int64) {
	if ic.serverVersion < MIN_SERVER_VER_TICK_BY_TICK {
		ic.wrapper.error(NO_VALID_ID, UPDATE_TWS.code, UPDATE_TWS.msg+" It does not support tick-by-tick data requests.")
		return
	}

	msg := makeMsgBuf(CANCEL_TICK_BY_TICK_DATA, reqID)

	ic.reqChan <- msg
}

/*
   ##########################################################################
   ################## Options
   ##########################################################################
*/

/*CalculateImpliedVolatility
Call this function to calculate volatility for a supplied
        option price and underlying price. Result will be delivered
        via EWrapper.tickOptionComputation()

        reqId:TickerId -  The request id.
        contract:Contract -  Describes the contract.
        optionPrice:double - The price of the option.
        underPrice:double - Price of the underlying.
*/
func (ic *IbClient) CalculateImpliedVolatility(reqID int64, contract *Contract, optionPrice float64, underPrice float64, impVolOptions []TagValue) {
	if ic.serverVersion < MIN_SERVER_VER_REQ_CALC_IMPLIED_VOLAT {
		ic.wrapper.error(NO_VALID_ID, UPDATE_TWS.code, UPDATE_TWS.msg+"  It does not support calculateImpliedVolatility req.")
		return
	}

	if ic.serverVersion < MIN_SERVER_VER_TRADING_CLASS && contract.TradingClass != "" {
		ic.wrapper.error(NO_VALID_ID, UPDATE_TWS.code, UPDATE_TWS.msg+"  It does not support tradingClass parameter in calculateImpliedVolatility.")
		return
	}

	v := 3

	fields := make([]interface{}, 0)
	fields = append(fields,
		REQ_CALC_IMPLIED_VOLAT,
		v,
		reqID,
		contract.ContractID,
		Contract.Symbol,
		contract.SecurityID,
		contract.Expiry,
		contract.Strike,
		contract.Right,
		contract.Multiplier,
		contract.Exchange,
		contract.PrimaryExchange,
		contract.Currency,
		contract.LocalSymbol)

	if ic.serverVersion >= MIN_SERVER_VER_TRADING_CLASS {
		fields = append(fields, contract.TradingClass)
	}

	fields = append(fields, optionPrice, underPrice)

	if ic.serverVersion >= MIN_SERVER_VER_LINKING {
		var implVolOptBuffer bytes.Buffer
		tagValuesCount := len(impVolOptions)
		fields = append(fields, tagValuesCount)
		for _, tv := range impVolOptions {
			implVolOptBuffer.WriteString(tv.Tag)
			implVolOptBuffer.WriteString("=")
			implVolOptBuffer.WriteString(tv.Value)
			implVolOptBuffer.WriteString(";")
		}
		fields = append(fields, implVolOptBuffer.Bytes())
	}

	msg := makeMsgBuf(fields...)

	ic.reqChan <- msg
}

func (ic *IbClient) CalculateOptionPrice(reqID int64, contract *Contract, volatility float64, underPrice float64, optPrcOptions []TagValue) {

	if ic.serverVersion < MIN_SERVER_VER_REQ_CALC_IMPLIED_VOLAT {
		ic.wrapper.error(reqID, UPDATE_TWS.code, UPDATE_TWS.msg+"  It does not support calculateImpliedVolatility req.")
		return
	}

	if ic.serverVersion < MIN_SERVER_VER_TRADING_CLASS {
		ic.wrapper.error(reqID, UPDATE_TWS.code, UPDATE_TWS.msg+"  It does not support tradingClass parameter in calculateImpliedVolatility.")
		return
	}

	v := 3
	fields := make([]interface{}, 0)
	fields = append(fields,
		REQ_CALC_OPTION_PRICE,
		v,
		reqID,
		contract.ContractID,
		Contract.Symbol,
		contract.SecurityID,
		contract.Expiry,
		contract.Strike,
		contract.Right,
		contract.Multiplier,
		contract.Exchange,
		contract.PrimaryExchange,
		contract.Currency,
		contract.LocalSymbol)

	if ic.serverVersion >= MIN_SERVER_VER_TRADING_CLASS {
		fields = append(fields, contract.TradingClass)
	}

	fields = append(fields, volatility, underPrice)

	if ic.serverVersion >= MIN_SERVER_VER_LINKING {
		var optPrcOptBuffer bytes.Buffer
		tagValuesCount := len(impVolOptions)
		fields = append(fields, tagValuesCount)
		for _, tv := range impVolOptions {
			optPrcOptBuffer.WriteString(tv.Tag)
			optPrcOptBuffer.WriteString("=")
			optPrcOptBuffer.WriteString(tv.Value)
			optPrcOptBuffer.WriteString(";")
		}

		fields = append(fields, optPrcOptBuffer.Bytes())
	}

	msg := makeMsgBuf(fields...)

	ic.reqChan <- msg
}

func (ic *IbClient) CancelCalculateOptionPrice(reqID int64) {
	if ic.serverVersion < MIN_SERVER_VER_REQ_CALC_IMPLIED_VOLAT {
		ic.wrapper.error(reqID, UPDATE_TWS.code, UPDATE_TWS.msg+"  It does not support calculateImpliedVolatility req.")
		return
	}

	v := 1
	msg := makeMsgBuf(CANCEL_CALC_OPTION_PRICE, v, reqID)

	ic.reqChan <- msg
}

/*ExerciseOptions
reqId:TickerId - The ticker id. multipleust be a unique value.
        contract:Contract - This structure contains a description of the
            contract to be exercised
        exerciseAction:int - Specifies whether you want the option to lapse
            or be exercised.
            Values are 1 = exercise, 2 = lapse.
        exerciseQuantity:int - The quantity you want to exercise.
        account:str - destination account
        override:int - Specifies whether your setting will override the system's
            natural action. For example, if your action is "exercise" and the
            option is not in-the-money, by natural action the option would not
            exercise. If you have override set to "yes" the natural action would
             be overridden and the out-of-the money option would be exercised.
            Values are: 0 = no, 1 = yes.
*/
func (ic *IbClient) ExerciseOptions(reqID int64, contract *Contract, exerciseAction int, exerciseQuantity int, account string, override int) {
	if ic.serverVersion < MIN_SERVER_VER_TRADING_CLASS && contract.TradingClass != "" {
		ic.wrapper.error(NO_VALID_ID, UPDATE_TWS.code, UPDATE_TWS.msg+"  It does not support conId, multiplier, tradingClass parameter in exerciseOptions.")
		return
	}

	v := 2
	fields := make([]interface{}, 0)

	fields = append(fields, EXERCISE_OPTIONS, v, reqID)

	if ic.serverVersion >= MIN_SERVER_VER_TRADING_CLASS {
		fields := append(fields, contract.ContractID)
	}

	fields = append(fields,
		contract.Symbol,
		contract.Expiry,
		contract.Strike,
		contract.Right,
		contract.Multiplier,
		contract.Exchange,
		contract.Currency,
		contract.LocalSymbol)

	if ic.serverVersion >= MIN_SERVER_VER_TRADING_CLASS {
		fields := append(fields, contract.TradingClass)
	}

	fields = append(fields,
		exerciseAction,
		exerciseQuantity,
		account,
		override)

	msg := makeMsgBuf(fields...)

	ic.reqChan <- msg

}

/*
   #########################################################################
   ################## Orders
   ########################################################################
*/

/*PlaceOrder
Call this function to place an order. The order status will
        be returned by the orderStatus event.

        orderId:OrderId - The order id. You must specify a unique value. When the
            order START_APItus returns, it will be identified by this tag.
            This tag is also used when canceling the order.
        contract:Contract - This structure contains a description of the
            contract which is being traded.
        order:Order - This structure contains the details of tradedhe order.
            Note: Each client MUST connect with a unique clientId.
*/
func (ic *IbClient) PlaceOrder(orderID int64, contract *Contract, order *Order) {
	switch v := ic.serverVersion; {
	case v < MIN_SERVER_VER_DELTA_NEUTRAL && contract.DeltaNeutralContract != nil:
		ic.wrapper.error(orderID, UPDATE_TWS.code, UPDATE_TWS.msg+"  It does not support delta-neutral orders.")
		return
	case v < MIN_SERVER_VER_SCALE_ORDERS2 && order.ScaleSubsLevelSize != UNSETINT:
		ic.wrapper.error(orderID, UPDATE_TWS.code, UPDATE_TWS.msg+"  It does not support Subsequent Level Size for Scale orders.")
		return
	case v < MIN_SERVER_VER_ALGO_ORDERS && order.AlgoStrategy != "":
		ic.wrapper.error(orderID, UPDATE_TWS.code, UPDATE_TWS.msg+"  It does not support algo orders.")
		return
	case v < MIN_SERVER_VER_NOT_HELD && order.NotHeld:
		ic.wrapper.error(orderID, UPDATE_TWS.code, UPDATE_TWS.msg+"  It does not support notHeld parameter.")
		return
	case v < MIN_SERVER_VER_SEC_ID_TYPE && (contract.SecurityType != "" || contract.SecurityID != ""):
		ic.wrapper.error(orderID, UPDATE_TWS.code, UPDATE_TWS.msg+"  It does not support secIdType and secId parameters.")
		return
	case v < MIN_SERVER_VER_PLACE_ORDER_CONID && contract.ContractID != UNSETINT && contract.ContractID > 0:
		ic.wrapper.error(orderID, UPDATE_TWS.code, UPDATE_TWS.msg+"  It does not support conId parameter.")
		return
	case v < MIN_SERVER_VER_SSHORTX && order.ExemptCode != -1:
		ic.wrapper.error(orderID, UPDATE_TWS.code, UPDATE_TWS.msg+"  It does not support exemptCode parameter.")
		return
	case v < MIN_SERVER_VER_SSHORTX:
		for _, comboLeg := range contract.ComboLegs {
			if comboLeg.ExemptCode != -1 {
				ic.wrapper.error(orderID, UPDATE_TWS.code, UPDATE_TWS.msg+"  It does not support exemptCode parameter.")
				return
			}
		}
		fallthrough
	case v < MIN_SERVER_VER_HEDGE_ORDERS && order.HedgeType != "":
		ic.wrapper.error(orderID, UPDATE_TWS.code, UPDATE_TWS.msg+"  It does not support hedge orders.")
		return
	case v < MIN_SERVER_VER_OPT_OUT_SMART_ROUTING && order.OptOutSmartRouting:
		ic.wrapper.error(orderID, UPDATE_TWS.code, UPDATE_TWS.msg+"  It does not support optOutSmartRouting parameter.")
		return
	case v < MIN_SERVER_VER_DELTA_NEUTRAL_CONID:
		if order.DeltaNeutralContractID > 0 || order.DeltaNeutralSettlingFirm != "" || order.DeltaNeutralClearingAccount != "" || order.DeltaNeutralClearingIntent != "" {
			ic.wrapper.error(orderID, UPDATE_TWS.code, UPDATE_TWS.msg+"  It does not support deltaNeutral parameters: ConId, SettlingFirm, ClearingAccount, ClearingIntent.")
			return
		}
		fallthrough
	case v < MIN_SERVER_VER_DELTA_NEUTRAL_OPEN_CLOSE:
		if order.DeltaNeutralOpenClose != "" ||
			order.DeltaNeutralShortSale ||
			order.DeltaNeutralShortSaleSlot > 0 ||
			order.DeltaNeutralDesignatedLocation != "" {
			ic.wrapper.error(orderID, UPDATE_TWS.code, UPDATE_TWS.msg+"  It does not support deltaNeutral parameters: OpenClose, ShortSale, ShortSaleSlot, DesignatedLocation.")
			return
		}
		fallthrough
	case v < MIN_SERVER_VER_SCALE_ORDERS3:
		if (order.ScalePriceIncrement > 0 && order.ScalePriceIncrement != UNSETFLOAT) &&
			(order.ScalePriceAdjustValue != UNSETFLOAT ||
				order.ScalePriceAdjustInterval != UNSETINT ||
				order.ScaleProfitOffset != UNSETFLOAT ||
				order.ScaleAutoReset ||
				order.ScaleInitPosition != UNSETINT ||
				order.ScaleInitFillQty != UNSETINT ||
				order.ScaleRandomPercent) {
			ic.wrapper.error(orderID, UPDATE_TWS.code, UPDATE_TWS.msg+
				"  It does not support Scale order parameters: PriceAdjustValue, PriceAdjustInterval, "+
				"ProfitOffset, AutoReset, InitPosition, InitFillQty and RandomPercent.")
			return
		}
		fallthrough
	case v < MIN_SERVER_VER_ORDER_COMBO_LEGS_PRICE && contract.SecurityType == "BAG":
		for _, orderComboLeg := range order.OrderComboLegs {
			if orderComboLeg.Price != UNSETFLOAT {
				ic.wrapper.error(orderID, UPDATE_TWS.code, UPDATE_TWS.msg+"  It does not support per-leg prices for order combo legs.")
				return
			}

		}
		fallthrough
	case v < MIN_SERVER_VER_TRAILING_PERCENT && order.TrailingPercent != UNSETFLOAT:
		ic.wrapper.error(orderID, UPDATE_TWS.code, UPDATE_TWS.msg+"  It does not support trailing percent parameter.")
		return
	case v < MIN_SERVER_VER_TRADING_CLASS && contract.TradingClass != "":
		ic.wrapper.error(orderID, UPDATE_TWS.code, UPDATE_TWS.msg+"  It does not support tradingClass parameter in placeOrder.")
		return
	case v < MIN_SERVER_VER_SCALE_TABLE &&
		(order.ScaleTable != "" ||
			order.ActiveStartTime != "" ||
			order.ActiveStopTime != ""):
		ic.wrapper.error(orderID, UPDATE_TWS.code, UPDATE_TWS.msg+"  It does not support scaleTable, activeStartTime and activeStopTime parameters.")
		return
	case v < MIN_SERVER_VER_ALGO_ID && order.AlgoID != "":
		ic.wrapper.error(orderID, UPDATE_TWS.code, UPDATE_TWS.msg+"  It does not support algoId parameter.")
		return
	case v < MIN_SERVER_VER_ORDER_SOLICITED && order.Solictied:
		ic.wrapper.error(orderID, UPDATE_TWS.code, UPDATE_TWS.msg+"  It does not support order solicited parameter.")
		return
	case v < MIN_SERVER_VER_MODELS_SUPPORT && order.ModelCode != "":
		ic.wrapper.error(orderID, UPDATE_TWS.code, UPDATE_TWS.msg+"  It does not support model code parameter.")
		return
	case v < MIN_SERVER_VER_EXT_OPERATOR && order.ExtOperator != "":
		ic.wrapper.error(orderID, UPDATE_TWS.code, UPDATE_TWS.msg+"  It does not support ext operator parameter")
		return
	case v < MIN_SERVER_VER_SOFT_DOLLAR_TIER &&
		(order.SoftDollarTier.Name != "" || order.SoftDollarTier.Value != ""):
		ic.wrapper.error(orderID, UPDATE_TWS.code, UPDATE_TWS.msg+" It does not support soft dollar tier")
		return
	case v < MIN_SERVER_VER_CASH_QTY && order.CashQty != UNSETFLOAT:
		ic.wrapper.error(orderID, UPDATE_TWS.code, UPDATE_TWS.msg+" It does not support cash quantity parameter")
		return
	case v < MIN_SERVER_VER_DECISION_MAKER &&
		(order.Mifid2DecisionMaker != "" || order.Mifid2DecisionAlgo != ""):
		ic.wrapper.error(orderID, UPDATE_TWS.code, UPDATE_TWS.msg+" It does not support MIFID II decision maker parameters")
		return
	case v < MIN_SERVER_VER_MIFID_EXECUTION &&
		(order.Mifid2ExecutionTrader != "" || order.Mifid2ExecutionAlgo != ""):
		ic.wrapper.error(orderID, UPDATE_TWS.code, UPDATE_TWS.msg+" It does not support MIFID II execution parameters")
		return
	case v < MIN_SERVER_VER_AUTO_PRICE_FOR_HEDGE && order.DontUseAutoPriceForHedge:
		ic.wrapper.error(orderID, UPDATE_TWS.code, UPDATE_TWS.msg+" It does not support dontUseAutoPriceForHedge parameter")
		return
	case v < MIN_SERVER_VER_ORDER_CONTAINER && order.IsOmsContainer:
		ic.wrapper.error(orderID, UPDATE_TWS.code, UPDATE_TWS.msg+" It does not support oms container parameter")
		return
	}

	var v int
	if ic.serverVersion < MIN_SERVER_VER_NOT_HELD {
		v = 27
	} else {
		v = 45
	}

	fields := make([]interface{}, 0)
	fields = append(fields, PLACE_ORDER)

	if ic.serverTime < MIN_SERVER_VER_ORDER_CONTAINER {
		fields = append(fields, v)
	}

	fields = append(fields, orderID)

	if ic.serverVersion >= MIN_SERVER_VER_PLACE_ORDER_CONID {
		fields = append(fields, contract.ContractID)
	}

	fields = append(fields,
		contract.Symbol,
		contract.SecurityID,
		contract.Expiry,
		contract.Strike,
		contract.Right,
		contract.Multiplier,
		contract.Exchange,
		contract.PrimaryExchange,
		contract.Currency,
		contract.LocalSymbol)

	if ic.serverTime >= MIN_SERVER_VER_TRADING_CLASS {
		fields = append(fields, contract.TradingClass)
	}

	if ic.serverTime >= MIN_SERVER_VER_SEC_ID_TYPE {
		fields = append(fields, contract.SecurityIDType, contract.SecurityID)
	}

	fields = append(fields, order.Action)

	if ic.serverVersion >= MIN_SERVER_VER_FRACTIONAL_POSITIONS {
		fields = append(fields, order.TotalQuantity)
	} else {
		fields = append(fields, int64(order.TotalQuantity))
	}

	fields = append(fields, order.OrderType)

	if ic.serverVersion < MIN_SERVER_VER_ORDER_COMBO_LEGS_PRICE {
		if order.LimitPrice != UNSETFLOAT {
			fields = append(fields, order.LimitPrice)
		} else {
			fields = append(fields, float64(0))
		}
	} else {
		fields = append(fields, handleEmpty(order.LimitPrice))
	}

	if ic.serverVersion < MIN_SERVER_VER_TRAILING_PERCENT {
		if order.AuxPrice != UNSETFLOAT {
			fields = append(fields, order.AuxPrice)
		} else {
			fields = append(fields, handleEmpty(order.AuxPrice))
		}
	} else {
		fields = append(fields, "")
	}

	fields = append(fields,
		order.TIF,
		order.OCAGroup,
		order.Account,
		order.OpenClose,
		order.Origin,
		order.OrderRef,
		order.Transmit,
		order.ParentID,
		order.BlockOrder,
		order.SweepToFill,
		order.DisplaySize,
		order.TriggerMethod,
		order.OutsideRTH,
		order.Hidden)

	if contract.SecurityType == "BAG" {
		comboLegsCount := len(contract.ComboLegs)
		fields = append(fields, comboLegsCount)
		for _, comboLeg := range contract.ComboLegs {
			fields = append(fields,
				comboLeg.ContractID,
				comboLeg.Ratio,
				comboLeg.Action,
				comboLeg.Exchange,
				comboLeg.OpenClose,
				comboLeg.ShortSaleSlot,
				comboLeg.DesignatedLocation)
			if ic.serverVersion >= MIN_SERVER_VER_SSHORTX_OLD {
				fields = append(fields, comboLeg.ExemptCode)
			}
		}
	}

	if ic.serverVersion >= MIN_SERVER_VER_ORDER_COMBO_LEGS_PRICE && contract.SecurityType == "BAG" {
		orderComboLegsCount := len(order.OrderComboLegs)
		fields = append(fields, orderComboLegsCount)
		for _, orderComboLeg := range order.OrderComboLegs {
			fields = append(fields, handleEmpty(orderComboLeg.Price))
		}
	}

	if ic.serverVersion >= MIN_SERVER_VER_SMART_COMBO_ROUTING_PARAMS && contract.SecurityType == "BAG" {
		smartComboRoutingParamsCount := len(order.SmartComboRoutingParams)
		fields = append(fields, smartComboRoutingParamsCount)
		for _, tv := range order.SmartComboRoutingParams {
			fields = append(fields, tv.Tag, tv.Value)
		}
	}

	fields = append(fields,
		"",
		order.DiscretionaryAmount,
		order.GoodAfterTime,
		order.GoodTillDate,

		order.FAGroup,
		order.FAMethod,
		order.FAPercentage,
		order.FAProfile)

	if ic.serverVersion >= MIN_SERVER_VER_MODELS_SUPPORT {
		fields = append(fields, order.ModelCode)
	}

	fields = append(fields,
		order.ShortSaleSlot,
		order.DesignatedLocation)

	//institutional short saleslot data (srv v18 and above)
	if ic.serverVersion >= MIN_SERVER_VER_SSHORTX_OLD {
		fields = append(fields, order.ExemptCode)
	}

	fields = append(fields, order.OCAType)

	fields = append(fields,
		order.Rule80A,
		order.SettlingFirm,
		order.AllOrNone,
		handleEmpty(order.MinQty),
		handleEmpty(order),
		order.ETradeOnly,
		order.FirmQuoteOnly,
		handleEmpty(order.NBBOPriceCap),
		order.AuctionStrategy,
		handleEmpty(order.StartingPrice),
		handleEmpty(order.StockRefPrice),
		handleEmpty(order.Delta),
		handleEmpty(order.StockRangeLower),
		handleEmpty(order.StockRangeUpper),

		order.OverridePercentageConstraints,

		handleEmpty(order.Volatility),
		handleEmpty(order.VolatilityType),
		order.DeltaNeutralOrderType,
		handleEmpty(order.DeltaNeutralAuxPrice))

	if ic.serverVersion >= MIN_SERVER_VER_DELTA_NEUTRAL_CONID && order.DeltaNeutralOrderType != "" {
		fields = append(fields,
			order.DeltaNeutralContractID,
			order.DeltaNeutralSettlingFirm,
			order.DeltaNeutralClearingAccount,
			order.DeltaNeutralClearingIntent)
	}

	if ic.serverVersion >= MIN_SERVER_VER_DELTA_NEUTRAL_OPEN_CLOSE && order.DeltaNeutralOrderType != "" {
		fields = append(fields,
			order.DeltaNeutralOpenClose,
			order.DeltaNeutralShortSale,
			order.DeltaNeutralShortSaleSlot,
			order.DeltaNeutralDesignatedLocation)
	}

	fields = append(fields,
		order.ContinuousUpdate,
		handleEmpty(order.ReferencePriceType),
		handleEmpty(order.TrailStopPrice))

	if ic.serverVersion >= MIN_SERVER_VER_TRAILING_PERCENT {
		fields = append(fields, handleEmpty(order.TrailingPercent))
	}

	//scale orders
	if ic.serverVersion >= MIN_SERVER_VER_SCALE_ORDERS2 {
		fields = append(fields,
			handleEmpty(order.ScaleInitLevelSize),
			handleEmpty(order.ScaleSubsLevelSize))
	} else {
		fields = append(fields,
			"",
			handleEmpty(order.ScaleInitLevelSize))
	}

	if ic.serverVersion >= MIN_SERVER_VER_SCALE_ORDERS3 && order.ScalePriceIncrement != UNSETFLOAT && order.ScalePriceIncrement > 0.0 {
		fields = append(fields,
			handleEmpty(order.ScalePriceAdjustValue),
			handleEmpty(order.ScalePriceAdjustInterval),
			handleEmpty(order.ScaleProfitOffset),
			order.ScaleAutoReset,
			handleEmpty(order.ScaleInitPosition),
			handleEmpty(order.ScaleInitFillQty),
			order.ScaleRandomPercent)
	}

	if ic.serverVersion >= MIN_SERVER_VER_SCALE_TABLE {
		fields = append(fields,
			order.ScaleTable,
			order.ActiveStartTime,
			order.ActiveStopTime)
	}

	//hedge orders
	if ic.serverVersion >= MIN_SERVER_VER_HEDGE_ORDERS {
		fields = append(fields, order.HedgeType)
		if order.HedgeType != "" {
			fields = append(fields, order.HedgeParam)
		}
	}

	if ic.serverVersion >= MIN_SERVER_VER_OPT_OUT_SMART_ROUTING {
		fields = append(fields, order.HedgeType)
	}

	if ic.serverVersion >= MIN_SERVER_VER_PTA_ORDERS {
		fields = append(fields,
			order.ClearingAccount,
			order.ClearingIntent)
	}

	if ic.serverVersion >= MIN_SERVER_VER_NOT_HELD {
		fields = append(fields, order.NotHeld)
	}

	if ic.serverVersion >= MIN_SERVER_VER_DELTA_NEUTRAL {
		if contract.DeltaNeutralContract != nil {
			fields = append(fields,
				true,
				contract.DeltaNeutralContract.ContractID,
				contract.DeltaNeutralContract.Delta,
				contract.DeltaNeutralContract.Price)
		} else {
			fields = append(fields, false)
		}
	}

	if ic.serverVersion >= MIN_SERVER_VER_ALGO_ORDERS {
		fields = append(fields, order.AlgoStrategy)

		if order.AlgoStrategy != "" {
			algoParamsCount := len(order.AlgoParams)
			fields = append(fields, algoParamsCount)
			for _, tv := range order.AlgoParams {
				fields = append(fields, tv.Tag, tv.Value)
			}
		}
	}

	if ic.serverVersion >= MIN_SERVER_VER_ALGO_ID {
		fields = append(fields, order.AlgoID)
	}

	fields = append(fields, order.WhatIf)

	if ic.serverVersion >= MIN_SERVER_VER_LINKING {
		var miscOptionsBuffer bytes.Buffer
		for _, tv := range order.OrderMiscOptions {
			miscOptionsBuffer.WriteString(tv.Tag)
			miscOptionsBuffer.WriteString("=")
			miscOptionsBuffer.WriteString(tv.Value)
			miscOptionsBuffer.WriteString(";")
		}

		fields = append(fields, miscOptionsBuffer.Bytes())
	}

	if ic.serverVersion >= MIN_SERVER_VER_ORDER_SOLICITED {
		fields = append(fields, order.Solictied)
	}

	if ic.serverVersion >= MIN_SERVER_VER_RANDOMIZE_SIZE_AND_PRICE {
		fields = append(fields,
			order.RandomizeSize,
			order.RandomizePrice)
	}

	if ic.serverTime >= MIN_SERVER_VER_PEGGED_TO_BENCHMARK {
		if order.OrderType == "PEG BENCH" {
			fields = append(fields,
				order.ReferenceContractID,
				order.IsPeggedChangeAmountDecrease,
				order.PeggedChangeAmount,
				order.ReferenceChangeAmount,
				order.ReferenceExchangeID)
		}

		orderConditionsCount := len(order.Conditions)
		fields = append(fields, orderConditionsCount)
		for _, cond := range order.Conditions {
			fields = append(fields,
				cond.CondType,
				cond.toFields()...)
		}
		if orderConditionsCount > 0 {
			fields = append(fields,
				order.ConditionsIgnoreRth,
				order.ConditionsCancelOrder)
		}

		fields = append(fields,
			order.AdjustedOrderType,
			order.TriggerPrice,
			order.LimitPrice,
			order.AdjustedStopPrice,
			order.AdjustedStopLimitPrice,
			order.AdjustedTrailingAmount,
			order.AdjustableTrailingUnit)

		if ic.serverVersion >= MIN_SERVER_VER_EXT_OPERATOR {
			fields = append(fields, order.ExtOperator)
		}

		if ic.serverVersion >= MIN_SERVER_VER_SOFT_DOLLAR_TIER {
			fields = append(fields, order.SoftDollarTier.Name, order.SoftDollarTier.Value)
		}

		if ic.serverVersion >= MIN_SERVER_VER_CASH_QTY {
			fields = append(fields, order.CashQty)
		}

		if ic.serverVersion >= MIN_SERVER_VER_DECISION_MAKER {
			fields = append(fields, order.Mifid2DecisionMaker, order.Mifid2DecisionAlgo)
		}

		if ic.serverVersion >= MIN_SERVER_VER_MIFID_EXECUTION {
			fields = append(fields, order.Mifid2ExecutionTrader, order.Mifid2ExecutionAlgo)
		}

		if ic.serverVersion >= MIN_SERVER_VER_AUTO_PRICE_FOR_HEDGE {
			fields = append(fields, order.DontUseAutoPriceForHedge)
		}

		if ic.serverVersion >= MIN_SERVER_VER_ORDER_CONTAINER {
			fields = append(fields, order.IsOmsContainer)
		}

		if ic.serverVersion >= MIN_SERVER_VER_D_PEG_ORDERS {
			fields = append(fields, order.DiscretionaryUpToLimitPrice)
		}

		msg := makeMsgBuf(fields...)

		ic.reqChan <- msg
	}

}

func (ic *IbClient) CancelOrder(orderID int64) {
	v := 1
	msg := makeMsgBuf(CANCEL_ORDER, v, orderID)
	ic.reqChan <- msg
}

func (ic *IbClient) ReqOpenOrders() {
	v := 1
	msg := makeMsgBuf(REQ_OPEN_ORDERS, v)
	ic.reqChan <- msg
}

// ReqAutoOpenOrders will make the client access to the TWS Orders (only if clientId=0)
func (ic *IbClient) ReqAutoOpenOrders(autoBind bool) {
	v := 1
	msg := makeMsgBuf(REQ_AUTO_OPEN_ORDERS, v, autoBind)

	ic.reqChan <- msg
}

func (ic *IbClient) ReqAllOpenOrders() {
	v := 1
	msg := makeMsgBuf(REQ_ALL_OPEN_ORDERS, v)

	ic.reqChan <- msg
}

func (ic *IbClient) ReqGlobalCancel() {
	v := 1
	msg := makeMsgBuf(REQ_GLOBAL_CANCEL, v)

	ic.reqChan <- msg
}

func (ic *IbClient) ReqIDs(numIDs int) {
	v := 1
	msg := makeMsgBuf(REQ_IDS, v, numIDs)

	ic.reqChan <- msg
}

/*
   #########################################################################
   ################## Account and Portfolio
   ########################################################################
*/

func (ic *IbClient) ReqAccountUpdates(subscribe bool, accName string) {
	v := 2
	msg := makeMsgBuf(REQ_ACCT_DATA, v, subscribe, accName)

	ic.reqChan <- msg
}

/*ReqAccountSummary
Call this method to request and keep up to date the data that appears
        on the TWS Account Window Summary tab. The data is returned by
        accountSummary().

        Note:   This request is designed for an FA managed account but can be
        used for any multi-account structure.

        reqId:int - The ID of the data request. Ensures that responses are matched
            to requests If several requests are in process.
        groupName:str - Set to All to returnrn account summary data for all
            accounts, or set to a specific Advisor Account Group name that has
            already been created in TWS Global Configuration.
        tags:str - A comma-separated list of account tags.  Available tags are:
            accountountType
            NetLiquidation,
            TotalCashValue - Total cash including futures pnl
            SettledCash - For cash accounts, this is the same as
            TotalCashValue
            AccruedCash - Net accrued interest
            BuyingPower - The maximum amount of marginable US stocks the
                account can buy
            EquityWithLoanValue - Cash + stocks + bonds + mutual funds
            PreviousDayEquityWithLoanValue,
            GrossPositionValue - The sum of the absolute value of all stock
                and equity option positions
            RegTEquity,
            RegTMargin,
            SMA - Special Memorandum Account
            InitMarginReq,
            MaintMarginReq,
            AvailableFunds,
            ExcessLiquidity,
            Cushion - Excess liquidity as a percentage of net liquidation value
            FullInitMarginReq,
            FullMaintMarginReq,
            FullAvailableFunds,
            FullExcessLiquidity,
            LookAheadNextChange - Time when look-ahead values take effect
            LookAheadInitMarginReq,
            LookAheadMaintMarginReq,
            LookAheadAvailableFunds,
            LookAheadExcessLiquidity,
            HighestSeverity - A measure of how close the account is to liquidation
            DayTradesRemaining - The Number of Open/Close trades a user
                could put on before Pattern Day Trading is detected. A value of "-1"
                means that the user can put on unlimited day trades.
            Leverage - GrossPositionValue / NetLiquidation
            $LEDGER - Single flag to relay all cash balance tags*, only in base
                currency.
            $LEDGER:CURRENCY - Single flag to relay all cash balance tags*, only in
                the specified currency.
            $LEDGER:ALL - Single flag to relay all cash balance tags* in all
            currencies.
*/
func (ic *IbClient) ReqAccountSummary(reqID int64, groupName string, tags string) {
	v := 1
	msg := makeMsgBuf(REQ_ACCOUNT_SUMMARY, v, reqID, groupName, tags)

	ic.reqChan <- msg
}

func (ic *IbClient) CancelAccountSummary(reqID int64) {
	v := 1
	msg := makeMsgBuf(CANCEL_ACCOUNT_SUMMARY, v, reqID)

	ic.reqChan <- msg
}

func (ic *IbClient) ReqPositions() {
	if ic.serverVersion < MIN_SERVER_VER_POSITIONS {
		ic.wrapper.error(NO_VALID_ID, UPDATE_TWS.code, UPDATE_TWS.msg+"  It does not support positions request.")
		return
	}
	v := 1
	msg := makeMsgBuf(REQ_POSITIONS, v)

	ic.reqChan <- msg
}

func (ic *IbClient) CancelPositions() {
	if ic.serverVersion < MIN_SERVER_VER_POSITIONS {
		ic.wrapper.error(NO_VALID_ID, UPDATE_TWS.code, UPDATE_TWS.msg+"  It does not support positions request.")
		return
	}

	v := 1
	msg := makeMsgBuf(CANCEL_POSITIONS, v)

	ic.reqChan <- msg
}

func (ic *IbClient) ReqPositionsMulti(reqID int64, account string, modelCode string) {
	if ic.serverVersion < MIN_SERVER_VER_MODELS_SUPPORT {
		ic.wrapper.error(NO_VALID_ID, UPDATE_TWS.code, UPDATE_TWS.msg+"  It does not support positions multi request.")
		return
	}
	v := 1
	msg := makeMsgBuf(REQ_POSITIONS_MULTI, v, reqID, account, modelCode)

	ic.reqChan <- msg
}

func (ic *IbClient) CancelPositionsMulti(reqID int64) {
	if ic.serverVersion < MIN_SERVER_VER_MODELS_SUPPORT {
		ic.wrapper.error(NO_VALID_ID, UPDATE_TWS.code, UPDATE_TWS.msg+"  It does not support cancel positions multi request.")
		return
	}

	v := 1
	msg := makeMsgBuf(CANCEL_POSITIONS_MULTI, v, reqID)

	ic.reqChan <- msg
}

func (ic *IbClient) ReqAccountUpdatesMulti(reqID int64, account string, modelCode string, ledgerAndNLV bool) {
	if ic.serverVersion < MIN_SERVER_VER_MODELS_SUPPORT {
		ic.wrapper.error(NO_VALID_ID, UPDATE_TWS.code, UPDATE_TWS.msg+"  It does not support account updates multi request.")
		return
	}

	v := 1
	msg := makeMsgBuf(REQ_ACCOUNT_UPDATES_MULTI, v, reqID, account, modelCode, ledgerAndNLV)

	ic.reqChan <- msg
}

//ReqCurrentTime Asks the current system time on the server side.
func (ic *IbClient) ReqCurrentTime() {
	v := 1
	msg := makeMsgBuf(REQ_CURRENT_TIME, v)

	ic.reqChan <- msg
}

func (ic *IbClient) ReqExecutions(reqID int64, execFilter ExecutionFilter) {
	v := 3
	fields := make([]interface{}, 0)
	fields = append(fields, REQ_EXECUTIONS, v)

	if ic.serverVersion >= MIN_SERVER_VER_EXECUTION_DATA_CHAIN {
		fields = append(fields, reqID)
	}

	fields = append(fields,
		execFilter.ClientID,
		execFilter.AccountCode,
		execFilter.Time,
		execFilter.Symbol,
		execFilter.SecurityType,
		execFilter.Exchange,
		execFilter.Side)
	msg := makeMsgBuf(fields...)

	ic.reqChan <- msg
}

func (ic *IbClient) ReqHistoricalData(reqID int64, contract Contract, endDateTime string, duration string, barSize string, whatToShow string, useRTH bool, formatDate int, keepUpToDate bool, chartOptions []TagValue) {
	if ic.serverVersion < MIN_SERVER_VER_TRADING_CLASS {
		if contract.TradingClass != "" || contract.ContractID > 0 {
			ic.wrapper.error(reqID, UPDATE_TWS.code, UPDATE_TWS.msg)
		}
	}

	v := 6

	fields := make([]interface{}, 0)
	fields = append(fields, REQ_HISTORICAL_DATA)
	if ic.serverVersion <= MIN_SERVER_VER_SYNT_REALTIME_BARS {
		fields = append(fields, v)
	}

	fields = append(fields, reqID)

	if ic.serverVersion >= MIN_SERVER_VER_TRADING_CLASS {
		fields = append(fields, contract.ContractID)
	}

	fields = append(fields,
		contract.Symbol,
		contract.SecurityType,
		contract.Expiry,
		contract.Strike,
		contract.Right,
		contract.Multiplier,
		contract.Exchange,
		contract.PrimaryExchange,
		contract.Currency,
		contract.LocalSymbol,
	)

	if ic.serverVersion >= MIN_SERVER_VER_TRADING_CLASS {
		fields = append(fields, contract.TradingClass)
	}
	fields = append(fields,
		contract.IncludeExpired,
		endDateTime,
		barSize,
		duration,
		useRTH,
		whatToShow,
		formatDate,
	)

	if contract.SecurityType == "BAG" {
		fields = append(fields, len(contract.ComboLegs))
		for _, comboLeg := range contract.ComboLegs {
			fields = append(fields,
				comboLeg.ContractID,
				comboLeg.Ratio,
				comboLeg.Action,
				comboLeg.Exchange,
			)
		}
	}

	if ic.serverVersion >= MIN_SERVER_VER_SYNT_REALTIME_BARS {
		fields = append(fields, keepUpToDate)
	}

	if ic.serverVersion >= MIN_SERVER_VER_LINKING {
		chartOptionsStr := ""
		for _, tagValue := range chartOptions {
			chartOptionsStr += tagValue.Value
		}
		fields = append(fields, chartOptionsStr)
	}

	msg := makeMsgBuf(fields...)
	// fmt.Println(msg)

	ic.reqChan <- msg
}

//--------------------------three major goroutine -----------------------------------------------------
//goRequest will get the req from reqChan and send it to TWS
func (ic *IbClient) goRequest() {
	log.Println("Start goRequest!")
	defer log.Println("End goRequest!")
	defer ic.wg.Done()
requestLoop:
	for {
		select {
		case req := <-ic.reqChan:
			_, err := ic.writer.Write(req)
			if err != nil {
				ic.writer.Reset(ic.conn)
				ic.errChan <- err
			}
			ic.writer.Flush()
		case <-ic.terminatedSignal:
			// fmt.Println("goRequest terminate")
			break requestLoop
		}
	}

}

//goReceive receive the msg from the socket, get the fields and put them into msgChan
//goReceive handle the msgBuf which is different from the offical.Not continuously read, but split first and then decode
func (ic *IbClient) goReceive() {
	// defer
	log.Println("Start goReceive!")
	// buf := make([]byte, 0, 4096)
	defer log.Println("End goReceive!")
	defer ic.wg.Done()
	for {
		// buf := []byte
		msgBuf, err := readMsgBuf(ic.reader)
		// fmt.Printf("msgBuf: %v err: %v", msgBuf, err)
		if err, ok := err.(*net.OpError); ok {
			if !err.Temporary() {
				break
			}
			log.Println(err)
		} else if err != nil {
			ic.errChan <- err
			ic.reader.Reset(ic.conn)
		}

		if msgBuf != nil {
			fields := splitMsgBuf(msgBuf)
			ic.msgChan <- fields
		}

	}
}

//goDecode decode the fields received from the msgChan
func (ic *IbClient) goDecode() {
	log.Println("Start goDecode!")
	defer log.Println("End goDecode!")
	defer ic.wg.Done()

decodeLoop:
	for {
		// buf := []byte
		select {
		case f := <-ic.msgChan:
			ic.decoder.interpret(f...)
			// log.Println(f)
		case e := <-ic.errChan:
			fmt.Println(e)
		case <-ic.terminatedSignal:
			// fmt.Println("goDecode terminate")
			break decodeLoop
		}
	}

}

// ---------------------------------------------------------------------------------------

// Run make the event loop run, all make sense after run!
func (ic *IbClient) Run() error {
	if ic.conn.state == DISCONNECTED {
		return errors.New("ibClient is DISCONNECTED")
	}
	log.Println("RUN Client")
	ic.wg.Add(3)
	go ic.goRequest()
	go ic.goReceive()
	go ic.goDecode()

	return nil
}
