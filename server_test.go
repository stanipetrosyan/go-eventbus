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
	server := NewServer("localhost:8082")
	go server.Listen()

	conn, err := net.Dial("tcp", "localhost:8082")
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
		assert.Equal(t, "Hello there", request.Message.Data)
		wg.Done()

	}()

	msg := CreateMessage().SetBody("Hello there")
	time.Sleep(time.Second * 2)
	server.Publish("my-channel", msg)
	wg.Wait()

}
