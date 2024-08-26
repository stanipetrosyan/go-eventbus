package goeventbus

type Publisher interface {
	Publish(message Message)
	Request(message Message, consumer func(context Context))
}

type defaultPublisher struct {
	ch      chan Message
	channel chan packet
}

func newPublisher(channel chan packet, ch chan Message) Publisher {
	return defaultPublisher{ch: ch, channel: channel}
}

func (p defaultPublisher) Publish(message Message) {
	p.channel <- newPublisherPacket(message)
}

func (p defaultPublisher) Request(message Message, consumer func(context Context)) {
	p.channel <- newPublisherPacket(message)

	go func() {
		for {
			message, ok := <-p.ch
			if !ok {
				return
			}

			consumer(newConsumerContextWithMessageAndChannel(message, p.channel))
		}
	}()
}
