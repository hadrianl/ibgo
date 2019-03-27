package ibgo

import (
	"fmt"
	"log"
	"time"
)

type IbWrapper interface {
	tickPrice(reqID int64, tickType int64, price float64, attrib TickAttrib)
	tickSize(reqID int64, tickType int64, size int64)
	orderStatus(orderID int64, status string, filled float64, remaining float64, avgFillPrice float64, permId int64, parentId int64, lastFillPrice float64, clientId int64, whyHeld string, mktCapPrice float64)
	nextValidId(reqId int64)
	managedAccounts(accountsList []Account)
	updateAccountValue(tag string, val string, currency string, accName string)
	updatePortfolio(contract *Contract, position float64, marketPrice float64, marketValue float64, averageCost float64, unrealizedPNL float64, realizedPNL float64, accName string)
	updateAccountTime(accTime time.Time)
	openOrder(orderId int64, contract *Contract, order *Order, orderState *OrderState)
	contractDetails(reqID int64, conDetails *ContractDetails)
	execDetails(reqID int64, contract *Contract, execution *Execution)
	connectAck()
	error(reqId int64, errCode int64, errString string)

	//wrap end
	accountDownloadEnd(accName string)
	currentTime(t time.Time)
}

type Wrapper struct {
	ic *IbClient
}

func (w Wrapper) connectAck() {

}

func (w Wrapper) nextValidId(reqId int64) {
	log.Printf("nextValidId: %v\n", reqId)

}

func (w Wrapper) managedAccounts(accountsList []Account) {
	log.Printf("managedAccounts: %v\n", accountsList)

}

func (w Wrapper) tickPrice(reqId int64, tickType int64, price float64, attrib TickAttrib) {
	log.Printf("reqId: %v tickType: %v price: %v\n", reqId, tickType, price)
}

func (w Wrapper) updateAccountTime(accTime time.Time) {
	log.Printf("updateAccountTime: %v", accTime)

}
func (w Wrapper) updateAccountValue(tag string, val string, currency string, accName string) {
	log.Printf("Account:%v <%v>:%v currency:%v", accName, tag, val, currency)

}

func (w Wrapper) accountDownloadEnd(accName string) {
	log.Printf("accountDownloadEnd: %v", accName)

}

func (w *Wrapper) accountUpdateMulti(reqId int, account Account, modelCode int, tag string, val float32, currency string) {

}

func (w *Wrapper) accountUpdateMultiEnd(reqId int) {

}

func (w *Wrapper) accountSummary() {

}

func (w *Wrapper) accountSummaryEnd() {

}

func (w Wrapper) updatePortfolio(contract *Contract, position float64, marketPrice float64, marketValue float64, averageCost float64, unrealizedPNL float64, realizedPNL float64, accName string) {
	log.Printf("contract: %v pos: %v marketPrice: %v averageCost: %v unrealizedPNL: %v realizedPNL: %v", contract, position, marketPrice, averageCost, unrealizedPNL, realizedPNL)
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
	log.Printf("orderId: %v contract: <%v> order: %v orderState: %v\n", orderID, contract.LocalSymbol, order.OrderID, orderState.Status)

}
func (w *Wrapper) openOrderEnd() {

}
func (w Wrapper) orderStatus(orderID int64, status string, filled float64, remaining float64, avgFillPrice float64, permId int64, parentId int64, lastFillPrice float64, clientId int64, whyHeld string, mktCapPrice float64) {
	log.Printf("orderId: %v status: %v filled: %v remaining: %v avgFillPrice: %v\n", orderID, status, filled, remaining, avgFillPrice)
}
func (w Wrapper) execDetails(reqID int64, contract *Contract, execution *Execution) {
	log.Printf("reqID: %v contract: %v execution: %v\n", reqID, contract, execution)
}

func (w *Wrapper) execDetailsEnd() {

}
func (w *Wrapper) commissionReport() {

}
func (w *Wrapper) orderBound() {

}
func (w Wrapper) contractDetails(reqID int64, conDetails *ContractDetails) {
	fmt.Printf("reqID: %v contractDetails: %v", reqID, conDetails)

}
func (w *Wrapper) contractDetailsEnd() {

}
func (w *Wrapper) symbolSamples() {

}
func (w *Wrapper) marketRule() {

}
func (w *Wrapper) realtimeBar() {

}
func (w *Wrapper) historicalData() {

}
func (w *Wrapper) historicalDataEnd() {

}
func (w *Wrapper) historicalDataUpdate() {

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
func (w Wrapper) tickSize(reqId int64, tickType int64, size int64) {
	log.Printf("reqId: %v tickType: %v size: %v\n", reqId, tickType, size)

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
func (w *Wrapper) updateMktDepth() {

}
func (w *Wrapper) updateMktDepthL2() {

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
func (w *Wrapper) updateNewsBulletin() {

}
func (w *Wrapper) receiveFA() {

}
func (w Wrapper) currentTime(t time.Time) {
	log.Printf("CurrentTime :%v\n", t)
}
func (w Wrapper) error(reqId int64, errCode int64, errString string) {
	log.Printf("reqId: %v errCode: %v errString: %v\n", reqId, errCode, errString)

}
