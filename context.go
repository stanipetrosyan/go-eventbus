package goeventbus

type Context interface {
	Result() Message
	Reply(message Message)
}

type defaultContext struct {
	message Message
	ch      chan packet
}

func (c defaultContext) Result() Message {
	return c.message
}

func (c defaultContext) Reply(message Message) {
	c.ch <- newSubscriberPacket(message)
}

func newConsumerContextWithMessageAndChannel(message Message, ch chan packet) Context {
	return defaultContext{message: message, ch: ch}
}
