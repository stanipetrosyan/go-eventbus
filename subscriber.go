package goeventbus

import (
	"errors"
)

type Subscriber interface {
	Listen(consumer func(context Context))
}

type defaultSubscriber struct {
	listenChannel <-chan Message
	sendChannel   chan<- packet
}

func (s defaultSubscriber) Listen(consumer func(context Context)) {
	go func() {
		for {
			message, ok := <-s.listenChannel
			if !ok {
				newContextWithError(errors.New("channel closed"))
				return
			}
			consumer(newContextWithMessageAndChannel(message, s.sendChannel))
		}
	}()
}

func newSubscriber(ch <-chan Message, channel chan packet) Subscriber {
	return defaultSubscriber{listenChannel: ch, sendChannel: channel}
}
