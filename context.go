package goeventbus

type Context interface {
	Result() Message
}

type defaultContext struct {
	message Message
}

func (c defaultContext) Result() Message {
	return c.message
}

func NewConsumerContextWithMessage(message Message) Context {
	return defaultContext{message: message}
}

type ConsumerContext struct {
	Ch      chan Message
	message Message
}

func (d *ConsumerContext) Reply(data any) {
	d.Ch <- Message{Data: data}
}

func (d ConsumerContext) Result() Message {
	return d.message
}

func (d *ConsumerContext) Handle(consume func(message Message)) {
	go func(ch chan Message) {
		for data := range ch {
			consume(data)
		}
	}(d.Ch)
}

func (d ConsumerContext) SetData(msg Message) ConsumerContext {
	d.message = msg
	return d
}

func NewConsumerContext(ch chan Message) ConsumerContext {
	return ConsumerContext{Ch: ch}
}

type InterceptorContext struct {
	message Message
	topic   *Topic
}

func (d *InterceptorContext) Next() {
	for _, item := range d.topic.GetChannels() {
		item <- d.message
	}
}

func (d *InterceptorContext) Result() Message {
	return d.message
}

func (d InterceptorContext) SetData(msg Message) InterceptorContext {
	d.message = msg
	return d
}

func NewInterceptorContext(topic *Topic) InterceptorContext {
	return InterceptorContext{topic: topic}
}
