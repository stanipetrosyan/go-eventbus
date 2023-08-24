package goeventbus

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

var wg sync.WaitGroup

func TestSubscribeHandler(t *testing.T) {
	var eventBus = NewEventBus()

	t.Run("should handle a message", func(t *testing.T) {
		wg.Add(1)

		eventBus.Subscribe("address", func(context DeliveryContext) {
			assert.Equal(t, "Hi There", context.Result().Data)
			wg.Done()
		})

		message := CreateMessage().SetBody("Hi There")
		eventBus.Publish("address", message)
		wg.Wait()
	})

	t.Run("should handle message for more of one handlers", func(t *testing.T) {
		wg.Add(2)

		eventBus.Subscribe("newaddress", func(context DeliveryContext) {
			assert.Equal(t, "Hi There", context.Result().Data)
			wg.Done()
		})

		eventBus.Subscribe("newaddress", func(context DeliveryContext) {
			assert.Equal(t, "Hi There", context.Result().Data)
			wg.Done()
		})

		message := CreateMessage().SetBody("Hi There")
		eventBus.Publish("newaddress", message)
		wg.Wait()
	})

}

func TestRequestReplyHandler(t *testing.T) {
	var eventBus = NewEventBus()

	t.Run("should send a request and reply", func(t *testing.T) {
		wg.Add(1)

		eventBus.Subscribe("address", func(context DeliveryContext) {
			context.Reply("Hello")
		})

		message := CreateMessage().SetBody("Hi There")
		eventBus.Request("address", message, func(context DeliveryContext) {
			context.Handle(func(message Message) {
				assert.Equal(t, "Hello", message.Data)
				wg.Done()
			})
		})
		wg.Wait()
	})
}

func TestInBoundInterceptorHandler(t *testing.T) {
	var eventBus = NewEventBus()

	t.Run("should pass interceptor handler", func(t *testing.T) {
		wg.Add(2)

		eventBus.Subscribe("address", func(context DeliveryContext) {
			assert.Equal(t, "Hi There", context.Result().Data)
			wg.Done()
		})

		eventBus.AddInBoundInterceptor("address", func(context DeliveryContext) {
			assert.Equal(t, "Hi There", context.Result().Data)
			wg.Done()
			context.Next()
		})

		message := CreateMessage().SetBody("Hi There")
		eventBus.Publish("address", message)
		wg.Wait()
	})
}

func TestMessageOptions(t *testing.T) {
	var eventBus = NewEventBus()

	wg.Add(1)
	eventBus.Subscribe("address", func(context DeliveryContext) {
		assert.Equal(t, "value", context.Result().Options.Header("key"))
		wg.Done()
	})

	message := CreateMessage().SetBody("Hi There").SetOptions(NewMessageOptions().AddHeader("key", "value"))
	eventBus.Publish("address", message)
	wg.Wait()
}
