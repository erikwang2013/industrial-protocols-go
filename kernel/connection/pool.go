// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package connection

import (
	"errors"
	"net"
	"sync"
)

// ErrPoolClosed 表示连接池已关闭。
var ErrPoolClosed = errors.New("connection: pool closed")

// ErrPoolExhausted 表示连接池已耗尽，无可用连接。
var ErrPoolExhausted = errors.New("connection: pool exhausted")

// ConnPool 基于 channel 的通用连接池。
// 使用工厂函数创建新连接，支持最大连接数限制和平滑关闭。
type ConnPool struct {
	mu      sync.Mutex
	conns   chan net.Conn
	factory func() (net.Conn, error)
	maxSize int
	active  int
	closed  bool
}

// NewConnPool 创建连接池。
func NewConnPool(factory func() (net.Conn, error), maxSize int) *ConnPool {
	if maxSize < 1 {
		maxSize = 1
	}
	return &ConnPool{
		conns:   make(chan net.Conn, maxSize),
		factory: factory,
		maxSize: maxSize,
	}
}

func (p *ConnPool) Get() (net.Conn, error) {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return nil, ErrPoolClosed
	}
	p.mu.Unlock()

	select {
	case conn, ok := <-p.conns:
		if ok && conn != nil {
			return conn, nil
		}
	default:
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return nil, ErrPoolClosed
	}

	if p.active >= p.maxSize {
		return nil, ErrPoolExhausted
	}

	conn, err := p.factory()
	if err != nil {
		return nil, err
	}
	p.active++
	return conn, nil
}

func (p *ConnPool) Put(conn net.Conn) {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		conn.Close()
		return
	}
	p.mu.Unlock()

	select {
	case p.conns <- conn:
	default:
		p.mu.Lock()
		p.active--
		p.mu.Unlock()
		conn.Close()
	}
}

func (p *ConnPool) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return nil
	}
	p.closed = true
	close(p.conns)
	for conn := range p.conns {
		conn.Close()
	}
	return nil
}
