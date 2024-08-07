package goeventbus

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"sync"
)

type Server interface {
	Listen() (Server, error)
	Publish(channel string, message Message)
}

type tcpServer struct {
	sync.RWMutex
	address string
	clients []net.Conn
}

func (s *tcpServer) Listen() (Server, error) {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	slog.Info("Server started", slog.String("host", s.address))

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		s.Lock()
		s.clients = append(s.clients, conn)
		s.Unlock()
	}
}

func (s *tcpServer) Publish(channel string, message Message) {
	var encoder *json.Encoder
	s.Lock()
	for _, client := range s.clients {
		encoder = json.NewEncoder(client)
		err := encoder.Encode(Request{Channel: channel, Message: message})

		if err != nil {
			fmt.Println("Error:", err)
		}
	}
	s.Unlock()
}

func NewServer(address string) Server {
	return &tcpServer{address: address, clients: []net.Conn{}}
}
