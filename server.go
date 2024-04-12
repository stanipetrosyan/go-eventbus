package goeventbus

import (
	"fmt"
	"net"
)

type Server interface {
	Listen()
}

type tcpServer struct {
	address, path string
}

func (s tcpServer) Listen() {

	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		handleClient(conn)
	}
}

func NewServer(address, path string) Server {
	return tcpServer{address: address, path: path}
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
