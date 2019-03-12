package ibgo

type IB struct {
	Client   *IbClient
	Wrapper  *IbWrapper
	host     string
	port     int
	clientId int8
}

func NewIB(host string, port int, clientId int8) *IB {
	ib := &IB{
		Client:   &IbClient{},
		Wrapper:  &IbWrapper{},
		host:     host,
		port:     port,
		clientId: clientId,
	}

	return ib
}

// Connect to TWS or Gateway
func (ib *IB) Connect(host string, port int, clientId int8) error {
	ib.host = host
	ib.port = port
	ib.clientId = clientId
	ib.Client.reset()
	err := ib.Client.Connect(host, port, clientId)

	return err
}
