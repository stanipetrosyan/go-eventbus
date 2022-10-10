package goeventbus

import (
	"sync"
)

type DataEvent struct {
	Data    interface{} // we have defined underlying data to be an interface which means it can be any value
	Address string
}

type DataChannel chan DataEvent

type DataChannelSlice []DataChannel

type EventBus interface {
	Subscribe(address string, ch DataChannel) // find method for don't use ch inside
	Publish(address string, data interface{})
}

type DefaultEventBus struct {
	subscribers map[string]DataChannelSlice
	rm          sync.RWMutex
}

func (e *DefaultEventBus) Subscribe(address string, ch DataChannel) {
	e.rm.Lock()
	if prev, found := e.subscribers[address]; found {
		e.subscribers[address] = append(prev, ch)
	} else {
		e.subscribers[address] = append([]DataChannel{}, ch)
	}
	e.rm.Unlock()
}

func (e *DefaultEventBus) Publish(address string, data interface{}) {
	e.rm.Lock()

	if chans, found := e.subscribers[address]; found {
		channels := append(DataChannelSlice{}, chans...)

		go func(data DataEvent, slices DataChannelSlice) {
			for _, ch := range slices {
				ch <- data
			}
		}(DataEvent{Data: data, Address: address}, channels)
	}
	e.rm.Unlock()
}

func NewEventBus() EventBus {
	return &DefaultEventBus{
		subscribers: map[string]DataChannelSlice{},
	}
}
