package goeventbus

type Processor interface {
	forward(message Message) bool
}

type defaultProcessor struct {
	predicate func(message Message) bool
}

func newProcessor() Processor {
	return defaultProcessor{func(message Message) bool { return true }}
}

func newProcessorWithPredicate(predicate func(message Message) bool) Processor {
	return defaultProcessor{predicate: predicate}
}

func (p defaultProcessor) forward(message Message) bool {
	return p.predicate(message)
}
