package goeventbus

type DeliveryContext interface {
	Reply(data any)
	Result() Message
	SetData(msg Message) DeliveryContext
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

func (d *DefaultDeliveryContext) SetData(msg Message) DeliveryContext {
	d.message = msg
	return d
}

func NewDeliveryContext(message Message, ch chan Message) DeliveryContext {
	return &DefaultDeliveryContext{message: message, ch: ch}
}
