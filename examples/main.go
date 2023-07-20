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

	eventbus.Subscribe("topic1", func(dc goeventbus.DeliveryContext) {
		printMessage(dc.Result())
	})

	eventbus.AddInBoundInterceptor("topic1", func(context goeventbus.DeliveryContext) {
		context.Next()
	})

	eventbus.Subscribe("topic2", func(dc goeventbus.DeliveryContext) {
		printMessage(dc.Result())
	})

	eventbus.Subscribe("topic3", func(dc goeventbus.DeliveryContext) {
		printMessage(dc.Result())
	})

	eventbus.SubscribeOnce("topic4", func(dc goeventbus.DeliveryContext) {
		printMessage(dc.Result())
	})

	go publishTo("topic1", "Hi topic 1")
	go publishTo("topic2", "Hi topic 2")
	go publishTo("topic3", "Hi topic 3")
	go publishTo("topic4", "Hi topic 4")

	wg.Wait()
}

func publishTo(address string, data string) {
	options := goeventbus.NewMessageOptions().AddHeader("header", "value")
	message := goeventbus.CreateMessage().SetBody(data).SetOptions(options)
	for {
		eventbus.Publish(address, message)
		time.Sleep(time.Second)
	}
}

func printMessage(data goeventbus.Message) {
	fmt.Printf("Message %s, Headers %s\n", data.Data, data.Options)
}
