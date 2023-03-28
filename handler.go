package goeventbus

import "sync"

type Handler struct {
	Ch       chan Message
	Consumer HandlerFunc
	Address  string
	closed   bool
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

		context := NewDeliveryContext(data)

		h.Consumer(context)

		if once {
			h.closed = true
			h.Close()
			wg.Done()
		}
	}
}
