[![Go Report Card](https://goreportcard.com/badge/github.com/stanipetrosyan/go-eventbus)](https://goreportcard.com/report/github.com/stanipetrosyan/go-eventbus)
[![codecov](https://codecov.io/gh/stanipetrosyan/go-eventbus/graph/badge.svg?token=YAGXYA64E6)](https://codecov.io/gh/stanipetrosyan/go-eventbus)
[![Go Reference](https://pkg.go.dev/badge/github.com/stanipetrosyan/go-eventbus.svg)](https://pkg.go.dev/github.com/stanipetrosyan/go-eventbus)
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
options := goeventbus.NewMessageOptions().AddHeader("header", "value")
message := goeventbus.CreateMessage().SetBody("Hi Topic").SetOptions(options)

eventbus.Subscribe(address, func(dc goeventbus.ConsumerContext) {
	fmt.Printf("Message %s\n", dc.Result().Data)
})

for {
	eventbus.Publish(address, message)
	time.Sleep(time.Second)
}
```

## Request/Reply messaging

```go

var eventbus = goeventbus.NewEventBus()

address := "topic"

eventbus.Subscribe(address, func(dc goeventbus.ConsumerContext) {
	fmt.Printf("Message %s\n", dc.Result().Data)
	dc.Reply("Hi from topic")
})
	
eventbus.Request(address, "Hi Topic", func(dc goeventbus.ConsumerContext) {
	dc.Handle(func(message Message) {
			fmt.Printf("Message %s\n", message.Data)
	})
})
```

## Message

For publishing, you need to create a Message object using this method. 

```go
message := goeventbus.CreateMessage().SetBody("Hi Topic")
```
Each message can have some options:

```go

options := goeventbus.NewMessageOptions().AddHeader("header", "value")
message := goeventbus.CreateMessage()

message.SetOptions(options)

eventBus.Publish("address", message)
```

## In Bound Interceptor 

```go

eventbus.AddInBoundInterceptor("topic1", func(context goeventbus.InterceptonContext) {
	if (some logic)
		context.Next()
})
```





