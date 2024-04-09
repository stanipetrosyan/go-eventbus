package goeventbus

type NetworkBus interface {
	Server()
	Client()
}

type defaultNetworkBus struct {
	localBus EventBus
	address  string
	path     string
}

func NewNetworkBus(bus EventBus, address, path string) NetworkBus {
	return defaultNetworkBus{localBus: bus, address: address, path: path}
}

func (b defaultNetworkBus) Server() {}

func (b defaultNetworkBus) Client() {}
