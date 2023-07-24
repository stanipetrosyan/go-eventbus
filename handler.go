package goeventbus

import "sync"

type HandlerType int

const (
	Consumer HandlerType = iota
	Interceptor
)

type Handler struct {
	Ch       chan Message
	Consumer func(context DeliveryContext)
	Context  DeliveryContext
	Address  string
	Closed   bool
	Once     bool
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

		h.Consumer(h.Context.SetData(data))

		if h.Once {
			h.Closed = true
			close(h.Ch)
			wg.Done()
		}
	}
}

func (h *Handler) NewHandler() *Handler {
	return &Handler{Closed: false}
}
