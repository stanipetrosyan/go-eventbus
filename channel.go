package goeventbus

type Channel interface {
	Register() Channel
	Publisher() Publisher
	Subscriber() Subscriber
}

type defaultChannel struct {
	address string
	ch      chan string
	chs     []chan string
}

func (c *defaultChannel) Listen() {
	for {
		data, ok := <-c.ch
		if !ok {
			return
		}

		for _, item := range c.chs {
			item <- data
		}
	}
}

func (c *defaultChannel) Register() Channel {
	return c
}

func (c *defaultChannel) Publisher() Publisher {
	return NewPublisher(c.ch)
}

func (c *defaultChannel) Subscriber() Subscriber {
	ch := make(chan string)
	c.chs = append(c.chs, ch)

	return NewSubscriber(ch)
}

func NewChannel(address string) Channel {
	ch := make(chan string)
	channel := defaultChannel{address: address, ch: ch, chs: []chan string{}}
	go channel.Listen()
	return &channel
}
