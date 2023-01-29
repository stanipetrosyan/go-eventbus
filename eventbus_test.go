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
	eventBus.Subscribe("address", func(data Message) {
		assert.Equal(t, "Hi There", data.Data)
		wg.Done()
	})
	eventBus.Publish("address", "Hi There", MessageOptions{})
	wg.Wait()

}

func TestMessageOptions(t *testing.T) {
	var eventBus = NewEventBus()

	wg.Add(1)
	eventBus.Subscribe("address", func(data Message) {
		assert.Equal(t, "value", data.Headers["key"])
		wg.Done()
	})

	options := NewMessageOptions()
	options.AddHeader("key", "value")

	eventBus.Publish("address", "Hi There", options)
	wg.Wait()
}
