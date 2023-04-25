package goeventbus

import (
	"sync"
)

type HandlerFunc func(DeliveryContext)

type EventBus interface {
	Subscribe(address string, consumer HandlerFunc)
	SubscribeOnce(address string, consumer HandlerFunc)
	Publish(address string, data any, options MessageOptions)
	Unsubscribe(address string)
	Request(address string, data any, options MessageOptions, consumer func(context DeliveryContext))
	AddInBoundInterceptor(address string, consumer func(context DeliveryContext))
}

type DefaultEventBus struct {
	handlers map[string][]*Handler
	topics   map[string]*Topic
	rm       sync.RWMutex
	wg       sync.WaitGroup
}

func (e *DefaultEventBus) Subscribe(address string, consumer HandlerFunc) {
	e.subscribe(address, consumer, false)
}

func (e *DefaultEventBus) SubscribeOnce(address string, consumer HandlerFunc) {
	e.subscribe(address, consumer, true)
}

func (e *DefaultEventBus) Publish(address string, data any, options MessageOptions) {
	e.rm.Lock()

	defer e.rm.Unlock()

	message := Message{Data: data, Headers: options.headers}

	topic, exists := e.topics[address]
	if !exists {
		return
	}

	for _, item := range topic.Handlers {
		go func(handler *Handler, data Message) {
			if !handler.closed {
				if len(handler.Interceptors) > 0 {
					handler.Interceptors[0].Ch <- data
					handler.Interceptors[0].Consumer(handler.Context)
				} else {
					handler.Ch <- data
				}
			}
		}(item, message)
	}

}

func (e *DefaultEventBus) Request(address string, data any, options MessageOptions, consumer func(context DeliveryContext)) {
	e.rm.Lock()

	message := Message{Data: data, Headers: options.headers}

	for _, item := range e.topics[address].Handlers {
		go func(handler *Handler, data Message) {
			if !handler.closed {
				handler.Ch <- data
				consumer(handler.Context)
			}
		}(item, message)
	}

	e.rm.Unlock()
}

func (e *DefaultEventBus) AddInBoundInterceptor(address string, consumer func(context DeliveryContext)) {
	for _, item := range e.handlers[address] {
		ch := make(chan Message)
		item.AddInterceptor(Interceptor{Ch: ch, Consumer: consumer})
	}
}

func (e *DefaultEventBus) Unsubscribe(address string) {
	e.removeHandler(address)
}

func (e *DefaultEventBus) subscribe(address string, consumer HandlerFunc, once bool) {
	ch := make(chan Message)
	context := DefaultDeliveryContext{ch: ch}
	handler := Handler{Ch: ch, Consumer: consumer, Context: &context, Address: address, closed: false}

	e.rm.Lock()
	_, exists := e.topics[address]
	if !exists {
		e.topics[address] = &Topic{Address: address, Handlers: []*Handler{}}
	}

	e.topics[address].Handlers = append(e.topics[address].Handlers, &handler)

	e.handlers[address] = append(e.handlers[address], &handler)
	e.rm.Unlock()

	go e.handle(&handler, once)
}

func (e *DefaultEventBus) handle(handler *Handler, once bool) {
	e.wg.Add(1)
	go handler.Handle(once, &e.wg)
	e.wg.Wait()
}

func (e *DefaultEventBus) removeHandler(address string) {
	e.rm.Lock()
	defer e.rm.Unlock()

	for _, handler := range e.topics[address].Handlers {
		handler.Close()
	}
	delete(e.topics, address)
}

func NewEventBus() EventBus {
	return &DefaultEventBus{
		handlers: map[string][]*Handler{},
		topics:   map[string]*Topic{},
	}
}
