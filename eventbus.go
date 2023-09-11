package goeventbus

import (
	"sync"
)

type EventBus interface {
	Subscribe(address string, callback func(context ConsumerContext))
	AddInBoundInterceptor(address string, callback func(context InterceptorContext))
	Publish(address string, message Message)
	Request(address string, message Message, callback func(context ConsumerContext))
}

type DefaultEventBus struct {
	topics map[string]*Topic
	rm     sync.RWMutex
	wg     sync.WaitGroup
}

func (e *DefaultEventBus) Subscribe(address string, callback func(context ConsumerContext)) {
	_, exists := e.topics[address]
	if !exists {
		e.topics[address] = NewTopic(address)
	}

	e.rm.Lock()
	handler := e.topics[address].AddConsumer(callback)
	e.rm.Unlock()
	go e.handle(Handler(handler))
}

func (e *DefaultEventBus) AddInBoundInterceptor(address string, callback func(context InterceptorContext)) {
	_, exists := e.topics[address]
	if !exists {
		e.topics[address] = NewTopic(address)
	}

	e.rm.Lock()
	handler := e.topics[address].AddInterceptor(callback)
	e.rm.Unlock()
	go e.handle(handler)
}

func (e *DefaultEventBus) handle(handler Handler) {
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
		if !handler.Closed() {
			handler.Chain() <- message
		}
	}
}

func (e *DefaultEventBus) Request(address string, message Message, callback func(context ConsumerContext)) {
	e.rm.Lock()
	defer e.rm.Unlock()

	topic, exists := e.topics[address]
	if !exists {
		return
	}
	for _, item := range topic.GetHandlers() {
		go func(handler Handler, data Message) {
			if !handler.Closed() {
				handler.Chain() <- data
				callback(NewConsumerContext(handler.Chain()))
			}
		}(item, message)
	}
}

func NewEventBus() EventBus {
	return &DefaultEventBus{
		topics: map[string]*Topic{},
	}
}
