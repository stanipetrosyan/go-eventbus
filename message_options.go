package goeventbus

type MessageOptions struct {
	headers Headers
}

func (op MessageOptions) SetHeaders(headers Headers) MessageOptions {
	op.headers = headers
	return op
}

func (op MessageOptions) Headers() Headers {
	return op.headers
}

func NewMessageOptions() MessageOptions {
	return MessageOptions{
		headers: Headers{},
	}
}

type Headers struct {
	headers map[string]string
}

func (h Headers) Add(key string, value string) Headers {
	h.headers[key] = value
	return h
}

func (h Headers) Header(key string) string {
	return h.headers[key]
}

func (h Headers) Contains(key string) bool {
	_, exist := h.headers[key]
	return exist
}

func NewHeaders() Headers {
	return Headers{headers: map[string]string{}}
}
