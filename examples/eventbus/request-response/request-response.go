package main

import (
	"fmt"
	"sync"
	"time"

	goeventbus "github.com/stanipetrosyan/go-eventbus"
)

var eventbus = goeventbus.NewEventBus()

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	eventbus.Channel("topic1").Subscriber().Listen(func(context goeventbus.Context) {
		printMessage(context.Result())
		message := goeventbus.NewMessageBuilder().SetPayload("Hello from subscriber").Build()
		context.Reply(message)
	})

	go requestTo("topic1", "Hi topic 1")

	wg.Wait()
}

func requestTo(address, data string) {
	options := goeventbus.NewMessageHeadersBuilder().SetHeader("header", "value").Build()
	message := goeventbus.NewMessageBuilder().SetPayload(data).SetHeaders(options).Build()
	publisher := eventbus.Channel(address).Publisher()

	for {
		publisher.Request(message, func(context goeventbus.Context) {
			printMessage(context.Result())
		})
		time.Sleep(time.Second)
	}
}

func printMessage(data goeventbus.Message) {
	fmt.Printf("Message %s, Headers %s\n", data.Extract(), data.ExtractHeaders())
}
