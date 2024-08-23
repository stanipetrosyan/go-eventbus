package goeventbus

import "log/slog"

type packet struct {
	from    string
	message Message
}

func newSubscriberPacket(message Message) packet {
	return packet{from: "subscriber", message: message}
}

func newPublisherPacket(message Message) packet {
	return packet{from: "publisher", message: message}
}

type Channel interface {
	Publisher() Publisher
	Subscriber() Subscriber
	Processor(predicate func(message Message) bool) Channel
}

type defaultChannel struct {
	address     string
	ch          chan packet
	subscribers []chan Message
	publishers  []chan Message
	processor   Processor
}

func (c *defaultChannel) Listen() {
	for {
		packet, ok := <-c.ch
		if !ok {
			slog.Error("Something went wrong during listening on channel")
			return
		}

		if packet.from == "subscriber" {
			for _, item := range c.publishers {
				item <- packet.message
			}
		}

		if c.processor.forward(packet.message) {
			for _, item := range c.subscribers {
				item <- packet.message
			}
		}

	}
}

func (c *defaultChannel) Publisher() Publisher {
	ch := make(chan Message)
	c.publishers = append(c.publishers, ch)
	slog.Info("Publisher created", slog.String("channel", c.address))

	return newPublisher(c.ch)
}

func (c *defaultChannel) Subscriber() Subscriber {
	ch := make(chan Message)
	c.subscribers = append(c.subscribers, ch)

	slog.Info("Subscriber created", slog.String("channel", c.address))

	return newSubscriber(ch, c.ch)
}

func (c *defaultChannel) Processor(predicate func(message Message) bool) Channel {
	c.processor = newProcessorWithPredicate(predicate)

	slog.Info("Processor created", slog.String("channel", c.address))
	return c
}

func newChannel(address string) Channel {
	ch := make(chan packet)
	channel := defaultChannel{
		address:     address,
		ch:          ch,
		subscribers: []chan Message{},
		publishers:  []chan Message{},
		processor:   newProcessor(),
	}
	go channel.Listen()

	return &channel
}
