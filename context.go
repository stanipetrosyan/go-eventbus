package goeventbus

type Context interface {
	// Returns the message received
	Result() Message
	// Respond to the publisher that send a Request
	Reply(message Message)
	// Returns error
	Error() error
}

type defaultContext struct {
	message Message
	ch      chan<- packet
	err     error
}

func (c defaultContext) Result() Message {
	return c.message
}

func (c defaultContext) Reply(message Message) {
	c.ch <- newSubscriberPacket(message)
}

func (c defaultContext) Error() error {
	return c.err
}

func newContextWithError(err error) Context {
	return defaultContext{err: err}
}

func newContextWithMessageAndChannel(message Message, ch chan<- packet) Context {
	return defaultContext{message: message, ch: ch, err: nil}
}
