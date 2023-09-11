package goeventbus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddNewConsumer(t *testing.T) {
	topic := NewTopic("topic")

	actual := NewConsumer("topic", func(context ConsumerContext) {})
	topic.AddConsumer(actual)

	handlers := topic.Consumers
	assert.Equal(t, handlers[0].Closed(), false)
}

func TestAddNewInterceptor(t *testing.T) {
	topic := NewTopic("topic")

	actual := NewInterceptor("topic", func(context InterceptorContext) {}, NewInterceptorContext([]chan Message{}))
	topic.AddInterceptor(actual)

	handlers := topic.Interceptors
	assert.Equal(t, handlers[0].Closed(), false)
}
