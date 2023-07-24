package goeventbus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddNewConsumer(t *testing.T) {
	topic := NewTopic("topic")

	actual := Handler{Address: "topic", Type: Consumer}
	topic.AddHandler(&actual)

	handlers := topic.GetHandlers()
	assert.Equal(t, handlers[0].Address, actual.Address)
	assert.Equal(t, handlers[0].Type, actual.Type)
}

func TestAddNewInterceptor(t *testing.T) {
	topic := NewTopic("topic")

	actual := Handler{Address: "topic", Type: Consumer}
	topic.AddHandler(&actual)

	handlers := topic.GetHandlers()
	assert.Equal(t, handlers[0].Address, actual.Address)
	assert.Equal(t, handlers[0].Type, actual.Type)
}
