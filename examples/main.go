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

	header := []string{"this is a header"}

	go publishTo("topic1", "Hi topic 1", header)
	go publishTo("topic2", "Hi topic 2", header)
	go publishTo("topic3", "Hi topic 3", header)

	wg.Wait()
}

func publishTo(address string, data string, headers []string) {
	for {
		eventbus.Publish(address, goeventbus.Message{Data: data, Headers: headers}, goeventbus.MessageOptions{})
		time.Sleep(time.Second)
	}
}

func printMessage(data goeventbus.Message) {
	fmt.Printf("Message %s, Headers %s\n", data.Data, data.Headers[0])
}
