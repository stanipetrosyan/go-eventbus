package goeventbus

type Topic struct {
	Address     string
	Handlers    []*Handler
	Interceptor []*Handler
}
