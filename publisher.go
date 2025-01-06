package goeventbus

import (
	"errors"
)

type Publisher interface {
	// Publish the message on the channel
	Publish(message Message)
	// Publish the message on the channel and execute consume when receive a reply from a subscriber
	Request(message Message, consumer func(context Context))
}

type defaultPublisher struct {
	listenChannel <-chan Message
	sendChannel   chan<- packet
}

func newPublisher(channel chan packet, ch <-chan Message) Publisher {
	return defaultPublisher{listenChannel: ch, sendChannel: channel}
}

func (p defaultPublisher) Publish(message Message) {
	p.sendChannel <- newPublisherPacket(message)
}

func (p defaultPublisher) Request(message Message, consumer func(context Context)) {
	p.sendChannel <- newPublisherPacket(message)

	go func() {
		for {
			message, ok := <-p.listenChannel
			if !ok {
				newContextWithError(errors.New("channel closed"))
				return
			}
			consumer(newContextWithMessageAndChannel(message, p.sendChannel))
		}
	}()
}
