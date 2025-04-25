package goeventbus

import "log/slog"

type Sender int

const (
	SUBSCRIBER Sender = iota
	PUBLISHER
	PROCESSOR
)

type Channel interface {
	// Create a Publisher for the channel. A publisher publish to all subscriber
	Publisher() Publisher

	// Create a Subscriber for the channel. A subscriber listen all channel messages
	Subscriber() Subscriber

	// Create a Processor for the channel. A processor forward the message if the predicate returns true.
	Processor() Processor
}

type defaultChannel struct {
	address     string
	ch          chan packet
	subscribers []chan Message
	publishers  []chan Message
	processor   chan Message
}

func (c *defaultChannel) listen() {
	for {
		packet, ok := <-c.ch
		if !ok {
			slog.Error("Something went wrong during listening on channel")
			return
		}

		switch packet.from {
		case SUBSCRIBER:
			{
				for _, item := range c.publishers {
					item <- packet.message
				}
			}

		case PUBLISHER:
			{
				if c.processor != nil {
					c.processor <- packet.message
				} else {
					for _, item := range c.subscribers {
						item <- packet.message
					}
				}

			}
		case PROCESSOR:
			{
				for _, item := range c.subscribers {
					item <- packet.message
				}
			}
		}

	}
}

func (c *defaultChannel) Publisher() Publisher {
	ch := make(chan Message)
	c.publishers = append(c.publishers, ch)
	slog.Info("Publisher created", slog.String("channel", c.address))

	return newPublisher(c.ch, ch)
}

func (c *defaultChannel) Subscriber() Subscriber {
	ch := make(chan Message)
	c.subscribers = append(c.subscribers, ch)

	slog.Info("Subscriber created", slog.String("channel", c.address))

	return newSubscriber(ch, c.ch)
}

func (c *defaultChannel) Processor() Processor {
	ch := make(chan Message)
	processor := newProcessor(ch, c.ch)

	c.processor = ch
	slog.Info("Processor created", slog.String("channel", c.address))
	return processor
}

func newChannel(address string) Channel {
	ch := make(chan packet)
	channel := defaultChannel{
		address:     address,
		ch:          ch,
		subscribers: []chan Message{},
		publishers:  []chan Message{},
		processor:   nil,
	}
	go channel.listen()

	return &channel
}
