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
