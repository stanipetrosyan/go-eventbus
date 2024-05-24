package goeventbus

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	var wg sync.WaitGroup

	wg.Add(1)
	server := NewServer("localhost:8082")
	go server.Listen()

	var conn net.Conn
	var err error

	for conn, err = net.Dial("tcp", "localhost:8082"); err != nil; {
		conn, err = net.Dial("tcp", "localhost:8082")
	}

	assert.Nil(t, err)

	go func() {
		var request Request
		d := json.NewDecoder(conn)
		err := d.Decode(&request)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		assert.Equal(t, "my-channel", request.Channel)
		assert.Equal(t, "Hello there", request.Message.Extract())
		wg.Done()

	}()

	msg := CreateMessage().SetBody("Hello there")
	server.Publish("my-channel", msg)
	wg.Wait()

}
