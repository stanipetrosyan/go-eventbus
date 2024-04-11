package goeventbus

type Client interface {
	Connect()
}

type tcpClient struct {
}

func (s tcpClient) Connect() {
}

func NewClient() Client {
	return tcpClient{}
}
