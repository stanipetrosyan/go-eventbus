package goeventbus

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessageHeaders_MarshalJSON(t *testing.T) {
	t.Run("should marshal headers to JSON", func(t *testing.T) {
		headers := NewMessageHeadersBuilder().
			SetHeader("key1", "value1").
			SetHeader("key2", "value2").
			Build()

		data, err := json.Marshal(headers)
		assert.Nil(t, err)
		assert.Contains(t, string(data), "key1")
		assert.Contains(t, string(data), "value1")
		assert.Contains(t, string(data), "key2")
		assert.Contains(t, string(data), "value2")
	})

	t.Run("should marshal empty headers to empty JSON object", func(t *testing.T) {
		headers := MessageHeaders{headers: nil}
		data, err := json.Marshal(headers)
		assert.Nil(t, err)
		assert.Equal(t, "{}", string(data))
	})

	t.Run("should marshal nil headers", func(t *testing.T) {
		headers := MessageHeaders{}
		data, err := json.Marshal(headers)
		assert.Nil(t, err)
		assert.Equal(t, "{}", string(data))
	})
}

func TestMessageHeaders_UnmarshalJSON(t *testing.T) {
	t.Run("should unmarshal JSON to headers", func(t *testing.T) {
		jsonData := `{"key1":"value1","key2":"value2"}`
		var headers MessageHeaders
		err := json.Unmarshal([]byte(jsonData), &headers)

		assert.Nil(t, err)
		assert.True(t, headers.Contains("key1"))
		assert.Equal(t, "value1", headers.Get("key1"))
		assert.True(t, headers.Contains("key2"))
		assert.Equal(t, "value2", headers.Get("key2"))
	})

	t.Run("should unmarshal empty JSON object", func(t *testing.T) {
		jsonData := `{}`
		var headers MessageHeaders
		err := json.Unmarshal([]byte(jsonData), &headers)

		assert.Nil(t, err)
		assert.False(t, headers.Contains("key1"))
	})
}

func TestMessageHeaders_Get_Contains_Nil(t *testing.T) {
	t.Run("should handle nil headers in Get", func(t *testing.T) {
		headers := MessageHeaders{}
		value := headers.Get("key")
		assert.Empty(t, value)
	})

	t.Run("should handle nil headers in Contains", func(t *testing.T) {
		headers := MessageHeaders{}
		exists := headers.Contains("key")
		assert.False(t, exists)
	})
}

func TestMessageHeadersBuilder(t *testing.T) {
	t.Run("should build headers with multiple values", func(t *testing.T) {
		headers := NewMessageHeadersBuilder().
			SetHeader("header1", "value1").
			SetHeader("header2", "value2").
			SetHeader("header3", "value3").
			Build()

		assert.True(t, headers.Contains("header1"))
		assert.Equal(t, "value1", headers.Get("header1"))
		assert.True(t, headers.Contains("header2"))
		assert.Equal(t, "value2", headers.Get("header2"))
		assert.True(t, headers.Contains("header3"))
		assert.Equal(t, "value3", headers.Get("header3"))
	})

	t.Run("should allow overwriting header values", func(t *testing.T) {
		headers := NewMessageHeadersBuilder().
			SetHeader("key", "value1").
			SetHeader("key", "value2").
			Build()

		assert.Equal(t, "value2", headers.Get("key"))
	})
}

