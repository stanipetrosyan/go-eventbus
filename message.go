package goeventbus

import "encoding/json"

type Message struct {
	Data    interface{}
	Options MessageOptions
}

func (m Message) ToJson() ([]byte, error) {
	return json.Marshal(m.Data)
}

func CreateMessage() Message {
	return Message{}
}

func (m Message) SetBody(data interface{}) Message {
	m.Data = data
	return m
}

func (m Message) SetOptions(options MessageOptions) Message {
	m.Options = options
	return m
}
