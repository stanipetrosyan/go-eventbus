package main

import (
	"fmt"
	"math/rand"
	"time"

	goeventbus "github.com/StaniPetrosyan/go-eventbus"
)

var eventbus = goeventbus.NewEventBus()

func main() {

	ch1 := make(chan goeventbus.DataEvent)
	ch2 := make(chan goeventbus.DataEvent)
	ch3 := make(chan goeventbus.DataEvent)

	eventbus.Subscribe("topic1", ch1)
	eventbus.Subscribe("topic2", ch2)
	eventbus.Subscribe("topic2", ch3)

	go publishTo("topic1", "Hi topic 1")
	go publishTo("topic2", "Welcome to topic 2")

	for {
		select {
		case d := <-ch1:
			go printDataEvent("ch1", d)
		case d := <-ch2:
			go printDataEvent("ch2", d)
		case d := <-ch3:
			go printDataEvent("ch3", d)
		}
	}
}

func publishTo(address string, data string) {
	for {
		eventbus.Publish(address, data)
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	}
}

func printDataEvent(ch string, data goeventbus.DataEvent) {
	fmt.Printf("Channel: %s; Address: %s; DataEvent %v\n", ch, data.Address, data.Data)
}
