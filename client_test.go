package goeventbus

import (
	"net"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	var wg sync.WaitGroup
	var eventbus EventBus = NewEventBus()

	eventbus.Channel("channel").Subscriber().Listen(func(context Context) {
		println(context.Result().Data)
		wg.Done()
	})

	wg.Add(1)

	listener, err := net.Listen("tcp", "localhost:8082")
	assert.Nil(t, err)

	go func() {
		for {
			conn, err := listener.Accept()
			assert.Nil(t, err)

			conn.Write([]byte("channel"))
		}
	}()

	client := NewClient("localhost:8082", "/", eventbus)
	go client.Connect()

	wg.Wait()
}
