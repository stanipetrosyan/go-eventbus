package goeventbus

type Topic struct {
	Address      string
	Consumers    []ConsumerHandler
	Interceptors []InterceptorHandler
}

func (t *Topic) AddConsumer(handler ConsumerHandler) ConsumerHandler {
	t.Consumers = append(t.Consumers, handler)
	return handler
}

func (t *Topic) AddInterceptor(interceptor InterceptorHandler) InterceptorHandler {
	t.Interceptors = append(t.Interceptors, interceptor)
	return interceptor
}

func (t *Topic) GetConsumers() []ConsumerHandler {
	return t.Consumers
}

func (t *Topic) GetInterceptors() []InterceptorHandler {
	return t.Interceptors
}

func (t *Topic) ExistInterceptor() bool {
	return len(t.Interceptors) > 0
}

func (t *Topic) GetChannels() []chan Message {
	chs := []chan Message{}

	for _, item := range t.Consumers {
		chs = append(chs, item.Chain())
	}

	return chs
}

func NewTopic(address string) *Topic {
	return &Topic{Address: address, Consumers: []ConsumerHandler{}, Interceptors: []InterceptorHandler{}}
}
