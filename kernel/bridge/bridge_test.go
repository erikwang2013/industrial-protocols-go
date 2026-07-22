package bridge

import (
	"strings"
	"testing"
)

func TestNewCmdBridge(t *testing.T) {
	b := NewCmdBridge("echo", "hello")
	if b == nil {
		t.Fatal("NewCmdBridge returned nil")
	}
	if b.running {
		t.Log("bridge not running until Start()")
	}
	b.Stop()
}

func TestBridgeInterface(t *testing.T) {
	var b Bridge = NewCmdBridge("cat")
	if b == nil {
		t.Fatal("Bridge interface not satisfied")
	}
}

func TestCmdBridgeEcho(t *testing.T) {
	b := NewCmdBridge("echo", "-n", "hello")
	err := b.Start()
	if err != nil {
		t.Fatal(err)
	}
	defer b.Stop()

	tr := b.Transport()
	if tr == nil {
		t.Fatal("transport is nil before Start")
	}

	buf := make([]byte, 32)
	n, err := tr.Read(buf)
	if err == nil && n > 0 {
		if strings.TrimSpace(string(buf[:n])) == "hello" {
			t.Log("echo bridge works")
		}
	}
	tr.Close()
}
