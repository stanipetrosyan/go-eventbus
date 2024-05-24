package goeventbus

type Request struct {
	Channel string  `json:"channel"`
	Message Message `json:"message"`
}
