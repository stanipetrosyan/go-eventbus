package goeventbus

import (
	"encoding/json"
	"log/slog"
	"net"
)

type Client interface {
	Connect()
}

type tcpClient struct {
	address  string
	eventbus EventBus
}

func (s *tcpClient) Connect() {
	var conn net.Conn
	var err error

	for conn, err = net.Dial("tcp", s.address); err != nil; {
		conn, err = net.Dial("tcp", s.address)
	}

	slog.Info("Connected to server", slog.String("address", conn.LocalAddr().String()))

	defer conn.Close()

	for {
		var msg request
		d := json.NewDecoder(conn)
		err = d.Decode(&msg)

		if err != nil {
			slog.Error(err.Error())
			return
		}

		channel := msg.Channel

		message := NewMessageBuilder().SetPayload(msg.Payload).SetHeaders(msg.Headers).Build()
		s.eventbus.Channel(channel).Publisher().Publish(message)

	}
}

func newClient(address string, eventbus EventBus) Client {
	return &tcpClient{address: address, eventbus: eventbus}
}
