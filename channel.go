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
	processor Processor
}

func (c *defaultChannel) Listen() {
	for {
		data, ok := <-c.ch
		if !ok {
			return
		}

		if c.processor.forward(data) {
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
	c.processor = NewProcessorWithPredicate(predicate)

	return c
}

func NewChannel(address string) Channel {
	ch := make(chan Message)
	channel := defaultChannel{address: address, ch: ch, chs: []chan Message{}, processor: NewProcessor()}
	go channel.Listen()

	return &channel
}
