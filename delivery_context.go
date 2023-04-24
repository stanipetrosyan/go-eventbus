package goeventbus

type DeliveryContext interface {
	Reply(data any)
	Handle(func(message Message))
	Result() Message
	SetData(msg Message) DeliveryContext
	Next()
}

type DefaultDeliveryContext struct {
	message Message
	ch      chan Message
}

func (d *DefaultDeliveryContext) Result() Message {
	return d.message
}

func (d *DefaultDeliveryContext) Reply(data any) {
	d.ch <- Message{Data: data}
}

func (d *DefaultDeliveryContext) Handle(consume func(message Message)) {
	for data := range d.ch {
		consume(data)
	}
}

func (d *DefaultDeliveryContext) SetData(msg Message) DeliveryContext {
	d.message = msg
	return d
}

func (d *DefaultDeliveryContext) Next() {
	d.ch <- d.message
}

func NewDeliveryContext(message Message, ch chan Message) DeliveryContext {
	return &DefaultDeliveryContext{message: message, ch: ch}
}
