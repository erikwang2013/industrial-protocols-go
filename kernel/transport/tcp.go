// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package transport

import (
	"net"
	"sync"
)

// TCPTransport implements Transport over a TCP connection.
type TCPTransport struct {
	conn   net.Conn
	addr   string
	mu     sync.Mutex
	closed bool
}

// DialTCP opens a TCP connection to address and returns a Transport.
func DialTCP(address string) (*TCPTransport, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	return &TCPTransport{conn: conn, addr: address}, nil
}

func (t *TCPTransport) Read(p []byte) (int, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.closed {
		return 0, net.ErrClosed
	}
	return t.conn.Read(p)
}

func (t *TCPTransport) Write(p []byte) (int, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.closed {
		return 0, net.ErrClosed
	}
	return t.conn.Write(p)
}

func (t *TCPTransport) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.closed {
		return nil
	}
	t.closed = true
	return t.conn.Close()
}

func (t *TCPTransport) Addr() string {
	return t.addr
}

func (t *TCPTransport) Alive() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return !t.closed
}
