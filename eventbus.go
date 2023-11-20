package goeventbus

type EventBus interface {
	Channel(adress string) Channel
}

type DefaultEventBus struct {
	channels map[string]Channel
}

func (e *DefaultEventBus) Channel(address string) Channel {
	_, exists := e.channels[address]
	if !exists {
		e.channels[address] = NewChannel(address)
	}
	return e.channels[address]
}

func NewEventBus() EventBus {
	return &DefaultEventBus{
		channels: map[string]Channel{},
	}
}
