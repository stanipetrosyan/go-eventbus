package goeventbus

type MessageHeaders struct {
	headers map[string]string
}

type MessageHeadersBuilder interface {
	SetHeader(key string, value string) MessageHeadersBuilder
	Build() MessageHeaders
}

type defaultMessageHeadersBuilder struct {
	messageHeaders MessageHeaders
}

func NewMessageHeadersBuilder() MessageHeadersBuilder {
	return &defaultMessageHeadersBuilder{messageHeaders: MessageHeaders{headers: map[string]string{}}}
}

func (hb *defaultMessageHeadersBuilder) SetHeader(key string, value string) MessageHeadersBuilder {
	hb.messageHeaders.headers[key] = value
	return hb
}

func (hb *defaultMessageHeadersBuilder) Build() MessageHeaders {
	return hb.messageHeaders
}

func (h MessageHeaders) Get(key string) string {
	return h.headers[key]
}

func (h MessageHeaders) Contains(key string) bool {
	_, exist := h.headers[key]
	return exist
}
