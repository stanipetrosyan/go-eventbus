package goeventbus

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestProcessor_BlocksMessages(t *testing.T) {
	t.Run("should block messages when processor does not call Next()", func(t *testing.T) {
		eventBus := NewEventBus()
		var received bool

		// Processor that doesn't call Next()
		eventBus.Channel("test-channel").Processor().Listen(func(context Context) {
			// Do not call Next(), so message should be blocked
			received = true
		})

		// Subscriber should never receive the message
		eventBus.Channel("test-channel").Subscriber().Listen(func(context Context) {
			t.Error("Subscriber should not receive message when processor blocks")
		})

		message := NewMessageBuilder().SetPayload("test").Build()
		eventBus.Channel("test-channel").Publisher().Publish(message)

		// Wait a bit to ensure message doesn't propagate
		time.Sleep(10 * time.Millisecond)
		assert.True(t, received, "Processor should have received the message")
		// Subscriber should not have received anything (no assertion error means success)
	})
}

func TestProcessor_FiltersMessages(t *testing.T) {
	t.Run("should filter messages based on condition", func(t *testing.T) {
		eventBus := NewEventBus()
		var wg sync.WaitGroup
		var mu sync.Mutex
		var receivedCount int

		wg.Add(1)

		// Processor that only forwards messages with specific header
		eventBus.Channel("filter-channel").Processor().Listen(func(context Context) {
			if context.Result().ExtractHeaders().Contains("allowed") {
				context.Next()
			}
		})

		eventBus.Channel("filter-channel").Subscriber().Listen(func(context Context) {
			mu.Lock()
			receivedCount++
			mu.Unlock()
			wg.Done()
		})

		// Message that should be blocked (no header)
		message1 := NewMessageBuilder().SetPayload("blocked").Build()
		eventBus.Channel("filter-channel").Publisher().Publish(message1)

		// Wait a bit for blocked message to be processed
		time.Sleep(10 * time.Millisecond)

		// Message that should pass
		headers := NewMessageHeadersBuilder().SetHeader("allowed", "true").Build()
		message2 := NewMessageBuilder().SetPayload("allowed").SetHeaders(headers).Build()
		eventBus.Channel("filter-channel").Publisher().Publish(message2)

		// Add timeout
		done := make(chan bool, 1)
		go func() {
			wg.Wait()
			done <- true
		}()

		select {
		case <-done:
			mu.Lock()
			count := receivedCount
			mu.Unlock()
			assert.Equal(t, 1, count, "Only one message should pass the filter")
		case <-time.After(2 * time.Second):
			t.Fatal("Test timed out waiting for message")
		}
	})
}

