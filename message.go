package goeventbus

type Message struct {
	Data    interface{}
	Headers []string
}

type MessageOptions struct {
	headers map[string]string
}

func (op MessageOptions) AddHeader(key string, value string) {
	op.headers[key] = value
}

func NewMessageOptions() MessageOptions {
	return MessageOptions{
		headers: map[string]string{},
	}
}
