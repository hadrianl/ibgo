package ibgo

type IB struct {
	Client   *IbClient
	Wrapper  Wrapper
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
	ibwrapper := Wrapper{}
	ibclient := NewIbClient(ibwrapper)
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
	// hsij9 := Contract{355299154, "HSI", "FUT", "20190429", 0, "?", "50", "HKFE", "HKD", "HSIJ9", "HSI", "", false, "", "", "", []ComboLeg{}, DeltaNeutralContract{}}
	// fmt.Println(hsij9)
	ib.Client.ReqCurrentTime()
	ib.Client.ReqAutoOpenOrders(true)
	ib.Client.ReqAccountUpdates(true, "")

	// ib.Client.ReqHistoricalData(ib.Client.GetReqID(), hsij9, "", "600 S", "1 min", "TRADES", false, 1, true, []TagValue{})
	// ef := ExecutionFilter{0, "", "DU1382837", "", "", "", ""}
	// ib.Client.ReqExecutions(699, ef)
}
