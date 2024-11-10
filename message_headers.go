package goeventbus

type Headers struct {
	headers map[string]string
}

type MessageHeaders struct {
	headers Headers
}

type MessageHeadersBuilder interface {
	SetHeader(key string, value string) MessageHeadersBuilder
	Build() MessageHeaders
}

type defaultMessageHeadersBuilder struct {
	messageHeaders MessageHeaders
}

func NewMessageHeadersBuilder() MessageHeadersBuilder {
	return &defaultMessageHeadersBuilder{messageHeaders: MessageHeaders{Headers{headers: map[string]string{}}}}
}

func (hb *defaultMessageHeadersBuilder) SetHeader(key string, value string) MessageHeadersBuilder {
	hb.messageHeaders.headers.headers[key] = value
	return hb
}

func (hb *defaultMessageHeadersBuilder) Build() MessageHeaders {
	return hb.messageHeaders
}

func NewMessageHeaders() MessageHeaders {
	return MessageHeaders{
		headers: Headers{},
	}
}

func (h MessageHeaders) Get(key string) string {
	return h.headers.headers[key]
}

func (h MessageHeaders) Contains(key string) bool {
	_, exist := h.headers.headers[key]
	return exist
}
