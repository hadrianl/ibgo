package IBAlgoTrade

import (
	"fmt"
	"github.com/hadrianl/ibgo/ibapi"
	// "runtime"
	"strings"
	"sync/atomic"
	"time"
)

// func init() {
// 	fmt.Print("GOMAXPROCS")
// 	runtime.GOMAXPROCS(4)
// }

type IB struct {
	Client   *ibapi.IbClient
	Wrapper  *GoWrapper
	host     string
	port     int
	clientID int64
	reqID    int64
}

func NewIB(host string, port int, clientID int64) *IB {
	ib := &IB{
		host:     host,
		port:     port,
		clientID: clientID,
	}
	wrapper := &GoWrapper{ib: ib}
	wrapper.reset()
	ibclient := ibapi.NewIbClient(wrapper)
	ib.Client = ibclient
	ib.Wrapper = wrapper

	return ib
}

// Connect to TWS or Gateway
func (ib *IB) Connect() error {

	if err := ib.Client.Connect(ib.host, ib.port, ib.clientID); err != nil {
		return err
	}

	if err := ib.Client.HandShake(); err != nil {
		return err
	}

	err := ib.Client.Run()
	return err
}

func (ib *IB) DisConnect() error {
	err := ib.Client.Disconnect()
	return err
}

func (ib *IB) GetReqID() int64 {
	return atomic.AddInt64(&ib.reqID, 1)
}

func (ib *IB) DoSomeTest() {
	hsik9 := ibapi.Contract{359142357, "HSI", "FUT", "20190530", 0, "?", "50", "HKFE", "HKD", "HSIK9", "HSI", "", false, "", "", "", nil, nil}
	// fmt.Println(hsij9)
	ib.Client.ReqCurrentTime()
	ib.Client.ReqAutoOpenOrders(true)
	ib.Client.ReqAccountUpdates(true, "")
	ib.Client.ReqPositions()
	// ib.Client.ReqIDs(5)
	time.Sleep(time.Second * 3)
	bars := ib.ReqHistoricalData(hsik9, "", "600 S", "1 min", "TRADES", false, 1, true, nil)
	ib.ReqOpenOrders()
	// fmt.Print(openOrders)
	// contractDetails := ib.ReqContractDetails(&hsik9)
	time.Sleep(time.Second * 3)
	// fmt.Println(bars)
	fmt.Println(bars)
	tags := []string{"AccountType,NetLiquidation,TotalCashValue,SettledCash,",
		"AccruedCash,BuyingPower,EquityWithLoanValue,",
		"PreviousEquityWithLoanValue,GrossPositionValue,ReqTEquity,",
		"ReqTMargin,SMA,InitMarginReq,MaintMarginReq,AvailableFunds,",
		"ExcessLiquidity,Cushion,FullInitMarginReq,FullMaintMarginReq,",
		"FullAvailableFunds,FullExcessLiquidity,LookAheadNextChange,",
		"LookAheadInitMarginReq,LookAheadMaintMarginReq,",
		"LookAheadAvailableFunds,LookAheadExcessLiquidity,",
		"HighestSeverity,DayTradesRemaining,Leverage,$LEDGER:ALL"}
	ib.Client.ReqAccountSummary(ib.GetReqID(), "All", strings.Join(tags, ""))
	ib.ReqPnL("DU1382837", "")
	// ib.Client.ReqOpenOrders()
	// ib.Client.ReqContractDetails(ib.GetReqID(), &hsik9)
	// ib.Client.ReqRealTimeBars(ib.GetReqID(), &hsik9, 5, "TRADES", false, nil)
	// ib.Client.ReqHistoricalData(ib.Client.GetReqID(), hsij9, "", "600 S", "1 min", "TRADES", false, 1, true, []TagValue{})
	// ef := ExecutionFilter{0, "", "DU1382837", "", "", "", ""}
	// ef := ExecutionFilter{}
	// ib.Client.ReqExecutions(ib.Client.GetReqID(), ef)
	// ib.Client.ReqMktData(1, hsij9, "", false, false, nil)
	// order := NewDefaultOrder()
	// order.LimitPrice = 30050
	// order.Action = "BUY"
	// order.OrderType = "LMT"
	// order.TotalQuantity = 1
	// ib.Client.ReqPnL(1, "DU1382837", "")
	// ib.Client.PlaceOrder(2271, &hsij9, order)
	// ib.ReqAccountSummary("", "")
}

// func (ib *IB) ReqAccountSummary(groupName string, tags string) {
// 	reqID := ib.GetReqID()
// 	ib.Wrapper.dataChanMap[reqID] = make(chan map[string]string, 5)
// 	ib.Client.ReqAccountSummary(reqID, groupName, tags)

// 	go func() {
// 		defer delete(ib.Wrapper.dataChanMap, reqID)
// 		for {
// 			select {
// 			case v, ok := <-ib.Wrapper.dataChanMap[reqID]:
// 				if ok {
// 					fmt.Println(v)
// 				} else {
// 					break
// 				}
// 			}
// 		}
// 	}()
// }

func (ib *IB) ReqPnL(account string, modelCode string) *PnL {
	reqID := ib.GetReqID()
	pnl := PnL{Account: account, ModelCode: modelCode}
	ib.Wrapper.PnLs[reqID] = &pnl
	ib.Client.ReqPnL(reqID, account, modelCode)
	return &pnl
}

func (ib *IB) ReqPnLSingle(account string, modelCode string, contractID int64) *PnLSingle {
	reqID := ib.GetReqID()
	pnl := PnL{Account: account, ModelCode: modelCode}
	pnlSingle := PnLSingle{ContractID: contractID, PnL: pnl}
	ib.Wrapper.PnLSingles[reqID] = &pnlSingle
	ib.Client.ReqPnLSingle(reqID, account, modelCode, contractID)
	return &pnlSingle
}

func (ib *IB) ReqHistoricalData(contract ibapi.Contract, endDateTime string, duration string, barSize string, whatToShow string, useRTH bool, formatDate int, keepUpToDate bool, chartOptions []ibapi.TagValue) *BarDataList {
	reqID := ib.GetReqID()
	bars := BarDataList{
		ReqID:          reqID,
		Contract:       contract,
		EndDateTime:    endDateTime,
		Duration:       duration,
		BarSizeSetting: barSize,
		WhatToShow:     whatToShow,
		UseRTH:         useRTH,
		FormatDate:     formatDate,
		KeepUpToDate:   keepUpToDate,
		ChartOptions:   chartOptions}

	ib.Wrapper.startReq(reqID, 10)

	if keepUpToDate {
		ib.Wrapper.startSubscription(reqID, 10, contract)
		barUpdateChan := ib.Wrapper.subDataChanMap[reqID]
		go func() {
		barUpdateLoop:
			for {
				select {
				case bar, ok := <-barUpdateChan:
					if !ok {
						break barUpdateLoop
					}
					newBar := *bar.(*ibapi.BarData)
					count := len(bars.BarList)
					if count == 0 {
						bars.BarList = append(bars.BarList, newBar)
					} else if bars.BarList[count-1].Date == newBar.Date {
						bars.BarList[count-1] = newBar
					} else {
						bars.BarList = append(bars.BarList, newBar)
					}
				}
			}
		}()
	}

	ib.Client.ReqHistoricalData(reqID, contract, endDateTime, duration, barSize, whatToShow, useRTH, formatDate, keepUpToDate, chartOptions)

	barChan := ib.Wrapper.dataChanMap[reqID]
barLoop:
	for {
		select {
		case bar, ok := <-barChan:
			if !ok {
				break barLoop
			}
			bars.BarList = append(bars.BarList, *bar.(*ibapi.BarData))
		}
	}

	return &bars
}

func (ib *IB) CancelHistoricalData(bars *BarDataList) {
	ib.Client.CancelHistoricalData(bars.ReqID)
	ib.Wrapper.endSubscription(bars.ReqID)
}

func (ib *IB) ReqOpenOrders() []ibapi.Order {
	reqID := int64(-ibapi.OPEN_ORDER)
	openOrders := []ibapi.Order{}
	ib.Wrapper.startReq(reqID, 10)
	ib.Client.ReqOpenOrders()
	openOrderChan := ib.Wrapper.dataChanMap[reqID]
openOrderLoop:
	for {
		select {
		case o, ok := <-openOrderChan:
			if !ok {
				break openOrderLoop
			}
			openOrders = append(openOrders, *o.(*ibapi.Order))

		}
	}

	return openOrders
}

func (ib *IB) ReqContractDetails(contract *ibapi.Contract) []ibapi.ContractDetails {
	reqID := ib.GetReqID()
	contractDetailsList := []ibapi.ContractDetails{}
	ib.Wrapper.startReq(reqID, 10)
	ib.Client.ReqContractDetails(reqID, contract)
	contractDetailsChan := ib.Wrapper.dataChanMap[reqID]
contractDetailsLoop:
	for {
		select {
		case contractDetails, ok := <-contractDetailsChan:
			if !ok {
				break contractDetailsLoop
			}
			contractDetailsList = append(contractDetailsList, *contractDetails.(*ibapi.ContractDetails))
		}
	}

	return contractDetailsList
}
