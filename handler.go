package goeventbus

type Handler struct {
	Ch      chan Message
	Consume func(data Message)
	Address string
}

func (h *Handler) Close() {
	close(h.Ch)
}
