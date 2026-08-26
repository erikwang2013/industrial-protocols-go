// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package transport

import (
	"io"
	"testing"
)

func TestPipeTransport_RoundTrip(t *testing.T) {
	pr1, pw1 := io.Pipe()
	pr2, pw2 := io.Pipe()

	pt := NewPipeTransport(pr2, pw1, "proc")
	if !pt.Alive() {
		t.Error("transport should be alive before close")
	}
	if pt.Addr() != "proc" {
		t.Errorf("expected addr 'proc', got %s", pt.Addr())
	}

	// Echo: whatever is written to the transport comes back through its reader.
	go io.Copy(pw2, pr1)

	if _, err := pt.Write([]byte("ping")); err != nil {
		t.Fatal(err)
	}
	buf := make([]byte, 64)
	n, err := pt.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	if string(buf[:n]) != "ping" {
		t.Errorf("expected 'ping', got %q", string(buf[:n]))
	}

	if err := pt.Close(); err != nil {
		t.Fatal(err)
	}
	if pt.Alive() {
		t.Error("transport should not be alive after close")
	}
	if err := pt.Close(); err != nil {
		t.Errorf("double close should return nil, got %v", err)
	}
}
