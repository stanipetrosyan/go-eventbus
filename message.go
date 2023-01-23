package goeventbus

import "encoding/json"

type Message struct {
	Data    interface{}
	Headers map[string]string
}

func (m Message) ToJson() ([]byte, error) {
	return json.Marshal(m.Data)
}

type MessageOptions struct {
	headers map[string]string
}

func (op MessageOptions) AddHeader(key string, value string) MessageOptions {
	op.headers[key] = value
	return op
}

func (op MessageOptions) SetHeader(headers map[string]string) MessageOptions {
	op.headers = headers
	return op
}

func NewMessageOptions() MessageOptions {
	return MessageOptions{
		headers: map[string]string{},
	}
}
