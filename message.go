package goeventbus

type Message struct {
	Data    interface{}
	Headers []string
}
