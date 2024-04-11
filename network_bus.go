package goeventbus

type NetworkBus interface {
	Server() Server
	Client() Client
}

type defaultNetworkBus struct {
	localBus EventBus
	address  string
	path     string
}

func NewNetworkBus(bus EventBus, address, path string) NetworkBus {
	return defaultNetworkBus{localBus: bus, address: address, path: path}
}

func (b defaultNetworkBus) Server() Server {
	return NewServer()
}

func (b defaultNetworkBus) Client() Client {
	return NewClient()
}
