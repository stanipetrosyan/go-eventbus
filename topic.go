package goeventbus

type Topic struct {
	Address      string
	Consumers    []ConsumerHandler
	Interceptors []InterceptorHandler
}

func (t *Topic) AddConsumer(callback func(context ConsumerContext)) ConsumerHandler {
	consumer := NewConsumer(t.Address, callback)
	t.Consumers = append(t.Consumers, consumer)

	return consumer
}

func (t *Topic) AddInterceptor(callback func(context InterceptorContext)) InterceptorHandler {
	context := NewInterceptorContext(t)
	interceptor := NewInterceptor(t.Address, callback, context)
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
