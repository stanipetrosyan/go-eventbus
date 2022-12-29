package goeventbus

import (
	"sync"
)

type Message struct {
	Data interface{}
}

type DataChannel chan Message

type Channel struct {
	Ch       chan Message
	Consumer func(data Message)
}

type EventBus interface {
	Subscribe(address string)
	Publish(address string, data interface{})
	On(address string, handle func(data Message))
	Unsubscribe(address string)
}

type DefaultEventBus struct {
	subscribers map[string]Channel
	rm          sync.RWMutex
}

func (e *DefaultEventBus) Subscribe(address string) {
	e.rm.Lock()

	ch := make(chan Message)
	e.subscribers[address] = Channel{Ch: ch, Consumer: func(data Message) { println("consumer not started") }}

	e.rm.Unlock()
}

func (e *DefaultEventBus) Publish(address string, data interface{}) {
	e.rm.Lock()

	found := e.subscribers[address]
	go func(data Message, ch Channel) {
		ch.Ch <- data
	}(Message{Data: data}, found)

	e.rm.Unlock()
}

func (e *DefaultEventBus) consume(address string) {
	var wg sync.WaitGroup
	wg.Add(len(e.subscribers))

	// for use address here, wg should be a top level declaretion
	for _, ch := range e.subscribers {
		go func(ch Channel) {
			for {
				data, ok := <-ch.Ch

				if !ok {
					wg.Done()
					return
				}
				ch.Consumer(data)
			}
		}(ch)
	}
	wg.Wait()

}

func (e *DefaultEventBus) On(address string, handle func(data Message)) {
	ch := e.subscribers[address]

	e.subscribers[address] = Channel{Ch: ch.Ch, Consumer: handle}

	go e.consume(address)
}

func (e *DefaultEventBus) Unsubscribe(address string) {
	e.rm.Lock()

	ch := e.subscribers[address]
	close(ch.Ch)
	delete(e.subscribers, address)

	e.rm.Unlock()
}

func NewEventBus() EventBus {
	return &DefaultEventBus{
		subscribers: map[string]Channel{},
	}
}
