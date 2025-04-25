[![Go Report Card](https://goreportcard.com/badge/github.com/stanipetrosyan/go-eventbus)](https://goreportcard.com/report/github.com/stanipetrosyan/go-eventbus)
[![codecov](https://codecov.io/gh/stanipetrosyan/go-eventbus/graph/badge.svg?token=YAGXYA64E6)](https://codecov.io/gh/stanipetrosyan/go-eventbus)
[![Go Reference](https://pkg.go.dev/badge/github.com/stanipetrosyan/go-eventbus.svg)](https://pkg.go.dev/github.com/stanipetrosyan/go-eventbus)
![workflow](https://github.com/StaniPetrosyan/go-eventbus/actions/workflows/test.yml/badge.svg)

# EventBus for Golang

## Description

This is a simple implementation of an event bus in golang. Actually support follwing pattern:

- [x] publish/subscribe
- [x] request/response

## Get Started

To start use eventbus in your project, you can run the following command.

```
go get github.com/stanipetrosyan/go-eventbus
```

And import
``` go
import (
	goeventbus "github.com/stanipetrosyan/go-eventbus"
)

```

## Publish/Subscribe

Simple example of publish/subscribe pattern.

```go

eventbus = goeventbus.NewEventBus()

address := "topic"
options := goeventbus.NewMessageHeadersBuilder().SetHeader("header", "value").Build()
message := goeventbus.NewMessageBuilder().SetPayload("Hi Topic").SetHeaders(options).Build()

eventbus.Channel(address).Subscriber().Listen(func(dc goeventbus.Context) {
	fmt.Printf("Message %s\n", dc.Result().Data)
})

eventbus.Channel(address).Publisher().Publish(message)
```

## Request/Response

Simple example of request/response pattern.

```go

eventbus = goeventbus.NewEventBus()

address := "topic"
message := goeventbus.NewMessageBuilder().SetPayload("Hi Topic").Build()

eventbus.Channel(address).Subscriber().Listen(func(context goeventbus.Context) {
	fmt.Printf("Message %s\n", context.Result().Extract())
	context.Reply("Hello from subscriber")
})

eventbus.Channel(address).Publisher().Request(message, func(context goeventbus.Context) {
	fmt.Printf("Message %s\n", context.Result().Extract())
})
```

## Message

For publishing, you need to create a Message object using this method.

```go
message := goeventbus.NewMessageBuilder().SetPayload("Hi Topic").SetHeaders(options).Build()
```
Each message can have some options:

```go

options := goeventbus.NewMessageHeadersBuilder().SetHeader("header", "value").Build()
message := goeventbus.NewMessageBuilder().setHeaders(options).Build()

eventBus.Channel("address").Publisher().Publish(message)
```

## Processor

A processor works like a middleware, in fact forwards messages only if context Next method is called. The entity works as Subscriber: Listen method accept a callback with context.

The processor intercept message from publisher to subscribers on a specific channel.

```go
eventbus.Channel("topic1").Processor().Listen(func(context goeventbus.Context) {
	if context.Result().ExtractHeaders().Contains("header") {
		context.Next()
	}
})
```

## Network Bus

A Network bus create a tcp connection between different services.

NetworkBus is a wrapper of local eventbus.

A simple server/client example is in `examples/networkbus` directory.
