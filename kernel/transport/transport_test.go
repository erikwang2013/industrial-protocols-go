// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package transport

import (
	"net"
	"testing"
)

func TestTCPTransportIntegration(t *testing.T) {
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
		buf := make([]byte, 64)
		n, _ := conn.Read(buf)
		conn.Write(buf[:n])
		conn.Close()
	}()

	tr, err := DialTCP(ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer tr.Close()

	if !tr.Alive() {
		t.Error("transport should be alive after dial")
	}

	data := []byte("hello")
	if _, err := tr.Write(data); err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 64)
	n, err := tr.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	if string(buf[:n]) != "hello" {
		t.Errorf("expected 'hello', got '%s'", string(buf[:n]))
	}

	if err := tr.Close(); err != nil {
		t.Fatal(err)
	}
	if tr.Alive() {
		t.Error("transport should not be alive after close")
	}
}
