package goeventbus

import (
	"encoding/json"
	"log/slog"
	"net"
	"time"
)

type Client interface {
	Connect()
}

type tcpClient struct {
	address  string
	eventbus EventBus
}

func (s *tcpClient) Connect() {
	for {
		conn, err := net.Dial("tcp", s.address)
		if err != nil {
			backoff := time.Second
			for {
				slog.Warn("Failed to connect, retrying...", slog.String("address", s.address), slog.String("error", err.Error()), slog.Duration("backoff", backoff))
				time.Sleep(backoff)
				conn, err = net.Dial("tcp", s.address)
				if err == nil {
					break
				}
				// Exponential backoff with max 30 seconds
				if backoff < 30*time.Second {
					backoff *= 2
				}
			}
		}

		slog.Info("Connected to server", slog.String("address", conn.RemoteAddr().String()))

		// Handle connection in a separate function that properly manages cleanup
		s.handleConnection(conn)
		
		// If we reach here, connection was closed, try to reconnect
		slog.Info("Connection closed, attempting to reconnect...")
	}
}

func (s *tcpClient) handleConnection(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			slog.Error("Error closing connection", slog.String("error", err.Error()))
		}
	}()

	for {
		var msg request
		d := json.NewDecoder(conn)
		err := d.Decode(&msg)

		if err != nil {
			slog.Error("Connection error, reconnecting...", slog.String("error", err.Error()))
			return // Return to reconnect
		}

		channel := msg.Channel

		message := NewMessageBuilder().SetPayload(msg.Payload).SetHeaders(msg.Headers).Build()
		s.eventbus.Channel(channel).Publisher().Publish(message)
	}
}

func newClient(address string, eventbus EventBus) Client {
	return &tcpClient{address: address, eventbus: eventbus}
}
