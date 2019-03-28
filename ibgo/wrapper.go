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
	connectAck()
	error(reqID int64, errCode int64, errString string)

	//wrap end
	accountDownloadEnd(accName string)
	openOrderEnd()
	currentTime(t time.Time)
}
