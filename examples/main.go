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

	eventbus.Subscribe("topic1", func(dc goeventbus.DeliveryContext) {
		printMessage(dc.Result())
	})

	eventbus.AddInBoundInterceptor("topic1", func(context goeventbus.DeliveryContext) {
		if context.Result().Data == "Hi topic 1" {
			context.Next()
		} else {
			println("Message not passed")
		}
	})

	eventbus.Subscribe("topic2", func(dc goeventbus.DeliveryContext) {
		printMessage(dc.Result())
	})

	eventbus.Subscribe("topic3", func(dc goeventbus.DeliveryContext) {
		printMessage(dc.Result())
	})

	eventbus.Subscribe("topic4", func(dc goeventbus.DeliveryContext) {
		printMessage(dc.Result())
	})

	eventbus.Subscribe("topic5", func(dc goeventbus.DeliveryContext) {
		printMessage(dc.Result())
		dc.Reply("Hi Publisher")
	})

	eventbus.SubscribeOnce("topic6", func(context goeventbus.DeliveryContext) {
		printMessage(context.Result())
	})

	go publishTo("topic1", "Hi topic 1")
	go publishTo("topic1", "Message to block")
	go publishTo("topic2", "Hi topic 2")
	go publishTo("topic3", "Hi topic 3")
	go publishTo("topic4", "Hi topic 4")
	go RequestTo("topic5", "Hi topic 5")
	go publishTo("topic6", "Hi topic 6")

	wg.Wait()
}

func RequestTo(address string, data string) {
	options := goeventbus.NewMessageOptions().AddHeader("header", "value")
	message := goeventbus.CreateMessage().SetBody(data).SetOptions(options)
	for {
		eventbus.Request(address, message, func(context goeventbus.DeliveryContext) {
			context.Handle(func(message goeventbus.Message) {
				printMessage(message)
			})
		})
		time.Sleep(time.Second)
	}
}

func publishTo(address string, data string) {
	options := goeventbus.NewMessageOptions().AddHeader("header", "value")
	message := goeventbus.CreateMessage().SetBody(data).SetOptions(options)
	for {
		eventbus.Publish(address, message)
		time.Sleep(time.Second * 2)
	}
}

func printMessage(data goeventbus.Message) {
	fmt.Printf("Message %s, Headers %s\n", data.Data, data.Options)
}
