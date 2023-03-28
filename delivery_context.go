package goeventbus

type DeliveryContext interface {
	//Next()
	//Reply()
	Result() Message
}

type DefaultDeliveryContext struct {
	message Message
}

func (d *DefaultDeliveryContext) Result() Message {
	return d.message
}

func NewDeliveryContext(message Message) DeliveryContext {
	return &DefaultDeliveryContext{message: message}
}
