package goeventbus

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubscriberHandler(t *testing.T) {
	var eventBus = NewEventBus()
	var wg sync.WaitGroup

	t.Run("should handle a message using channels", func(t *testing.T) {
		wg.Add(1)

		eventBus.Channel("my-channel").Subscriber().Listen(func(context Context) {
			assert.Equal(t, "Hi There", context.Result().Data)
			wg.Done()
		})

		message := CreateMessage().SetBody("Hi There")
		eventBus.Channel("my-channel").Publisher().Publish(message)
		wg.Wait()
	})

	t.Run("should handle message for more of one handlers", func(t *testing.T) {
		wg.Add(2)

		eventBus.Channel("newaddress").Subscriber().Listen(func(context Context) {
			assert.Equal(t, "Hi There", context.Result().Data)
			wg.Done()
		})

		eventBus.Channel("newaddress").Subscriber().Listen(func(context Context) {
			assert.Equal(t, "Hi There", context.Result().Data)
			wg.Done()
		})

		message := CreateMessage().SetBody("Hi There")
		eventBus.Channel("newaddress").Publisher().Publish(message)
		wg.Wait()
	})
}

func TestProcessorHandler(t *testing.T) {
	var eventBus = NewEventBus()
	var wg sync.WaitGroup

	t.Run("should pass processor handler", func(t *testing.T) {
		wg.Add(2)

		eventBus.Channel("my-channel").Subscriber().Listen(func(context Context) {
			assert.Equal(t, "Hi There", context.Result().Data)
			wg.Done()
		})

		eventBus.Channel("my-channel").Processor(func(message Message) bool {
			wg.Done()
			return message.Data == "Hi There"
		})

		message := CreateMessage().SetBody("Hi There")
		eventBus.Channel("my-channel").Publisher().Publish(message)
		wg.Wait()
	})

	t.Run("should pass message to handler created afted interceptor", func(t *testing.T) {
		wg.Add(2)

		eventBus.Channel("my-channel").Processor(func(message Message) bool {
			wg.Done()
			return message.Data == "Hi There"
		})

		eventBus.Channel("my-channel").Subscriber().Listen(func(context Context) {
			assert.Equal(t, "Hi There", context.Result().Data)
			wg.Done()
		})

		message := CreateMessage().SetBody("Hi There")
		eventBus.Channel("my-channel").Publisher().Publish(message)
		wg.Wait()
	})
}

func TestMessageOptions(t *testing.T) {
	var eventBus = NewEventBus()
	var wg sync.WaitGroup

	wg.Add(1)
	eventBus.Channel("address").Subscriber().Listen(func(context Context) {
		assert.Equal(t, "value", context.Result().Options.Header("key"))
		wg.Done()
	})

	message := CreateMessage().SetBody("Hi There").SetOptions(NewMessageOptions().AddHeader("key", "value"))
	eventBus.Channel("address").Publisher().Publish(message)
	wg.Wait()
}
