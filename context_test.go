package goeventbus

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContext_Error(t *testing.T) {
	t.Run("should return error when context has error", func(t *testing.T) {
		err := errors.New("test error")
		ctx := newContextWithError(err)
		assert.Equal(t, err, ctx.Error())
		assert.NotNil(t, ctx.Error())
	})

	t.Run("should return nil error when context has no error", func(t *testing.T) {
		message := NewMessageBuilder().SetPayload("test").Build()
		ch := make(chan packet)
		ctx := newContextWithMessageAndChannel(message, ch)
		assert.Nil(t, ctx.Error())
	})
}

func TestContext_Map(t *testing.T) {
	t.Run("should map message to new message", func(t *testing.T) {
		message1 := NewMessageBuilder().SetPayload("original").Build()
		message2 := NewMessageBuilder().SetPayload("mapped").Build()
		ch := make(chan packet)

		ctx := newContextWithMessageAndChannel(message1, ch)
		assert.Equal(t, "original", ctx.Result().Extract())

		mappedCtx := ctx.Map(message2)
		assert.Equal(t, "mapped", mappedCtx.Result().Extract())
		// Original context should also be modified (pointer receiver)
		assert.Equal(t, "mapped", ctx.Result().Extract())
	})
}

func TestContext_Reply(t *testing.T) {
	t.Run("should send reply through channel", func(t *testing.T) {
		message := NewMessageBuilder().SetPayload("request").Build()
		reply := NewMessageBuilder().SetPayload("reply").Build()
		ch := make(chan packet, 1)

		ctx := newContextWithMessageAndChannel(message, ch)
		ctx.Reply(reply)

		packet := <-ch
		assert.Equal(t, SUBSCRIBER, packet.from)
		assert.Equal(t, "reply", packet.message.Extract())
	})
}

func TestContext_Next(t *testing.T) {
	t.Run("should forward message through channel", func(t *testing.T) {
		message := NewMessageBuilder().SetPayload("test").Build()
		ch := make(chan packet, 1)

		ctx := newContextWithMessageAndChannel(message, ch)
		ctx.Next()

		packet := <-ch
		assert.Equal(t, PROCESSOR, packet.from)
		assert.Equal(t, "test", packet.message.Extract())
	})
}

