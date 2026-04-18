// Package events provides an in-process event bus with async subscriber support.
package events

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"sync"
)

// Event is the base interface for all events.
type Event interface {
	EventName() string
}

// Handler is a function that handles an event.
type Handler func(ctx context.Context, event Event) error

// Subscriber registers event handlers.
type Subscriber interface {
	Subscribe(bus *EventBus)
}

// EventBus is a simple in-process event bus.
type EventBus struct {
	mu       sync.RWMutex
	handlers map[string][]Handler
	logger   *slog.Logger
	async    bool
}

// NewEventBus creates a new event bus.
func NewEventBus(logger *slog.Logger) *EventBus {
	return &EventBus{
		handlers: make(map[string][]Handler),
		logger:   logger,
		async:    true,
	}
}

// SetSync makes the event bus synchronous (useful for testing).
func (eb *EventBus) SetSync() *EventBus {
	eb.async = false
	return eb
}

// On registers a handler for an event name.
func (eb *EventBus) On(eventName string, handler Handler) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	eb.handlers[eventName] = append(eb.handlers[eventName], handler)
}

// OnEvent registers a typed handler for an event.
func OnEvent[T Event](eb *EventBus, handler func(ctx context.Context, event T) error) {
	var zero T
	eventName := zero.EventName()
	eb.On(eventName, func(ctx context.Context, event Event) error {
		typed, ok := event.(T)
		if !ok {
			return fmt.Errorf("events: unexpected event type %s for %s", reflect.TypeOf(event), eventName)
		}
		return handler(ctx, typed)
	})
}

// Emit publishes an event to all registered handlers.
func (eb *EventBus) Emit(ctx context.Context, event Event) error {
	eb.mu.RLock()
	handlers := eb.handlers[event.EventName()]
	eb.mu.RUnlock()

	if len(handlers) == 0 {
		return nil
	}

	if eb.async {
		for _, h := range handlers {
			go func(handler Handler) {
				if err := handler(ctx, event); err != nil {
					eb.logger.Error("event handler error",
						"event", event.EventName(),
						"error", err.Error(),
					)
				}
			}(h)
		}
		return nil
	}

	// Synchronous execution.
	for _, h := range handlers {
		if err := h(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

// EmitAsync publishes an event asynchronously regardless of bus setting.
func (eb *EventBus) EmitAsync(ctx context.Context, event Event) {
	eb.mu.RLock()
	handlers := eb.handlers[event.EventName()]
	eb.mu.RUnlock()

	for _, h := range handlers {
		go func(handler Handler) {
			if err := handler(ctx, event); err != nil {
				eb.logger.Error("async event handler error",
					"event", event.EventName(),
					"error", err.Error(),
				)
			}
		}(h)
	}
}

// HasHandlers returns true if there are handlers for the event.
func (eb *EventBus) HasHandlers(eventName string) bool {
	eb.mu.RLock()
	defer eb.mu.RUnlock()
	return len(eb.handlers[eventName]) > 0
}

// HandlerCount returns the number of handlers for an event.
func (eb *EventBus) HandlerCount(eventName string) int {
	eb.mu.RLock()
	defer eb.mu.RUnlock()
	return len(eb.handlers[eventName])
}

// Clear removes all handlers.
func (eb *EventBus) Clear() {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	eb.handlers = make(map[string][]Handler)
}
