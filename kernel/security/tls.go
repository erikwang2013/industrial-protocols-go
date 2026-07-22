// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package security

import (
	"crypto/tls"
	"crypto/x509"
	"net"
	"os"
	"time"

	"github.com/erikwang2013/industrial-protocols-go/kernel/transport"
)

// TLSConfig builds a tls.Config from PEM-encoded certificate files.
// caFile is optional; pass an empty string for server-only authentication.
func TLSConfig(certFile, keyFile, caFile string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	cfg := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}
	if caFile != "" {
		caCert, err := os.ReadFile(caFile)
		if err != nil {
			return nil, err
		}
		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(caCert)
		cfg.RootCAs = pool
	}
	return cfg, nil
}

// netConnAdapter wraps a transport.Transport to satisfy net.Conn so it can
// be passed to tls.Client.
type netConnAdapter struct {
	transport.Transport
}

func (a *netConnAdapter) LocalAddr() net.Addr                { return stringAddr(a.Addr()) }
func (a *netConnAdapter) RemoteAddr() net.Addr               { return stringAddr(a.Addr()) }
func (a *netConnAdapter) SetDeadline(t time.Time) error      { return nil }
func (a *netConnAdapter) SetReadDeadline(t time.Time) error  { return nil }
func (a *netConnAdapter) SetWriteDeadline(t time.Time) error { return nil }

type stringAddr string

func (s stringAddr) Network() string { return "tcp" }
func (s stringAddr) String() string  { return string(s) }

// secureTransport decorates an underlying Transport with TLS encryption.
type secureTransport struct {
	inner transport.Transport
	conn  *tls.Conn
}

// NewSecureTransport performs a TLS client handshake over inner and returns
// a Transport that encrypts all reads and writes.
func NewSecureTransport(inner transport.Transport, cfg *tls.Config) (transport.Transport, error) {
	adapter := &netConnAdapter{inner}
	conn := tls.Client(adapter, cfg)
	if err := conn.Handshake(); err != nil {
		return nil, err
	}
	return &secureTransport{inner: inner, conn: conn}, nil
}

func (s *secureTransport) Read(p []byte) (int, error)  { return s.conn.Read(p) }
func (s *secureTransport) Write(p []byte) (int, error)  { return s.conn.Write(p) }
func (s *secureTransport) Close() error                 { return s.conn.Close() }
func (s *secureTransport) Addr() string                 { return s.inner.Addr() }
func (s *secureTransport) Alive() bool                  { return s.inner.Alive() }

var _ transport.Transport = (*secureTransport)(nil)
