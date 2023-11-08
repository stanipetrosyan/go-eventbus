package goeventbus

type Publisher interface {
	Publish(message Message)
}

type defaultPublisher struct {
	ch chan Message
}

func NewPublisher(ch chan Message) Publisher {
	return defaultPublisher{ch: ch}
}

func (p defaultPublisher) Publish(message Message) {
	p.ch <- message
}
