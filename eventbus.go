package goeventbus

import (
	"sync"
)

type EventBus interface {
	Subscribe(address string, consumer func(context DeliveryContext))
	SubscribeOnce(address string, consumer func(context DeliveryContext))
	Publish(address string, message Message)
	Request(address string, message Message, consumer func(context DeliveryContext))
	AddInBoundInterceptor(address string, consumer func(context DeliveryContext))
}

type DefaultEventBus struct {
	topics map[string]*Topic
	rm     sync.RWMutex
	wg     sync.WaitGroup
}

func (e *DefaultEventBus) Subscribe(address string, consumer func(context DeliveryContext)) {
	e.subscribe(address, consumer, false)
}

func (e *DefaultEventBus) SubscribeOnce(address string, consumer func(context DeliveryContext)) {
	e.subscribe(address, consumer, true)
}

func (e *DefaultEventBus) Publish(address string, message Message) {
	e.rm.Lock()
	defer e.rm.Unlock()

	topic, exists := e.topics[address]
	if !exists {
		return
	}
	for _, handler := range topic.GetHandlers() {
		if !handler.closed {
			handler.Ch <- message
		}
	}
}

func (e *DefaultEventBus) Request(address string, message Message, consumer func(context DeliveryContext)) {
	e.rm.Lock()
	defer e.rm.Unlock()

	topic, exists := e.topics[address]
	if !exists {
		return
	}
	for _, item := range topic.Handlers {
		go func(handler *Handler, data Message) {
			if !handler.closed {
				handler.Ch <- data
				consumer(handler.Context)
			}
		}(item, message)
	}

}

func (e *DefaultEventBus) AddInBoundInterceptor(address string, consumer func(context DeliveryContext)) {
	_, exists := e.topics[address]
	if !exists {
		e.topics[address] = NewTopic(address)
	}

	ch := make(chan Message)
	context := NewDeliveryContext(e.topics[address].GetChannels())
	handler := Handler{Ch: ch, Consumer: consumer, Context: context, Address: address, closed: false}

	e.rm.Lock()
	e.topics[address].AddInterceptor(handler)
	e.rm.Unlock()

	go e.handle(&handler, false)
}

func (e *DefaultEventBus) subscribe(address string, consumer func(context DeliveryContext), once bool) {
	_, exists := e.topics[address]
	if !exists {
		e.topics[address] = NewTopic(address)
	}

	ch := make(chan Message)
	context := NewDeliveryContext([]chan Message{ch})
	handler := Handler{Ch: ch, Consumer: consumer, Context: context, Address: address, closed: false}

	e.rm.Lock()
	e.topics[address].AddHandler(handler)
	e.rm.Unlock()
	go e.handle(&handler, once)
}

func (e *DefaultEventBus) handle(handler *Handler, once bool) {
	e.wg.Add(1)
	go handler.Handle(once, &e.wg)
	e.wg.Wait()
}

func NewEventBus() EventBus {
	return &DefaultEventBus{
		topics: map[string]*Topic{},
	}
}
