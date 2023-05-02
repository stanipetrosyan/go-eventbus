[![Go Report Card](https://goreportcard.com/badge/github.com/StaniPetrosyan/go-eventbus)](https://goreportcard.com/report/github.com/StaniPetrosyan/go-eventbus)
![workflow](https://github.com/StaniPetrosyan/go-eventbus/actions/workflows/test.yml/badge.svg)

# EventBus for Golang

## Description

This is a simple implementation of an event bus in golang. Actually support:
* publish/subscribe messaging.
* request/reply messaging

## Get Started

To start use eventbus in your project, you can run the following command. 

```
go get github.com/StaniPetrosyan/go-eventbus
```

And import 
``` go
import (
	goeventbus "github.com/StaniPetrosyan/go-eventbus"
)

```

## Publish/Subscribe

```go

var eventbus = goeventbus.NewEventBus()

address := "topic"

eventbus.Subscribe(address, func(dc goeventbus.DeliveryContext) {
	fmt.Printf("Message %s\n", dc.Result().Data)
})

for {
	eventbus.Publish(address, "Hi Topic", MessageOptions{})
	time.Sleep(time.Second)
}
```

If you want handle once: 
```go
var eventbus = goeventbus.NewEventBus()

address := "topic"

eventbus.SubscribeOnce(address, func(dc goeventbus.DeliveryContext) {
	fmt.Printf("This Message %s\n will be printed once time", dc.Result().Data)
})

eventbus.Publish(address, "Hi Topic", MessageOptions{})
```

## Request/Reply messaging

```go

var eventbus = goeventbus.NewEventBus()

address := "topic"

eventbus.Subscribe(address, func(dc goeventbus.DeliveryContext) {
	fmt.Printf("Message %s\n", dc.Result().Data)
	dc.Reply("Hi from topic")
})
	
eventbus.Request(address, "Hi Topic", func(dc goeventbus.DeliveryContext) {
	dc.Handle(func(message Message) {
			fmt.Printf("Message %s\n", message.Data)
	})
})
```

## Options

When publish a message, you can add message options like the following:

```go

// define new message option object
options := NewMessageOptions()

//add new header
options.AddHeader("key", "value")

eventBus.Publish("address", "Hi There", options)
```

## In Bound Interceptor 

```go

eventbus.AddInBoundInterceptor("topic1", func(context goeventbus.DeliveryContext) {
	context.Next()
})
```





