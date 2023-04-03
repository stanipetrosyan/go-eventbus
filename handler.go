package goeventbus

import "sync"

type Handler struct {
	Ch       chan Message
	Consumer HandlerFunc
	Context  DeliveryContext
	Address  string
	closed   bool
}

func (h *Handler) Close() {
	close(h.Ch)
}

func (h *Handler) Handle(once bool, wg *sync.WaitGroup) {
	context := NewDeliveryContext(Message{}, h.Ch)

	for {
		data, ok := <-h.Ch
		if !ok {
			wg.Done()
			return
		}

		h.Consumer(context.SetData(data))

		if once {
			h.closed = true
			h.Close()
			wg.Done()
		}
	}
}
