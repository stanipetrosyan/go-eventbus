package goeventbus

import "sync"

type Interceptor struct {
	Ch       chan Message
	Consumer HandlerFunc
}

type Handler struct {
	Ch           chan Message
	Consumer     HandlerFunc
	Context      DeliveryContext
	Address      string
	closed       bool
	Interceptors []Interceptor
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

func (h *Handler) AddInterceptor(interceptor Interceptor) {
	h.Interceptors = append(h.Interceptors, interceptor)
}
