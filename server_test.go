package goeventbus

import (
	"encoding/json"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServer_Listen(t *testing.T) {
	var wg sync.WaitGroup

	wg.Add(1)
	server := newServer("localhost:8082")
	srv, err := server.Listen()
	assert.Nil(t, err)
	assert.NotNil(t, srv)

	tcpSrv := srv.(*tcpServer)

	var conn net.Conn
	var dialErr error

	// Wait for server to be ready
	for i := 0; i < 10; i++ {
		conn, dialErr = net.Dial("tcp", "localhost:8082")
		if dialErr == nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	assert.Nil(t, dialErr)
	defer conn.Close()

	go func() {
		defer wg.Done()
		var request request
		d := json.NewDecoder(conn)
		err := d.Decode(&request)
		if err != nil {
			return
		}
		assert.Equal(t, "my-channel", request.Channel)
		assert.Equal(t, "Hello there", request.Payload)
	}()

	time.Sleep(50 * time.Millisecond) // Give time for connection to be registered

	msg := NewMessageBuilder().SetPayload("Hello there").Build()
	server.Publish("my-channel", msg)

	// Add timeout to prevent hanging
	done := make(chan bool, 1)
	go func() {
		wg.Wait()
		done <- true
	}()

	select {
	case <-done:
		// Success
	case <-time.After(2 * time.Second):
		t.Fatal("Test timed out")
	}

	// Cleanup
	tcpSrv.Stop()
}

func TestServer_Stop(t *testing.T) {
	t.Run("should stop server and close all connections", func(t *testing.T) {
		server := newServer("localhost:8084")
		srv, err := server.Listen()
		assert.Nil(t, err)

		tcpSrv := srv.(*tcpServer)

		// Connect a client
		conn, err := net.Dial("tcp", "localhost:8084")
		assert.Nil(t, err)
		defer conn.Close()

		// Wait a bit for connection to be registered
		time.Sleep(50 * time.Millisecond)

		// Stop server
		err = tcpSrv.Stop()
		assert.Nil(t, err)

		// Verify connection is closed
		time.Sleep(50 * time.Millisecond)
		conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
		_, err = conn.Read(make([]byte, 1))
		assert.NotNil(t, err, "Connection should be closed")
	})

	t.Run("should handle stop on server that never started", func(t *testing.T) {
		server := newServer("localhost:8085")
		tcpSrv := server.(*tcpServer)
		err := tcpSrv.Stop()
		assert.Nil(t, err)
	})
}

func TestServer_RemoveDeadClients(t *testing.T) {
	t.Run("should remove dead clients when publish fails", func(t *testing.T) {
		server := newServer("localhost:8086")
		srv, err := server.Listen()
		assert.Nil(t, err)

		tcpSrv := srv.(*tcpServer)
		defer tcpSrv.Stop()

		// Connect and immediately close
		conn, err := net.Dial("tcp", "localhost:8086")
		assert.Nil(t, err)
		conn.Close()

		// Wait a bit for connection to be registered
		time.Sleep(50 * time.Millisecond)

		// Try to publish - should remove dead connection
		msg := NewMessageBuilder().SetPayload("test").Build()
		tcpSrv.Publish("test-channel", msg)

		// Wait for cleanup
		time.Sleep(100 * time.Millisecond)

		tcpSrv.RLock()
		clientCount := len(tcpSrv.clients)
		tcpSrv.RUnlock()

		assert.Equal(t, 0, clientCount, "Dead clients should be removed")
	})
}
