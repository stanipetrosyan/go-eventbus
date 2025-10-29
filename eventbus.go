package goeventbus

import (
	"log/slog"
	"sync"
)

type EventBus interface {
	// Returns a Channel with that address. If it doesn't exist, it returns a new one
	Channel(address string) Channel
}

type defaultEventBus struct {
	channels sync.Map
}

func (e *defaultEventBus) Channel(address string) Channel {
	channel, loaded := e.channels.LoadOrStore(address, newChannel(address))

	if !loaded {
		slog.Info("Channel created", slog.String("name", address))
	}
	
	// Type assertion with safety check
	ch, ok := channel.(Channel)
	if !ok {
		// This should never happen, but handle gracefully
		slog.Error("Invalid channel type in sync.Map, recreating channel", slog.String("name", address))
		newCh := newChannel(address)
		e.channels.Store(address, newCh)
		return newCh
	}
	return ch
}

// Create a new eventbus with default parameters
func NewEventBus() EventBus {
	return &defaultEventBus{
		channels: sync.Map{},
	}
}
