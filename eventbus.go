package goeventbus

type EventBus interface {
	Channel(adress string) Channel
}

type defaultEventBus struct {
	channels map[string]Channel
}

func (e *defaultEventBus) Channel(address string) Channel {
	_, exists := e.channels[address]
	if !exists {
		e.channels[address] = NewChannel(address)
	}
	return e.channels[address]
}

func NewEventBus() EventBus {
	return &defaultEventBus{
		channels: map[string]Channel{},
	}
}
