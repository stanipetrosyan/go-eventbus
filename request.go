package goeventbus

type request struct {
	Channel string         `json:"channel"`
	Payload any            `json:"payload"`
	Headers MessageHeaders `json:"headers"`
}
