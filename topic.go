package goeventbus

type Topic struct {
	Address     string
	Handlers    []*Handler
	Interceptor []*Handler
}

func (t *Topic) AddInterceptor(interceptor Handler) *Handler {
	t.Interceptor = append(t.Interceptor, &interceptor)
	return &interceptor
}

func (t *Topic) GetChannels() []chan Message {
	chs := []chan Message{}

	for _, item := range t.Handlers {
		chs = append(chs, item.Ch)
	}

	return chs
}
