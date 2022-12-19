package main

import (
	"fmt"
	"sync"
	"time"

	goeventbus "github.com/StaniPetrosyan/go-eventbus"
)

var eventbus = goeventbus.NewEventBus()

func main() {

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		time.Sleep(10 * time.Second)
		eventbus.Unsubscribe("topic1")
		wg.Done()
	}()

	eventbus.Subscribe("topic1")
	eventbus.Subscribe("topic2")
	eventbus.Subscribe("topic3")

	go publishTo("topic1", "Hi topic 1")
	go publishTo("topic2", "Hi topic 2")

	eventbus.On("topic2", func(data goeventbus.DataEvent) {
		printDataEvent(data)
	})

	eventbus.On("topic1", func(data goeventbus.DataEvent) {
		printDataEvent(data)
	})

	eventbus.On("topic3", func(data goeventbus.DataEvent) {
		printDataEvent(data)
	})

	wg.Wait()

}

func publishTo(address string, data string) {
	for {
		println("invio")
		eventbus.Publish(address, data)
		time.Sleep(time.Second)
	}
}

func printDataEvent(data goeventbus.DataEvent) {
	fmt.Printf("Address: %s; DataEvent %v\n", data.Address, data.Data)
}
