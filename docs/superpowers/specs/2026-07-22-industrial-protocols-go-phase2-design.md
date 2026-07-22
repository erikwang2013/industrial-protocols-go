# Industrial Protocols Go — Phase 2 Design Spec

Date: 2026-07-22
Reference: https://github.com/erikwang2013/industrial-protocols
Phase 1 Spec: docs/superpowers/specs/2026-07-22-industrial-protocols-go-design.md

## Overview

Phase 2 补齐 PHP 参考项目中 Go 版本缺失的所有模块：8 个内核模块 + 11 个协议实现。

Phase 1 已交付：kernel (errors/codec/protocol/transport/pipeline)、40 协议 stub、Modbus/MQTT/OPC UA 三协议完整实现。

## Gap Analysis

### Kernel modules (8 new subpackages)

| Module | PHP Reference | Go Approach |
|--------|-------------|-------------|
| ConnectionManager | 3 strategies (Lazy/Eager/Pooled) | channel-based pool + strategy pattern |
| ConfigRepository | File config | YAML/JSON loader |
| GatewayEngine | Protocol transform rules | Rule-based engine |
| Bridge Layer | External process bridging | stdin/stdout io.ReadWriter |
| Vendor Profiles | Vendor parameter templates | Registry + preset structs |
| Event System | PSR-14 | channel-based EventBus |
| Metrics | Prometheus | Recorder interface (Prometheus + noop) |
| Security | TLS/certs | crypto/tls wrapper for Transport |

### Protocol implementations (11 upgrade from stub)

| Category | Protocol | Transport | Complexity | Key features |
|----------|----------|-----------|------------|--------------|
| **Medium** | BACnet | UDP 47808 | Medium | Who-Is/I-Am discovery, ReadProperty |
| **Medium** | EtherNet/IP | TCP 44818 | Medium | ENIP session + CIP Read Tag |
| **Medium** | PROFINET | UDP 34964 | Medium | DCP discovery + Record Data Read |
| **Medium** | CC-Link | RS-485 | Medium | Master-slave polling, CRC-16/XMODEM |
| **Low** | HART | Serial FSK | Low | HART command frames |
| **Low** | HART-IP | TCP/UDP 5094 | Low | HART over TCP, reuse HART codec |
| **Low** | DALI | Serial | Low | Short address commands |
| **Low** | LIN | UART | Low | Master/slave frame + checksum |
| **Low** | K-Line | Serial 10400 | Medium | ISO 9141/14230 init, OBD-II |
| **High** | DNP3 | TCP/Serial 20000 | High | Class 0 poll, time sync, app-layer frag |
| **High** | IEC 61850 | TCP 102 | High | MMS stack + data type system |

## Added File Structure

```
kernel/
├── connection/
│   ├── manager.go              # ConnectionManager
│   ├── strategy.go              # Lazy/Eager/Pooled strategies
│   ├── pool.go                  # channel-based ConnPool
│   └── manager_test.go
│
├── config/
│   ├── config.go                # ConfigRepository (YAML/JSON)
│   ├── config_test.go
│   └── testdata/example.yaml
│
├── gateway/
│   ├── engine.go                # GatewayEngine
│   ├── rule.go                  # Rule definition
│   └── engine_test.go
│
├── bridge/
│   ├── bridge.go                # Bridge interface + CmdBridge
│   └── bridge_test.go
│
├── vendor/
│   ├── vendor.go                # VendorProfile + Registry
│   ├── siemens.go               # Siemens S7 presets
│   ├── rockwell.go              # Rockwell/AB presets
│   └── vendor_test.go
│
├── event/
│   ├── event.go                 # Event types + Bus (channel)
│   └── event_test.go
│
├── metrics/
│   ├── metrics.go               # Recorder interface + Prometheus + noop
│   └── metrics_test.go
│
└── security/
    ├── tls.go                   # TLS config factory + SecureTransport
    └── tls_test.go
```

Protocols follow existing patterns (each is an independent Go module with protocol_name.go, test, optional discover.go/client.go).

## Core Module Design

### ConnectionManager

```go
type ConnStrategy interface {
    Dial(network, addr string) (net.Conn, error)
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

type ConnectionManager struct {
    configs  map[string]*DeviceConfig
    pools    map[string]*ConnPool
    strategy ConnStrategy
    events   event.Bus
    mu       sync.RWMutex
}

func (m *ConnectionManager) Register(cfg *DeviceConfig) error
func (m *ConnectionManager) Acquire(name string) (transport.Transport, error)
func (m *ConnectionManager) Release(t transport.Transport)
func (m *ConnectionManager) Health(name string) HealthStatus
```

### ConnPool

```go
type ConnPool struct {
    mu      sync.Mutex
    conns   chan net.Conn
    factory func() (net.Conn, error)
    maxSize int
}

func NewConnPool(factory func() (net.Conn, error), maxSize int) *ConnPool
func (p *ConnPool) Get() (net.Conn, error)
func (p *ConnPool) Put(conn net.Conn)
func (p *ConnPool) Close() error
```

### ConfigRepository

```go
type ConfigRepository struct {
    devices map[string]*connection.DeviceConfig
}

func Load(path string) (*ConfigRepository, error)
func (r *ConfigRepository) Get(name string) (*connection.DeviceConfig, bool)
```

### GatewayEngine

```go
type Rule struct {
    Src kernel.Protocol
    Dst kernel.Protocol
    Map func(srcResp *kernel.Response) *kernel.Request
}

type GatewayEngine struct {
    rules []Rule
}

func (e *GatewayEngine) Add(rule Rule)
func (e *GatewayEngine) Transform(ctx context.Context, src, dst transport.Transport, data *kernel.Request) (*kernel.Response, error)
```

### Bridge

```go
type Bridge interface {
    Start() error
    Stop() error
    Transport() transport.Transport
}

type CmdBridge struct { ... }

func NewCmdBridge(cmd string, args ...string) *CmdBridge
```

### Vendor

```go
type VendorProfile struct {
    Make     string
    Model    string
    Defaults map[string]any
}

type VendorRegistry struct { ... }

func (r *VendorRegistry) Register(p VendorProfile)
func (r *VendorRegistry) Find(make, model string) (*VendorProfile, bool)
```

### Event

```go
type Event struct {
    Type      string
    Device    string
    Timestamp time.Time
    Data      map[string]any
}

type Bus chan Event

func NewBus(buffer int) Bus
func (b Bus) Subscribe() <-chan Event
func (b Bus) Emit(e Event)
```

### Metrics

```go
type Recorder interface {
    Inc(name string, labels map[string]string)
    Observe(name string, value float64, labels map[string]string)
}

func PrometheusRecorder(namespace string) Recorder
func NoopRecorder() Recorder
```

### Security

```go
func TLSConfig(certFile, keyFile, caFile string) (*tls.Config, error)
func NewSecureTransport(base transport.Transport, config *tls.Config) transport.Transport
```

## Protocol Implementation Plan

| Protocol | Complexity | Features |
|----------|-----------|----------|
| BACnet | Medium | Who-Is/I-Am discovery, ReadProperty, WriteProperty |
| EtherNet/IP | Medium | ENIP RegisterSession, CIP Read Tag Service |
| PROFINET | Medium | DCP Identify, Record Data Read/Write |
| CC-Link | Medium | Master-slave polling, CRC-16/XMODEM |
| HART | Low | Command 0/3 (read), short/long frame codec |
| HART-IP | Low | HART over TCP, reuse HART codec |
| DALI | Low | 16-bit forward frame, 8-bit backward frame |
| LIN | Low | Master header + slave response, classic checksum |
| K-Line | Medium | Fast init 5-baud wakeup, ISO 9141/14230 codec |
| DNP3 | High | Transport/Application layer, Class 0 poll, time sync |
| IEC 61850 | High | MMS Initiate/Conclude, Read/Write variable, basic data types |

## Decisions Summary

| Dimension | Decision |
|-----------|----------|
| ConnectionPool | channel-based, sync.Pool-style |
| Config format | YAML primary, JSON secondary |
| Gateway | Rule engine (not middleware) |
| Bridge | CmdBridge via stdin/stdout |
| Event bus | Pure channel, zero deps |
| Metrics | Interface with Prometheus + noop |
| Security | crypto/tls wrapper, factory pattern |
| All protocols | Complete for low-complexity, core features for medium, framework for high |

## Total Scope

| Category | Count | New files (est.) |
|----------|-------|-------------------|
| Kernel modules | 8 subpackages | ~20 .go files |
| Protocol implementations | 11 upgrades | ~30-40 .go files |
| **Total** | | **~50-60 Go source files** |
