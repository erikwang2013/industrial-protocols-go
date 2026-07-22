// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package connection

import (
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
