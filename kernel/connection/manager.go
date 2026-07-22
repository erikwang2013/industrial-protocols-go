// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package connection

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/erikwang2013/industrial-protocols-go/kernel/transport"
)

// ConnectionManager 管理设备连接的生命周期。
// 支持三种策略：Lazy（按需连接）、Eager（注册时连接）、Pooled（连接池复用）。
type ConnectionManager struct {
	configs  map[string]*DeviceConfig
	pools    map[string]*ConnPool
	strategy ConnStrategy
	mu       sync.RWMutex
}

// NewManager 创建连接管理器。strategy 为 nil 时默认使用 LazyStrategy。
func NewManager(strategy ConnStrategy) *ConnectionManager {
	if strategy == nil {
		strategy = &LazyStrategy{}
	}
	return &ConnectionManager{
		configs:  make(map[string]*DeviceConfig),
		pools:    make(map[string]*ConnPool),
		strategy: strategy,
	}
}

// Register 注册一个设备配置。
func (m *ConnectionManager) Register(cfg *DeviceConfig) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if cfg.Name == "" {
		return fmt.Errorf("connection: device name required")
	}
	if cfg.Network == "" {
		cfg.Network = "tcp"
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 5 * time.Second
	}

	if _, ok := m.configs[cfg.Name]; ok {
		return fmt.Errorf("connection: device %q already registered", cfg.Name)
	}

	if cfg.Pool != nil && cfg.Pool.MaxSize > 0 {
		deviceAddr := cfg.Addr
		deviceNetwork := cfg.Network
		deviceTimeout := cfg.Timeout
		pool := NewConnPool(func() (net.Conn, error) {
			return m.strategy.Dial(deviceNetwork, deviceAddr, deviceTimeout)
		}, cfg.Pool.MaxSize)
		m.pools[cfg.Name] = pool
	}

	m.configs[cfg.Name] = cfg
	return nil
}

// Acquire 获取一个就绪的 Transport。池化策略时从连接池获取，其他策略时新建连接。
func (m *ConnectionManager) Acquire(name string) (transport.Transport, error) {
	m.mu.RLock()
	cfg, ok := m.configs[name]
	m.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("connection: unknown device %q", name)
	}

	if pool, ok := m.pools[name]; ok {
		conn, err := pool.Get()
		if err != nil {
			return nil, fmt.Errorf("connection: pool acquire %q: %w", name, err)
		}
		return &poolTransport{conn: conn, pool: pool, addr: cfg.Addr}, nil
	}

	conn, err := m.strategy.Dial(cfg.Network, cfg.Addr, cfg.Timeout)
	if err != nil {
		return nil, fmt.Errorf("connection: dial %q: %w", name, err)
	}
	return &simpleTransport{conn: conn, addr: cfg.Addr, alive: true}, nil
}

// Release 归还 Transport。池化策略时放回池中，其他策略时关闭连接。
func (m *ConnectionManager) Release(t transport.Transport) {
	if pt, ok := t.(*poolTransport); ok {
		pt.pool.Put(pt.conn)
		return
	}
	t.Close()
}

// Health 健康检查，返回设备的健康状态和延迟。
func (m *ConnectionManager) Health(name string) HealthStatus {
	start := time.Now()
	tr, err := m.Acquire(name)
	if err != nil {
		return HealthStatus{Healthy: false}
	}
	defer m.Release(tr)
	return HealthStatus{Healthy: true, LatencyMs: float64(time.Since(start).Microseconds()) / 1000.0}
}

// Shutdown 关闭所有连接池。
func (m *ConnectionManager) Shutdown() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, pool := range m.pools {
		pool.Close()
	}
}

type poolTransport struct {
	conn net.Conn
	pool *ConnPool
	addr string
}

func (t *poolTransport) Read(p []byte) (int, error)  { return t.conn.Read(p) }
func (t *poolTransport) Write(p []byte) (int, error) { return t.conn.Write(p) }
func (t *poolTransport) Close() error                { t.pool.Put(t.conn); return nil }
func (t *poolTransport) Addr() string                { return t.addr }
func (t *poolTransport) Alive() bool                 { return t.conn != nil }

type simpleTransport struct {
	conn  net.Conn
	addr  string
	alive bool
}

func (t *simpleTransport) Read(p []byte) (int, error)  { return t.conn.Read(p) }
func (t *simpleTransport) Write(p []byte) (int, error) { return t.conn.Write(p) }
func (t *simpleTransport) Close() error                { t.alive = false; return t.conn.Close() }
func (t *simpleTransport) Addr() string                { return t.addr }
func (t *simpleTransport) Alive() bool                 { return t.alive }
