package goeventbus

import (
	"fmt"
	"sync"
)

type DataEvent struct {
	Data    interface{}
	Address string
}

type DataChannel chan DataEvent

type Channel struct {
	Ch       chan DataEvent
	Consumer func(data DataEvent)
}

type EventBus interface {
	Subscribe(address string)
	Publish(address string, data interface{})
	On(address string, handle func(data DataEvent))
	Unsubscribe(address string)
}

type DefaultEventBus struct {
	subscribers map[string]Channel
	rm          sync.RWMutex
}

func (e *DefaultEventBus) Subscribe(address string) {
	e.rm.Lock()

	ch := make(chan DataEvent)
	e.subscribers[address] = Channel{Ch: ch, Consumer: func(data DataEvent) { println("consumer not started") }}

	e.rm.Unlock()
}

func (e *DefaultEventBus) Publish(address string, data interface{}) {
	e.rm.Lock()

	found := e.subscribers[address]
	go func(data DataEvent, ch Channel) {
		ch.Ch <- data
	}(DataEvent{Data: data, Address: address}, found)

	e.rm.Unlock()
}

func (e *DefaultEventBus) consume() {
	var wg sync.WaitGroup
	wg.Add(len(e.subscribers))

	for address, ch := range e.subscribers {
		println(address)

		go func(ch Channel) {
			for d := range ch.Ch {
				fmt.Printf("%s\n", d.Data)
				ch.Consumer(d)
			}
		}(ch)
	}
	wg.Wait()

}

func (e *DefaultEventBus) On(address string, handle func(data DataEvent)) {
	ch := e.subscribers[address]

	e.subscribers[address] = Channel{Ch: ch.Ch, Consumer: handle}

	go e.consume()
}

func (e *DefaultEventBus) Unsubscribe(address string) {
	e.rm.Lock()

	println("removed address")
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
