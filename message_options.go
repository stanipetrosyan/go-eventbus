package goeventbus

type MessageOptions struct {
	headers Headers
}

// Set the headers of the message. For create a new headers see: NewHeaders()
func (op MessageOptions) SetHeaders(headers Headers) MessageOptions {
	op.headers = headers
	return op
}

// Returns headers of message
func (op MessageOptions) Headers() Headers {
	return op.headers
}

// Retuns empty message options
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

// Returns empty headers
func NewHeaders() Headers {
	return Headers{headers: map[string]string{}}
}
