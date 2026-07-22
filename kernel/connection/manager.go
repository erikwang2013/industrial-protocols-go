package connection

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/erikwang2013/industrial-protocols-go/kernel/transport"
)

type ConnectionManager struct {
	configs  map[string]*DeviceConfig
	pools    map[string]*ConnPool
	strategy ConnStrategy
	mu       sync.RWMutex
}

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

func (m *ConnectionManager) Release(t transport.Transport) {
	if pt, ok := t.(*poolTransport); ok {
		pt.pool.Put(pt.conn)
		return
	}
	t.Close()
}

func (m *ConnectionManager) Health(name string) HealthStatus {
	start := time.Now()
	tr, err := m.Acquire(name)
	if err != nil {
		return HealthStatus{Healthy: false}
	}
	defer m.Release(tr)
	return HealthStatus{Healthy: true, LatencyMs: float64(time.Since(start).Microseconds()) / 1000.0}
}

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
