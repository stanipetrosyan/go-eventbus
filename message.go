package goeventbus

type Message struct {
	Data           interface{}
	MessageOptions MessageOptions
}

// Returns an empty message
func CreateMessage() Message {
	return Message{}
}

// Set the body of the message. Paramater can be any type
func (m Message) SetBody(data any) Message {
	m.Data = data
	return m
}

// Set the options of the message. For create a new options see: NewMessageOptions()
func (m Message) SetOptions(options MessageOptions) Message {
	m.MessageOptions = options
	return m
}

// Returns data of the message
func (m Message) Extract() any {
	return m.Data
}

// Returns options of the message
func (m Message) Options() MessageOptions {
	return m.MessageOptions
}
