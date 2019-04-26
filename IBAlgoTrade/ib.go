package IBAlgoTrade

import (
	"sync/atomic"
	"time"

	"github.com/hadrianl/ibgo/ibapi"
)

type IB struct {
	Client   *ibapi.IbClient
	Wrapper  ibapi.IbWrapper
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
	wrapper := GoWrapper{ib: ib}
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
	// hsij9 := Contract{355299154, "HSI", "FUT", "20190429", 0, "?", "50", "HKFE", "HKD", "HSIJ9", "HSI", "", false, "", "", "", nil, nil}
	// fmt.Println(hsij9)
	// ib.Client.ReqCurrentTime()
	// ib.Client.ReqAutoOpenOrders(true)
	// ib.Client.ReqAccountUpdates(true, "")

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
	time.Sleep(time.Second * 3)
	ib.Client.ReqPnL(1, "DU1382837", "")
	// ib.Client.PlaceOrder(2271, &hsij9, order)
}

func (ib *IB) ReqAccountSummary(groupName string, tags string)
