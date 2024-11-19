package goeventbus

import (
	"errors"
)

type Subscriber interface {
	Listen(consumer func(context Context))
}

type defaultSubscriber struct {
	ch      <-chan Message
	channel chan packet
}

func (s defaultSubscriber) Listen(consumer func(context Context)) {
	go func() {
		select {
		case message, ok := <-s.ch:
			if !ok {
				newContextWithError(errors.New("channel closed"))
				return
			}
			consumer(newContextWithMessageAndChannel(message, s.channel))
		}
	}()
}

func newSubscriber(ch <-chan Message, channel chan packet) Subscriber {
	return defaultSubscriber{ch: ch, channel: channel}
}
