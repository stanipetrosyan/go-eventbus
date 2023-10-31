package goeventbus

type Channel interface {
	Register() Channel
	Publisher() Channel
	Subscriber() Channel
}
