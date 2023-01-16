package goeventbus

import (
	"sync"
)

type EventBus interface {
	Subscribe(address string)
	Publish(address string, data any, options MessageOptions)
	On(address string, consumer func(data Message))
	Unsubscribe(address string)
	handle(handler Handler)
}

type DefaultEventBus struct {
	handlers map[string]Handler
	rm       sync.RWMutex
	wg       sync.WaitGroup
}

func (e *DefaultEventBus) Subscribe(address string) {
	e.rm.Lock()

	ch := make(chan Message)
	e.handlers[address] = Handler{Ch: ch, Consume: func(data Message) {}}

	e.rm.Unlock()
}

func (e *DefaultEventBus) Publish(address string, data any, options MessageOptions) {
	e.rm.Lock()

	message := Message{Data: data, Headers: options.headers}
	found := e.handlers[address]

	go func(data Message, ch Handler) {
		ch.Ch <- data
	}(message, found)

	e.rm.Unlock()
}

func (e *DefaultEventBus) On(address string, consumer func(data Message)) {
	ch := e.handlers[address]

	e.handlers[address] = Handler{Ch: ch.Ch, Consume: consumer}

	go e.handle(e.handlers[address])
}

func (e *DefaultEventBus) handle(handler Handler) {
	e.wg.Add(1)

	go func(Handler Handler) {
		for {
			data, ok := <-handler.Ch

			if !ok {
				e.wg.Done()
				return
			}
			Handler.Consume(data)
		}
	}(handler)

	e.wg.Wait()
}

func (e *DefaultEventBus) Unsubscribe(address string) {
	e.rm.Lock()

	ch := e.handlers[address]
	close(ch.Ch)
	delete(e.handlers, address)

	e.rm.Unlock()
}

func NewEventBus() EventBus {
	return &DefaultEventBus{
		handlers: map[string]Handler{},
	}
}
