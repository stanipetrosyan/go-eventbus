package goeventbus

type Publisher interface {
	Publish()
}

type defaultPublisher struct {
	ch chan string
}

func NewPublisher(ch chan string) Publisher {
	return defaultPublisher{ch: ch}
}

func (p defaultPublisher) Publish() {
	p.ch <- "Hello World"
}
