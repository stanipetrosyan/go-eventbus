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

	eventBus.Request("address", "Hi there", MessageOptions{}, func(context DeliveryContext) {
		context.Handle(func(message Message) {
			assert.Equal(t, "Hello", message.Data)
			wg.Done()
		})
	})
	wg.Wait()
}

func TestInBoundInterceptorHandler(t *testing.T) {
	var eventBus = NewEventBus()

	wg.Add(2)

	eventBus.Subscribe("address", func(context DeliveryContext) {
		assert.Equal(t, "Hi there", context.Result().Data)
		wg.Done()
	})

	eventBus.AddInBoundInterceptor("address", func(context DeliveryContext) {
		assert.Equal(t, "Hi there", context.Result().Data)
		wg.Done()
		context.Next()
	})

	eventBus.Publish("address", "Hi there", MessageOptions{})
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
