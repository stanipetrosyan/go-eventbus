package goeventbus

import (
	"fmt"
	"net"
)

type Server interface {
	Listen() (Server, error)
	Publish(channel string)
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

		handleClient(conn)
	}
}

func (s *tcpServer) Publish(channel string) {
	for _, client := range s.clients {
		_, err := client.Write([]byte(channel))
		if err != nil {
			fmt.Println("Error:", err)
		}
	}
}

func NewServer(address, path string) Server {
	return &tcpServer{address: address, path: path, clients: []net.Conn{}}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)

	for {
		// Read data from the client
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// Process and use the data (here, we'll just print it)
		fmt.Printf("Received: %s\n", buffer[:n])
	}
}
