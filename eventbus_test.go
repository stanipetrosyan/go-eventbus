package goeventbus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubscribeHandler(t *testing.T) {
	var eventBus = NewEventBus()

	eventBus.Subscribe("address")

	eventBus.On("address", func(data Message) {
		assert.Equal(t, "Hi There", data.Data)
	})

	eventBus.Publish("address", "Hi There", MessageOptions{})
}

func TestMessageOptions(t *testing.T) {
	var eventBus = NewEventBus()

	eventBus.Subscribe("address")

	eventBus.On("address", func(data Message) {
		assert.Equal(t, "key", data.Headers[0])
	})

	options := NewMessageOptions()
	options.AddHeader("key", "value")

	eventBus.Publish("address", "Hi There", options)
}
