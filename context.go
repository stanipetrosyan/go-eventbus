package goeventbus

type Context interface {
	Result() Message
}

type ConsumerContext struct {
	Ch      chan Message
	message Message
}

func (d *ConsumerContext) Reply(data any) {
	d.Ch <- Message{Data: data}
}

func (d *ConsumerContext) Result() Message {
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
	chs     []chan Message
	message Message
}

func (d *InterceptorContext) Next() {
	for _, item := range d.chs {
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

func NewInterceptorContext(ch []chan Message) InterceptorContext {
	return InterceptorContext{chs: ch}
}
