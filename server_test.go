package goeventbus

import (
	"fmt"
	"net"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	var wg sync.WaitGroup

	wg.Add(1)
	server := NewServer("localhost:8082", "/")
	go server.Listen()

	conn, err := net.Dial("tcp", "localhost:8082")
	assert.Nil(t, err)

	//defer conn.Close()

	buffer := make([]byte, 1024)

	go func() {
		println("reading")
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		/*if err == io.EOF {
			t.FailNow()
		} */

		println(string(buffer[:n]))
		assert.Equal(t, "my-channel", string(buffer[:n]))
		wg.Done()

	}()

	server.Publish("my-channel")
	wg.Wait()

}
