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
	closed   bool
	Type     HandlerType
}

func (h *Handler) Close() {
	close(h.Ch)
}

func (h *Handler) Handle(once bool, wg *sync.WaitGroup) {
	for {
		data, ok := <-h.Ch
		if !ok {
			wg.Done()
			return
		}

		h.Consumer(h.Context.SetData(data))

		if once {
			h.closed = true
			h.Close()
			wg.Done()
		}
	}
}

func (h *Handler) NewHandler() *Handler {
	return &Handler{closed: false}
}
