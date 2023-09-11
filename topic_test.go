package goeventbus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddNewConsumer(t *testing.T) {
	topic := NewTopic("topic")

	topic.AddConsumer(func(context ConsumerContext) {})

	handlers := topic.Consumers
	assert.Equal(t, handlers[0].Closed(), false)
}

func TestAddNewInterceptor(t *testing.T) {
	topic := NewTopic("topic")

	topic.AddInterceptor(func(context InterceptorContext) {})

	handlers := topic.Interceptors
	assert.Equal(t, handlers[0].Closed(), false)
}
