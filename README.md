[![Go Report Card](https://goreportcard.com/badge/github.com/stanipetrosyan/go-eventbus)](https://goreportcard.com/report/github.com/stanipetrosyan/go-eventbus)
[![codecov](https://codecov.io/gh/stanipetrosyan/go-eventbus/graph/badge.svg?token=YAGXYA64E6)](https://codecov.io/gh/stanipetrosyan/go-eventbus)
[![Go Reference](https://pkg.go.dev/badge/github.com/stanipetrosyan/go-eventbus.svg)](https://pkg.go.dev/github.com/stanipetrosyan/go-eventbus)
![workflow](https://github.com/StaniPetrosyan/go-eventbus/actions/workflows/test.yml/badge.svg)

# EventBus for Golang

## Description

This is a simple implementation of an event bus in golang. Actually support:
* publish/subscribe messaging.

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

var eventbus = goeventbus.NewEventBus()

address := "topic"
options := goeventbus.NewMessageOptions().AddHeader("header", "value")
message := goeventbus.CreateMessage().SetBody("Hi Topic").SetOptions(options)

eventbus.Channel(address).Subscriber().Listen(func(dc goeventbus.Context) {
	fmt.Printf("Message %s\n", dc.Result().Data)
})

eventbus.Channel(address).Publisher().Publish(message)
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

eventBus.Channel("address").Publisher().Publish(message)
```

## Processor

A processor works like a middleware, in fact forwards messages only if the predicate is satisfied. The method accepts a function with message and return must return a boolean.

```go

eventbus.Channel("topic1").Processor(func(message goeventbus.Message) bool {
	return logic
})
```





