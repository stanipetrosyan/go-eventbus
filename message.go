package goeventbus

import "encoding/json"

type Message struct {
	Data    interface{}
	Headers map[string]string
}

func (m Message) ToJson() ([]byte, error) {
	return json.Marshal(m.Data)
}
