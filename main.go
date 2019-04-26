package main

import (
<<<<<<< HEAD
	"time"

	"github.com/hadrianl/ibgo/ibapi"
	log "github.com/sirupsen/logrus"
)

type IB struct {
	Client   *ibapi.IbClient
	Wrapper  ibapi.IbWrapper
	host     string
	port     int
	clientID int64
}

func NewIB(host string, port int, clientID int64) *IB {
	ib := &IB{
		host:     host,
		port:     port,
		clientID: clientID,
	}
	ibwrapper := ibapi.Wrapper{}
	ibclient := ibapi.NewIbClient(ibwrapper)
	ib.Client = ibclient
	ib.Wrapper = ibwrapper

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

func main() {
	log.SetLevel(log.InfoLevel)
	ib := NewIB("127.0.0.1", 7497, 0)
	if err := ib.Connect(); err != nil {
		log.Panicf("Connect failed: %v", err)
	}
	ib.DoSomeTest()
	time.Sleep(time.Second * 20)
	if err := ib.DisConnect(); err != nil {
		log.Panicf("DisConnect failed: %v", err)
	}
	// time.Sleep(time.Second * 10)
	// if err := ib.Connect(); err != nil {
	// 	log.Panicf("Connect failed: %v", err)
	// }
	// ib.DoSomeTest()
	// time.Sleep(time.Second * 10)
	// if err := ib.DisConnect(); err != nil {
	// 	log.Panicf("DisConnect failed: %v", err)
	// }
=======
	"fmt"
	"strconv"
)

// type testingStruct struct {
// 	A  int64
// 	BS *B
// }
// type B struct {
// 	C int64
// }

func main() {
	// f := ""
	// b := []byte(f)
	// c := []byte(" ")
	// fmt.Println(len(bytes.Split(b, c)))
	// fmt.Println("string:", f, "byte:", b, c)
	// var t []testingStruct
	// fmt.Println(t)
	// fmt.Println(len(t))
	// fmt.Println(t == nil)
	// t = []testingStruct{}
	// fmt.Println(t == nil)
	// fmt.Println(new(testingStruct) == nil)
	// fmt.Println(&testingStruct{} == nil)
	// var ts *testingStruct
	// fmt.Println(ts == nil)
	// var t *testingStruct
	// fmt.Println(t == nil)
	// fmt.Println(t.BS)
	// fmt.Println(t.BS == nil)
	i, err := strconv.ParseInt(string([]byte{'1'}), 10, 64)
	fmt.Println(i, err)
>>>>>>> f2840998567b04c698134653f69f15c8f6fa7b40
}
