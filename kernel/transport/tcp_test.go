// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package transport

import (
	"errors"
	"net"
	"testing"
)

func TestTCPTransport_DialError(t *testing.T) {
	if _, err := DialTCP("not-an-address"); err == nil {
		t.Error("expected error for invalid address")
	}
}

func TestTCPTransport_ClosedBehavior(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		conn.Close()
	}()

	tr, err := DialTCP(ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	if tr.Addr() != ln.Addr().String() {
		t.Errorf("expected addr %s, got %s", ln.Addr().String(), tr.Addr())
	}
	if !tr.Alive() {
		t.Error("transport should be alive after dial")
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

	if _, err := tr.Read(make([]byte, 1)); !errors.Is(err, net.ErrClosed) {
		t.Errorf("read after close: expected net.ErrClosed, got %v", err)
	}
	if _, err := tr.Write([]byte("x")); !errors.Is(err, net.ErrClosed) {
		t.Errorf("write after close: expected net.ErrClosed, got %v", err)
	}
}
