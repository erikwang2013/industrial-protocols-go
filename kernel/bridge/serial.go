package bridge

import (
	"io"

	"github.com/erikwang2013/industrial-protocols-go/kernel/transport"
)

// SerialBridge is a Bridge implementation for serial port (UART/RS-232/RS-485) communication.
type SerialBridge struct {
	rw      io.ReadWriteCloser
	dev     string
	baud    int
	running bool
}

// NewSerialBridge creates a new SerialBridge for the given device and baud rate.
func NewSerialBridge(dev string, baud int) *SerialBridge {
	return &SerialBridge{dev: dev, baud: baud}
}

// Start marks the bridge as running. The actual port must be opened externally
// and injected via SetPort before Start is called.
func (b *SerialBridge) Start() error {
	b.running = true
	return nil
}

// Stop marks the bridge as stopped and closes the underlying port if open.
func (b *SerialBridge) Stop() error {
	b.running = false
	if b.rw != nil {
		return b.rw.Close()
	}
	return nil
}

// Transport returns a transport.Transport backed by the serial port I/O.
// Returns nil if no port has been set.
func (b *SerialBridge) Transport() transport.Transport {
	if b.rw == nil {
		return nil
	}
	return transport.NewPipeTransport(b.rw, b.rw, b.dev)
}

// SetPort injects an io.ReadWriteCloser that represents the open serial port.
func (b *SerialBridge) SetPort(rw io.ReadWriteCloser) { b.rw = rw }
