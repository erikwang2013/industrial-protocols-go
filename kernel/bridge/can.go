//go:build linux

package bridge

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/erikwang2013/industrial-protocols-go/kernel/transport"
)

// CAN_RAW protocol constant not exported by syscall.
const canRaw = 1 // syscall.CAN_RAW

// CANFrame represents a single CAN bus frame.
type CANFrame struct {
	ID   uint32
	Data []byte
	Ext  bool
}

// CANBridge is a Bridge implementation for Linux SocketCAN interfaces.
type CANBridge struct {
	iface   string
	fd      int
	running bool
}

// NewCANBridge creates a new CANBridge for the given SocketCAN interface (e.g. "vcan0", "can0").
func NewCANBridge(iface string) *CANBridge {
	return &CANBridge{iface: iface}
}

// Start opens a CAN_RAW socket and binds it to the named interface.
func (b *CANBridge) Start() error {
	fd, err := syscall.Socket(syscall.AF_CAN, syscall.SOCK_RAW, canRaw)
	if err != nil {
		return fmt.Errorf("can: socket: %w", err)
	}
	var addr [16]byte
	copy(addr[:], b.iface)
	if _, _, errno := syscall.Syscall(syscall.SYS_BIND, uintptr(fd), uintptr(unsafe.Pointer(&addr[0])), 0); errno != 0 {
		return fmt.Errorf("can: bind %s: %w", b.iface, errno)
	}
	b.fd = fd
	b.running = true
	return nil
}

// Stop closes the socket and marks the bridge as stopped.
func (b *CANBridge) Stop() error {
	b.running = false
	if b.fd > 0 {
		return syscall.Close(b.fd)
	}
	return nil
}

// Transport returns a transport.Transport backed by the CAN socket.
func (b *CANBridge) Transport() transport.Transport {
	return &canTransport{fd: b.fd, iface: b.iface}
}

type canTransport struct {
	fd    int
	iface string
}

func (t *canTransport) Read(p []byte) (int, error)  { return syscall.Read(t.fd, p) }
func (t *canTransport) Write(p []byte) (int, error) { return syscall.Write(t.fd, p) }
func (t *canTransport) Close() error                { return syscall.Close(t.fd) }
func (t *canTransport) Addr() string                { return t.iface }
func (t *canTransport) Alive() bool                 { return t.fd > 0 }
