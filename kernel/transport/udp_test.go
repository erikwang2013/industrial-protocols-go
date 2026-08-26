// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package transport

import (
	"errors"
	"net"
	"testing"
)

func TestUDPTransport_Echo(t *testing.T) {
	ln, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	go func() {
		buf := make([]byte, 64)
		for {
			n, addr, err := ln.ReadFrom(buf)
			if err != nil {
				return
			}
			ln.WriteTo(buf[:n], addr)
		}
	}()

	tr, err := DialUDP(ln.LocalAddr().String())
	if err != nil {
		t.Fatal(err)
	}

	if tr.Addr() != ln.LocalAddr().String() {
		t.Errorf("expected addr %s, got %s", ln.LocalAddr().String(), tr.Addr())
	}
	if !tr.Alive() {
		t.Error("transport should be alive after dial")
	}

	if _, err := tr.Write([]byte("hi")); err != nil {
		t.Fatal(err)
	}
	buf := make([]byte, 64)
	n, err := tr.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	if string(buf[:n]) != "hi" {
		t.Errorf("expected 'hi', got %q", string(buf[:n]))
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
	if _, err := tr.Read(buf); !errors.Is(err, net.ErrClosed) {
		t.Errorf("read after close: expected net.ErrClosed, got %v", err)
	}
	if _, err := tr.Write(buf); !errors.Is(err, net.ErrClosed) {
		t.Errorf("write after close: expected net.ErrClosed, got %v", err)
	}
}

func TestUDPTransport_DialError(t *testing.T) {
	if _, err := DialUDP("not-an-address"); err == nil {
		t.Error("expected error for invalid address")
	}
}
