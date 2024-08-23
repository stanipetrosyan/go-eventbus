package goeventbus

type Subscriber interface {
	Listen(consumer func(context Context))
}

type defaultSubscriber struct {
	ch      chan Message
	channel chan packet
}

func (s defaultSubscriber) Listen(consumer func(context Context)) {
	go func() {
		for {
			message, ok := <-s.ch
			if !ok {
				return
			}

			consumer(newConsumerContextWithMessageAndChannel(message, s.channel))
		}
	}()
}

func newSubscriber(ch chan Message, channel chan packet) Subscriber {
	return defaultSubscriber{ch: ch, channel: channel}
}
