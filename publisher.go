package goeventbus

type Publisher interface {
	Publish(message Message)
	Request(message Message, consumer func(context Context))
}

type defaultPublisher struct {
	ch chan packet
}

func newPublisher(ch chan packet) Publisher {
	return defaultPublisher{ch: ch}
}

func (p defaultPublisher) Publish(message Message) {
	p.ch <- newPublisherPacket(message)
}

func (p defaultPublisher) Request(message Message, consumer func(context Context)) {
	p.ch <- newPublisherPacket(message)

	go func() {
		for {
			packet, ok := <-p.ch
			if !ok {
				return
			}

			consumer(newConsumerContextWithMessageAndChannel(packet.message, p.ch))
		}
	}()
}
