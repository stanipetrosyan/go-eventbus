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

func (t *Topic) GetHandlers() []*Handler {
	if len(t.Interceptor) > 0 {
		return t.Interceptor
	}

	return t.Handlers
}

func (t *Topic) GetChannels() []chan Message {
	chs := []chan Message{}

	for _, item := range t.Handlers {
		chs = append(chs, item.Ch)
	}

	return chs
}
