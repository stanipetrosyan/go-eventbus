package goeventbus

import "sync"

type Handler interface {
	Handle(wg *sync.WaitGroup)
	Chain() chan Message
	Closed() bool
}

type ConsumerHandler struct {
	ch       chan Message
	callback func(context ConsumerContext)
	context  ConsumerContext
	closed   bool
}

func (h ConsumerHandler) Handle(wg *sync.WaitGroup) {
	for {
		data, ok := <-h.ch
		if !ok {
			h.closed = true
			wg.Done()
			return
		}

		h.callback(h.Context().SetData(data))
	}
}

func (h ConsumerHandler) Chain() chan Message {
	return h.ch
}

func (h ConsumerHandler) Closed() bool {
	return h.closed
}

func (h ConsumerHandler) Context() ConsumerContext {
	return h.context
}

func NewConsumer(address string, callback func(context ConsumerContext)) ConsumerHandler {
	ch := make(chan Message)

	return ConsumerHandler{
		ch:       ch,
		callback: callback,
		closed:   false,
		context:  NewConsumerContext(ch),
	}
}

type InterceptorHandler struct {
	ch       chan Message
	callback func(context InterceptorContext)
	context  InterceptorContext
	closed   bool
}

func (h InterceptorHandler) Handle(wg *sync.WaitGroup) {
	for {
		data, ok := <-h.ch
		if !ok {
			h.closed = true
			wg.Done()
			return
		}

		h.callback(h.Context().SetData(data))
	}
}

func (h InterceptorHandler) Chain() chan Message {
	return h.ch
}

func (h InterceptorHandler) Closed() bool {
	return h.closed
}

func (h InterceptorHandler) Context() InterceptorContext {
	return h.context
}

func NewInterceptor(address string, callback func(context InterceptorContext), context InterceptorContext) InterceptorHandler {
	ch := make(chan Message)

	return InterceptorHandler{
		ch:       ch,
		callback: callback,
		closed:   false,
		context:  context,
	}
}
