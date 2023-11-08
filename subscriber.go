package goeventbus

type Subscriber interface {
	Listen(consumer func(message Message))
}

type defaultSubscriber struct {
	ch chan Message
}

func (s defaultSubscriber) Listen(consumer func(message Message)) {
	go func() {
		for {
			data, ok := <-s.ch
			if !ok {
				return
			}

			consumer(data)
		}
	}()
}

func NewSubscriber(ch chan Message) Subscriber {
	return defaultSubscriber{ch: ch}
}
