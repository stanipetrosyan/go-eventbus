package goeventbus

import (
	"sync"
)

type DataEvent struct {
	Data    interface{}
	Address string
}

type DataChannel chan DataEvent

type DataChannelSlice []DataChannel

type EventBus interface {
	Subscribe(address string)
	Publish(address string, data interface{})
	On(address string, handle func(data DataEvent))
}

type DefaultEventBus struct {
	subscribers map[string]DataChannel
	rm          sync.RWMutex
}

func (e *DefaultEventBus) Subscribe(address string) {
	e.rm.Lock()

	ch := make(chan DataEvent)
	e.subscribers[address] = ch

	e.rm.Unlock()
}

func (e *DefaultEventBus) Publish(address string, data interface{}) {
	e.rm.Lock()

	found := e.subscribers[address]

	go func(data DataEvent, ch DataChannel) {
		ch <- data
	}(DataEvent{Data: data, Address: address}, found)

	e.rm.Unlock()
}

func (e *DefaultEventBus) On(address string, handle func(data DataEvent)) {
	for {
		d := <-e.subscribers[address]
		println(d.Address)
		handle(d)
	}
}

func NewEventBus() EventBus {
	return &DefaultEventBus{
		subscribers: map[string]DataChannel{},
	}
}
