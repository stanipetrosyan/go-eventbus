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

			consumer(NewConsumerContextWithMessage(data))
		}
	}()
}

func NewSubscriber(ch chan Message) Subscriber {
	return defaultSubscriber{ch: ch}
}
