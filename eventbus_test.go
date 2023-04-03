package goeventbus

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

var wg sync.WaitGroup

func TestSubscribeHandler(t *testing.T) {
	var eventBus = NewEventBus()

	wg.Add(1)

	eventBus.Subscribe("address", func(context DeliveryContext) {
		assert.Equal(t, "Hi There", context.Result().Data)
		wg.Done()
	})
	eventBus.Publish("address", "Hi There", MessageOptions{})
	wg.Wait()
}

func TestTwiceSubscribe(t *testing.T) {
	var eventBus = NewEventBus()

	wg.Add(2)

	eventBus.Subscribe("address", func(context DeliveryContext) {
		assert.Equal(t, "Hi There", context.Result().Data)
		wg.Done()
	})

	eventBus.Subscribe("address", func(context DeliveryContext) {
		assert.Equal(t, "Hi There", context.Result().Data)
		wg.Done()
	})

	eventBus.Publish("address", "Hi There", MessageOptions{})
	wg.Wait()
}

func TestRequestReplyHandler(t *testing.T) {
	var eventBus = NewEventBus()

	wg.Add(1)

	eventBus.Subscribe("address", func(context DeliveryContext) {
		context.Reply("Hello")
	})

	eventBus.Request("address", "Hi there", func(context DeliveryContext) {
		println("sono nella request con")
		println(context)
		assert.Equal(t, "Hello", context.Result().Data)
		wg.Done()
	})
	wg.Wait()
}

func TestMessageOptions(t *testing.T) {
	var eventBus = NewEventBus()

	wg.Add(1)
	eventBus.Subscribe("address", func(context DeliveryContext) {
		assert.Equal(t, "value", context.Result().Headers["key"])
		wg.Done()
	})

	options := NewMessageOptions()
	options.AddHeader("key", "value")

	eventBus.Publish("address", "Hi There", options)
	wg.Wait()
}
