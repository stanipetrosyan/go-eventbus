package goeventbus

import "errors"

type Processor interface {
	Listen(consumer func(context Context))
}

type defaultProcessor struct {
	listenChannel <-chan Message
	sendChannel   chan<- packet
}

func newProcessor(ch <-chan Message, channel chan packet) Processor {
	return defaultProcessor{listenChannel: ch, sendChannel: channel}
}

func (p defaultProcessor) Listen(consumer func(context Context)) {
	go func() {
		for {
			message, ok := <-p.listenChannel
			if !ok {
				newContextWithError(errors.New("channel closed"))
				return
			}
			consumer(newContextWithMessageAndChannel(message, p.sendChannel))
		}
	}()
}
