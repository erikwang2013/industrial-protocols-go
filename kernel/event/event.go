// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package event

import (
	"sync"
	"time"
)

// Event 事件数据结构。
type Event struct {
	Type      string
	Device    string
	Timestamp time.Time
	Data      map[string]any
}

// Bus 事件总线。基于 channel 实现，支持多订阅者 fan-out 广播。
type Bus struct {
	mu   sync.RWMutex
	subs []chan Event
	buf  int
}

// NewBus 创建带缓冲的事件总线。
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
