package main

import (
	"fmt"
	"math/rand"
	"time"

	goeventbus "github.com/StaniPetrosyan/go-eventbus"
)

var eventbus = goeventbus.NewEventBus()

func main() {

	eventbus.Subscribe("topic1")
	eventbus.Subscribe("topic2")
	eventbus.Subscribe("topic3")

	go publishTo("topic1", "Hi topic 1")
	go publishTo("topic2", "Welcome to topic 2")
	go publishTo("topic3", "Welcome to topic 3")

	eventbus.On("topic1", func(data goeventbus.DataEvent) {
		printDataEvent(data)
	})

	eventbus.On("topic2", func(data goeventbus.DataEvent) {
		printDataEvent(data)
	})

	eventbus.On("topic3", func(data goeventbus.DataEvent) {
		printDataEvent(data)
	})

}

func publishTo(address string, data string) {
	for {
		eventbus.Publish(address, data)
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	}
}

func printDataEvent(data goeventbus.DataEvent) {
	fmt.Printf("Address: %s; DataEvent %v\n", data.Address, data.Data)
}
