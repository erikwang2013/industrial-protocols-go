// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package transport

import (
	"net"
	"testing"
)

func TestNewSerial_NilPort(t *testing.T) {
	if _, err := NewSerial(nil, "/dev/ttyS0"); err == nil {
		t.Error("expected error for nil serial port")
	}
}

func TestSerialTransport_Echo(t *testing.T) {
	c1, c2 := net.Pipe()
	go func() {
		buf := make([]byte, 64)
		for {
			n, err := c2.Read(buf)
			if err != nil {
				return
			}
			c2.Write(buf[:n])
		}
	}()

	tr, err := NewSerial(c1, "/dev/ttyS0")
	if err != nil {
		t.Fatal(err)
	}
	if tr.Addr() != "/dev/ttyS0" {
		t.Errorf("expected addr /dev/ttyS0, got %s", tr.Addr())
	}
	if !tr.Alive() {
		t.Error("transport should be alive before close")
	}

	if _, err := tr.Write([]byte("AT")); err != nil {
		t.Fatal(err)
	}
	buf := make([]byte, 64)
	n, err := tr.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	if string(buf[:n]) != "AT" {
		t.Errorf("expected 'AT', got %q", string(buf[:n]))
	}

	if err := tr.Close(); err != nil {
		t.Fatal(err)
	}
	if err := tr.Close(); err != nil {
		t.Errorf("double close should return nil, got %v", err)
	}
	if tr.Alive() {
		t.Error("transport should not be alive after close")
	}
	if _, err := tr.Read(buf); err == nil || err.Error() != "transport: serial port closed" {
		t.Errorf("read after close: expected closed error, got %v", err)
	}
	if _, err := tr.Write(buf); err == nil || err.Error() != "transport: serial port closed" {
		t.Errorf("write after close: expected closed error, got %v", err)
	}
}
