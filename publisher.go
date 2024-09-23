package goeventbus

type Publisher interface {
	// Publish the message on the channel
	Publish(message Message)
	// Publish the message on the channel and execute consume when receive a reply from a subscriber
	Request(message Message, consume func(context Context))
}

type defaultPublisher struct {
	ch      <-chan Message
	channel chan packet
}

func newPublisher(channel chan packet, ch <-chan Message) Publisher {
	return defaultPublisher{ch: ch, channel: channel}
}

func (p defaultPublisher) Publish(message Message) {
	p.channel <- newPublisherPacket(message)
}

func (p defaultPublisher) Request(message Message, consume func(context Context)) {
	p.channel <- newPublisherPacket(message)

	go func() {
		for {
			message, ok := <-p.ch
			if !ok {
				return
			}

			consume(newConsumerContextWithMessageAndChannel(message, p.channel))
		}
	}()
}
