package goeventbus

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEventBus_Channel(t *testing.T) {
	t.Run("should create new channel on first access", func(t *testing.T) {
		bus := NewEventBus()
		channel := bus.Channel("new-channel")
		assert.NotNil(t, channel)
	})

	t.Run("should return same channel for same address", func(t *testing.T) {
		bus := NewEventBus()
		channel1 := bus.Channel("same-channel")
		channel2 := bus.Channel("same-channel")
		assert.Equal(t, channel1, channel2)
	})

	t.Run("should create different channels for different addresses", func(t *testing.T) {
		bus := NewEventBus()
		channel1 := bus.Channel("channel-1")
		channel2 := bus.Channel("channel-2")
		assert.NotEqual(t, channel1, channel2)
	})
}

func TestEventBus_ConcurrentAccess(t *testing.T) {
	t.Run("should handle concurrent channel creation", func(t *testing.T) {
		bus := NewEventBus()
		var wg sync.WaitGroup
		channels := make([]Channel, 100)

		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				channels[index] = bus.Channel("concurrent-channel")
			}(i)
		}

		wg.Wait()

		// All should be the same channel instance
		firstChannel := channels[0]
		for i := 1; i < 100; i++ {
			assert.Equal(t, firstChannel, channels[i], "All concurrent accesses should return the same channel")
		}
	})
}

