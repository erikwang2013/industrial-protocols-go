// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package transport

import (
	"errors"
)

// ReadWriteCloser is the subset of io.ReadWriteCloser needed for a serial port.
type ReadWriteCloser interface {
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	Close() error
}

// SerialTransport implements Transport over a serial connection.
type SerialTransport struct {
	rw     ReadWriteCloser
	addr   string
	closed bool
}

// NewSerial wraps an existing ReadWriteCloser (e.g. a serial port) as a Transport.
func NewSerial(rwc ReadWriteCloser, addr string) (*SerialTransport, error) {
	if rwc == nil {
		return nil, errors.New("transport: nil serial port")
	}
	return &SerialTransport{rw: rwc, addr: addr}, nil
}

func (t *SerialTransport) Read(p []byte) (int, error) {
	if t.closed {
		return 0, errors.New("transport: serial port closed")
	}
	return t.rw.Read(p)
}

func (t *SerialTransport) Write(p []byte) (int, error) {
	if t.closed {
		return 0, errors.New("transport: serial port closed")
	}
	return t.rw.Write(p)
}

func (t *SerialTransport) Close() error {
	if t.closed {
		return nil
	}
	t.closed = true
	return t.rw.Close()
}

func (t *SerialTransport) Addr() string {
	return t.addr
}

func (t *SerialTransport) Alive() bool {
	return !t.closed
}
