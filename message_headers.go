package goeventbus

import "encoding/json"

type MessageHeaders struct {
	headers map[string]string
}

// MarshalJSON implements json.Marshaler interface
func (h MessageHeaders) MarshalJSON() ([]byte, error) {
	if h.headers == nil {
		return json.Marshal(map[string]string{})
	}
	return json.Marshal(h.headers)
}

// UnmarshalJSON implements json.Unmarshaler interface
func (h *MessageHeaders) UnmarshalJSON(data []byte) error {
	if h.headers == nil {
		h.headers = make(map[string]string)
	}
	return json.Unmarshal(data, &h.headers)
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
	if h.headers == nil {
		return ""
	}
	return h.headers[key]
}

func (h MessageHeaders) Contains(key string) bool {
	if h.headers == nil {
		return false
	}
	_, exist := h.headers[key]
	return exist
}
