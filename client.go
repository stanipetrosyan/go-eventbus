package goeventbus

import (
	"encoding/json"
	"fmt"
	"net"
)

type Client interface {
	Connect()
}

type tcpClient struct {
	address  string
	eventbus EventBus
}

type Request struct {
	Channel string  `json:"channel"`
	Message Message `json:"message"`
}

func (s *tcpClient) Connect() {
	var conn net.Conn
	var err error

	for conn, err = net.Dial("tcp", s.address); err != nil; {
		conn, err = net.Dial("tcp", s.address)
	}

	defer conn.Close()

	for {
		var msg Request
		d := json.NewDecoder(conn)
		err = d.Decode(&msg)

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		channel := msg.Channel

		message := msg.Message
		s.eventbus.Channel(channel).Publisher().Publish(message)

	}
}

func NewClient(address string, eventbus EventBus) Client {
	return &tcpClient{address: address, eventbus: eventbus}
}
