package main

import (
	goeventbus "github.com/stanipetrosyan/go-eventbus"
)

var eventbus = goeventbus.NewEventBus()

func main() {
	network := goeventbus.NewNetworkBus(eventbus, "localhost", "/bus")
	network.Server("localhost", "/bus").Listen()
}
