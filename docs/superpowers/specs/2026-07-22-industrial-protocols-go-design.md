# Industrial Protocols Go — Design Spec

Date: 2026-07-22
Reference: https://github.com/erikwang2013/industrial-protocols

## Overview

Go 语言工业通信协议集，采用 **分层 + 中间件** 架构。40 个协议以 Go workspace 多模块 monorepo 组织，用户按需 `go get`。

## Architecture

```
┌─────────────────────────────────────────────┐
│              User Application                │
├─────────────────────────────────────────────┤
│  Pipeline Middleware                          │
│  Retry · Timeout · CircuitBreaker · Logger   │
├─────────────────────────────────────────────┤
│  Codec (Encode/Decode)         Transport     │
│  每协议实现                      (TCP/UDP/Serial) │
├─────────────────────────────────────────────┤
│  40 Protocol Packages                        │
└─────────────────────────────────────────────┘
```

Kernel 只定义接口 + 提供可复用横切关注点。协议只需实现 `Protocol` + `Codec`。

## Project Structure

```
industrial-protocols-go/
├── go.work
├── kernel/                       # github.com/erikwang2013/industrial-protocols-go/kernel
│   ├── go.mod
│   ├── protocol.go               # Protocol 接口
│   ├── codec.go                  # Codec, Request, Response
│   ├── transport/
│   │   ├── transport.go          # Transport 接口
│   │   ├── tcp.go                # TCP 实现
│   │   ├── udp.go                # UDP 实现
│   │   └── serial.go             # 串口实现
│   ├── pipeline/
│   │   ├── pipeline.go           # Chain, Middleware, Handler
│   │   ├── retry.go              # linear / exponential / jitter
│   │   ├── timeout.go
│   │   ├── logger.go
│   │   └── breaker.go            # closed → open → half-open
│   └── errors.go                 # Sentinel errors + ProtocolError
│
├── protocols/
│   ├── ethernet/                  # 工业以太网 (5)
│   │   ├── modbus/
│   │   │   ├── go.mod
│   │   │   ├── modbus.go
│   │   │   ├── tcp.go
│   │   │   ├── rtu.go
│   │   │   └── modbus_test.go
│   │   ├── bacnet/
│   │   ├── ethernetip/
│   │   ├── opcua/
│   │   └── profinet/
│   ├── fieldbus/                  # 现场总线 (11)
│   │   ├── hart/ cclink/ dnp3/ iec61850/ profibus/
│   │   ├── canopen/ devicenet/ foundationfieldbus/
│   │   ├── asinterface/ iolink/ cclinkie/
│   ├── iot/                       # IoT/消息 (2)
│   │   ├── mqtt/ hartip/
│   ├── automotive/                # 汽车总线 (5)
│   │   ├── lin/ kline/ flexray/ saej1850/ most/
│   ├── building/                  # 楼宇/照明 (2)
│   │   ├── lonworks/ dali/
│   ├── bridge/                    # 硬件桥接 stub-only (13)
│   │   ├── ethercat/ powerlink/ sercos/ sercos1/
│   │   ├── controlnet/ interbus/ worldfip/ lightbus/
│   │   ├── modbusplus/ isa100/ wirelesshart/ safej1850/
│   └── system/                    # 系统总线 stub-only (3)
│       ├── pci/ vme/ cpci/
│
├── examples/                      # 使用示例
├── _tools/Makefile
├── README.md
└── LICENSE
```

## Core Interfaces

### Protocol

```go
type Protocol interface {
    Name()        string      // "modbus"
    Variants()    []string    // ["tcp", "rtu", "ascii"]
    DefaultPort() int         // 502
    NewCodec(variant string) (Codec, error)
}
```

### Codec

```go
type Codec interface {
    Encode(req *Request) ([]byte, error)
    Decode(data []byte) (*Response, error)
}
```

### Request / Response

```go
type Request struct {
    Function string         // "read_coils", "publish"
    Address  string         // "40001", "topic/foo"
    Count    int
    Data     []byte
    Metadata map[string]any // 协议特有字段
}

type Response struct {
    Address  string
    Data     []byte
    Metadata map[string]any
}
```

### Transport

```go
type Transport interface {
    io.ReadWriter
    io.Closer
    Addr()  string
    Alive() bool
}
```

内置实现：`*TCPTransport`, `*UDPTransport`, `*SerialTransport`。

## Pipeline

```go
type Handler func(ctx context.Context, req *Request) (*Response, error)
type Middleware func(next Handler) Handler

func Chain(middlewares ...Middleware) Middleware
```

内置中间件：

| Middleware | 说明 |
|------------|------|
| `Timeout(d)` | 上下文超时控制 |
| `Retry(max, backoff)` | linear / exponential / jitter 三种退避 |
| `CircuitBreaker(threshold, cooldown)` | 三态熔断 |
| `Logger(logger)` | 结构化日志 |

## Error Handling

```go
var (
    ErrTimeout         = errors.New("protocol: timeout")
    ErrCircuitOpen     = errors.New("protocol: circuit breaker open")
    ErrRetryExhausted  = errors.New("protocol: all retries exhausted")
    ErrTransportClosed = errors.New("protocol: transport closed")
    ErrInvalidAddress  = errors.New("protocol: invalid address")
)

type ProtocolError struct {
    Code    string // 协议原生错误码
    Message string
    Raw     []byte // 调试用原始帧
}
```

## Protocol Implementation Categories

| Category | Count | Phase 1 | Phase 2+ |
|----------|-------|---------|----------|
| Industrial Ethernet | 5 | Modbus, OPC UA, BACnet | EtherNet/IP, PROFINET |
| Fieldbus | 11 | HART, DNP3, IEC 61850 | CC-Link, PROFIBUS, CANopen, etc. |
| IoT/Messaging | 2 | MQTT | HART-IP |
| Automotive | 5 | LIN, K-Line | FlexRay, SAE J1850, MOST |
| Building/Lighting | 2 | DALI | LonWorks |
| Hardware Bridge | 13 | — (stub only) | EtherCAT, POWERLINK, SERCOS, etc. |
| System Bus | 3 | — (stub only) | PCI, VME, CPCI |

Phase 1 priority: Modbus TCP/RTU, MQTT, OPC UA (covers 80% of industrial + IoT use cases).

## Module Policy

- Go workspace (`go.work`) aggregates all modules
- Each protocol has its own `go.mod`: `github.com/erikwang2013/industrial-protocols-go/protocols/ethernet/modbus`
- Kernel: `github.com/erikwang2013/industrial-protocols-go/kernel`
- All modules require `go 1.20`

## Testing Strategy

| Layer | Type | Approach |
|-------|------|----------|
| Transport | Integration (`-tags=integration`) | Real TCP loopback / serial loopback |
| Codec | Unit | Table-driven: `[]byte` in → `Response` out; `Request` in → `[]byte` out |
| Pipeline | Unit | Mock Handler + table-driven for retry counts, timeout, breaker state transitions |

## Decisions Summary

| Dimension | Decision |
|-----------|----------|
| Package strategy | Go workspace monorepo, per-protocol `go.mod` |
| Delivery | All 40 protocol stubs first, then fill implementations |
| Go version | 1.20+ |
| Architecture | Layered + middleware |
| Core interfaces | Protocol + Codec + Transport (4 types) |
| Middleware | Retry / Timeout / CircuitBreaker / Logger |
| Connection management | Transport-level transparent pooling, no global manager |
| Error handling | Sentinel errors + `ProtocolError` |
| Testing | Table-driven unit + integration tagged |
