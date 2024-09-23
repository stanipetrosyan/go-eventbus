package goeventbus

import (
	"log/slog"
	"sync"
)

type EventBus interface {
	// Returns a Channel with that address. If it doesn't exist, it returns a new one
	Channel(adress string) Channel
}

type defaultEventBus struct {
	channels sync.Map
}

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
