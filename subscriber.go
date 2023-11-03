package goeventbus

type Subscriber interface {
	Listen(consumer func())
}

type defaultSubscriber struct {
	ch chan string
}

func (s defaultSubscriber) Listen(consumer func()) {
	go func() {
		for {
			data, ok := <-s.ch
			if !ok {
				return
			}

			println(data)

			consumer()
		}
	}()
}

func NewSubscriber(ch chan string) Subscriber {
	return defaultSubscriber{ch: ch}
}
