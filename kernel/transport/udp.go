// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package transport

import (
	"net"
	"sync"
)

// UDPTransport implements Transport over a UDP connection.
type UDPTransport struct {
	conn   net.Conn
	addr   string
	mu     sync.Mutex
	closed bool
}

// DialUDP opens a UDP connection to address and returns a Transport.
func DialUDP(address string) (*UDPTransport, error) {
	conn, err := net.Dial("udp", address)
	if err != nil {
		return nil, err
	}
	return &UDPTransport{conn: conn, addr: address}, nil
}

func (t *UDPTransport) Read(p []byte) (int, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.closed {
		return 0, net.ErrClosed
	}
	return t.conn.Read(p)
}

func (t *UDPTransport) Write(p []byte) (int, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.closed {
		return 0, net.ErrClosed
	}
	return t.conn.Write(p)
}

func (t *UDPTransport) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.closed {
		return nil
	}
	t.closed = true
	return t.conn.Close()
}

func (t *UDPTransport) Addr() string {
	return t.addr
}

func (t *UDPTransport) Alive() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return !t.closed
}
