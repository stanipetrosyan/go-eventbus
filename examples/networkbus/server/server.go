package main

import (
	"time"

	goeventbus "github.com/stanipetrosyan/go-eventbus"
)

var eventbus = goeventbus.NewEventBus()

func main() {

	network := goeventbus.NewNetworkBus(eventbus, "localhost:9000", "/bus")
	server := network.Server()
	go server.Listen()

	for {
		message := goeventbus.CreateMessage().SetBody("Hello World!")
		server.Publish("hello", message)
		time.Sleep(time.Second * 2)
	}
}
