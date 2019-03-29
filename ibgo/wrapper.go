package ibgo

import (
	"time"
)

// IbWrapper contain the funcs to handle the msg from TWS or Gateway
type IbWrapper interface {
	tickPrice(reqID int64, tickType int64, price float64, attrib TickAttrib)
	tickSize(reqID int64, tickType int64, size int64)
	orderStatus(orderID int64, status string, filled float64, remaining float64, avgFillPrice float64, permID int64, parentID int64, lastFillPrice float64, clientID int64, whyHeld string, mktCapPrice float64)
	nextValidID(reqID int64)
	managedAccounts(accountsList []Account)
	updateAccountValue(tag string, val string, currency string, accName string)
	updatePortfolio(contract *Contract, position float64, marketPrice float64, marketValue float64, averageCost float64, unrealizedPNL float64, realizedPNL float64, accName string)
	updateAccountTime(accTime time.Time)
	openOrder(orderID int64, contract *Contract, order *Order, orderState *OrderState)
	contractDetails(reqID int64, conDetails *ContractDetails)
	execDetails(reqID int64, contract *Contract, execution *Execution)

	updateMktDepth(reqID int64, position int64, operation int64, side int64, price float64, size int64)
	updateMktDepthL2(reqID int64, position int64, marketMaker string, operation int64, side int64, price float64, size int64, isSmartDepth bool)
	updateNewsBulletin(msgID int64, msgType int64, newsMessage string, originExch string)
	historicalData(reqID int64, bar *BarData)
	receiveFA(faData int64, cxml string)
	historicalDataUpdate(reqID int64, bar *BarData)
	bondContractDetails(reqID int64, conDetails *ContractDetails)
	scannerParameters(xml string)
	scannerData(reqID int64, rank int64, conDetails *ContractDetails, distance string, benchmark string, projection string, legs string)
	scannerDataEnd(reqID int64)
	tickOptionComputation(reqID int64, tickType int64, impliedVol float64, delta float64, optPrice float64, pvDiviedn float64, gamma float64, vega float64, theta float64, undPrice float64)
	tickGeneric(reqID int64, tickType int64, value float64)
	tickString(reqID int64, tickType int64, value string)
	tickEFP(reqID int64, tickType int64, basisPoints float64, formattedBasisPoints string, totalDividends float64, holdDays int64, futureLastTradeDate string, dividendImpact float64, dividendsToLastTradeDate float64)

	connectAck()
	error(reqID int64, errCode int64, errString string)

	//wrap end
	accountDownloadEnd(accName string)
	openOrderEnd()
	historicalDataEnd(reqID int64, startDateStr string, endDateStr string)
	execDetailsEnd(reqID int64)
	currentTime(t time.Time)
}
