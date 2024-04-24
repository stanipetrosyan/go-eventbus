package goeventbus

import (
	"fmt"
	"net"
)

type Client interface {
	Connect()
}

type tcpClient struct {
	address  string
	path     string
	eventbus EventBus
}

func (s *tcpClient) Connect() {
	var conn net.Conn
	var err error

	for conn, err = net.Dial("tcp", s.address); err != nil; {
		conn, err = net.Dial("tcp", s.address)
	}

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		channel := string(buffer[:n])

		message := CreateMessage().SetBody("Test")
		s.eventbus.Channel(channel).Publisher().Publish(message)

	}
}

func NewClient(address, path string, eventbus EventBus) Client {
	return &tcpClient{address: address, path: path, eventbus: eventbus}
}
