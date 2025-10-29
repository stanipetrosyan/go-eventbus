package goeventbus

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestChannel_ConcurrentSubscribers(t *testing.T) {
	t.Run("should handle concurrent subscribers", func(t *testing.T) {
		eventBus := NewEventBus()
		var wg sync.WaitGroup
		const subscriberCount = 10
		const messageCount = 5

		wg.Add(subscriberCount * messageCount)

		// Create multiple subscribers
		for i := 0; i < subscriberCount; i++ {
			eventBus.Channel("concurrent-channel").Subscriber().Listen(func(context Context) {
				assert.Nil(t, context.Error())
				wg.Done()
			})
		}

		// Publish multiple messages
		for i := 0; i < messageCount; i++ {
			message := NewMessageBuilder().SetPayload(i).Build()
			eventBus.Channel("concurrent-channel").Publisher().Publish(message)
		}

		// Add timeout
		done := make(chan bool, 1)
		go func() {
			wg.Wait()
			done <- true
		}()

		select {
		case <-done:
			// Success
		case <-time.After(2 * time.Second):
			t.Fatal("Test timed out")
		}
	})
}

func TestChannel_ConcurrentPublishers(t *testing.T) {
	t.Run("should handle concurrent publishers", func(t *testing.T) {
		eventBus := NewEventBus()
		var wg sync.WaitGroup
		const publisherCount = 10
		var receivedMessages sync.Map

		wg.Add(publisherCount)

		// Single subscriber
		eventBus.Channel("publisher-channel").Subscriber().Listen(func(context Context) {
			msg := context.Result().Extract()
			receivedMessages.Store(msg, true)
			wg.Done()
		})

		// Multiple concurrent publishers
		var publishWg sync.WaitGroup
		for i := 0; i < publisherCount; i++ {
			publishWg.Add(1)
			go func(value int) {
				defer publishWg.Done()
				message := NewMessageBuilder().SetPayload(value).Build()
				eventBus.Channel("publisher-channel").Publisher().Publish(message)
			}(i)
		}

		// Add timeout
		done := make(chan bool, 1)
		go func() {
			publishWg.Wait()
			wg.Wait()
			done <- true
		}()

		select {
		case <-done:
			// Verify all messages were received
			for i := 0; i < publisherCount; i++ {
				_, ok := receivedMessages.Load(i)
				assert.True(t, ok, "Message %d should have been received", i)
			}
		case <-time.After(2 * time.Second):
			t.Fatal("Test timed out")
		}
	})
}

func TestChannel_ConcurrentPublisherAndSubscriberCreation(t *testing.T) {
	t.Run("should handle concurrent creation of publishers and subscribers", func(t *testing.T) {
		eventBus := NewEventBus()
		var wg sync.WaitGroup
		const goroutineCount = 20

		wg.Add(goroutineCount)

		for i := 0; i < goroutineCount; i++ {
			go func(index int) {
				defer wg.Done()
				if index%2 == 0 {
					// Create subscriber
					eventBus.Channel("test-channel").Subscriber().Listen(func(context Context) {})
				} else {
					// Create publisher
					eventBus.Channel("test-channel").Publisher()
				}
			}(i)
		}

		wg.Wait()
		// Test passes if no race condition or panic occurs
	})
}

func TestRequest_FirstResponse(t *testing.T) {
	t.Run("should only handle first response in Request", func(t *testing.T) {
		eventBus := NewEventBus()
		var responseCount int32
		var wg sync.WaitGroup
		wg.Add(1)

		// Multiple subscribers that reply
		for i := 0; i < 3; i++ {
			eventBus.Channel("request-channel").Subscriber().Listen(func(context Context) {
				context.Reply(NewMessageBuilder().SetPayload("response").Build())
			})
		}

		message := NewMessageBuilder().SetPayload("request").Build()
		eventBus.Channel("request-channel").Publisher().Request(message, func(context Context) {
			assert.Nil(t, context.Error())
			atomic.AddInt32(&responseCount, 1)
			wg.Done()
		})

		wg.Wait()
		time.Sleep(10 * time.Millisecond) // Wait for other responses
		assert.Equal(t, int32(1), atomic.LoadInt32(&responseCount), "Only first response should be handled")
	})
}

