// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package bridge

import (
	"io"
	"os"
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

func TestCmdBridge_StartError(t *testing.T) {
	b := NewCmdBridge("/nonexistent-cmd-xyz")
	if err := b.Start(); err == nil {
		t.Error("expected error for missing binary")
	}
}

func TestCmdBridge_TransportNilBeforeStart(t *testing.T) {
	b := NewCmdBridge("echo")
	if tr := b.Transport(); tr != nil {
		t.Error("expected nil transport before Start")
	}
}

func TestCmdBridge_StopWithoutStart(t *testing.T) {
	b := NewCmdBridge("echo")
	if err := b.Stop(); err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

// TestCmdBridge_RoundTrip spawns the test binary as a child process that
// echoes stdin to stdout, and verifies the PipeTransport round-trip.
func TestCmdBridge_RoundTrip(t *testing.T) {
	b := NewCmdBridge(os.Args[0], "-test.run=TestBridgeHelperProcess", "--")
	b.cmd.Env = append(os.Environ(), "BRIDGE_HELPER=1")
	if err := b.Start(); err != nil {
		t.Fatal(err)
	}
	defer b.Stop()

	tr := b.Transport()
	if tr == nil {
		t.Fatal("transport is nil after Start")
	}
	if !tr.Alive() {
		t.Error("transport should be alive after Start")
	}

	if _, err := tr.Write([]byte("ping")); err != nil {
		t.Fatal(err)
	}
	buf := make([]byte, 64)
	n, err := tr.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	if string(buf[:n]) != "ping" {
		t.Errorf("expected 'ping', got %q", string(buf[:n]))
	}
}

// TestBridgeHelperProcess is the child-process echo helper. It only runs when
// BRIDGE_HELPER=1 is set by TestCmdBridge_RoundTrip.
func TestBridgeHelperProcess(t *testing.T) {
	if os.Getenv("BRIDGE_HELPER") != "1" {
		return
	}
	io.Copy(os.Stdout, os.Stdin)
	os.Exit(0)
}
