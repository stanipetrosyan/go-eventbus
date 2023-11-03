package goeventbus

type Channel interface {
	Register() Channel
	Publisher() Publisher
	Subscriber() Subscriber
}

type defaultChannel struct {
	address string
	ch      chan string
}

func (c defaultChannel) Register() Channel {
	return c
}

func (c defaultChannel) Publisher() Publisher {
	println("publisher created")
	return NewPublisher(c.ch)
}

func (c defaultChannel) Subscriber() Subscriber {
	println("subscriber created")
	return NewSubscriber(c.ch)
}

func NewChannel(address string) Channel {
	ch := make(chan string)
	return defaultChannel{address: address, ch: ch}
}
