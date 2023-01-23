package goeventbus

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Fuzz struct {
	Field string `json:"field"`
}

func TestStringToJson(t *testing.T) {

	message := Message{Data: "Hi There"}

	js, err := message.ToJson()

	assert.Nil(t, err)

	var actual string

	json.Unmarshal(js, &actual)

	assert.Equal(t, "Hi There", actual)
}

func TestStructToJson(t *testing.T) {

	message := Message{Data: Fuzz{"Value"}}

	js, err := message.ToJson()

	assert.Nil(t, err)

	var actual Fuzz

	json.Unmarshal(js, &actual)

	assert.Equal(t, Fuzz{"Value"}, actual)
}
