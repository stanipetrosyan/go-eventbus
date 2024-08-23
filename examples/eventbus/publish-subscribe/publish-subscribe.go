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

	eventbus.Channel("topic1").Subscriber().Listen(func(dc goeventbus.Context) {
		printMessage(dc.Result())
	})

	eventbus.Channel("topic1").Processor(func(message goeventbus.Message) bool {
		return message.Options().Headers().Contains("header")
	})

	eventbus.Channel("topic2").Subscriber().Listen(func(dc goeventbus.Context) {
		printMessage(dc.Result())
	})

	eventbus.Channel("topic3").Subscriber().Listen(func(dc goeventbus.Context) {
		printMessage(dc.Result())
	})

	go publishTo("topic1", "Hi topic 1")
	go publishTo("topic2", "Hi topic 2")
	go publishTo("topic3", "Hi topic 3")

	wg.Wait()
}

func publishTo(address string, data string) {
	options := goeventbus.NewMessageOptions().SetHeaders(goeventbus.NewHeaders().Add("header", "value"))
	message := goeventbus.CreateMessage().SetBody(data).SetOptions(options)
	publisher := eventbus.Channel(address).Publisher()

	for {
		publisher.Publish(message)
		time.Sleep(time.Second)
	}
}

func printMessage(data goeventbus.Message) {
	fmt.Printf("Message %s, Headers %s\n", data.Extract(), data.Options().Headers())
}
