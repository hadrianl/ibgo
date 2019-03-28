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
	ib.Client.ReqCurrentTime()
	ib.Client.ReqAutoOpenOrders(true)
	ib.Client.ReqAccountUpdates(true, "")
	ib.Client.ReqExecutions(ib.Client.GetReqID(), ExecutionFilter{})
}
