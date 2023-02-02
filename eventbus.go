package goeventbus

import (
	"sync"
)

type EventBus interface {
	Subscribe(address string, consumer func(data Message))
	SubscribeOnce(address string, consumer func(data Message))
	Publish(address string, data any, options MessageOptions)
	Unsubscribe(address string)
}

type DefaultEventBus struct {
	handlers map[string][]Handler // Use a pointer ?
	rm       sync.RWMutex
	wg       sync.WaitGroup
}

func (e *DefaultEventBus) Subscribe(address string, consumer func(data Message)) {
	e.subscribe(address, consumer, false)
}

func (e *DefaultEventBus) SubscribeOnce(address string, consumer func(data Message)) {
	e.subscribe(address, consumer, true)
}

func (e *DefaultEventBus) Publish(address string, data any, options MessageOptions) {
	e.rm.Lock()

	message := Message{Data: data, Headers: options.headers}

	for _, item := range e.handlers[address] {
		go func(handler Handler, data Message) {
			handler.Ch <- data
		}(item, message)
	}

	e.rm.Unlock()
}

func (e *DefaultEventBus) Unsubscribe(address string) {
	e.removeHandler(address)
}

func (e *DefaultEventBus) subscribe(address string, consumer func(data Message), once bool) {
	ch := make(chan Message)
	handler := Handler{Ch: ch, Consume: consumer, Address: address}

	e.rm.Lock()
	e.handlers[address] = append(e.handlers[address], handler)
	e.rm.Unlock()

	go e.handle(handler, once)
}

// maybe a Handler function ?
func (e *DefaultEventBus) handle(handler Handler, once bool) {
	e.wg.Add(1)

	go func(handler Handler) {
		for {
			data, ok := <-handler.Ch
			if !ok {
				e.wg.Done()
				return
			}

			handler.Consume(data)

			if once {
				e.removeHandler(handler.Address)
				e.wg.Done()
			}
		}
	}(handler)

	e.wg.Wait()
}

func (e *DefaultEventBus) removeHandler(address string) {
	e.rm.Lock()
	defer e.rm.Unlock()

	for _, handler := range e.handlers[address] {
		handler.Close()
	}
	delete(e.handlers, address)
}

func NewEventBus() EventBus {
	return &DefaultEventBus{
		handlers: map[string][]Handler{},
	}
}
