package goeventbus

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	var wg sync.WaitGroup

	wg.Add(1)
	server := newServer("localhost:8082")
	go server.Listen()

	var conn net.Conn
	var err error

	for conn, err = net.Dial("tcp", "localhost:8082"); err != nil; {
		conn, err = net.Dial("tcp", "localhost:8082")
	}

	assert.Nil(t, err)

	go func() {
		var request request
		d := json.NewDecoder(conn)
		err := d.Decode(&request)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		assert.Equal(t, "my-channel", request.Channel)
		assert.Equal(t, "Hello there", request.Payload)
		wg.Done()
	}()

	time.Sleep(time.Millisecond)

	msg := NewMessageBuilder().SetPayload("Hello there").Build()
	server.Publish("my-channel", msg)
	wg.Wait()
}
