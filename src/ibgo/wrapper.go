package ibgo

import "time"

type IbWrapper struct {
}

type Account struct {
}

func (w *IbWrapper) connectAck() {

}

func (w *IbWrapper) nextValidId(reqId int) {

}

func (w *IbWrapper) managedAccounts(accountsList []Account) {

}

func (w *IbWrapper) updateAccountTime(timestamp time.Time) {

}
func (w *IbWrapper) updateAccountValue(tag string, val float32, currency string, account Account) {

}

func (w *IbWrapper) accountDownloadEnd(_account Account) {

}

func (w *IbWrapper) accountUpdateMulti(reqId int, account Account, modelCode int, tag string, val float32, currency string) {

}

func (w *IbWrapper) accountUpdateMultiEnd(reqId int) {

}

func (w *IbWrapper) accountSummary() {

}

func (w *IbWrapper) accountSummaryEnd() {

}

func (w *IbWrapper) updatePortfolio() {

}
func (w *IbWrapper) position() {

}
func (w *IbWrapper) positionEnd() {

}
func (w *IbWrapper) pnl() {

}
func (w *IbWrapper) pnlSingle() {

}
func (w *IbWrapper) openOrder() {

}
func (w *IbWrapper) openOrderEnd() {

}
func (w *IbWrapper) orderStatus() {

}
func (w *IbWrapper) execDetails() {

}
func (w *IbWrapper) execDetailsEnd() {

}
func (w *IbWrapper) commissionReport() {

}
func (w *IbWrapper) orderBound() {

}
func (w *IbWrapper) contractDetails() {

}
func (w *IbWrapper) contractDetailsEnd() {

}
func (w *IbWrapper) symbolSamples() {

}
func (w *IbWrapper) marketRule() {

}
func (w *IbWrapper) realtimeBar() {

}
func (w *IbWrapper) historicalData() {

}
func (w *IbWrapper) historicalDataEnd() {

}
func (w *IbWrapper) historicalDataUpdate() {

}
func (w *IbWrapper) headTimestamp() {

}
func (w *IbWrapper) historicalTicks() {

}
func (w *IbWrapper) historicalTicksBidAsk() {

}
func (w *IbWrapper) historicalTicksLast() {

}
func (w *IbWrapper) priceSizeTick() {

}
func (w *IbWrapper) tickSize() {

}
func (w *IbWrapper) tickSnapshotEnd() {

}
func (w *IbWrapper) tickByTickAllLast() {

}
func (w *IbWrapper) tickByTickBidAsk() {

}
func (w *IbWrapper) tickByTickMidPoint() {

}
func (w *IbWrapper) tickString() {

}
func (w *IbWrapper) tickGeneric() {

}
func (w *IbWrapper) tickReqParams() {

}
func (w *IbWrapper) mktDepthExchanges() {

}
func (w *IbWrapper) updateMktDepth() {

}
func (w *IbWrapper) updateMktDepthL2() {

}
func (w *IbWrapper) tickOptionComputation() {

}
func (w *IbWrapper) fundamentalData() {

}
func (w *IbWrapper) scannerParameters() {

}
func (w *IbWrapper) scannerData() {

}
func (w *IbWrapper) scannerDataEnd() {

}
func (w *IbWrapper) histogramData() {

}
func (w *IbWrapper) securityDefinitionOptionParameter() {

}
func (w *IbWrapper) securityDefinitionOptionParameterEnd() {

}
func (w *IbWrapper) newsProviders() {

}
func (w *IbWrapper) tickNews() {

}
func (w *IbWrapper) newsArticle() {

}
func (w *IbWrapper) historicalNews() {

}
func (w *IbWrapper) historicalNewsEnd() {

}
func (w *IbWrapper) updateNewsBulletin() {

}
func (w *IbWrapper) receiveFA() {

}
func (w *IbWrapper) currentTime() {

}
func (w *IbWrapper) error() {

}
