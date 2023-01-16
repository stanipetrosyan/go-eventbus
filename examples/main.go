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
		time.Sleep(4 * time.Second)
		eventbus.Unsubscribe("topic1")
		println("Unsubscribed topic1 handler")
		//wg.Done()
	}()

	eventbus.Subscribe("topic1")
	eventbus.Subscribe("topic2")
	eventbus.Subscribe("topic3")

	eventbus.On("topic1", func(data goeventbus.Message) {
		printMessage(data)
	})

	eventbus.On("topic2", func(data goeventbus.Message) {
		printMessage(data)
	})

	eventbus.On("topic3", func(data goeventbus.Message) {
		printMessage(data)
	})

	go publishTo("topic1", "Hi topic 1")
	go publishTo("topic2", "Hi topic 2")
	go publishTo("topic3", "Hi topic 3")

	wg.Wait()
}

func publishTo(address string, data string) {
	options := goeventbus.NewMessageOptions().AddHeader("header", "value")
	for {
		eventbus.Publish(address, data, options)
		time.Sleep(time.Second)
	}
}

func printMessage(data goeventbus.Message) {
	fmt.Printf("Message %s, Headers %s\n", data.Data, data.Headers["header"])
}
