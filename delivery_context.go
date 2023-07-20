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
	chs     []chan Message
}

func (d *DefaultDeliveryContext) Result() Message {
	return d.message
}

func (d *DefaultDeliveryContext) Reply(data any) {
	for _, item := range d.chs {
		item <- Message{Data: data}
	}
}

func (d *DefaultDeliveryContext) Handle(consume func(message Message)) {

	for _, ch := range d.chs {
		go func(ch chan Message) {
			for data := range ch {
				consume(data)
			}
		}(ch)
	}

}

func (d *DefaultDeliveryContext) SetData(msg Message) DeliveryContext {
	d.message = msg
	return d
}

func (d *DefaultDeliveryContext) Next() {
	for _, item := range d.chs {
		item <- d.message
	}
}

func NewDeliveryContext(chs []chan Message) DeliveryContext {
	return &DefaultDeliveryContext{chs: chs}
}
