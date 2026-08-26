// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package security

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/erikwang2013/industrial-protocols-go/kernel/transport"
)

func TestTLSConfig_MissingFiles(t *testing.T) {
	_, err := TLSConfig("/nonexistent.crt", "/nonexistent.key", "")
	if err == nil {
		t.Error("expected error for missing cert files")
	}
}

// generateCert creates an in-memory self-signed TLS certificate for localhost.
func generateCert(t *testing.T) tls.Certificate {
	t.Helper()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "localhost"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
	}

	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}

	return tls.Certificate{
		Certificate: [][]byte{der},
		PrivateKey:  key,
	}
}

func TestSecureTransport_Echo(t *testing.T) {
	cert := generateCert(t)
	cfg := &tls.Config{Certificates: []tls.Certificate{cert}, MinVersion: tls.VersionTLS12}

	ln, err := tls.Listen("tcp", "127.0.0.1:0", cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	errCh := make(chan error, 1)
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			errCh <- err
			return
		}
		buf := make([]byte, 64)
		n, _ := conn.Read(buf)
		conn.Write(buf[:n])
		conn.Close()
		errCh <- nil
	}()

	tr, err := transport.DialTCP(ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	tlsCfg := &tls.Config{InsecureSkipVerify: true, MinVersion: tls.VersionTLS12}
	sec, err := NewSecureTransport(tr, tlsCfg)
	if err != nil {
		t.Fatal(err)
	}
	defer sec.Close()

	if !sec.Alive() {
		t.Error("secure transport should be alive after handshake")
	}

	data := []byte("hello-tls")
	if _, err := sec.Write(data); err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 64)
	n, err := sec.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	if string(buf[:n]) != "hello-tls" {
		t.Errorf("expected 'hello-tls', got '%s'", string(buf[:n]))
	}

	if err := sec.Close(); err != nil {
		t.Fatal(err)
	}

	// Wait for the background accept goroutine.
	select {
	case err := <-errCh:
		if err != nil {
			t.Error(err)
		}
	case <-time.After(time.Second):
		t.Error("timeout waiting for accept goroutine")
	}
}

// writePEMFiles writes the generated certificate and key to temp files.
func writePEMFiles(t *testing.T) (certFile, keyFile string) {
	t.Helper()
	cert := generateCert(t)
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert.Certificate[0]})
	keyDER, err := x509.MarshalECPrivateKey(cert.PrivateKey.(*ecdsa.PrivateKey))
	if err != nil {
		t.Fatal(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER})

	dir := t.TempDir()
	certFile = filepath.Join(dir, "cert.pem")
	keyFile = filepath.Join(dir, "key.pem")
	if err := os.WriteFile(certFile, certPEM, 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(keyFile, keyPEM, 0600); err != nil {
		t.Fatal(err)
	}
	return certFile, keyFile
}

func TestTLSConfig_ValidFiles(t *testing.T) {
	certFile, keyFile := writePEMFiles(t)
	cfg, err := TLSConfig(certFile, keyFile, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.Certificates) != 1 {
		t.Errorf("expected 1 certificate, got %d", len(cfg.Certificates))
	}
	if cfg.MinVersion != tls.VersionTLS12 {
		t.Errorf("expected MinVersion TLS1.2, got %d", cfg.MinVersion)
	}
	if cfg.RootCAs != nil {
		t.Error("expected nil RootCAs when caFile is empty")
	}
}

func TestTLSConfig_WithCA(t *testing.T) {
	certFile, keyFile := writePEMFiles(t)
	// Reuse the self-signed certificate as its own CA.
	cfg, err := TLSConfig(certFile, keyFile, certFile)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.RootCAs == nil {
		t.Error("expected RootCAs pool when caFile is provided")
	}
}

func TestTLSConfig_MissingCA(t *testing.T) {
	certFile, keyFile := writePEMFiles(t)
	if _, err := TLSConfig(certFile, keyFile, "/nonexistent-ca.pem"); err == nil {
		t.Error("expected error for missing CA file")
	}
}

func TestTLSConfig_InvalidPEM(t *testing.T) {
	dir := t.TempDir()
	certFile := filepath.Join(dir, "cert.pem")
	keyFile := filepath.Join(dir, "key.pem")
	if err := os.WriteFile(certFile, []byte("not a pem"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(keyFile, []byte("not a pem"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := TLSConfig(certFile, keyFile, ""); err == nil {
		t.Error("expected error for invalid PEM data")
	}
}

// stubTransport is a minimal transport.Transport for metadata tests.
type stubTransport struct{ addr string }

func (s *stubTransport) Read(p []byte) (int, error)  { return 0, nil }
func (s *stubTransport) Write(p []byte) (int, error) { return 0, nil }
func (s *stubTransport) Close() error                { return nil }
func (s *stubTransport) Addr() string                { return s.addr }
func (s *stubTransport) Alive() bool                 { return true }

func TestNetConnAdapter(t *testing.T) {
	a := &netConnAdapter{&stubTransport{addr: "1.2.3.4:502"}}
	if a.LocalAddr().Network() != "tcp" || a.LocalAddr().String() != "1.2.3.4:502" {
		t.Errorf("unexpected LocalAddr: %v", a.LocalAddr())
	}
	if a.RemoteAddr().String() != "1.2.3.4:502" {
		t.Errorf("unexpected RemoteAddr: %v", a.RemoteAddr())
	}
	if err := a.SetDeadline(time.Time{}); err != nil {
		t.Errorf("SetDeadline: %v", err)
	}
	if err := a.SetReadDeadline(time.Time{}); err != nil {
		t.Errorf("SetReadDeadline: %v", err)
	}
	if err := a.SetWriteDeadline(time.Time{}); err != nil {
		t.Errorf("SetWriteDeadline: %v", err)
	}
}

func TestSecureTransport_Metadata(t *testing.T) {
	s := &secureTransport{inner: &stubTransport{addr: "10.0.0.1:502"}}
	if s.Addr() != "10.0.0.1:502" {
		t.Errorf("expected addr 10.0.0.1:502, got %s", s.Addr())
	}
	if !s.Alive() {
		t.Error("expected alive")
	}
}

func TestNewSecureTransport_HandshakeFailure(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		conn.Close()
	}()

	tr, err := transport.DialTCP(ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer tr.Close()

	tlsCfg := &tls.Config{InsecureSkipVerify: true, MinVersion: tls.VersionTLS12}
	if _, err := NewSecureTransport(tr, tlsCfg); err == nil {
		t.Error("expected handshake failure against a non-TLS server")
	}
}
