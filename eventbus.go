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
	Publish(address string, data Message)
	On(address string, handle func(data Message))
	Unsubscribe(address string)
	handle(address string)
}

type DefaultEventBus struct {
	Handlers map[string]Handler
	rm       sync.RWMutex
	wg       sync.WaitGroup
}

func (e *DefaultEventBus) Subscribe(address string) {
	e.rm.Lock()

	ch := make(chan Message)
	e.Handlers[address] = Handler{Ch: ch, Consume: func(data Message) {}}

	e.rm.Unlock()
}

func (e *DefaultEventBus) Publish(address string, data Message) {
	e.rm.Lock()

	found := e.Handlers[address]
	go func(data Message, ch Handler) {
		ch.Ch <- data
	}(data, found)

	e.rm.Unlock()
}

func (e *DefaultEventBus) handle(address string) {
	e.wg.Add(1)

	ch := e.Handlers[address]

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
	ch := e.Handlers[address]

	e.Handlers[address] = Handler{Ch: ch.Ch, Consume: handle}

	go e.handle(address)
}

func (e *DefaultEventBus) Unsubscribe(address string) {
	e.rm.Lock()

	ch := e.Handlers[address]
	close(ch.Ch)
	delete(e.Handlers, address)

	e.rm.Unlock()
}

func NewEventBus() EventBus {
	return &DefaultEventBus{
		Handlers: map[string]Handler{},
	}
}
