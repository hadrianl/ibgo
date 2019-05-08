package IBAlgoTrade

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	. "github.com/hadrianl/ibgo/ibapi"
	log "github.com/sirupsen/logrus"
)

type Subscriber interface {
	// GetReqID() int64
}

type GoWrapper struct {
	ib               *IB
	Logger           *log.Logger
	mu               sync.RWMutex
	dataChanMap      map[int64]chan interface{}
	subDataChanMap   map[int64]chan interface{}
	reqID2Subscriber map[int64]Subscriber
	reqID2Contract   map[int64]Contract
	AccValues        map[string]*AccountValues
	AccSummary       map[string]*AccountSummary
	Portfolio        map[string]map[int64]*PortfolioItem
	Positions        map[string]map[int64]*Position

	Trades map[[2]int64]Trade
	Fills  map[string]Fill

	NewTicks      []NewsTick
	NewsBulletins map[int64]NewsBulletin

	Tickers map[int64]Ticker

	PnLs       map[int64]*PnL
	PnLSingles map[int64]*PnLSingle
}

func (w *GoWrapper) reset() {
	w.dataChanMap = make(map[int64]chan interface{})
	w.subDataChanMap = make(map[int64]chan interface{})
	w.AccValues = make(map[string]*AccountValues)
	w.AccSummary = make(map[string]*AccountSummary)
	w.AccSummary["All"] = &AccountSummary{Account: "All", TagValues: make(map[string][2]string)}
	w.Portfolio = make(map[string]map[int64]*PortfolioItem)
	w.Positions = make(map[string]map[int64]*Position)
	w.PnLs = make(map[int64]*PnL)
	w.Trades = make(map[[2]int64]Trade)
}

func (w *GoWrapper) startReq(reqID int64, chanSize int) {
	w.dataChanMap[reqID] = make(chan interface{}, chanSize)
}

func (w *GoWrapper) endReq(reqID int64) {
	if dataChan, ok := w.dataChanMap[reqID]; ok {
		close(dataChan)
		delete(w.dataChanMap, reqID)
	}

}

func (w *GoWrapper) startSubscription(reqID int64, chanSize int, contract Contract) {
	w.subDataChanMap[reqID] = make(chan interface{}, chanSize)
}

func (w *GoWrapper) endSubscription(reqID int64) {
	if subChan, ok := w.subDataChanMap[reqID]; ok {
		close(subChan)
		delete(w.subDataChanMap, reqID)
	}
}

func OrderKey(clientID, orderID, permID int64) [2]int64 {
	if orderID <= 0 {
		return [2]int64{-1, permID}
	}

	return [2]int64{clientID, orderID}
}

func (w *GoWrapper) ConnectAck() {
	log.Printf("<connectAck>...")
}

func (w *GoWrapper) NextValidID(reqID int64) {
	log.WithField("reqID", reqID).Printf("<nextValidID>: %v.", reqID)
	atomic.StoreInt64(&w.ib.reqID, reqID)
}

func (w *GoWrapper) ManagedAccounts(accountsList []string) {
	log.Printf("<managedAccounts>: %v.", accountsList)
	w.mu.RLock()
	defer w.mu.RUnlock()
	for _, acc := range accountsList {
		if _, ok := w.AccSummary[acc]; !ok {
			w.AccSummary[acc] = &AccountSummary{Account: acc, TagValues: make(map[string][2]string)}
		}

		if _, ok := w.AccValues[acc]; !ok {
			w.AccValues[acc] = &AccountValues{Account: acc, TagValues: make(map[string][3]string)}
		}

		if _, ok := w.Portfolio[acc]; !ok {
			w.Portfolio[acc] = map[int64]*PortfolioItem{}
		}
	}

}

func (w *GoWrapper) TickPrice(reqID int64, tickType int64, price float64, attrib TickAttrib) {
	log.WithField("reqID", reqID).Printf("<tickPrice>: tickType: %v price: %v.", tickType, price)
}

func (w *GoWrapper) UpdateAccountTime(accTime time.Time) {
	log.Printf("<updateAccountTime>: %v", accTime)

}

func (w *GoWrapper) UpdateAccountValue(tag string, value string, currency string, account string) {
	log.WithFields(log.Fields{"account": account, tag: value, "currency": currency}).Print("<updateAccountValue>")
	w.AccValues[account].TagValues[tag] = [3]string{value, currency, ""}
}

func (w *GoWrapper) AccountDownloadEnd(accName string) {
	log.Printf("<accountDownloadEnd>: %v", accName)
}

func (w *GoWrapper) AccountUpdateMulti(reqID int64, account string, modelCode string, tag string, value string, currency string) {
	log.WithFields(log.Fields{"reqID": reqID, "account": account, tag: value, "currency": currency, "modelCode": modelCode}).Print("<accountUpdateMulti>")
	w.AccValues[account].TagValues[tag] = [3]string{value, currency, modelCode}
}

func (w *GoWrapper) AccountUpdateMultiEnd(reqID int64) {
	log.WithField("reqID", reqID).Print("<accountUpdateMultiEnd>")
	w.endReq(reqID)
}

func (w *GoWrapper) AccountSummary(reqID int64, account string, tag string, value string, currency string) {
	log.WithFields(log.Fields{"reqID": reqID, "account": account, tag: value, "currency": currency}).Print("<accountSummary>")
	w.mu.Lock()
	defer w.mu.Unlock()
	w.AccSummary[account].TagValues[tag] = [2]string{value, currency}
	// c := w.dataChanMap[reqID]
	// if c != nil {
	// 	c <- map[string]string{"account": account, "tag": tag, "value": value, "currency": currency}
	// }
}

func (w *GoWrapper) AccountSummaryEnd(reqID int64) {
	log.WithField("reqID", reqID).Print("<accountSummaryEnd>")
	// c := w.dataChanMap[reqID]
	// if c != nil {
	// 	close(c)
	// }
}

func (w *GoWrapper) VerifyMessageAPI(apiData string) {
	log.Printf("<verifyMessageAPI>: apiData: %v", apiData)
}

func (w *GoWrapper) VerifyCompleted(isSuccessful bool, err string) {
	log.Printf("<verifyCompleted>: isSuccessful: %v error: %v", isSuccessful, err)
}

func (w *GoWrapper) VerifyAndAuthMessageAPI(apiData string, xyzChallange string) {
	log.Printf("<verifyCompleted>: apiData: %v xyzChallange: %v", apiData, xyzChallange)
}

func (w *GoWrapper) VerifyAndAuthCompleted(isSuccessful bool, err string) {
	log.Printf("<verifyAndAuthCompleted>: isSuccessful: %v error: %v", isSuccessful, err)
}

func (w *GoWrapper) DisplayGroupList(reqID int64, groups string) {
	log.WithField("reqID", reqID).Printf("<displayGroupList>: groups: %v", groups)
}

func (w *GoWrapper) DisplayGroupUpdated(reqID int64, contractInfo string) {
	log.WithField("reqID", reqID).Printf("<displayGroupUpdated>: contractInfo: %v", contractInfo)
}

func (w *GoWrapper) PositionMulti(reqID int64, account string, modelCode string, contract *Contract, position float64, avgCost float64) {
	log.WithField("reqID", reqID).Printf("<positionMulti>: account: %v modelCode: %v contract: <%v> position: %v avgCost: %v", account, modelCode, contract, position, avgCost)
}

func (w *GoWrapper) PositionMultiEnd(reqID int64) {
	log.WithField("reqID", reqID).Print("<positionMultiEnd>")
}

func (w *GoWrapper) UpdatePortfolio(contract *Contract, position float64, marketPrice float64, marketValue float64, averageCost float64, unrealizedPNL float64, realizedPNL float64, account string) {
	log.Printf("<updatePortfolio>: contract: %v pos: %v marketPrice: %v averageCost: %v unrealizedPNL: %v realizedPNL: %v", contract.LocalSymbol, position, marketPrice, averageCost, unrealizedPNL, realizedPNL)
	portfolioItem := PortfolioItem{*contract, position, marketPrice, marketValue, averageCost, unrealizedPNL, realizedPNL, account}
	conID := contract.ContractID
	portfolioItemMap, ok := w.Portfolio[account]
	if !ok {
		portfolioItemMap = make(map[int64]*PortfolioItem)
		w.Portfolio[account] = portfolioItemMap
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	if position != 0 {
		portfolioItemMap[conID] = &portfolioItem
	} else {
		delete(portfolioItemMap, conID)
	}

}

func (w *GoWrapper) Position(account string, contract *Contract, position float64, avgCost float64) {
	log.Printf("<Position>: account: %v, contract: %v position: %v, avgCost: %v", account, contract, position, avgCost)
	p := Position{account, *contract, position, avgCost}
	positionMap, ok := w.Positions[account]
	if !ok {
		positionMap = make(map[int64]*Position)
		w.Positions[account] = positionMap
	}
	conID := contract.ContractID
	if position != 0 {
		positionMap[conID] = &p
	} else {
		delete(positionMap, conID)
	}

}

func (w *GoWrapper) PositionEnd() {
	log.Printf("<positionEnd>")
}

func (w *GoWrapper) Pnl(reqID int64, dailyPnL float64, unrealizedPnL float64, realizedPnL float64) {
	log.WithField("reqID", reqID).Printf("<pnl>: dailyPnL: %v unrealizedPnL: %v realizedPnL: %v", dailyPnL, unrealizedPnL, realizedPnL)
	if pnl, ok := w.PnLs[reqID]; ok {
		pnl.DailyPnL = dailyPnL
		pnl.UnrealizePnL = unrealizedPnL
		pnl.RealizePnL = realizedPnL
	}
}

func (w *GoWrapper) PnlSingle(reqID int64, position int64, dailyPnL float64, unrealizedPnL float64, realizedPnL float64, value float64) {
	log.WithField("reqID", reqID).Printf("<pnl>:  position: %v dailyPnL: %v unrealizedPnL: %v realizedPnL: %v value: %v", position, dailyPnL, unrealizedPnL, realizedPnL, value)
	if pnlSingle, ok := w.PnLSingles[reqID]; ok {
		pnlSingle.Position = position
		pnlSingle.PnL.DailyPnL = dailyPnL
		pnlSingle.PnL.UnrealizePnL = unrealizedPnL
		pnlSingle.PnL.RealizePnL = realizedPnL
		pnlSingle.Value = value
	}
}

func (w *GoWrapper) OpenOrder(orderID int64, contract *Contract, order *Order, orderState *OrderState) {
	log.WithField("orderID", orderID).Printf("<openOrder>: orderId: %v contract: <%v> order: %v orderState: %v.", orderID, contract, order.OrderID, orderState.Status)
	if order.WhatIf {

	} else {
		key := OrderKey(order.ClientID, order.OrderID, order.PermID)
		trade, ok := w.Trades[key]
		if ok {
			trade.Order = order //XXX:
		} else {
			w.Trades[key] = Trade{Contract: contract,
				Order:       order,
				OrderStatus: &OrderStatus{Status: orderState.Status}}
			log.WithField("openOrder", key).Info(w.Trades[key])
		}

		if openOrderChan, ok := w.dataChanMap[int64(-OPEN_ORDER)]; ok {
			openOrderChan <- order
		} else {
			// push openOrderEvent
		}

	}
}

func (w *GoWrapper) OpenOrderEnd() {
	log.Printf("<openOrderEnd>")
	w.endReq(int64(-OPEN_ORDER))
}

func (w *GoWrapper) OrderStatus(orderID int64, status string, filled float64, remaining float64, avgFillPrice float64, permID int64, parentID int64, lastFillPrice float64, clientID int64, whyHeld string, mktCapPrice float64) {
	log.WithField("orderID", orderID).Printf("<orderStatus>: orderId: %v status: %v filled: %v remaining: %v avgFillPrice: %v.", orderID, status, filled, remaining, avgFillPrice)
	key := OrderKey(clientID, orderID, permID)
	trade, ok := w.Trades[key]
	if ok {

		oldStatus := trade.OrderStatus.Status
		trade.OrderStatus.Status = status
		trade.OrderStatus.Filled = filled
		trade.OrderStatus.Remaining = remaining
		trade.OrderStatus.AverageFillPrice = avgFillPrice
		trade.OrderStatus.PermID = permID
		trade.OrderStatus.LastFillPrice = lastFillPrice
		trade.OrderStatus.ClientID = clientID
		trade.OrderStatus.WhyHeld = whyHeld
		trade.OrderStatus.MarketCappedPrice = mktCapPrice
		trade.OrderStatus.LastLiquidity = 0.0

		logLen := len(trade.Log)
		var msg string

		if status == "Submitted" || logLen > 0 || trade.Log[logLen-1].Message == "Modify" {
			msg = "Modified"
		}

		logEntry := TradeLogEntry{Status: status, Message: msg}
		trade.Log = append(trade.Log, logEntry)

		if status != oldStatus {
			switch status {
			case Filled:
			case Cancelled:
			}
		}
	} else {
		log.WithFields(log.Fields{"orderID": orderID, "clientID": clientID}).Error("orderStatus: No order found")
	}
}

func (w *GoWrapper) ExecDetails(reqID int64, contract *Contract, execution *Execution) {
	log.WithField("reqID", reqID).Printf("<execDetails>: contract: %v execution: %v.", contract, execution)

	key := OrderKey(execution.ClientID, execution.OrderID, execution.PermID)
	trade, hasTrade := w.Trades[key]

	var c Contract
	if hasTrade || contract.ContractID == trade.Contract.ContractID {
		c = *trade.Contract
	} else {
		c = *contract
	}

	// execTime := execution.Time
	execTime, err := time.Parse("2006-01-02 15:04:05", execution.Time)
	if err != nil {
		fmt.Printf("parse time err: %v execTime: %v", err, execution.Time)
		panic("parse time err")
	}
	fill := Fill{Contract: c, Execution: *execution, CommissionReport: CommissionReport{}, Time: execTime}

	if _, ok := w.Fills[execution.ExecID]; !ok {
		w.Fills[execution.ExecID] = fill
		if hasTrade {
			trade.Fills = append(trade.Fills, fill)
			logEntry := TradeLogEntry{Time: execTime, Status: trade.OrderStatus.Status, Message: fmt.Sprintf("Fill %v@%v", execution.Shares, execution.Price)}
			trade.Log = append(trade.Log, logEntry)
		}
	}
}

func (w *GoWrapper) ExecDetailsEnd(reqID int64) {
	log.WithField("reqID", reqID).Print("<execDetailsEnd>")
}

func (w *GoWrapper) DeltaNeutralValidation(reqID int64, deltaNeutralContract DeltaNeutralContract) {
	log.WithField("reqID", reqID).Printf("<deltaNeutralValidation>: deltaNeutralContract: %v", deltaNeutralContract)
}

func (w *GoWrapper) CommissionReport(commissionReport CommissionReport) {
	log.Printf("<commissionReport>: commissionReport: %v", commissionReport)
	if commissionReport.Yield == UNSETFLOAT {
		commissionReport.Yield = 0
	}

	if commissionReport.RealizedPNL == UNSETFLOAT {
		commissionReport.RealizedPNL = 0
	}

	if fill, hasFill := w.Fills[commissionReport.ExecId]; hasFill {
		log.WithField("execID", commissionReport.ExecId).Info(commissionReport)
		key := OrderKey(fill.Execution.ClientID, fill.Execution.OrderID, fill.Execution.PermID)
		if trade, hasTrade := w.Trades[key]; hasTrade {
			_ = trade
			// push the commissionReportEvent
		} else {

		}
	} else {
		log.WithField("commissionReport", commissionReport).Error("No execution found for commissionReport")
	}

}

func (w *GoWrapper) OrderBound(reqID int64, apiClientID int64, apiOrderID int64) {
	log.WithField("reqID", reqID).Printf("<orderBound>: apiClientID: %v apiOrderID: %v", apiClientID, apiOrderID)
}

func (w *GoWrapper) ContractDetails(reqID int64, conDetails *ContractDetails) {
	log.WithField("reqID", reqID).Printf("<contractDetails>: contractDetails: %v", conDetails)
	if dataChan, ok := w.dataChanMap[reqID]; ok {
		dataChan <- conDetails
	}

}

func (w *GoWrapper) ContractDetailsEnd(reqID int64) {
	log.WithField("reqID", reqID).Print("<contractDetailsEnd>")
	w.endReq(reqID)
}

func (w *GoWrapper) BondContractDetails(reqID int64, conDetails *ContractDetails) {
	log.WithField("reqID", reqID).Printf("<bondContractDetails>: contractDetails: %v", conDetails)
	if dataChan, ok := w.dataChanMap[reqID]; ok {
		dataChan <- conDetails
	}
}

func (w *GoWrapper) SymbolSamples(reqID int64, contractDescriptions []ContractDescription) {
	log.WithField("reqID", reqID).Printf("<symbolSamples>: contractDescriptions: %v", contractDescriptions)
}

func (w *GoWrapper) SmartComponents(reqID int64, smartComps []SmartComponent) {
	log.WithField("reqID", reqID).Printf("<smartComponents>: smartComponents: %v", smartComps)
}

func (w *GoWrapper) MarketRule(marketRuleID int64, priceIncrements []PriceIncrement) {
	log.WithField("marketRuleID", marketRuleID).Printf("<marketRule>: marketRuleID: %v priceIncrements: %v", marketRuleID, priceIncrements)
}

func (w *GoWrapper) RealtimeBar(reqID int64, time int64, open float64, high float64, low float64, close float64, volume int64, wap float64, count int64) {
	log.WithField("reqID", reqID).Printf("<realtimeBar>: time: %v [O: %v H: %v, L: %v C: %v] volume: %v wap: %v count: %v", time, open, high, low, close, volume, wap, count)
}

func (w *GoWrapper) HistoricalData(reqID int64, bar *BarData) {
	log.WithField("reqID", reqID).Printf("<historicalData>: bar: %v", bar)
	w.dataChanMap[reqID] <- bar
}

func (w *GoWrapper) HistoricalDataEnd(reqID int64, startDateStr string, endDateStr string) {
	log.WithField("reqID", reqID).Printf("<historicalDataEnd>: startDate: %v endDate: %v", startDateStr, endDateStr)
	w.endReq(reqID)
}

func (w *GoWrapper) HistoricalDataUpdate(reqID int64, bar *BarData) {
	log.WithField("reqID", reqID).Printf("<historicalDataUpdate>: bar: %v", bar)
	w.subDataChanMap[reqID] <- bar
}

func (w *GoWrapper) HeadTimestamp(reqID int64, headTimestamp string) {
	log.WithField("reqID", reqID).Printf("<headTimestamp>: headTimestamp: %v", headTimestamp)
}

func (w *GoWrapper) HistoricalTicks(reqID int64, ticks []HistoricalTick, done bool) {
	log.WithField("reqID", reqID).Printf("<historicalTicks>:  done: %v", done)
}

func (w *GoWrapper) HistoricalTicksBidAsk(reqID int64, ticks []HistoricalTickBidAsk, done bool) {
	log.WithField("reqID", reqID).Printf("<historicalTicksBidAsk>: done: %v", done)
}

func (w *GoWrapper) HistoricalTicksLast(reqID int64, ticks []HistoricalTickLast, done bool) {
	log.WithField("reqID", reqID).Printf("<historicalTicksLast>: done: %v", done)
}

func (w *GoWrapper) TickSize(reqID int64, tickType int64, size int64) {
	log.WithField("reqID", reqID).Printf("<tickSize>:  tickType: %v size: %v.", tickType, size)
}

func (w *GoWrapper) TickSnapshotEnd(reqID int64) {
	log.WithField("reqID", reqID).Print("<tickSnapshotEnd>")
}

func (w *GoWrapper) MarketDataType(reqID int64, marketDataType int64) {
	log.WithField("reqID", reqID).Printf("<marketDataType>: marketDataType: %v", marketDataType)
}

func (w *GoWrapper) TickByTickAllLast(reqID int64, tickType int64, time int64, price float64, size int64, tickAttribLast TickAttribLast, exchange string, specialConditions string) {
	log.WithField("reqID", reqID).Printf("<tickByTickAllLast>:tickType: %v time: %v price: %v size: %v", tickType, time, price, size)
}

func (w *GoWrapper) TickByTickBidAsk(reqID int64, time int64, bidPrice float64, askPrice float64, bidSize int64, askSize int64, tickAttribBidAsk TickAttribBidAsk) {
	log.WithField("reqID", reqID).Printf("<tickByTickBidAsk>: time: %v bidPrice: %v askPrice: %v bidSize: %v askSize: %v", time, bidPrice, askPrice, bidSize, askSize)
}

func (w *GoWrapper) TickByTickMidPoint(reqID int64, time int64, midPoint float64) {
	log.WithField("reqID", reqID).Printf("<tickByTickMidPoint>: time: %v midPoint: %v ", time, midPoint)
}

func (w *GoWrapper) TickString(reqID int64, tickType int64, value string) {
	log.WithField("reqID", reqID).Printf("<tickString>: tickType: %v value: %v.", tickType, value)
}

func (w *GoWrapper) TickGeneric(reqID int64, tickType int64, value float64) {
	log.WithField("reqID", reqID).Printf("<tickGeneric>:tickType: %v value: %v.", tickType, value)
}

func (w *GoWrapper) TickEFP(reqID int64, tickType int64, basisPoints float64, formattedBasisPoints string, totalDividends float64, holdDays int64, futureLastTradeDate string, dividendImpact float64, dividendsToLastTradeDate float64) {
	log.WithField("reqID", reqID).Printf("<tickEFP>: tickType: %v basisPoints: %v.", tickType, basisPoints)
}

func (w *GoWrapper) TickReqParams(tickerID int64, minTick float64, bboExchange string, snapshotPermissions int64) {
	log.WithField("tickerID", tickerID).Printf("<tickReqParams>: tickerId: %v", tickerID)
}
func (w *GoWrapper) MktDepthExchanges(depthMktDataDescriptions []DepthMktDataDescription) {
	log.Printf("<mktDepthExchanges>: depthMktDataDescriptions: %v", depthMktDataDescriptions)
}

/*Returns the order book.

tickerId -  the request's identifier
position -  the order book's row being updated
operation - how to refresh the row:
	0 = insert (insert this new order into the row identified by 'position')
	1 = update (update the existing order in the row identified by 'position')
	2 = delete (delete the existing order at the row identified by 'position').
side -  0 for ask, 1 for bid
price - the order's price
size -  the order's size*/
func (w *GoWrapper) UpdateMktDepth(reqID int64, position int64, operation int64, side int64, price float64, size int64) {
	log.WithField("reqID", reqID).Printf("<updateMktDepth>: position: %v operation: %v side: %v price: %v size: %v", position, operation, side, price, size)
}

func (w *GoWrapper) UpdateMktDepthL2(reqID int64, position int64, marketMaker string, operation int64, side int64, price float64, size int64, isSmartDepth bool) {
	log.WithField("reqID", reqID).Printf("<updateMktDepthL2>: position: %v marketMaker: %v operation: %v side: %v price: %v size: %v isSmartDepth: %v", position, marketMaker, operation, side, price, size, isSmartDepth)
}

func (w *GoWrapper) TickOptionComputation(reqID int64, tickType int64, impliedVol float64, delta float64, optPrice float64, pvDiviedn float64, gamma float64, vega float64, theta float64, undPrice float64) {
	log.WithField("reqID", reqID).Printf("<tickOptionComputation>: tickType: %v ", tickType)
}

func (w *GoWrapper) FundamentalData(reqID int64, data string) {
	log.WithField("reqID", reqID).Printf("<fundamentalData>:data: %v", data)
}

func (w *GoWrapper) ScannerParameters(xml string) {
	log.Printf("<scannerParameters>: xml: %v", xml)

}

func (w *GoWrapper) ScannerData(reqID int64, rank int64, conDetails *ContractDetails, distance string, benchmark string, projection string, legs string) {
	log.WithField("reqID", reqID).Printf("<scannerData>: rank: %v", rank)
}

func (w *GoWrapper) ScannerDataEnd(reqID int64) {
	log.WithField("reqID", reqID).Print("<scannerDataEnd>")
}

func (w *GoWrapper) HistogramData(reqID int64, histogram []HistogramData) {
	log.WithField("reqID", reqID).Printf("<histogramData>: histogram: %v", histogram)
}

func (w *GoWrapper) RerouteMktDataReq(reqID int64, contractID int64, exchange string) {
	log.WithField("reqID", reqID).Printf("<rerouteMktDataReq>: contractID: %v exchange: %v", contractID, exchange)
}

func (w *GoWrapper) RerouteMktDepthReq(reqID int64, contractID int64, exchange string) {
	log.WithField("reqID", reqID).Printf("<rerouteMktDepthReq>: contractID: %v", contractID)
}

func (w *GoWrapper) SecurityDefinitionOptionParameter(reqID int64, exchange string, underlyingContractID int64, tradingClass string, multiplier string, expirations []string, strikes []float64) {
	log.WithField("reqID", reqID).Printf("<securityDefinitionOptionParameter>: underlyingContractID: %v", underlyingContractID)
}

func (w *GoWrapper) SecurityDefinitionOptionParameterEnd(reqID int64) {
	log.WithField("reqID", reqID).Print("<securityDefinitionOptionParameterEnd>")
}

func (w *GoWrapper) SoftDollarTiers(reqID int64, tiers []SoftDollarTier) {
	log.WithField("reqID", reqID).Printf("<softDollarTiers>: tiers: %v", tiers)
}

func (w *GoWrapper) FamilyCodes(famCodes []FamilyCode) {
	log.Printf("<familyCodes>: familyCodes: %v", famCodes)
}

func (w *GoWrapper) NewsProviders(newsProviders []NewsProvider) {
	log.Printf("<newsProviders>: newsProviders: %v", newsProviders)
}

func (w *GoWrapper) TickNews(tickerID int64, timeStamp int64, providerCode string, articleID string, headline string, extraData string) {
	log.WithField("tickerID", tickerID).Printf("<tickNews>: tickerID: %v timeStamp: %v", tickerID, timeStamp)
}

func (w *GoWrapper) NewsArticle(reqID int64, articleType int64, articleText string) {
	log.WithField("reqID", reqID).Printf("<newsArticle>: articleType: %v articleText: %v", articleType, articleText)
}

func (w *GoWrapper) HistoricalNews(reqID int64, time string, providerCode string, articleID string, headline string) {
	log.WithField("reqID", reqID).Printf("<historicalNews>: time: %v providerCode: %v articleID: %v, headline: %v", time, providerCode, articleID, headline)
}

func (w *GoWrapper) HistoricalNewsEnd(reqID int64, hasMore bool) {
	log.WithField("reqID", reqID).Printf("<historicalNewsEnd>: hasMore: %v", hasMore)
}

func (w *GoWrapper) UpdateNewsBulletin(msgID int64, msgType int64, newsMessage string, originExch string) {
	log.WithField("msgID", msgID).Printf("<updateNewsBulletin>: msgID: %v", msgID)
}

func (w *GoWrapper) ReceiveFA(faData int64, cxml string) {
	log.Printf("<receiveFA>: faData: %v", faData)

}

func (w *GoWrapper) CurrentTime(t time.Time) {
	log.Printf("<currentTime>: %v.", t)
}

func (w *GoWrapper) Error(reqID int64, errCode int64, errString string) {
	log.WithFields(log.Fields{"reqID": reqID, "errCode": errCode, "errString": errString}).Error("GoWrapper Error!")
}
