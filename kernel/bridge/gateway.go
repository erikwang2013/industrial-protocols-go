package bridge

import (
	"net"

	"github.com/erikwang2013/industrial-protocols-go/kernel/transport"
)

// GatewayBridge connects to TCP-based gateway hardware.
type GatewayBridge struct {
	conn    net.Conn
	addr    string
	running bool
}

// NewGatewayBridge creates a GatewayBridge for the given TCP address.
func NewGatewayBridge(address string) *GatewayBridge {
	return &GatewayBridge{addr: address}
}

// Start dials the TCP gateway and marks the bridge as running.
func (b *GatewayBridge) Start() error {
	conn, err := net.Dial("tcp", b.addr)
	if err != nil {
		return err
	}
	b.conn = conn
	b.running = true
	return nil
}

// Stop closes the underlying connection.
func (b *GatewayBridge) Stop() error {
	b.running = false
	if b.conn != nil {
		return b.conn.Close()
	}
	return nil
}

// Transport returns a transport that wraps the TCP connection.
// It returns nil if Start has not been called.
func (b *GatewayBridge) Transport() transport.Transport {
	if b.conn == nil {
		return nil
	}
	conn := b.conn
	return &gatewayTransport{conn: conn, addr: b.addr}
}

// gatewayTransport adapts a net.Conn to the transport.Transport interface.
type gatewayTransport struct {
	conn net.Conn
	addr string
}

func (t *gatewayTransport) Read(p []byte) (int, error)  { return t.conn.Read(p) }
func (t *gatewayTransport) Write(p []byte) (int, error) { return t.conn.Write(p) }
func (t *gatewayTransport) Close() error                { return t.conn.Close() }
func (t *gatewayTransport) Addr() string                { return t.addr }
func (t *gatewayTransport) Alive() bool                 { return true }
