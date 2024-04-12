package goeventbus

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	go NewServer("localhost:8080", "/").Listen()

	conn, err := net.Dial("tcp", "localhost:8080")
	assert.Nil(t, err)

	defer conn.Close()

	data := []byte("Hello, Server!")
	_, err = conn.Write(data)
	assert.Nil(t, err)
}
