package goeventbus

type Channel interface {
	Publisher() Publisher
	Subscriber() Subscriber
	Processor(predicate func(message Message) bool) Channel
}

type defaultChannel struct {
	address   string
	ch        chan Message
	chs       []chan Message
	predicate func(message Message) bool
}

func (c *defaultChannel) Listen() {
	for {
		data, ok := <-c.ch
		if !ok {
			return
		}

		if c.predicate(data) {
			for _, item := range c.chs {
				item <- data
			}
		}

	}
}

func (c *defaultChannel) Publisher() Publisher {
	return NewPublisher(c.ch)
}

func (c *defaultChannel) Subscriber() Subscriber {
	ch := make(chan Message)
	c.chs = append(c.chs, ch)

	return NewSubscriber(ch)
}

func (c *defaultChannel) Processor(predicate func(message Message) bool) Channel {
	c.predicate = predicate

	return c
}

func NewChannel(address string) Channel {
	ch := make(chan Message)
	predicate := func(message Message) bool { return true }
	channel := defaultChannel{address: address, ch: ch, chs: []chan Message{}, predicate: predicate}
	go channel.Listen()

	return &channel
}
