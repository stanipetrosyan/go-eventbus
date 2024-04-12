package goeventbus

type NetworkBus interface {
	Server(address, path string) Server
	Client(address, path string) Client
}

type defaultNetworkBus struct {
	localBus EventBus
	address  string
	path     string
}

func NewNetworkBus(bus EventBus, address, path string) NetworkBus {
	return defaultNetworkBus{localBus: bus, address: address, path: path}
}

func (b defaultNetworkBus) Server(address, path string) Server {
	return NewServer(address, path)
}

func (b defaultNetworkBus) Client(address, path string) Client {
	return NewClient()
}
