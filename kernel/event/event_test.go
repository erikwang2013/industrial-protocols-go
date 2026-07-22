// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package event

import (
	"testing"
	"time"
)

func TestBusEmitSubscribe(t *testing.T) {
	bus := NewBus(10)
	ch := bus.Subscribe()

	bus.Emit(Event{Type: "connect", Device: "plc-001", Timestamp: time.Now()})

	select {
	case e := <-ch:
		if e.Type != "connect" {
			t.Errorf("expected connect, got %s", e.Type)
		}
		if e.Device != "plc-001" {
			t.Errorf("expected plc-001, got %s", e.Device)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout waiting for event")
	}
}

func TestBusMultipleSubscribers(t *testing.T) {
	bus := NewBus(5)
	ch1 := bus.Subscribe()
	ch2 := bus.Subscribe()

	bus.Emit(Event{Type: "error", Device: "x"})

	for _, ch := range []<-chan Event{ch1, ch2} {
		select {
		case e := <-ch:
			if e.Type != "error" {
				t.Errorf("expected error, got %s", e.Type)
			}
		case <-time.After(100 * time.Millisecond):
			t.Fatal("timeout")
		}
	}
}
