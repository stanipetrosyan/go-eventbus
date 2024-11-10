package goeventbus

type Message struct {
	Payload interface{}
	Headers MessageHeaders
}

type MessageBuilder interface {
	SetPayload(payload any) MessageBuilder
	SetHeaders(headers MessageHeaders) MessageBuilder
	Build() Message
}

type defaultMessageBuilder struct {
	message Message
}

func NewMessageBuilder() MessageBuilder {
	return &defaultMessageBuilder{message: Message{}}
}

func (mb *defaultMessageBuilder) SetPayload(payload any) MessageBuilder {
	mb.message.Payload = payload
	return mb
}

func (mb *defaultMessageBuilder) SetHeaders(headers MessageHeaders) MessageBuilder {
	mb.message.Headers = headers
	return mb
}

func (mb *defaultMessageBuilder) Build() Message {
	return mb.message
}

// Returns data of the message
func (m Message) Extract() any {
	return m.Payload
}

func (m Message) ExtractHeaders() MessageHeaders {
	return m.Headers
}
