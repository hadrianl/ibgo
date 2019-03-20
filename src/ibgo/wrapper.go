package ibgo

import (
	"fmt"
	"time"
)

type IbWrapper interface {
	nextValidId(reqId int)
	managedAccounts(accountsList []Account)
	connectAck()
	error(reqId int, errCode int, errString string)
	currentTime(t time.Time)
}

type Wrapper struct {
	ic *IbClient
}

func (w Wrapper) connectAck() {

}

func (w Wrapper) nextValidId(reqId int) {
	fmt.Println("nextValidId:", reqId)

}

func (w Wrapper) managedAccounts(accountsList []Account) {
	fmt.Println("managedAccounts:", accountsList)

}

func (w *Wrapper) updateAccountTime(timestamp time.Time) {

}
func (w *Wrapper) updateAccountValue(tag string, val float32, currency string, account Account) {

}

func (w *Wrapper) accountDownloadEnd(_account Account) {

}

func (w *Wrapper) accountUpdateMulti(reqId int, account Account, modelCode int, tag string, val float32, currency string) {

}

func (w *Wrapper) accountUpdateMultiEnd(reqId int) {

}

func (w *Wrapper) accountSummary() {

}

func (w *Wrapper) accountSummaryEnd() {

}

func (w *Wrapper) updatePortfolio() {

}
func (w *Wrapper) position() {

}
func (w *Wrapper) positionEnd() {

}
func (w *Wrapper) pnl() {

}
func (w *Wrapper) pnlSingle() {

}
func (w *Wrapper) openOrder() {

}
func (w *Wrapper) openOrderEnd() {

}
func (w *Wrapper) orderStatus() {

}
func (w *Wrapper) execDetails() {

}
func (w *Wrapper) execDetailsEnd() {

}
func (w *Wrapper) commissionReport() {

}
func (w *Wrapper) orderBound() {

}
func (w *Wrapper) contractDetails() {

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
func (w *Wrapper) tickSize() {

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

}
func (w Wrapper) error(reqId int, errCode int, errString string) {
	fmt.Printf("reqId: %v errCode: %v errString: %v\n", reqId, errCode, errString)

}
