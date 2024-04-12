package main

import (
	goeventbus "github.com/stanipetrosyan/go-eventbus"
)

var eventbus = goeventbus.NewEventBus()

func main() {
	network := goeventbus.NewNetworkBus(eventbus, "localhost", "/bus")

	network.Client("localhost", "/bus").Connect()
}
