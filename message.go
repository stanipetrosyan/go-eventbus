package goeventbus

type Message struct {
	Data           interface{}
	MessageOptions MessageOptions
}

func CreateMessage() Message {
	return Message{}
}

func (m Message) SetBody(data any) Message {
	m.Data = data
	return m
}

func (m Message) SetOptions(options MessageOptions) Message {
	m.MessageOptions = options
	return m
}

func (m Message) Extract() any {
	return m.Data
}

func (m Message) Options() MessageOptions {
	return m.MessageOptions
}
