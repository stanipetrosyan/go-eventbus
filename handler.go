package goeventbus

import "sync"

type HandlerType int

const (
	Consumer HandlerType = iota
	Interceptor
)

type Handler struct {
	Ch       chan Message
	Callback func(context DeliveryContext)
	Context  DeliveryContext
	Address  string
	Closed   bool
	Type     HandlerType
}

func (h *Handler) Close() {
	close(h.Ch)
	h.Closed = true
}

func (h *Handler) Handle(wg *sync.WaitGroup) {
	for {
		data, ok := <-h.Ch
		if !ok {
			wg.Done()
			return
		}

		h.Callback(h.Context.SetData(data))
	}
}

func (h *Handler) NewHandler() *Handler {
	return &Handler{Closed: false}
}
