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
