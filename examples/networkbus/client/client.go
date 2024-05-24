package main

import (
	"fmt"

	goeventbus "github.com/stanipetrosyan/go-eventbus"
)

var eventbus = goeventbus.NewEventBus()

func main() {

	network := goeventbus.NewNetworkBus(eventbus, "localhost:9000")
	eventbus.Channel("hello").Subscriber().Listen(func(context goeventbus.Context) {
		printMessage(context.Result())
	})

	network.Client().Connect()
}

func printMessage(data goeventbus.Message) {
	fmt.Printf("Message %s, Headers %s\n", data.Extract(), data.Options().Headers())
}
