package services

import (
	"messages/core/entity"
	"sync"
	"time"
)

const maxEvents = 1000

type EventStore interface {
	Publish(event entity.ConsumerEvent)
	Subscribe() chan entity.ConsumerEvent
	Unsubscribe(ch chan entity.ConsumerEvent)
	Clear()
	History() []entity.ConsumerEvent
}

type EventStoreImpl struct {
	mu        sync.RWMutex
	history   []entity.ConsumerEvent
	clients   []chan entity.ConsumerEvent
	clearedAt time.Time
}

func NewEventStore() *EventStoreImpl {
	return &EventStoreImpl{
		history: make([]entity.ConsumerEvent, 0, maxEvents),
		clients: make([]chan entity.ConsumerEvent, 0),
	}
}

func (s *EventStoreImpl) Publish(event entity.ConsumerEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Discard events that were buffered before the last Clear
	if !s.clearedAt.IsZero() && event.ReceivedAt.Before(s.clearedAt) {
		return
	}

	if len(s.history) >= maxEvents {
		s.history = s.history[1:]
	}
	s.history = append(s.history, event)

	for _, ch := range s.clients {
		select {
		case ch <- event:
		default:
			// drop event if client is too slow
		}
	}
}

func (s *EventStoreImpl) Subscribe() chan entity.ConsumerEvent {
	ch := make(chan entity.ConsumerEvent, 64)
	s.mu.Lock()
	s.clients = append(s.clients, ch)
	s.mu.Unlock()
	return ch
}

func (s *EventStoreImpl) Unsubscribe(ch chan entity.ConsumerEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, c := range s.clients {
		if c == ch {
			s.clients = append(s.clients[:i], s.clients[i+1:]...)
			close(ch)
			return
		}
	}
}

func (s *EventStoreImpl) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.history = s.history[:0]
	s.clearedAt = time.Now()
}

func (s *EventStoreImpl) History() []entity.ConsumerEvent {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]entity.ConsumerEvent, len(s.history))
	copy(result, s.history)
	return result
}
