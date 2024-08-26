package goeventbus

type request struct {
	Channel string  `json:"channel"`
	Message Message `json:"message"`
}
