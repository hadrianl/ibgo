package ibgo

import (
	"fmt"
	"log"
	"time"
)

type Wrapper struct {
}

func (w Wrapper) connectAck() {
	log.Printf("<connectAck>...")
}

func (w Wrapper) nextValidID(reqID int64) {
	log.Printf("<nextValidId>: %v\n", reqID)

}

func (w Wrapper) managedAccounts(accountsList []string) {
	log.Printf("<managedAccounts>: %v\n", accountsList)

}

func (w Wrapper) tickPrice(reqID int64, tickType int64, price float64, attrib TickAttrib) {
	log.Printf("<tickPrice>: reqID: %v tickType: %v price: %v\n", reqID, tickType, price)
}

func (w Wrapper) updateAccountTime(accTime time.Time) {
	log.Printf("<updateAccountTime>: %v", accTime)

}

func (w Wrapper) updateAccountValue(tag string, value string, currency string, account string) {
	log.Printf("<updateAccountValue>: account:%v [%v]:%v currency:%v", account, tag, value, currency)

}

func (w Wrapper) accountDownloadEnd(accName string) {
	log.Printf("<accountDownloadEnd>: %v", accName)
}

func (w Wrapper) accountUpdateMulti(reqID int64, account string, modelCode string, tag string, value string, currency string) {
	log.Printf("<accountUpdateMulti>: reqID: %v account:%v modelCode: %v [%v]:%v currency:%v", reqID, account, modelCode, tag, value, currency)
}

func (w Wrapper) accountUpdateMultiEnd(reqID int64) {
	log.Printf("<accountUpdateMultiEnd>: reqID: %v", reqID)
}

func (w Wrapper) accountSummary(reqID int64, account string, tag string, value string, currency string) {
	log.Printf("<accountSummary>: reqID: %v account:%v [%v]:%v currency:%v", reqID, account, tag, value, currency)
}

func (w Wrapper) accountSummaryEnd(reqID int64) {
	log.Printf("<accountSummaryEnd>: reqID: %v", reqID)
}

func (w Wrapper) verifyMessageAPI(apiData string) {
	log.Printf("<verifyMessageAPI>: apiData: %v", apiData)
}

func (w Wrapper) verifyCompleted(isSuccessful bool, err string) {
	log.Printf("<verifyCompleted>: isSuccessful: %v error: %v", isSuccessful, err)
}

func (w Wrapper) verifyAndAuthMessageAPI(apiData string, xyzChallange string) {
	log.Printf("<verifyCompleted>: apiData: %v xyzChallange: %v", apiData, xyzChallange)
}

func (w Wrapper) verifyAndAuthCompleted(isSuccessful bool, err string) {
	log.Printf("<verifyAndAuthCompleted>: isSuccessful: %v error: %v", isSuccessful, err)
}

func (w Wrapper) displayGroupList(reqID int64, groups string) {
	log.Printf("<displayGroupList>: reqID: %v groups: %v", reqID, groups)
}

func (w Wrapper) displayGroupUpdated(reqID int64, contractInfo string) {
	log.Printf("<displayGroupUpdated>: reqID: %v contractInfo: %v", reqID, contractInfo)
}

func (w Wrapper) positionMulti(reqID int64, account string, modelCode string, contract *Contract, position float64, avgCost float64) {
	log.Printf("<positionMulti>: reqID: %v account: %v modelCode: %v contract: <%v> position: %v avgCost: %v", reqID, account, modelCode, contract, position, avgCost)
}

func (w Wrapper) positionMultiEnd(reqID int64) {
	log.Printf("<positionMultiEnd>: reqID: %v", reqID)
}

func (w Wrapper) updatePortfolio(contract *Contract, position float64, marketPrice float64, marketValue float64, averageCost float64, unrealizedPNL float64, realizedPNL float64, accName string) {
	log.Printf("<updatePortfolio>: contract: %v pos: %v marketPrice: %v averageCost: %v unrealizedPNL: %v realizedPNL: %v", contract.LocalSymbol, position, marketPrice, averageCost, unrealizedPNL, realizedPNL)
}

func (w Wrapper) position(account string, contract *Contract, position float64, avgCost float64) {
	log.Printf("<updatePortfolio>: account: %v, contract: %v position: %v, avgCost: %v", account, contract, position, avgCost)
}

func (w Wrapper) positionEnd() {
	log.Printf("<positionEnd>:...")
}

func (w Wrapper) pnl(reqID int64, dailyPnL float64, unrealizedPnL float64, realizedPnL float64) {
	log.Printf("<pnl>: reqID: %v dailyPnL: %v unrealizedPnL: %v realizedPnL: %v", reqID, dailyPnL, unrealizedPnL, realizedPnL)
}

func (w Wrapper) pnlSingle(reqID int64, position int64, dailyPnL float64, unrealizedPnL float64, realizedPnL float64, value float64) {
	log.Printf("<pnl>: reqID: %v position: %v dailyPnL: %v unrealizedPnL: %v realizedPnL: %v value: %v", reqID, position, dailyPnL, unrealizedPnL, realizedPnL, value)
}

func (w Wrapper) openOrder(orderID int64, contract *Contract, order *Order, orderState *OrderState) {
	log.Printf("<openOrder>: orderId: %v contract: <%v> order: %v orderState: %v\n", orderID, contract, order.OrderID, orderState.Status)

}

func (w Wrapper) openOrderEnd() {
	log.Printf("<openOrderEnd>...")

}

func (w Wrapper) orderStatus(orderID int64, status string, filled float64, remaining float64, avgFillPrice float64, permID int64, parentID int64, lastFillPrice float64, clientID int64, whyHeld string, mktCapPrice float64) {
	log.Printf("<orderStatus>: orderId: %v status: %v filled: %v remaining: %v avgFillPrice: %v\n", orderID, status, filled, remaining, avgFillPrice)
}

func (w Wrapper) execDetails(reqID int64, contract *Contract, execution *Execution) {
	log.Printf("<execDetails>: reqID: %v contract: %v execution: %v\n", reqID, contract, execution)
}

func (w Wrapper) execDetailsEnd(reqID int64) {
	log.Printf("<execDetailsEnd>: reqID: %v", reqID)
}

func (w Wrapper) deltaNeutralValidation(reqID int64, deltaNeutralContract DeltaNeutralContract) {
	log.Printf("<deltaNeutralValidation>: reqID: %v deltaNeutralContract: %v", reqID, deltaNeutralContract)
}

func (w Wrapper) commissionReport(commissionReport CommissionReport) {
	log.Printf("<commissionReport>: commissionReport: %v", commissionReport)
}

func (w Wrapper) orderBound(reqID int64, apiClientID int64, apiOrderID int64) {
	log.Printf("<orderBound>: reqID: %v apiClientID: %v apiOrderID: %v", reqID, apiClientID, apiOrderID)
}

func (w Wrapper) contractDetails(reqID int64, conDetails *ContractDetails) {
	fmt.Printf("<contractDetails>: reqID: %v contractDetails: %v", reqID, conDetails)

}

func (w Wrapper) contractDetailsEnd(reqID int64) {
	fmt.Printf("<contractDetailsEnd>: reqID: %v", reqID)
}

func (w Wrapper) bondContractDetails(reqID int64, conDetails *ContractDetails) {
	log.Printf("<bondContractDetails>: reqID: %v contractDetails: %v", reqID, conDetails)
}

func (w Wrapper) symbolSamples(reqID int64, contractDescriptions []ContractDescription) {
	log.Printf("<symbolSamples>: reqID: %v contractDescriptions: %v", reqID, contractDescriptions)
}

func (w Wrapper) smartComponents(reqID int64, smartComps []SmartComponent) {
	log.Printf("<smartComponents>: reqID: %v smartComponents: %v", reqID, smartComps)
}

func (w Wrapper) marketRule(marketRuleID int64, priceIncrements []PriceIncrement) {
	log.Printf("<marketRule>: marketRuleID: %v priceIncrements: %v", marketRuleID, priceIncrements)
}

func (w Wrapper) realtimeBar(reqID int64, time int64, open float64, high float64, low float64, close float64, volume int64, wap float64, count int64) {
	log.Printf("<realtimeBar>: reqID: %v time: %v [O: %v H: %v, L: %v C: %v] volume: %v wap: %v count: %v", reqID, time, open, high, low, close, volume, wap, count)
}

func (w Wrapper) historicalData(reqID int64, bar *BarData) {
	log.Printf("<historicalData>: reqID: %v bar: %v", reqID, bar)

}

func (w Wrapper) historicalDataEnd(reqID int64, startDateStr string, endDateStr string) {
	log.Printf("<historicalDataEnd>: reqID: %v startDate: %v endDate: %v", reqID, startDateStr, endDateStr)
}

func (w Wrapper) historicalDataUpdate(reqID int64, bar *BarData) {
	log.Printf("<historicalDataUpdate>: reqID: %v bar: %v", reqID, bar)
}

func (w Wrapper) headTimestamp(reqID int64, headTimestamp string) {
	log.Printf("<headTimestamp>: reqID: %v headTimestamp: %v", reqID, headTimestamp)
}

func (w Wrapper) historicalTicks(reqID int64, ticks []HistoricalTick, done bool) {
	log.Printf("<historicalTicks>: reqID: %v done: %v", reqID, done)
}

func (w Wrapper) historicalTicksBidAsk(reqID int64, ticks []HistoricalTickBidAsk, done bool) {
	log.Printf("<historicalTicksBidAsk>: reqID: %v done: %v", reqID, done)
}

func (w Wrapper) historicalTicksLast(reqID int64, ticks []HistoricalTickLast, done bool) {
	log.Printf("<historicalTicksLast>: reqID: %v done: %v", reqID, done)
}

func (w *Wrapper) priceSizeTick() {

}

func (w Wrapper) tickSize(reqID int64, tickType int64, size int64) {
	log.Printf("<tickSize>: reqID: %v tickType: %v size: %v\n", reqID, tickType, size)

}

func (w Wrapper) tickSnapshotEnd(reqID int64) {
	log.Printf("<tickSnapshotEnd>: reqID: %v", reqID)
}

func (w Wrapper) marketDataType(reqID int64, marketDataType int64) {
	log.Printf("<marketDataType>: reqID: %v marketDataType: %v", reqID, marketDataType)
}

func (w Wrapper) tickByTickAllLast(reqID int64, tickType int64, time int64, price float64, size int64, tickAttribLast TickAttribLast, exchange string, specialConditions string) {
	log.Printf("<tickByTickAllLast>: reqID: %v tickType: %v time: %v price: %v size: %v", reqID, tickType, time, price, size)
}

func (w Wrapper) tickByTickBidAsk(reqID int64, time int64, bidPrice float64, askPrice float64, bidSize int64, askSize int64, tickAttribBidAsk TickAttribBidAsk) {
	log.Printf("<tickByTickBidAsk>: reqID: %v time: %v bidPrice: %v askPrice: %v bidSize: %v askSize: %v", reqID, time, bidPrice, askPrice, bidSize, askSize)
}

func (w Wrapper) tickByTickMidPoint(reqID int64, time int64, midPoint float64) {
	log.Printf("<tickByTickMidPoint>: reqID: %v time: %v midPoint: %v ", reqID, time, midPoint)
}

func (w Wrapper) tickString(reqID int64, tickType int64, value string) {
	log.Printf("<tickString>: reqID: %v tickType: %v value: %v\n", reqID, tickType, value)

}

func (w Wrapper) tickGeneric(reqID int64, tickType int64, value float64) {
	log.Printf("<tickGeneric>: reqID: %v tickType: %v value: %v\n", reqID, tickType, value)
}

func (w Wrapper) tickEFP(reqID int64, tickType int64, basisPoints float64, formattedBasisPoints string, totalDividends float64, holdDays int64, futureLastTradeDate string, dividendImpact float64, dividendsToLastTradeDate float64) {
	log.Printf("<tickEFP>: reqID: %v tickType: %v basisPoints: %v\n", reqID, tickType, basisPoints)
}

func (w Wrapper) tickReqParams(tickerID int64, minTick float64, bboExchange string, snapshotPermissions int64) {
	log.Printf("<tickReqParams>: tickerId: %v", tickerID)
}
func (w Wrapper) mktDepthExchanges(depthMktDataDescriptions []DepthMktDataDescription) {
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
func (w Wrapper) updateMktDepth(reqID int64, position int64, operation int64, side int64, price float64, size int64) {
	log.Printf("<updateMktDepth>: reqID: %v", reqID)
}

func (w Wrapper) updateMktDepthL2(reqID int64, position int64, marketMaker string, operation int64, side int64, price float64, size int64, isSmartDepth bool) {
	log.Printf("<updateMktDepthL2>: reqID: %v", reqID)
}

func (w Wrapper) tickOptionComputation(reqID int64, tickType int64, impliedVol float64, delta float64, optPrice float64, pvDiviedn float64, gamma float64, vega float64, theta float64, undPrice float64) {
	log.Printf("<tickOptionComputation>: reqID: %v ", reqID)
}

func (w Wrapper) fundamentalData(reqID int64, data string) {
	log.Printf("<fundamentalData>: reqID: %v data: %v", reqID, data)
}

func (w Wrapper) scannerParameters(xml string) {
	log.Printf("<scannerParameters>: xml: %v", xml)

}

func (w Wrapper) scannerData(reqID int64, rank int64, conDetails *ContractDetails, distance string, benchmark string, projection string, legs string) {
	log.Printf("<scannerData>: reqID: %v", reqID)
}

func (w Wrapper) scannerDataEnd(reqID int64) {
	log.Printf("<scannerDataEnd>: reqID: %v", reqID)
}

func (w Wrapper) histogramData(reqID int64, histogram []HistogramData) {
	log.Printf("<histogramData>: reqID: %v", reqID)
}

func (w Wrapper) rerouteMktDataReq(reqID int64, contractID int64, exchange string) {
	log.Printf("<rerouteMktDataReq>: reqID: %v contractID: %v", reqID, contractID)
}

func (w Wrapper) rerouteMktDepthReq(reqID int64, contractID int64, exchange string) {
	log.Printf("<rerouteMktDepthReq>: reqID: %v contractID: %v", reqID, contractID)
}

func (w Wrapper) securityDefinitionOptionParameter(reqID int64, exchange string, underlyingContractID int64, tradingClass string, multiplier string, expirations []string, strikes []float64) {
	log.Printf("<securityDefinitionOptionParameter>: reqID: %v", reqID)
}

func (w Wrapper) securityDefinitionOptionParameterEnd(reqID int64) {
	log.Printf("<securityDefinitionOptionParameterEnd>: reqID: %v", reqID)
}

func (w Wrapper) softDollarTiers(reqID int64, tiers []SoftDollarTier) {
	log.Printf("<softDollarTiers>: reqID: %v", reqID)
}

func (w Wrapper) familyCodes(famCodes []FamilyCode) {
	log.Printf("<familyCodes>: familyCodes: %v", famCodes)
}

func (w Wrapper) newsProviders(newsProviders []NewsProvider) {
	log.Printf("<newsProviders>: newsProviders: %v", newsProviders)
}

func (w Wrapper) tickNews(tickerID int64, timeStamp int64, providerCode string, articleID string, headline string, extraData string) {
	log.Printf("<tickNews>: tickerID: %v timeStamp: %v", tickerID, timeStamp)
}

func (w Wrapper) newsArticle(reqID int64, articleType int64, articleText string) {
	log.Printf("<newsArticle>: reqID: %v", reqID)
}

func (w Wrapper) historicalNews(reqID int64, time string, providerCode string, articleID string, headline string) {
	log.Printf("<historicalNews>: reqID: %v", reqID)
}

func (w Wrapper) historicalNewsEnd(reqID int64, hasMore bool) {
	log.Printf("<historicalNewsEnd>: reqID: %v hasMore: %v", reqID, hasMore)
}

func (w Wrapper) updateNewsBulletin(msgID int64, msgType int64, newsMessage string, originExch string) {
	log.Printf("<updateNewsBulletin>: msgID: %v", msgID)

}

func (w Wrapper) receiveFA(faData int64, cxml string) {
	log.Printf("<receiveFA>: faData: %v", faData)

}

func (w Wrapper) currentTime(t time.Time) {
	log.Printf("<currentTime>: %v\n", t)
}

func (w Wrapper) error(reqID int64, errCode int64, errString string) {
	log.Printf("<error>: reqID: %v errCode: %v errString: %v\n", reqID, errCode, errString)
}
