[![Go Report Card](https://goreportcard.com/badge/github.com/StaniPetrosyan/go-eventbus)](https://goreportcard.com/report/github.com/StaniPetrosyan/go-eventbus)
![workflow](https://github.com/StaniPetrosyan/go-eventbus/actions/workflows/test.yml/badge.svg)

# EventBus for Golang

## Description

This is a simple implementation of an event bus in golang. Actually support just publish/subscribe messaging.

## Simple Usage

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

Let's see a simple example 

```go

var eventbus = goeventbus.NewEventBus()

address := "topic"

eventbus.Subscribe(address, func(data goeventbus.Message) {
	fmt.Printf("Message %s\n", data.Data)
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

eventbus.SubscribeOnce(address, func(data goeventbus.Message) {
	fmt.Printf("This Message %s\n will be printed once time", data.Data)
})

eventbus.Publish(address, "Hi Topic", MessageOptions{})
```

### Options

When publish a message, you can add message options like the following:

```go

// define new message option object
options := NewMessageOptions()

//add new header
options.AddHeader("key", "value")

eventBus.Publish("address", "Hi There", options)
```



