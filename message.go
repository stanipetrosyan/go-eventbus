package goeventbus

type Message struct {
	payload interface{}
	headers MessageHeaders
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
	mb.message.payload = payload
	return mb
}

func (mb *defaultMessageBuilder) SetHeaders(headers MessageHeaders) MessageBuilder {
	mb.message.headers = headers
	return mb
}

func (mb *defaultMessageBuilder) Build() Message {
	return mb.message
}

// Returns data of the message
func (m Message) Extract() any {
	return m.payload
}

// Returns headers of the message
func (m Message) ExtractHeaders() MessageHeaders {
	return m.headers
}
