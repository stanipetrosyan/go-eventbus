package goeventbus

import (
	"encoding/json"
	"net"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	var wg sync.WaitGroup
	var eventbus EventBus = NewEventBus()

	eventbus.Channel("channel").Subscriber().Listen(func(context Context) {
		assert.Equal(t, "Hello there", context.Result().Data)
		wg.Done()
	})

	wg.Add(1)

	listener, err := net.Listen("tcp", "localhost:8083")
	assert.Nil(t, err)

	go func() {
		for {
			conn, err := listener.Accept()
			assert.Nil(t, err)

			msg := Request{Channel: "channel", Message: CreateMessage().SetBody("Hello there")}
			json.NewEncoder(conn).Encode(msg)

		}
	}()

	client := NewClient("localhost:8083", eventbus)
	go client.Connect()

	wg.Wait()
}
