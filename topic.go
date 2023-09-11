package goeventbus

type Topic struct {
	Address      string
	Consumers    []Handler
	Interceptors []Handler
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

func (t *Topic) GetHandlers() []Handler {
	if len(t.Interceptors) > 0 {
		return t.Interceptors
	}
	return t.Consumers
}

func (t *Topic) GetChannels() []chan Message {
	chs := []chan Message{}

	for _, item := range t.Consumers {
		chs = append(chs, item.Chain())

	}

	return chs
}
func NewTopic(address string) *Topic {
	return &Topic{Address: address, Consumers: []Handler{}, Interceptors: []Handler{}}
}
