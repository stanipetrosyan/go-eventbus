package goeventbus

import (
	"sync"
)

type EventBus interface {
	Subscribe(address string, callback func(context DeliveryContext))
	AddInBoundInterceptor(address string, callback func(context DeliveryContext))
	Publish(address string, message Message)
	Request(address string, message Message, callback func(context DeliveryContext))
}

type DefaultEventBus struct {
	topics map[string]*Topic
	rm     sync.RWMutex
	wg     sync.WaitGroup
}

func (e *DefaultEventBus) Subscribe(address string, callback func(context DeliveryContext)) {
	_, exists := e.topics[address]
	if !exists {
		e.topics[address] = NewTopic(address)
	}

	handler := NewConsumer(address, callback)
	channels := []chan Message{handler.Ch}
	handler = handler.SetContext(NewDeliveryContext(channels))

	e.rm.Lock()
	e.topics[address].AddHandler(handler)
	e.rm.Unlock()
	go e.handle(handler)
}

func (e *DefaultEventBus) AddInBoundInterceptor(address string, callback func(context DeliveryContext)) {
	_, exists := e.topics[address]
	if !exists {
		e.topics[address] = NewTopic(address)
	}

	channels := e.topics[address].GetChannels()
	handler := NewConsumer(address, callback).SetContext(NewDeliveryContext(channels))

	e.rm.Lock()
	e.topics[address].AddHandler(handler)
	e.rm.Unlock()
	go e.handle(handler)
}

func (e *DefaultEventBus) handle(handler *Handler) {
	e.wg.Add(1)
	go handler.Handle(&e.wg)
	e.wg.Wait()
}

func (e *DefaultEventBus) Publish(address string, message Message) {
	e.rm.Lock()
	defer e.rm.Unlock()

	topic, exists := e.topics[address]
	if !exists {
		return
	}
	for _, handler := range topic.GetHandlers() {
		if !handler.Closed {
			handler.Ch <- message
		}
	}
}

func (e *DefaultEventBus) Request(address string, message Message, callback func(context DeliveryContext)) {
	e.rm.Lock()
	defer e.rm.Unlock()

	topic, exists := e.topics[address]
	if !exists {
		return
	}
	for _, item := range topic.GetHandlers() {
		go func(handler *Handler, data Message) {
			if !handler.Closed {
				handler.Ch <- data
				callback(handler.Context)
			}
		}(item, message)
	}
}

func NewEventBus() EventBus {
	return &DefaultEventBus{
		topics: map[string]*Topic{},
	}
}
