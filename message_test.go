package goeventbus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessageBuilder(t *testing.T) {
	t.Run("should build message with payload", func(t *testing.T) {
		message := NewMessageBuilder().SetPayload("test payload").Build()
		assert.Equal(t, "test payload", message.Extract())
	})

	t.Run("should build message with headers", func(t *testing.T) {
		headers := NewMessageHeadersBuilder().SetHeader("key", "value").Build()
		message := NewMessageBuilder().SetHeaders(headers).Build()
		assert.True(t, message.ExtractHeaders().Contains("key"))
		assert.Equal(t, "value", message.ExtractHeaders().Get("key"))
	})

	t.Run("should build message with payload and headers", func(t *testing.T) {
		headers := NewMessageHeadersBuilder().SetHeader("header1", "value1").Build()
		message := NewMessageBuilder().
			SetPayload("payload").
			SetHeaders(headers).
			Build()

		assert.Equal(t, "payload", message.Extract())
		assert.True(t, message.ExtractHeaders().Contains("header1"))
		assert.Equal(t, "value1", message.ExtractHeaders().Get("header1"))
	})

	t.Run("should handle empty message", func(t *testing.T) {
		message := NewMessageBuilder().Build()
		assert.Nil(t, message.Extract())
		assert.False(t, message.ExtractHeaders().Contains("key"))
	})
}

