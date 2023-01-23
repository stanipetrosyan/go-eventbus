# EventBus for Golang

## Description

This is a simple implementation of an event bus in golang. Actually support just publish/subscribe messaging.

## Simple Usage

To start use eventbus in your project, you can run the following command. 

```
go get github.com/StaniPetrosyan/go-eventbus
```

Let's see a simple example 

```go

import (
	"fmt"
	"time"

	goeventbus "github.com/StaniPetrosyan/go-eventbus"
)

var eventbus = goeventbus.NewEventBus()

func main() {
	address := "topic"

	eventbus.Subscribe(address)

	eventbus.On("topic", func(data goeventbus.Message) {
		fmt.Printf("Message %s\n", data.Data)
	})

	for {
		eventbus.Publish(address, "Hi Topic", MessageOptions{})
		time.Sleep(time.Second)
	}
}

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



