package goeventbus

type Topic struct {
	Address     string
	Handlers    []*Handler
	Interceptor []*Handler
}

func (t *Topic) AddInterceptor(interceptor Handler) {
	t.Interceptor = append(t.Interceptor, &interceptor)
}
