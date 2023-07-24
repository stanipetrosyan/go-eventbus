package goeventbus

import (
	"sync"
)

type EventBus interface {
	Subscribe(address string, consumer func(context DeliveryContext))
	SubscribeOnce(address string, consumer func(context DeliveryContext))
	AddInBoundInterceptor(address string, consumer func(context DeliveryContext))
	Publish(address string, message Message)
	Request(address string, message Message, consumer func(context DeliveryContext))
}

type DefaultEventBus struct {
	topics map[string]*Topic
	rm     sync.RWMutex
	wg     sync.WaitGroup
}

func (e *DefaultEventBus) Subscribe(address string, consumer func(context DeliveryContext)) {
	e.subscribe(address, consumer, false, false)
}

func (e *DefaultEventBus) SubscribeOnce(address string, consumer func(context DeliveryContext)) {
	e.subscribe(address, consumer, true, false)
}

func (e *DefaultEventBus) AddInBoundInterceptor(address string, consumer func(context DeliveryContext)) {
	e.subscribe(address, consumer, false, true)
}

func (e *DefaultEventBus) subscribe(address string, consumer func(context DeliveryContext), once bool, interceptor bool) {
	_, exists := e.topics[address]
	if !exists {
		e.topics[address] = NewTopic(address)
	}

	ch := make(chan Message)

	var channels []chan Message = []chan Message{ch}
	var handlerType HandlerType = Consumer

	if interceptor {
		channels = e.topics[address].GetChannels()
		handlerType = Interceptor
	}

	context := NewDeliveryContext(channels)
	handler := Handler{Ch: ch, Consumer: consumer, Context: context, Address: address, Closed: false, Once: once, Type: handlerType}

	e.rm.Lock()
	e.topics[address].AddHandler(&handler)
	e.rm.Unlock()
	go e.handle(&handler)
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

func (e *DefaultEventBus) Request(address string, message Message, consumer func(context DeliveryContext)) {
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
				consumer(handler.Context)
			}
		}(item, message)
	}
}

func NewEventBus() EventBus {
	return &DefaultEventBus{
		topics: map[string]*Topic{},
	}
}
