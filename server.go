package goeventbus

import (
	"context"
	"encoding/json"
	"log/slog"
	"net"
	"sync"
)

type Server interface {
	Listen() (Server, error)
	Publish(channel string, message Message)
	Stop() error
}

type tcpServer struct {
	sync.RWMutex
	address  string
	clients  []net.Conn
	listener net.Listener
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

func (s *tcpServer) Listen() (Server, error) {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	s.Lock()
	s.listener = listener
	s.ctx, s.cancel = context.WithCancel(context.Background())
	s.Unlock()

	slog.Info("Server started", slog.String("host", s.address))

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		defer func() {
			if err := listener.Close(); err != nil {
				slog.Error("Error closing listener", slog.String("error", err.Error()))
			}
		}()
		for {
			conn, err := listener.Accept()
			if err != nil {
				select {
				case <-s.ctx.Done():
					slog.Info("Server stopping, no longer accepting connections")
					return
				default:
					slog.Error("Error accepting connection", slog.String("error", err.Error()))
					continue
				}
			}

			slog.Info("Client Connected", slog.String("host", conn.RemoteAddr().String()))

			s.Lock()
			s.clients = append(s.clients, conn)
			s.Unlock()

			// Start goroutine to handle client disconnection on context cancellation
			s.wg.Add(1)
			go func(c net.Conn) {
				defer s.wg.Done()
				<-s.ctx.Done()
				s.removeClient(c)
			}(conn)
		}
	}()

	return s, nil
}


func (s *tcpServer) Publish(channel string, message Message) {
	s.RLock()
	clients := make([]net.Conn, len(s.clients))
	copy(clients, s.clients)
	s.RUnlock()

	var encoder *json.Encoder
	for _, client := range clients {
		encoder = json.NewEncoder(client)
		err := encoder.Encode(request{Channel: channel, Payload: message.Extract(), Headers: message.ExtractHeaders()})

		if err != nil {
			slog.Error("Error publishing to client", slog.String("error", err.Error()))
			// Remove dead connection
			s.removeClient(client)
		}
	}
}

func (s *tcpServer) removeClient(conn net.Conn) {
	s.Lock()
	defer s.Unlock()
	for i, client := range s.clients {
		if client == conn {
			s.clients = append(s.clients[:i], s.clients[i+1:]...)
			conn.Close()
			break
		}
	}
}

func (s *tcpServer) Stop() error {
	s.Lock()
	defer s.Unlock()

	if s.cancel != nil {
		s.cancel()
	}

	if s.listener != nil {
		if err := s.listener.Close(); err != nil {
			slog.Error("Error closing listener", slog.String("error", err.Error()))
		}
	}

	// Close all client connections
	for _, client := range s.clients {
		client.Close()
	}
	s.clients = []net.Conn{}

	s.wg.Wait()
	slog.Info("Server stopped")
	return nil
}

func newServer(address string) Server {
	return &tcpServer{address: address, clients: []net.Conn{}}
}
