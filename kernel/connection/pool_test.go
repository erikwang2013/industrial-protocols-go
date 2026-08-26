// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package connection

import (
	"errors"
	"net"
	"testing"
	"time"
)

func TestConnPoolGetPut(t *testing.T) {
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

	factory := func() (net.Conn, error) {
		return net.Dial("tcp", ln.Addr().String())
	}
	p := NewConnPool(factory, 2)
	defer p.Close()

	c1, err := p.Get()
	if err != nil {
		t.Fatal(err)
	}
	p.Put(c1)

	c2, err := p.Get()
	if err != nil {
		t.Fatal(err)
	}
	p.Put(c2)
}

func TestConnPoolFull(t *testing.T) {
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

	factory := func() (net.Conn, error) {
		return net.Dial("tcp", ln.Addr().String())
	}
	p := NewConnPool(factory, 1)
	defer p.Close()

	c1, err := p.Get()
	if err != nil {
		t.Fatal(err)
	}
	_, err = p.Get()
	if err == nil {
		t.Error("expected error when pool is exhausted")
	}
	p.Put(c1)
	c2, err := p.Get()
	if err != nil {
		t.Error("should get connection after put")
	}
	p.Put(c2)
}

func TestConnPoolClose(t *testing.T) {
	ln, _ := net.Listen("tcp", "127.0.0.1:0")
	defer ln.Close()
	go func() {
		c, _ := ln.Accept()
		time.Sleep(10 * time.Millisecond)
		c.Close()
	}()

	factory := func() (net.Conn, error) {
		return net.Dial("tcp", ln.Addr().String())
	}
	p := NewConnPool(factory, 1)
	c, _ := p.Get()
	p.Put(c)
	p.Close()

	_, err := p.Get()
	if err == nil {
		t.Error("expected error from closed pool")
	}
}

func TestConnPool_FactoryError(t *testing.T) {
	p := NewConnPool(func() (net.Conn, error) {
		return nil, errors.New("dial failed")
	}, 1)
	defer p.Close()

	_, err := p.Get()
	if err == nil || err.Error() != "dial failed" {
		t.Errorf("expected factory error, got %v", err)
	}
}

// recordingConn records Close calls.
type recordingConn struct {
	net.Conn
	closed bool
}

func (c *recordingConn) Close() error {
	c.closed = true
	return c.Conn.Close()
}

func TestConnPool_PutAfterCloseClosesConn(t *testing.T) {
	c1, _ := net.Pipe()
	rec := &recordingConn{Conn: c1}
	p := NewConnPool(func() (net.Conn, error) { return rec, nil }, 1)

	conn, err := p.Get()
	if err != nil {
		t.Fatal(err)
	}
	p.Close()
	p.Put(conn) // put into closed pool must close the conn
	if !rec.closed {
		t.Error("connection should be closed when returned to a closed pool")
	}
}

func TestConnPool_MaxSizeBelowOne(t *testing.T) {
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

	factory := func() (net.Conn, error) {
		return net.Dial("tcp", ln.Addr().String())
	}
	p := NewConnPool(factory, 0) // clamped to 1
	defer p.Close()

	c1, err := p.Get()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := p.Get(); !errors.Is(err, ErrPoolExhausted) {
		t.Errorf("expected ErrPoolExhausted, got %v", err)
	}
	p.Put(c1)
}

func TestEagerStrategy_Dial(t *testing.T) {
	ln, _ := net.Listen("tcp", "127.0.0.1:0")
	defer ln.Close()

	s := &EagerStrategy{}
	conn, err := s.Dial("tcp", ln.Addr().String(), time.Second)
	if err != nil {
		t.Fatal(err)
	}
	conn.Close()
}

func TestPooledStrategy_DialFailure(t *testing.T) {
	ln, _ := net.Listen("tcp", "127.0.0.1:0")
	addr := ln.Addr().String()
	ln.Close() // free the port so the dial is refused

	s := &PooledStrategy{DefaultPoolSize: 4}
	if s.DefaultPoolSize != 4 {
		t.Errorf("expected DefaultPoolSize 4, got %d", s.DefaultPoolSize)
	}
	if _, err := s.Dial("tcp", addr, 100*time.Millisecond); err == nil {
		t.Error("expected dial failure")
	}
}
