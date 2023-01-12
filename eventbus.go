package goeventbus

import (
	"sync"
)

type Handler struct {
	Ch      chan Message
	Consume func(data Message)
}

type EventBus interface {
	Subscribe(address string)
	Publish(address string, data any, options MessageOptions)
	On(address string, handle func(data Message))
	Unsubscribe(address string)
	handle(address string)
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

	message := Message{Data: data}
	found := e.handlers[address]

	go func(data Message, ch Handler) {
		ch.Ch <- data
	}(message, found)

	e.rm.Unlock()
}

func (e *DefaultEventBus) handle(address string) {
	e.wg.Add(1)

	ch := e.handlers[address]

	go func(Handler Handler) {
		for {
			data, ok := <-ch.Ch

			if !ok {
				e.wg.Done()
				return
			}
			Handler.Consume(data)
		}
	}(ch)

	e.wg.Wait()
}

func (e *DefaultEventBus) On(address string, handle func(data Message)) {
	ch := e.handlers[address]

	e.handlers[address] = Handler{Ch: ch.Ch, Consume: handle}

	go e.handle(address)
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
