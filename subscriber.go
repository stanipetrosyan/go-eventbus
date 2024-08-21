package goeventbus

type Subscriber interface {
	Listen(consumer func(context Context))
}

type defaultSubscriber struct {
	ch chan Message
}

func (s defaultSubscriber) Listen(consumer func(context Context)) {
	go func() {
		for {
			data, ok := <-s.ch
			if !ok {
				return
			}

			consumer(newConsumerContextWithMessage(data))
		}
	}()
}

func newSubscriber(ch chan Message) Subscriber {
	return defaultSubscriber{ch: ch}
}
