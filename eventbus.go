package goeventbus

import (
	"fmt"
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
	Unsubscribe(address string)
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
	ch := e.subscribers[address]
	for d := range ch {
		fmt.Println(d.Address)
		handle(d)
	}
}

func (e *DefaultEventBus) Unsubscribe(address string) {
	e.rm.Lock()

	ch := e.subscribers[address]
	close(ch)
	delete(e.subscribers, address)

	e.rm.Unlock()
}

func NewEventBus() EventBus {
	return &DefaultEventBus{
		subscribers: map[string]DataChannel{},
	}
}
