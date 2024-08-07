package goeventbus

import "log/slog"

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
			slog.Error("Something went wrong during listening on channel")
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
	slog.Info("Publisher created", slog.String("channel", c.address))

	return NewPublisher(c.ch)
}

func (c *defaultChannel) Subscriber() Subscriber {
	ch := make(chan Message)
	c.chs = append(c.chs, ch)

	slog.Info("Subscriber created", slog.String("channel", c.address))

	return NewSubscriber(ch)
}

func (c *defaultChannel) Processor(predicate func(message Message) bool) Channel {
	c.processor = NewProcessorWithPredicate(predicate)

	slog.Info("Processor created", slog.String("channel", c.address))
	return c
}

func NewChannel(address string) Channel {
	ch := make(chan Message)
	channel := defaultChannel{address: address, ch: ch, chs: []chan Message{}, processor: NewProcessor()}
	go channel.Listen()

	return &channel
}
