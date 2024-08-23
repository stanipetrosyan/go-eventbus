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
		message := goeventbus.CreateMessage().SetBody("Hello from subscriber")
		context.Reply(message)
	})

	go requestTo("topic1", "Hi topic 1")

	wg.Wait()
}

func requestTo(address, data string) {
	options := goeventbus.NewMessageOptions().SetHeaders(goeventbus.NewHeaders().Add("header", "value"))
	message := goeventbus.CreateMessage().SetBody(data).SetOptions(options)
	publisher := eventbus.Channel(address).Publisher()

	for {
		publisher.Request(message, func(context goeventbus.Context) {
			printMessage(context.Result())
		})
		time.Sleep(time.Second)
	}
}

func printMessage(data goeventbus.Message) {
	fmt.Printf("Message %s, Headers %s\n", data.Extract(), data.Options().Headers())
}
