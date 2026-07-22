// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package security

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"net"
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
