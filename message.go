package goeventbus

type Message struct {
	Data    interface{}
	Options MessageOptions
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
