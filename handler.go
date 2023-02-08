package goeventbus

import "sync"

type Handler struct {
	Ch       chan Message
	Consume  func(data Message)
	Address  string
	flagOnce bool
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

		h.Consume(data)

		if once {
			h.flagOnce = true
			h.Close()
			wg.Done()
		}
	}
}
