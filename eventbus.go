package goeventbus

import (
	"log/slog"
	"sync"
)

type EventBus interface {
	Channel(adress string) Channel
}

type defaultEventBus struct {
	channels sync.Map
}

// Returns a Channel with that address. If it doesn't exist, it returns a new one
func (e *defaultEventBus) Channel(address string) Channel {
	channel, _ := e.channels.LoadOrStore(address, newChannel(address))

	slog.Info("Channel created", slog.String("name", address))
	return channel.(Channel)
}

// Create a new eventbus with default parameters
func NewEventBus() EventBus {
	return &defaultEventBus{
		channels: sync.Map{},
	}
}
