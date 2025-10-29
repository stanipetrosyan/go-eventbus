package goeventbus

import (
	"log/slog"
	"sync"
)

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
	mu          sync.RWMutex
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
				c.mu.RLock()
				publishers := make([]chan Message, len(c.publishers))
				copy(publishers, c.publishers)
				c.mu.RUnlock()
				// Non-blocking broadcast using goroutines
				for _, item := range publishers {
					go func(ch chan Message) {
						ch <- packet.message
					}(item)
				}
			}

		case PUBLISHER:
			{
				c.mu.RLock()
				processor := c.processor
				subscribers := make([]chan Message, len(c.subscribers))
				copy(subscribers, c.subscribers)
				c.mu.RUnlock()
				if processor != nil {
					processor <- packet.message
				} else {
					// Non-blocking broadcast using goroutines
					for _, item := range subscribers {
						go func(ch chan Message) {
							ch <- packet.message
						}(item)
					}
				}

			}
		case PROCESSOR:
			{
				c.mu.RLock()
				subscribers := make([]chan Message, len(c.subscribers))
				copy(subscribers, c.subscribers)
				c.mu.RUnlock()
				// Non-blocking broadcast using goroutines
				for _, item := range subscribers {
					go func(ch chan Message) {
						ch <- packet.message
					}(item)
				}
			}
		}

	}
}

func (c *defaultChannel) Publisher() Publisher {
	ch := make(chan Message)
	c.mu.Lock()
	c.publishers = append(c.publishers, ch)
	c.mu.Unlock()
	slog.Info("Publisher created", slog.String("channel", c.address))

	return newPublisher(c.ch, ch)
}

func (c *defaultChannel) Subscriber() Subscriber {
	ch := make(chan Message)
	c.mu.Lock()
	c.subscribers = append(c.subscribers, ch)
	c.mu.Unlock()

	slog.Info("Subscriber created", slog.String("channel", c.address))

	return newSubscriber(ch, c.ch)
}

func (c *defaultChannel) Processor() Processor {
	ch := make(chan Message)
	processor := newProcessor(ch, c.ch)

	c.mu.Lock()
	c.processor = ch
	c.mu.Unlock()
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
