package ibgo

import (
	"fmt"
	"log"
	"time"
)

type Wrapper struct {
}

func (w Wrapper) connectAck() {

}

func (w Wrapper) nextValidID(reqID int64) {
	log.Printf("<nextValidId>: %v\n", reqID)

}

func (w Wrapper) managedAccounts(accountsList []Account) {
	log.Printf("<managedAccounts>: %v\n", accountsList)

}

func (w Wrapper) tickPrice(reqID int64, tickType int64, price float64, attrib TickAttrib) {
	log.Printf("<tickPrice>: reqID: %v tickType: %v price: %v\n", reqID, tickType, price)
}

func (w Wrapper) updateAccountTime(accTime time.Time) {
	log.Printf("<updateAccountTime>: %v", accTime)

}
func (w Wrapper) updateAccountValue(tag string, val string, currency string, accName string) {
	// log.Printf("<updateAccountValue>: accName:%v [%v]:%v currency:%v", accName, tag, val, currency)

}

func (w Wrapper) accountDownloadEnd(accName string) {
	log.Printf("<accountDownloadEnd>: %v", accName)

}

func (w *Wrapper) accountUpdateMulti(reqID int, account Account, modelCode int, tag string, val float32, currency string) {

}

func (w *Wrapper) accountUpdateMultiEnd(reqID int) {

}

func (w *Wrapper) accountSummary() {

}

func (w *Wrapper) accountSummaryEnd() {

}

func (w Wrapper) updatePortfolio(contract *Contract, position float64, marketPrice float64, marketValue float64, averageCost float64, unrealizedPNL float64, realizedPNL float64, accName string) {
	log.Printf("<updatePortfolio>: contract: %v pos: %v marketPrice: %v averageCost: %v unrealizedPNL: %v realizedPNL: %v", contract.LocalSymbol, position, marketPrice, averageCost, unrealizedPNL, realizedPNL)
}
func (w *Wrapper) position() {

}
func (w *Wrapper) positionEnd() {

}
func (w *Wrapper) pnl() {

}
func (w *Wrapper) pnlSingle() {

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
func (w *Wrapper) commissionReport() {

}
func (w *Wrapper) orderBound() {

}
func (w Wrapper) contractDetails(reqID int64, conDetails *ContractDetails) {
	fmt.Printf("<contractDetails>: reqID: %v contractDetails: %v", reqID, conDetails)

}
func (w *Wrapper) contractDetailsEnd() {

}
func (w *Wrapper) symbolSamples() {

}
func (w *Wrapper) marketRule() {

}
func (w *Wrapper) realtimeBar() {

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
func (w *Wrapper) headTimestamp() {

}
func (w *Wrapper) historicalTicks() {

}
func (w *Wrapper) historicalTicksBidAsk() {

}
func (w *Wrapper) historicalTicksLast() {

}
func (w *Wrapper) priceSizeTick() {

}
func (w Wrapper) tickSize(reqID int64, tickType int64, size int64) {
	log.Printf("<tickSize>: reqID: %v tickType: %v size: %v\n", reqID, tickType, size)

}
func (w *Wrapper) tickSnapshotEnd() {

}
func (w *Wrapper) tickByTickAllLast() {

}
func (w *Wrapper) tickByTickBidAsk() {

}
func (w *Wrapper) tickByTickMidPoint() {

}
func (w *Wrapper) tickString() {

}
func (w *Wrapper) tickGeneric() {

}
func (w *Wrapper) tickReqParams() {

}
func (w *Wrapper) mktDepthExchanges() {

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
	log.Printf("<updateMktDepth>: reqID:%v", reqID)
}
func (w Wrapper) updateMktDepthL2(reqID int64, position int64, marketMaker string, operation int64, side int64, price float64, size int64, isSmartDepth bool) {
	log.Printf("<updateMktDepthL2>: reqID:%v", reqID)
}
func (w *Wrapper) tickOptionComputation() {

}
func (w *Wrapper) fundamentalData() {

}
func (w *Wrapper) scannerParameters() {

}
func (w *Wrapper) scannerData() {

}
func (w *Wrapper) scannerDataEnd() {

}
func (w *Wrapper) histogramData() {

}
func (w *Wrapper) securityDefinitionOptionParameter() {

}
func (w *Wrapper) securityDefinitionOptionParameterEnd() {

}
func (w *Wrapper) newsProviders() {

}
func (w *Wrapper) tickNews() {

}
func (w *Wrapper) newsArticle() {

}
func (w *Wrapper) historicalNews() {

}
func (w *Wrapper) historicalNewsEnd() {

}
func (w Wrapper) updateNewsBulletin(msgID int64, msgType int64, newsMessage string, originExch string) {
	log.Printf("<updateNewsBulletin>: msgID:%v", msgID)

}
func (w Wrapper) receiveFA(faData int64, cxml string) {
	log.Printf("<receiveFA>: faData:%v", faData)

}
func (w Wrapper) currentTime(t time.Time) {
	log.Printf("<currentTime>: %v\n", t)
}
func (w Wrapper) error(reqID int64, errCode int64, errString string) {
	log.Printf("<error>: reqID: %v errCode: %v errString: %v\n", reqID, errCode, errString)

}
