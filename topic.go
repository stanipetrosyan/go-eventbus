package goeventbus

type Topic struct {
	Address      string
	Consumers    []*Handler
	Interceptors []*Handler
}

func (t *Topic) AddHandler(handler *Handler) {
	switch handler.Type {
	case Consumer:
		t.addConsumer(handler)
	case Interceptor:
		t.addInterceptor(handler)
	}
}

func (t *Topic) addConsumer(handler *Handler) *Handler {
	t.Consumers = append(t.Consumers, handler)
	return handler
}

func (t *Topic) addInterceptor(interceptor *Handler) *Handler {
	t.Interceptors = append(t.Interceptors, interceptor)
	return interceptor
}

func (t *Topic) GetHandlers() []*Handler {
	if len(t.Interceptors) > 0 {
		return t.Interceptors
	}

	return t.Consumers
}

func (t *Topic) GetChannels() []chan Message {
	chs := []chan Message{}

	for _, item := range t.Consumers {
		chs = append(chs, item.Ch)
	}

	return chs
}

func NewTopic(address string) *Topic {
	return &Topic{Address: address, Consumers: []*Handler{}, Interceptors: []*Handler{}}
}
