package goeventbus

type Context interface {
	Result() Message
}

type defaultContext struct {
	message Message
}

func (c defaultContext) Result() Message {
	return c.message
}

func newConsumerContextWithMessage(message Message) Context {
	return defaultContext{message: message}
}
