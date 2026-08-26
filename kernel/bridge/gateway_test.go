// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package bridge

import (
	"net"
	"testing"
)

func TestGatewayBridgeEcho(t *testing.T) {
	ln, _ := net.Listen("tcp", "127.0.0.1:0")
	defer ln.Close()
	go func() {
		conn, _ := ln.Accept()
		buf := make([]byte, 64)
		n, _ := conn.Read(buf)
		conn.Write(buf[:n])
	}()

	b := NewGatewayBridge(ln.Addr().String())
	if err := b.Start(); err != nil {
		t.Fatal(err)
	}
	defer b.Stop()

	tr := b.Transport()
	tr.Write([]byte("test"))
	buf := make([]byte, 64)
	n, _ := tr.Read(buf)
	if string(buf[:n]) != "test" {
		t.Errorf("expected 'test', got '%s'", string(buf[:n]))
	}
}

func TestGatewayBridgeInterface(t *testing.T) {
	var b Bridge = NewGatewayBridge("127.0.0.1:9999")
	if b == nil {
		t.Fatal("Bridge interface not satisfied")
	}
}

func TestGatewayBridge_StartError(t *testing.T) {
	b := NewGatewayBridge("not-an-address")
	if err := b.Start(); err == nil {
		t.Error("expected error for invalid address")
	}
}

func TestGatewayBridge_TransportNilBeforeStart(t *testing.T) {
	b := NewGatewayBridge("127.0.0.1:9999")
	if tr := b.Transport(); tr != nil {
		t.Error("expected nil transport before Start")
	}
}

func TestGatewayBridge_StopWithoutStart(t *testing.T) {
	b := NewGatewayBridge("127.0.0.1:9999")
	if err := b.Stop(); err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func TestGatewayBridge_TransportMetadata(t *testing.T) {
	ln, _ := net.Listen("tcp", "127.0.0.1:0")
	defer ln.Close()
	go func() {
		conn, _ := ln.Accept()
		conn.Close()
	}()

	b := NewGatewayBridge(ln.Addr().String())
	if err := b.Start(); err != nil {
		t.Fatal(err)
	}
	defer b.Stop()

	tr := b.Transport()
	if tr.Addr() != ln.Addr().String() {
		t.Errorf("expected addr %s, got %s", ln.Addr().String(), tr.Addr())
	}
	if !tr.Alive() {
		t.Error("transport should be alive after Start")
	}
	if err := tr.Close(); err != nil {
		t.Fatal(err)
	}
}
