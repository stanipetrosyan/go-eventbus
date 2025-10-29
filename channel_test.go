package goeventbus

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubscriberHandler(t *testing.T) {
	t.Run("should handle a message", func(t *testing.T) {
		eventBus := NewEventBus()
		var wg sync.WaitGroup
		wg.Add(1)

		eventBus.Channel("subscriber-channel-1").Subscriber().Listen(func(context Context) {
			assert.Nil(t, context.Error())
			assert.Equal(t, "Hi There", context.Result().Extract())
			wg.Done()
		})

		message := NewMessageBuilder().SetPayload("Hi There").Build()
		eventBus.Channel("subscriber-channel-1").Publisher().Publish(message)
		wg.Wait()
	})

	t.Run("should handle message with many subscribers", func(t *testing.T) {
		eventBus := NewEventBus()
		var wg sync.WaitGroup
		wg.Add(2)

		eventBus.Channel("subscriber-channel-2").Subscriber().Listen(func(context Context) {
			assert.Nil(t, context.Error())
			assert.Equal(t, "Hi There", context.Result().Extract())
			wg.Done()
		})

		eventBus.Channel("subscriber-channel-2").Subscriber().Listen(func(context Context) {
			assert.Nil(t, context.Error())
			assert.Equal(t, "Hi There", context.Result().Extract())
			wg.Done()
		})

		message := NewMessageBuilder().SetPayload("Hi There").Build()
		eventBus.Channel("subscriber-channel-2").Publisher().Publish(message)
		wg.Wait()
	})
}

func TestRequestHandler(t *testing.T) {
	t.Run("should a subscriber response to a publisher", func(t *testing.T) {
		eventBus := NewEventBus()
		var wg sync.WaitGroup
		wg.Add(1)

		eventBus.Channel("request-channel-1").Subscriber().Listen(func(context Context) {
			assert.Nil(t, context.Error())
			assert.Equal(t, "Hi There", context.Result().Extract())
			context.Reply(NewMessageBuilder().SetPayload("Hello there!").Build())
		})

		message := NewMessageBuilder().SetPayload("Hi There").Build()
		eventBus.Channel("request-channel-1").Publisher().Request(message, func(context Context) {
			assert.Nil(t, context.Error())
			assert.Equal(t, "Hello there!", context.Result().Extract())
			wg.Done()
		})
		wg.Wait()
	})
}

func TestProcessorHandler(t *testing.T) {
	t.Run("should pass processor", func(t *testing.T) {
		eventBus := NewEventBus()
		var wg sync.WaitGroup
		wg.Add(2)

		eventBus.Channel("processor-channel-1").Subscriber().Listen(func(context Context) {
			assert.Nil(t, context.Error())
			assert.Equal(t, "Hi There", context.Result().Extract())
			wg.Done()
		})

		eventBus.Channel("processor-channel-1").Processor().Listen(func(context Context) {
			if context.Result().Extract() == "Hi There" {
				wg.Done()
				context.Next()
			}
		})

		message := NewMessageBuilder().SetPayload("Hi There").Build()
		eventBus.Channel("processor-channel-1").Publisher().Publish(message)
		wg.Wait()
	})

	t.Run("should re map message in processor", func(t *testing.T) {
		eventBus := NewEventBus()
		var wg sync.WaitGroup
		wg.Add(2)

		eventBus.Channel("processor-channel-2").Subscriber().Listen(func(context Context) {
			assert.Nil(t, context.Error())
			assert.Equal(t, "Hellooo", context.Result().Extract())
			wg.Done()
		})

		eventBus.Channel("processor-channel-2").Processor().Listen(func(context Context) {
			if context.Result().Extract() == "Hi There" {
				wg.Done()
				newMessage := NewMessageBuilder().SetPayload("Hellooo").Build()

				context.Map(newMessage).Next()
			}
		})

		message := NewMessageBuilder().SetPayload("Hi There").Build()
		eventBus.Channel("processor-channel-2").Publisher().Publish(message)
		wg.Wait()
	})
}

func TestMessageOptions(t *testing.T) {
	var eventBus = NewEventBus()
	var wg sync.WaitGroup

	wg.Add(1)
	eventBus.Channel("address").Subscriber().Listen(func(context Context) {
		assert.Nil(t, context.Error())
		assert.True(t, context.Result().ExtractHeaders().Contains("key"))
		assert.Equal(t, "value", context.Result().ExtractHeaders().Get("key"))
		wg.Done()
	})

	options := NewMessageHeadersBuilder().SetHeader("key", "value").Build()
	message := NewMessageBuilder().SetPayload("Hi There").SetHeaders(options).Build()
	eventBus.Channel("address").Publisher().Publish(message)
	wg.Wait()
}
