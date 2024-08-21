package goeventbus

import "log/slog"

type EventBus interface {
	Channel(adress string) Channel
}

type defaultEventBus struct {
	channels map[string]Channel
}

// Returns a Channel with that address. If it doesn't exist, it returns a new one
func (e *defaultEventBus) Channel(address string) Channel {
	_, exists := e.channels[address]
	if !exists {
		e.channels[address] = newChannel(address)
	}

	slog.Info("Channel created", slog.String("name", address))
	return e.channels[address]
}

// Create a new eventbus with default parameters
func NewEventBus() EventBus {
	return &defaultEventBus{
		channels: map[string]Channel{},
	}
}
