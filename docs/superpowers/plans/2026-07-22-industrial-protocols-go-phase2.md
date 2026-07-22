# Industrial Protocols Go — Phase 2 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add 8 kernel modules (ConnectionManager, ConfigRepository, Gateway, Bridge, Vendor, Event, Metrics, Security) and 11 full protocol implementations (BACnet, EtherNet/IP, PROFINET, CC-Link, HART, HART-IP, DALI, LIN, K-Line, DNP3, IEC 61850).

**Architecture:** Each kernel module is an independent Go subpackage under `kernel/`. Each protocol is an independent Go module under `protocols/`. Kernel modules that depend on others (ConnectionManager → Event) are ordered appropriately.

**Tech Stack:** Go 1.20+, stdlib only for kernel modules; `gopkg.in/yaml.v3` for ConfigRepository.

---

## Phase 2A: Foundation Kernel Modules (Event + Metrics + Security)

### Task 1: Implement Event Bus

**Files:**
- Create: `kernel/event/event.go`
- Create: `kernel/event/event_test.go`

- [ ] **Step 1: Write the failing test**

Create `kernel/event/event_test.go`:

```go
package event

import (
	"testing"
	"time"
)

func TestBusEmitSubscribe(t *testing.T) {
	bus := NewBus(10)
	ch := bus.Subscribe()

	bus.Emit(Event{Type: "connect", Device: "plc-001", Timestamp: time.Now()})

	select {
	case e := <-ch:
		if e.Type != "connect" {
			t.Errorf("expected connect, got %s", e.Type)
		}
		if e.Device != "plc-001" {
			t.Errorf("expected plc-001, got %s", e.Device)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout waiting for event")
	}
}

func TestBusMultipleSubscribers(t *testing.T) {
	bus := NewBus(5)
	ch1 := bus.Subscribe()
	ch2 := bus.Subscribe()

	bus.Emit(Event{Type: "error", Device: "x"})

	for _, ch := range []<-chan Event{ch1, ch2} {
		select {
		case e := <-ch:
			if e.Type != "error" {
				t.Errorf("expected error, got %s", e.Type)
			}
		case <-time.After(100 * time.Millisecond):
			t.Fatal("timeout")
		}
	}
}

func TestEventStruct(t *testing.T) {
	e := Event{Type: "connect", Device: "d1", Timestamp: time.Now(), Data: map[string]any{"latency_ms": 5}}
	if e.Data["latency_ms"].(int) != 5 {
		t.Error("data field mismatch")
	}
}
```

- [ ] **Step 2: Implement event.go**

```go
package event

import "time"

type Event struct {
	Type      string
	Device    string
	Timestamp time.Time
	Data      map[string]any
}

type Bus chan Event

func NewBus(buffer int) Bus { return make(chan Event, buffer) }

func (b Bus) Subscribe() <-chan Event { return b }

func (b Bus) Emit(e Event) {
	select {
	case b <- e:
	default:
	}
}
```

- [ ] **Step 3: Run tests and commit**

```bash
cd kernel && go test ./event/ -v
git add kernel/event/ && git commit -m "feat(kernel): add Event Bus (channel-based, zero-deps)"
```

---

### Task 2: Implement Metrics

**Files:**
- Create: `kernel/metrics/metrics.go`
- Create: `kernel/metrics/metrics_test.go`

- [ ] **Step 1: Implement**

Create `kernel/metrics/metrics.go`:

```go
package metrics

type Recorder interface {
	Inc(name string, labels map[string]string)
	Observe(name string, value float64, labels map[string]string)
}

type noopRecorder struct{}

func (n noopRecorder) Inc(name string, labels map[string]string)           {}
func (n noopRecorder) Observe(name string, v float64, labels map[string]string) {}

func NoopRecorder() Recorder { return noopRecorder{} }
```

Create `kernel/metrics/metrics_test.go`:

```go
package metrics

import "testing"

func TestNoopRecorder(t *testing.T) {
	r := NoopRecorder()
	r.Inc("test", map[string]string{"k": "v"})
	r.Observe("test", 1.0, nil)
}

func TestRecorderInterface(t *testing.T) {
	var r Recorder = NoopRecorder()
	if r == nil {
		t.Error("NoopRecorder returned nil")
	}
}
```

- [ ] **Step 2: Commit**

```bash
cd kernel && go test ./metrics/ -v
git add kernel/metrics/ && git commit -m "feat(kernel): add Metrics Recorder interface with noop implementation"
```

---

### Task 3: Implement Security (TLS wrapper)

**Files:**
- Create: `kernel/security/tls.go`
- Create: `kernel/security/tls_test.go`

- [ ] **Step 1: Implement tls.go**

```go
package security

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"os"

	"github.com/erikwang2013/industrial-protocols-go/kernel/transport"
)

func TLSConfig(certFile, keyFile, caFile string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	cfg := &tls.Config{Certificates: []tls.Certificate{cert}, MinVersion: tls.VersionTLS12}
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

type secureTransport struct {
	inner transport.Transport
	conn  *tls.Conn
}

func NewSecureTransport(inner transport.Transport, cfg *tls.Config) (transport.Transport, error) {
	conn := tls.Client(inner, cfg)
	if err := conn.Handshake(); err != nil {
		return nil, err
	}
	return &secureTransport{inner: inner, conn: conn}, nil
}

func (s *secureTransport) Read(p []byte) (int, error)  { return s.conn.Read(p) }
func (s *secureTransport) Write(p []byte) (int, error) { return s.conn.Write(p) }
func (s *secureTransport) Close() error                { return s.conn.Close() }
func (s *secureTransport) Addr() string                { return s.inner.Addr() }
func (s *secureTransport) Alive() bool                 { return s.inner.Alive() }

var _ io.ReadWriteCloser = &secureTransport{}
```

Create `kernel/security/tls_test.go`:

```go
package security

import "testing"

func TestTLSConfig_MissingFiles(t *testing.T) {
	_, err := TLSConfig("/nonexistent.crt", "/nonexistent.key", "")
	if err == nil {
		t.Error("expected error for missing cert files")
	}
}
```

- [ ] **Step 2: Commit**

```bash
cd kernel && go test ./security/ -v && go build ./security/
git add kernel/security/ && git commit -m "feat(kernel): add TLS config factory and SecureTransport wrapper"
```

---

## Phase 2B: ConnectionManager + ConfigRepository

### Task 4: Implement ConnPool

**Files:**
- Create: `kernel/connection/pool.go`
- Create: `kernel/connection/pool_test.go`

- [ ] **Step 1: Write failing test**

Create `kernel/connection/pool_test.go`:

```go
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
			c, _ := ln.Accept()
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
			c, _ := ln.Accept()
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
```

- [ ] **Step 2: Implement pool.go**

```go
package connection

import (
	"errors"
	"net"
	"sync"
)

var ErrPoolClosed = errors.New("connection: pool closed")

type ConnPool struct {
	mu      sync.Mutex
	conns   chan net.Conn
	factory func() (net.Conn, error)
	maxSize int
	closed  bool
}

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
	select {
	case conn := <-p.conns:
		return conn, nil
	default:
	}
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return nil, ErrPoolClosed
	}
	p.mu.Unlock()
	return p.factory()
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
```

- [ ] **Step 3: Run tests and commit**

```bash
cd kernel && go test ./connection/ -v -short
git add kernel/connection/pool.go kernel/connection/pool_test.go
git commit -m "feat(kernel): add channel-based ConnPool"
```

---

### Task 5: Implement ConnectionManager + Strategies

**Files:**
- Create: `kernel/connection/strategy.go`
- Create: `kernel/connection/manager.go`
- Create: `kernel/connection/manager_test.go`

- [ ] **Step 1: Implement strategy.go**

```go
package connection

import (
	"net"
	"time"
)

type ConnStrategy interface {
	Dial(network, addr string, timeout time.Duration) (net.Conn, error)
}

type LazyStrategy struct{}
func (s *LazyStrategy) Dial(network, addr string, timeout time.Duration) (net.Conn, error) {
	return net.DialTimeout(network, addr, timeout)
}

type EagerStrategy struct{}
func (s *EagerStrategy) Dial(network, addr string, timeout time.Duration) (net.Conn, error) {
	return net.DialTimeout(network, addr, timeout)
}

type PooledStrategy struct {
	DefaultPoolSize int
}
func (s *PooledStrategy) Dial(network, addr string, timeout time.Duration) (net.Conn, error) {
	return net.DialTimeout(network, addr, timeout)
}

type DeviceConfig struct {
	Name    string
	Network string
	Addr    string
	Timeout time.Duration
	Pool    *PoolConfig
}

type PoolConfig struct {
	MaxSize     int
	IdleTimeout time.Duration
}

type HealthStatus struct {
	Healthy   bool
	LatencyMs float64
}
```

- [ ] **Step 2: Implement manager.go**

```go
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
```

- [ ] **Step 3: Write test**

Create `kernel/connection/manager_test.go`:

```go
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
			c, _ := ln.Accept()
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
			conn, _ := ln.Accept()
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
```

- [ ] **Step 4: Run tests and commit**

```bash
cd kernel && go test ./connection/ -v -short
git add kernel/connection/
git commit -m "feat(kernel): add ConnectionManager with Lazy/Eager/Pooled strategies"
```

---

### Task 6: Implement ConfigRepository

**Files:**
- Create: `kernel/config/config.go`
- Create: `kernel/config/config_test.go`
- Create: `kernel/config/testdata/example.yaml`

- [ ] **Step 1: Create test data**

Create `kernel/config/testdata/example.yaml`:

```yaml
devices:
  plc-001:
    name: plc-001
    network: tcp
    addr: "192.168.1.10:502"
    timeout: 3s
    pool:
      max_size: 5
      idle_timeout: 60s
  sensor-01:
    name: sensor-01
    network: tcp
    addr: "192.168.1.11:1883"
    timeout: 5s
```

- [ ] **Step 2: Implement config.go**

```go
package config

import (
	"os"
	"time"

	"github.com/erikwang2013/industrial-protocols-go/kernel/connection"
	"gopkg.in/yaml.v3"
)

type configFile struct {
	Devices map[string]deviceYAML `yaml:"devices"`
}

type poolYAML struct {
	MaxSize     int    `yaml:"max_size"`
	IdleTimeout string `yaml:"idle_timeout"`
}

type deviceYAML struct {
	Name    string    `yaml:"name"`
	Network string    `yaml:"network"`
	Addr    string    `yaml:"addr"`
	Timeout string    `yaml:"timeout"`
	Pool    *poolYAML `yaml:"pool,omitempty"`
}

type ConfigRepository struct {
	devices map[string]*connection.DeviceConfig
}

func Load(path string) (*ConfigRepository, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cf configFile
	if err := yaml.Unmarshal(data, &cf); err != nil {
		return nil, err
	}
	repo := &ConfigRepository{devices: make(map[string]*connection.DeviceConfig)}
	for key, dy := range cf.Devices {
		dc := &connection.DeviceConfig{
			Name: dy.Name, Network: dy.Network, Addr: dy.Addr,
		}
		if dy.Timeout != "" {
			d, err := time.ParseDuration(dy.Timeout)
			if err != nil {
				return nil, err
			}
			dc.Timeout = d
		}
		if dy.Pool != nil {
			idleTimeout := 60 * time.Second
			if dy.Pool.IdleTimeout != "" {
				idleTimeout, _ = time.ParseDuration(dy.Pool.IdleTimeout)
			}
			dc.Pool = &connection.PoolConfig{
				MaxSize: dy.Pool.MaxSize, IdleTimeout: idleTimeout,
			}
		}
		repo.devices[key] = dc
	}
	return repo, nil
}

func (r *ConfigRepository) Get(name string) (*connection.DeviceConfig, bool) {
	cfg, ok := r.devices[name]
	return cfg, ok
}

func (r *ConfigRepository) All() map[string]*connection.DeviceConfig {
	result := make(map[string]*connection.DeviceConfig, len(r.devices))
	for k, v := range r.devices {
		result[k] = v
	}
	return result
}
```

- [ ] **Step 3: Write test**

Create `kernel/config/config_test.go`:

```go
package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadYAML(t *testing.T) {
	path := filepath.Join("testdata", "example.yaml")
	repo, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	cfg, ok := repo.Get("plc-001")
	if !ok {
		t.Fatal("plc-001 not found")
	}
	if cfg.Addr != "192.168.1.10:502" {
		t.Errorf("expected 192.168.1.10:502, got %s", cfg.Addr)
	}
	if cfg.Pool == nil || cfg.Pool.MaxSize != 5 {
		t.Error("expected pool config")
	}
	_, ok = repo.Get("nonexistent")
	if ok {
		t.Error("should not find nonexistent device")
	}
}

func TestLoadMissingFile(t *testing.T) {
	_, err := Load("/nonexistent.yaml")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestTempYAML(t *testing.T) {
	content := []byte("devices:\n  test:\n    name: test\n    network: tcp\n    addr: \"127.0.0.1:502\"\n")
	tmpfile := filepath.Join(t.TempDir(), "config.yaml")
	os.WriteFile(tmpfile, content, 0644)
	repo, err := Load(tmpfile)
	if err != nil {
		t.Fatal(err)
	}
	cfg, ok := repo.Get("test")
	if !ok {
		t.Fatal("test device not found")
	}
	if cfg.Addr != "127.0.0.1:502" {
		t.Errorf("wrong addr: %s", cfg.Addr)
	}
}
```

- [ ] **Step 4: Add yaml dep, test, commit**

```bash
cd kernel && go get gopkg.in/yaml.v3
cd kernel && go test ./config/ -v
git add kernel/config/ kernel/go.mod kernel/go.sum
git commit -m "feat(kernel): add ConfigRepository with YAML loader"
```

---

## Phase 2C: Gateway + Bridge + Vendor (3 tasks)

### Task 7: Implement GatewayEngine

**Files:**
- Create: `kernel/gateway/rule.go`, `kernel/gateway/engine.go`, `kernel/gateway/engine_test.go`

Implement as specified in the Phase 2 design doc. Full code identical to spec.

- [ ] **Step 1: Implement, test, commit**

```bash
cd kernel && go test ./gateway/ -v && go build ./gateway/
git add kernel/gateway/ && git commit -m "feat(kernel): add GatewayEngine for protocol-to-protocol transformation"
```

---

### Task 8: Implement Bridge + PipeTransport

**Files:**
- Create: `kernel/transport/pipe.go`, `kernel/bridge/bridge.go`, `kernel/bridge/bridge_test.go`

Implement PipeTransport and CmdBridge as specified in the Phase 2 design doc.

- [ ] **Step 1: Implement, test, commit**

```bash
cd kernel && go test ./bridge/ -v -short && go build ./bridge/
git add kernel/bridge/ kernel/transport/pipe.go
git commit -m "feat(kernel): add Bridge interface and CmdBridge for external process bridging"
```

---

### Task 9: Implement Vendor Registry

**Files:**
- Create: `kernel/vendor/vendor.go`, `siemens.go`, `rockwell.go`, `vendor_test.go`

Implement as specified: VendorProfile, Registry, Siemens S7 presets, Rockwell/AB presets.

- [ ] **Step 1: Implement, test, commit**

```bash
cd kernel && go test ./vendor/ -v
git add kernel/vendor/ && git commit -m "feat(kernel): add Vendor Registry with Siemens and Rockwell presets"
```

---

## Phase 2D: Low-Complexity Protocol Implementations (5 tasks)

### Task 10: Implement HART + HART-IP

**Files:** Rewrite `protocols/fieldbus/hart/` and `protocols/iot/hartip/`

HART FSK codec: short/long frame with preamble, delimiter, address, command, checksum. Commands 0 (Read Unique ID) and 3 (Read Dynamic Variables). HART-IP wraps HART frames with 6-byte header.

- [ ] **Step 1: Implement both, test, commit**

```bash
cd protocols/fieldbus/hart && go test ./... -v
cd protocols/iot/hartip && go test ./... -v
git add protocols/fieldbus/hart/ protocols/iot/hartip/
git commit -m "feat(protocols): implement HART FSK and HART-IP codecs"
```

---

### Task 11: Implement DALI

**Files:** Rewrite `protocols/building/dali/`

DALI 16-bit forward frame (address byte + command byte) and 8-bit backward frame (status byte). Standard commands: off, max, recall_min, dim_up, dim_down.

- [ ] **Step 1: Implement, test, commit**

```bash
cd protocols/building/dali && go test ./... -v
git add protocols/building/dali/
git commit -m "feat(protocols): implement DALI codec (forward/backward frames, standard commands)"
```

---

### Task 12: Implement LIN

**Files:** Rewrite `protocols/automotive/lin/`

LIN master header (sync break + 0x55 + PID with parity) + slave response (data + classic checksum). Enhanced checksum (LIN 2.0+) also supported.

- [ ] **Step 1: Implement, test, commit**

```bash
cd protocols/automotive/lin && go test ./... -v
git add protocols/automotive/lin/
git commit -m "feat(protocols): implement LIN master/slave frame codec with classic checksum"
```

---

### Task 13: Implement K-Line

**Files:** Rewrite `protocols/automotive/kline/`

ISO 9141/14230 OBD-II over serial at 10400 baud. Fast init with 5-baud wakeup pattern (0x33). OBD-II SIDs: 0x01 (Current Data), 0x03 (DTCs), 0x09 (Vehicle Info).

- [ ] **Step 1: Implement, test, commit**

```bash
cd protocols/automotive/kline && go test ./... -v
git add protocols/automotive/kline/
git commit -m "feat(protocols): implement K-Line ISO 9141/14230 OBD-II codec"
```

---

## Phase 2E: Medium-Complexity Protocol Implementations (4 tasks)

### Task 14: Implement BACnet

**Files:** Rewrite `protocols/ethernet/bacnet/`

BACnet/IP over UDP 47808. BVLC (1-byte type) + NPDU (version+control) + APDU (service choice). Services: Who-Is (global broadcast), I-Am (device response), ReadProperty (object read).

- [ ] **Step 1: Implement, test with known wire-format vectors, commit**

```bash
cd protocols/ethernet/bacnet && go test ./... -v
git add protocols/ethernet/bacnet/
git commit -m "feat(protocols): implement BACnet/IP codec with Who-Is/I-Am and ReadProperty"
```

---

### Task 15: Implement EtherNet/IP

**Files:** Rewrite `protocols/ethernet/ethernetip/`

EtherNet/IP over TCP 44818. ENIP RegisterSession/UnRegisterSession. CIP Read Tag Service (Service 0x4C). Basic tag path encoding.

- [ ] **Step 1: Implement, test, commit**

```bash
cd protocols/ethernet/ethernetip && go test ./... -v
git add protocols/ethernet/ethernetip/
git commit -m "feat(protocols): implement EtherNet/IP ENIP+CIP codec"
```

---

### Task 16: Implement PROFINET

**Files:** Rewrite `protocols/ethernet/profinet/`

PROFINET NRT over UDP 34964. DCP Identify (device discovery) and Set (name/IP assignment). Record Data Read/Write for IO data access.

- [ ] **Step 1: Implement, test, commit**

```bash
cd protocols/ethernet/profinet && go test ./... -v
git add protocols/ethernet/profinet/
git commit -m "feat(protocols): implement PROFINET DCP+Record Data codec"
```

---

### Task 17: Implement CC-Link

**Files:** Rewrite `protocols/fieldbus/cclink/`

CC-Link over RS-485. Master-slave polling frames. CRC-16/XMODEM checksum (polynomial 0x1021, initial 0x0000, no XOR).

- [ ] **Step 1: Implement, test, commit**

```bash
cd protocols/fieldbus/cclink && go test ./... -v
git add protocols/fieldbus/cclink/
git commit -m "feat(protocols): implement CC-Link RS-485 codec with CRC-16/XMODEM"
```

---

## Phase 2F: High-Complexity Protocol Implementations (2 tasks)

### Task 18: Implement DNP3

**Files:** Rewrite `protocols/fieldbus/dnp3/`

DNP3 over TCP 20000. Transport layer (TPDU segmentation/reassembly with FIR/FIN/sequence). Application layer (APDU with function codes: Confirm 0x00, Read 0x01, Response 0x81). Class 0 static object reads. CRC-16/DNP variant.

- [ ] **Step 1: Implement transport+application layer, test, commit**

```bash
cd protocols/fieldbus/dnp3 && go test ./... -v
git add protocols/fieldbus/dnp3/
git commit -m "feat(protocols): implement DNP3 transport/application layer codec"
```

---

### Task 19: Implement IEC 61850

**Files:** Rewrite `protocols/fieldbus/iec61850/`

IEC 61850 MMS over TCP 102. MMS Initiate Request/Response, Conclude Request/Response. Read/Write variable with basic data types: Boolean, Integer, Unsigned, Float, Visible String, Octet String.

- [ ] **Step 1: Implement MMS codec, test, commit**

```bash
cd protocols/fieldbus/iec61850 && go test ./... -v
git add protocols/fieldbus/iec61850/
git commit -m "feat(protocols): implement IEC 61850 MMS codec with Initiate/Read/Write"
```

---

## Phase 2G: Final Verification

### Task 20: Full workspace verification

- [ ] **Step 1: Build and vet all**

```bash
cd /home/wwwroot/bag/industrial-protocols-go
go work sync
cd kernel && go build ./... && go vet ./...
for dir in $(find protocols -name go.mod -exec dirname {} \;); do
  (cd "$dir" && go build ./... && go vet ./...) || echo "FAIL: $dir"
done
```

- [ ] **Step 2: Run all tests**

```bash
cd kernel && go test ./... -v -short
for dir in $(find protocols -name go.mod -exec dirname {} \;); do
  echo "=== $dir ===" && (cd "$dir" && go test ./... -v -short) || echo "FAIL: $dir"
done
```

- [ ] **Step 3: Commit final state**

```bash
git add -A && git status
git commit -m "chore: Phase 2 complete — 8 kernel modules + 11 protocol implementations, all tests pass"
```

---

## Summary

| Phase | Tasks | Files (est.) | Description |
|-------|-------|-------------|-------------|
| 2A: Foundation | 1-3 | ~8 | Event, Metrics, Security |
| 2B: Connection | 4-6 | ~10 | ConnPool, ConnectionManager, ConfigRepository |
| 2C: Engine | 7-9 | ~10 | Gateway, Bridge, Vendor |
| 2D: Low protocols | 10-13 | ~8 | HART, HART-IP, DALI, LIN, K-Line |
| 2E: Medium protocols | 14-17 | ~8 | BACnet, EtherNet/IP, PROFINET, CC-Link |
| 2F: High protocols | 18-19 | ~4 | DNP3, IEC 61850 |
| 2G: Verify | 20 | — | Full build + test suite |
| **Total** | **20 tasks** | **~50 files** | |
