package goeventbus

import (
	"encoding/json"
	"fmt"
	"net"
)

type Server interface {
	Listen() (Server, error)
	Publish(channel string, message Message)
}

type tcpServer struct {
	address string
	path    string
	clients []net.Conn
}

func (s *tcpServer) Listen() (Server, error) {

	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		s.clients = append(s.clients, conn)
	}
}

func (s *tcpServer) Publish(channel string, message Message) {
	var encoder *json.Encoder
	for _, client := range s.clients {
		encoder = json.NewEncoder(client)
		err := encoder.Encode(Request{Channel: channel, Message: message})

		if err != nil {
			fmt.Println("Error:", err)
		}
	}
}

func NewServer(address, path string) Server {
	return &tcpServer{address: address, path: path, clients: []net.Conn{}}
}
