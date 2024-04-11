package goeventbus

type Server interface {
	Listen()
}

type tcpServer struct {
}

func (s tcpServer) Listen() {
}

func NewServer() Server {
	return tcpServer{}
}
