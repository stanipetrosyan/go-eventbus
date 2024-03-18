package goeventbus

type Processor interface {
	forward(message Message) bool
}

type defaultProcessor struct {
	predicate func(message Message) bool
}

func NewProcessor() Processor {
	return defaultProcessor{func(message Message) bool { return true }}
}

func NewProcessorWithPredicate(predicate func(message Message) bool) Processor {
	return defaultProcessor{predicate: predicate}
}

func (p defaultProcessor) forward(message Message) bool {
	return p.predicate(message)
}
