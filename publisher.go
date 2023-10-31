package goeventbus

type Publisher struct {
}

func New() Publisher {
	return Publisher{}
}

func (p *Publisher) Publish() {
}
