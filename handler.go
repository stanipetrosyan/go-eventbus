package goeventbus

import "sync"

type Handler interface {
	Handle(wg *sync.WaitGroup)
	Chain() chan Message
	Closed() bool
}

type ConsumerHandler struct {
	Ch       chan Message
	Callback func(context ConsumerContext)
	context  ConsumerContext
	closed   bool
}

func (h ConsumerHandler) Handle(wg *sync.WaitGroup) {
	for {
		data, ok := <-h.Ch
		if !ok {
			h.closed = true
			wg.Done()
			return
		}

		h.Callback(h.Context().SetData(data))
	}
}

func (h ConsumerHandler) Chain() chan Message {
	return h.Ch
}

func (h ConsumerHandler) Closed() bool {
	return h.closed
}

func (h ConsumerHandler) Context() ConsumerContext {
	return h.context
}

func (h ConsumerHandler) SetContext(context ConsumerContext) ConsumerHandler {
	h.context = context
	return h
}

func NewConsumer(address string, callback func(context ConsumerContext)) ConsumerHandler {
	ch := make(chan Message)

	return ConsumerHandler{
		Ch:       ch,
		Callback: callback,
		closed:   false,
	}
}

type InterceptorHandler struct {
	Ch       chan Message
	Callback func(context InterceptorContext)
	context  InterceptorContext
	closed   bool
}

func (h InterceptorHandler) Handle(wg *sync.WaitGroup) {
	for {
		data, ok := <-h.Ch
		if !ok {
			h.closed = true
			wg.Done()
			return
		}

		h.Callback(h.Context().SetData(data))
	}
}

func (h InterceptorHandler) Chain() chan Message {
	return h.Ch
}

func (h InterceptorHandler) Closed() bool {
	return h.closed
}

func (h InterceptorHandler) Context() InterceptorContext {
	return h.context
}

func (h InterceptorHandler) SetContext(context InterceptorContext) InterceptorHandler {
	h.context = context
	return h
}

func NewInterceptor(address string, callback func(context InterceptorContext)) InterceptorHandler {
	ch := make(chan Message)

	return InterceptorHandler{
		Ch:       ch,
		Callback: callback,
		closed:   false,
	}
}
