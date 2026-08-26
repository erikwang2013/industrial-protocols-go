// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package connection

import (
	"net"
	"testing"
	"time"
)

func TestManagerRegisterAndHealth_Lazy(t *testing.T) {
	ln, _ := net.Listen("tcp", "127.0.0.1:0")
	defer ln.Close()
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			c.Close()
		}
	}()

	mgr := NewManager(&LazyStrategy{})
	err := mgr.Register(&DeviceConfig{
		Name: "test", Network: "tcp", Addr: ln.Addr().String(), Timeout: 1 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	h := mgr.Health("test")
	if !h.Healthy {
		t.Error("expected healthy device")
	}

	err = mgr.Register(&DeviceConfig{Name: "test", Network: "tcp", Addr: "x"})
	if err == nil {
		t.Error("expected duplicate registration error")
	}

	h2 := mgr.Health("nonexistent")
	if h2.Healthy {
		t.Error("expected unhealthy for unknown device")
	}
}

func TestManagerPooledAcquireRelease(t *testing.T) {
	ln, _ := net.Listen("tcp", "127.0.0.1:0")
	defer ln.Close()
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			buf := make([]byte, 64)
			conn.Read(buf)
			conn.Close()
		}
	}()

	mgr := NewManager(&PooledStrategy{})
	err := mgr.Register(&DeviceConfig{
		Name: "plc", Network: "tcp", Addr: ln.Addr().String(), Timeout: 1 * time.Second,
		Pool: &PoolConfig{MaxSize: 2},
	})
	if err != nil {
		t.Fatal(err)
	}

	tr1, err := mgr.Acquire("plc")
	if err != nil {
		t.Fatal(err)
	}
	mgr.Release(tr1)

	tr2, err := mgr.Acquire("plc")
	if err != nil {
		t.Fatal(err)
	}
	mgr.Release(tr2)

	mgr.Shutdown()
}

func TestManagerRegister_EmptyName(t *testing.T) {
	mgr := NewManager(nil)
	if err := mgr.Register(&DeviceConfig{}); err == nil {
		t.Error("expected error for empty device name")
	}
}

func TestManagerRegister_Defaults(t *testing.T) {
	mgr := NewManager(nil)
	if err := mgr.Register(&DeviceConfig{Name: "d"}); err != nil {
		t.Fatal(err)
	}
	cfg := mgr.configs["d"]
	if cfg.Network != "tcp" {
		t.Errorf("expected default network tcp, got %s", cfg.Network)
	}
	if cfg.Timeout != 5*time.Second {
		t.Errorf("expected default timeout 5s, got %v", cfg.Timeout)
	}
}

func TestManagerAcquire_UnknownDevice(t *testing.T) {
	mgr := NewManager(nil)
	if _, err := mgr.Acquire("nope"); err == nil {
		t.Error("expected error for unknown device")
	}
}

func TestManagerAcquire_DialFailure(t *testing.T) {
	ln, _ := net.Listen("tcp", "127.0.0.1:0")
	addr := ln.Addr().String()
	ln.Close() // free the port so the dial is refused

	mgr := NewManager(nil)
	mgr.Register(&DeviceConfig{Name: "dead", Network: "tcp", Addr: addr, Timeout: 100 * time.Millisecond})
	if _, err := mgr.Acquire("dead"); err == nil {
		t.Error("expected dial failure")
	}
}

func TestManagerRelease_NonPool(t *testing.T) {
	ln, _ := net.Listen("tcp", "127.0.0.1:0")
	defer ln.Close()
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			c.Close()
		}
	}()

	mgr := NewManager(&LazyStrategy{})
	mgr.Register(&DeviceConfig{Name: "dev", Network: "tcp", Addr: ln.Addr().String(), Timeout: time.Second})
	tr, err := mgr.Acquire("dev")
	if err != nil {
		t.Fatal(err)
	}
	mgr.Release(tr)
	if tr.Alive() {
		t.Error("transport should be closed after Release")
	}
}

func TestManagerShutdown_Idempotent(t *testing.T) {
	mgr := NewManager(nil)
	mgr.Register(&DeviceConfig{Name: "p", Network: "tcp", Addr: "127.0.0.1:1", Pool: &PoolConfig{MaxSize: 2}})
	mgr.Shutdown()
	mgr.Shutdown() // must not panic or block
}

// countingStrategy counts Dial calls to verify pool reuse.
type countingStrategy struct {
	dials int
}

func (s *countingStrategy) Dial(network, addr string, timeout time.Duration) (net.Conn, error) {
	s.dials++
	return net.DialTimeout(network, addr, timeout)
}

func TestManagerPooled_ReusesConnection(t *testing.T) {
	ln, _ := net.Listen("tcp", "127.0.0.1:0")
	defer ln.Close()
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			c.Close()
		}
	}()

	strategy := &countingStrategy{}
	mgr := NewManager(strategy)
	mgr.Register(&DeviceConfig{
		Name: "plc", Network: "tcp", Addr: ln.Addr().String(), Timeout: time.Second,
		Pool: &PoolConfig{MaxSize: 2},
	})

	tr1, err := mgr.Acquire("plc")
	if err != nil {
		t.Fatal(err)
	}
	mgr.Release(tr1) // back to the pool

	tr2, err := mgr.Acquire("plc")
	if err != nil {
		t.Fatal(err)
	}
	if !tr2.Alive() {
		t.Error("expected pooled transport alive")
	}
	mgr.Release(tr2)

	if strategy.dials != 1 {
		t.Errorf("expected 1 dial with pool reuse, got %d", strategy.dials)
	}
	mgr.Shutdown()
}

// echoServer accepts connections and echoes back whatever it reads.
func echoServer(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		go func(c net.Conn) {
			defer c.Close()
			buf := make([]byte, 64)
			for {
				n, err := c.Read(buf)
				if err != nil {
					return
				}
				c.Write(buf[:n])
			}
		}(conn)
	}
}

func TestManagerPooled_TransportIO(t *testing.T) {
	ln, _ := net.Listen("tcp", "127.0.0.1:0")
	defer ln.Close()
	go echoServer(ln)

	mgr := NewManager(&LazyStrategy{})
	mgr.Register(&DeviceConfig{
		Name: "plc", Network: "tcp", Addr: ln.Addr().String(), Timeout: time.Second,
		Pool: &PoolConfig{MaxSize: 2},
	})
	defer mgr.Shutdown()

	tr, err := mgr.Acquire("plc")
	if err != nil {
		t.Fatal(err)
	}
	if tr.Addr() != ln.Addr().String() {
		t.Errorf("expected addr %s, got %s", ln.Addr().String(), tr.Addr())
	}
	if !tr.Alive() {
		t.Error("expected transport alive")
	}
	if _, err := tr.Write([]byte("ping")); err != nil {
		t.Fatal(err)
	}
	buf := make([]byte, 64)
	n, err := tr.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	if string(buf[:n]) != "ping" {
		t.Errorf("expected 'ping', got %q", string(buf[:n]))
	}
	if err := tr.Close(); err != nil { // poolTransport.Close returns the conn to the pool
		t.Fatal(err)
	}
}

func TestManager_SimpleTransportIO(t *testing.T) {
	ln, _ := net.Listen("tcp", "127.0.0.1:0")
	defer ln.Close()
	go echoServer(ln)

	mgr := NewManager(&LazyStrategy{})
	mgr.Register(&DeviceConfig{
		Name: "dev", Network: "tcp", Addr: ln.Addr().String(), Timeout: time.Second,
	})

	tr, err := mgr.Acquire("dev")
	if err != nil {
		t.Fatal(err)
	}
	if tr.Addr() != ln.Addr().String() {
		t.Errorf("expected addr %s, got %s", ln.Addr().String(), tr.Addr())
	}
	if _, err := tr.Write([]byte("hi")); err != nil {
		t.Fatal(err)
	}
	buf := make([]byte, 64)
	n, err := tr.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	if string(buf[:n]) != "hi" {
		t.Errorf("expected 'hi', got %q", string(buf[:n]))
	}
	mgr.Release(tr)
	if tr.Alive() {
		t.Error("expected transport closed after Release")
	}
}
