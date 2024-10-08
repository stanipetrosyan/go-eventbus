package goeventbus

type NetworkBus interface {
	Server() Server
	Client() Client
}

type defaultNetworkBus struct {
	localBus EventBus
	address  string
}

func NewNetworkBus(bus EventBus, address string) NetworkBus {
	return defaultNetworkBus{localBus: bus, address: address}
}

func (b defaultNetworkBus) Server() Server {
	return newServer(b.address)
}

func (b defaultNetworkBus) Client() Client {
	return newClient(b.address, b.localBus)
}
