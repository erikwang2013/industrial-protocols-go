package event

import (
	"sync"
	"time"
)

// Event is a device- or system-level message dispatched through the Bus.
type Event struct {
	Type      string
	Device    string
	Timestamp time.Time
	Data      map[string]any
}

// Bus is a channel-based, zero-dependency event bus that supports
// multiple subscribers with non-blocking broadcast emission.
type Bus struct {
	mu   sync.RWMutex
	subs []chan Event
	buf  int
}

// NewBus creates a Bus where each subscriber channel is buffered
// with the given capacity.
func NewBus(buffer int) *Bus {
	return &Bus{buf: buffer}
}

// Subscribe registers a new subscriber and returns a read-only
// channel on which events will be delivered.
func (b *Bus) Subscribe() <-chan Event {
	b.mu.Lock()
	defer b.mu.Unlock()
	ch := make(chan Event, b.buf)
	b.subs = append(b.subs, ch)
	return ch
}

// Emit broadcasts e to all subscribers. If a subscriber's buffer
// is full the event is dropped for that subscriber (non-blocking).
func (b *Bus) Emit(e Event) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, ch := range b.subs {
		select {
		case ch <- e:
		default:
		}
	}
}
