# Industrial Protocols Go

Go 语言工业网络通信协议集 —— 分层 + 中间件架构，覆盖 40 种工业协议，14 个完整实现 + 26 个骨架。

> 参考 PHP 实现: [github.com/erikwang2013/industrial-protocols](https://github.com/erikwang2013/industrial-protocols)

---

## 项目设计

### 架构分层

```
┌──────────────────────────────────────────────────┐
│                User Application                   │
├──────────────────────────────────────────────────┤
│  Pipeline Middleware                               │
│  Retry · Timeout · CircuitBreaker · Logger        │
├──────────────────────────────────────────────────┤
│  Kernel                                           │
│  ConnectionManager · ConfigRepository              │
│  GatewayEngine · Bridge · Vendor · Event          │
│  Metrics · Security (TLS)                         │
├──────────────────────────────────────────────────┤
│  Codec (Encode/Decode)           Transport        │
│  每协议独立模块                    (TCP/UDP/Serial)  │
│                                                │
│  14 Full Impl   +   26 Stub                      │
└──────────────────────────────────────────────────┘
```

### 设计思路

**微内核 + 协议 SDK。** 内核只定义接口和横切关注点（连接池、重试、熔断、事件、指标），不包含任何具体协议实现。每个协议是独立的 Go module，按需引入。

**分层解耦。** 四层抽象：

| 层 | 职责 | 核心类型 |
|----|------|---------|
| **Transport** | 底层通信信道 | `Transport` 接口 — `TCPTransport`, `UDPTransport`, `PipeTransport` |
| **Codec** | 协议编解码 | `Codec` 接口 — `Encode(*Request) ([]byte, error)` / `Decode([]byte) (*Response, error)` |
| **Pipeline** | 横切中间件 | `Middleware` — `Timeout`, `Retry`, `CircuitBreaker`, `Logger` |
| **Kernel** | 设备生命周期 | `ConnectionManager`, `ConfigRepository` |

**协议作者只需实现 `Protocol` + `Codec` 两个接口。** 传输层、连接池、重试、超时、熔断全部从 kernel 复用。

### 核心接口

```go
// kernel/protocol.go — 每个协议模块必须实现
type Protocol interface {
    Name()        string              // "modbus"
    Variants()    []string            // ["tcp", "rtu", "ascii"]
    DefaultPort() int                 // 502
    NewCodec(variant string) (Codec, error)
}

// kernel/codec.go — 编解码
type Codec interface {
    Encode(req *Request) ([]byte, error)
    Decode(data []byte) (*Response, error)
}

// kernel/transport/transport.go — 传输信道
type Transport interface {
    io.ReadWriter
    io.Closer
    Addr()  string
    Alive() bool
}

// kernel/pipeline/pipeline.go — 中间件
type Handler    func(ctx context.Context, req *Request) (*Response, error)
type Middleware func(next Handler) Handler
```

### 中间件链

```
Request → [Timeout] → [Retry] → [CircuitBreaker] → [Logger] → Codec.Encode → Transport.Write
Response ← Codec.Decode ← Transport.Read
```

内置中间件：

| 中间件 | 构造 | 说明 |
|--------|------|------|
| `Timeout` | `Timeout(d time.Duration)` | context 超时控制 |
| `Retry` | `Retry(max int, backoff BackoffFunc)` | linear / exponential / jitter 三种退避 |
| `CircuitBreaker` | `NewCircuitBreaker(threshold int, cooldown time.Duration)` | closed → open → half-open 三态熔断 |
| `Logger` | `Logger(logger *log.Logger)` | 结构化请求日志 |

### 内核模块

| 模块 | 路径 | 说明 |
|------|------|------|
| ConnectionManager | `kernel/connection/` | 设备注册、连接池、健康检查，支持 Lazy/Eager/Pooled 三种策略 |
| ConfigRepository | `kernel/config/` | YAML/JSON 设备配置加载 |
| GatewayEngine | `kernel/gateway/` | 跨协议转换规则引擎（Modbus→MQTT 等） |
| Bridge | `kernel/bridge/` | 外部进程桥接（stdin/stdout 通信） |
| Vendor | `kernel/vendor/` | 厂商参数预设（Siemens S7, Rockwell AB 等） |
| Event | `kernel/event/` | channel-based 事件总线 |
| Metrics | `kernel/metrics/` | 指标采集接口（Prometheus + noop） |
| Security | `kernel/security/` | TLS 传输层安全包装 |

---

## 协议支持

### 工业以太网（5/5 已完成）

| 协议 | 传输 | 端口 | 功能 |
|------|------|------|------|
| **Modbus** | TCP, RTU | 502 | FC 01/02/03/04/05/06/16, CRC16, MBAP |
| **BACnet** | UDP | 47808 | Who-Is/I-Am 设备发现, ReadProperty |
| **EtherNet/IP** | TCP | 44818 | ENIP RegisterSession + CIP Read Tag |
| **OPC UA** | TCP | 4840 | Binary HEL + OpenSecureChannel |
| **PROFINET** | UDP | 34964 | DCP Identify/Set + Record Data Read/Write |

### 现场总线（4/11 已完成）

| 协议 | 传输 | 端口 | 功能 |
|------|------|------|------|
| **HART** | Serial FSK | — | 短帧/长帧, Command 0/3, XOR 校验 |
| **CC-Link** | RS-485 | — | 主从轮询, CRC-16/XMODEM |
| **DNP3** | TCP, Serial | 20000 | 传输层分段/重组, Class 0 轮询, CRC-16/DNP |
| **IEC 61850** | TCP | 102 | MMS Initiate/Conclude, Read/Write, BER-TLV |
| PROFIBUS | Serial | — | *stub — 需 CP 5611 硬件* |
| CANopen | CAN | — | *stub — 需 CAN 接口* |
| DeviceNet | CAN | — | *stub — 需 DeviceNet 扫描器* |
| Foundation Fieldbus | Serial | — | *stub — 需 FF 接口* |
| AS-Interface | Serial | — | *stub — 需 ASi 网关* |
| IO-Link | Serial | — | *stub — 需 IO-Link Master* |
| CC-Link IE | Ethernet | — | *stub — 需网关* |

### IoT / 消息（2/2 已完成）

| 协议 | 传输 | 端口 | 功能 |
|------|------|------|------|
| **MQTT** | TCP, WS | 1883 | 3.1.1 CONNECT/PUBLISH/SUBSCRIBE/PING |
| **HART-IP** | TCP, UDP | 5094 | HART over TCP, 复用 HART 编解码 |

### 汽车总线（2/5 已完成）

| 协议 | 传输 | 波特率 | 功能 |
|------|------|--------|------|
| **LIN** | UART | — | 主从帧, PID 校验, Classic/Enhanced Checksum |
| **K-Line** | Serial | 10400 | ISO 9141/14230, 5-baud Fast Init, OBD-II SID 01/03/09 |
| FlexRay | Serial | — | *stub — 需 FlexRay 控制器* |
| SAE J1850 | Serial | — | *stub — 需 J1850 接口* |
| MOST | Optical | — | *stub — 需 MOST 接口* |

### 楼宇 / 照明（1/2 已完成）

| 协议 | 传输 | 功能 |
|------|------|------|
| **DALI** | Serial | 16-bit 前向帧, 8-bit 后向帧, 标准命令 (Off/Max/Dim) |
| LonWorks | Serial | *stub — 需 Neuron 芯片* |

### 硬件桥接（0/12 — 全部 stub）

EtherCAT, POWERLINK, SERCOS III, SERCOS I/II, ControlNet, Interbus, WorldFIP, Lightbus, Modbus Plus, ISA100, WirelessHART, SAE J1850 — *需特定硬件/FPGA/网关*

### 系统总线（0/3 — 全部 stub）

PCI/PCIe, VME/VPX, CompactPCI — *需内核驱动/桥接模块*

---

## 使用指南

### 安装

```bash
go get github.com/erikwang2013/industrial-protocols-go/kernel
go get github.com/erikwang2013/industrial-protocols-go/protocols/ethernet/modbus
go get github.com/erikwang2013/industrial-protocols-go/protocols/iot/mqtt
```

### Modbus TCP 读写

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/kernel/pipeline"
    "github.com/erikwang2013/industrial-protocols-go/kernel/transport"
    "github.com/erikwang2013/industrial-protocols-go/protocols/ethernet/modbus"
)

func main() {
    tr, _ := transport.DialTCP("192.168.1.10:502")
    defer tr.Close()

    codec, _ := modbus.NewCodec("tcp")
    handler := func(ctx context.Context, req *kernel.Request) (*kernel.Response, error) {
        req.Metadata["unit_id"] = byte(1)
        raw, _ := codec.Encode(req)
        tr.Write(raw)
        buf := make([]byte, 256)
        n, _ := tr.Read(buf)
        return codec.Decode(buf[:n])
    }

    wrapped := pipeline.Chain(
        pipeline.Timeout(3*time.Second),
        pipeline.Retry(3, pipeline.ExponentialBackoff(100*time.Millisecond)),
    )(handler)

    resp, err := wrapped(context.Background(), &kernel.Request{
        Function: "read_holding_registers",
        Address:  "40001",
        Count:    2,
    })
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Registers: %x", resp.Data)
}
```

### 使用 ConnectionManager + 连接池

```go
mgr := connection.NewManager(&connection.PooledStrategy{})

// 注册设备（从代码）
mgr.Register(&connection.DeviceConfig{
    Name: "plc-001", Network: "tcp", Addr: "192.168.1.10:502",
    Timeout: 3 * time.Second,
    Pool:    &connection.PoolConfig{MaxSize: 5},
})

// 或从 YAML 配置加载
repo, _ := config.Load("devices.yaml")
for _, cfg := range repo.All() {
    mgr.Register(cfg)
}

// 获取连接 → 使用 → 归还
tr, _ := mgr.Acquire("plc-001")
// ... 协议读写 ...
mgr.Release(tr)

// 健康检查
status := mgr.Health("plc-001")
```

### MQTT 发布订阅

```go
tr, _ := transport.DialTCP("broker.emqx.io:1883")
codec, _ := mqtt.New().NewCodec("tcp")

// 发送 CONNECT
connect, _ := codec.Encode(&kernel.Request{
    Function: "connect",
    Metadata: map[string]any{"client_id": "goclient"},
})
tr.Write(connect)
buf := make([]byte, 256)
n, _ := tr.Read(buf)
codec.Decode(buf[:n]) // CONNACK

// 发布消息
pub, _ := codec.Encode(&kernel.Request{
    Function: "publish",
    Address:  "sensor/temperature",
    Data:     []byte("25.5"),
})
tr.Write(pub)

// 订阅主题
sub, _ := codec.Encode(&kernel.Request{
    Function: "subscribe",
    Address:  "sensor/#",
})
tr.Write(sub)
```

### 使用 Vendor 预设

```go
reg := vendor.NewRegistry()
vendor.RegisterSiemens(reg)
vendor.RegisterRockwell(reg)

profile, _ := reg.Find("Siemens", "S7-1200")
// profile.Defaults["protocol"] = "modbus"
// profile.Defaults["port"]     = 502
// profile.Defaults["endian"]   = "big"
```

### 协议间网关转换

```go
engine := gateway.New()
engine.Add(gateway.Rule{
    Src: modbus.NewProtocol(),
    Dst: mqtt.New(),
    Map: func(srcResp *kernel.Response) *kernel.Request {
        return &kernel.Request{
            Function: "publish",
            Address:  "plc/registers",
            Data:     srcResp.Data,
        }
    },
})

srcTr, _ := transport.DialTCP("192.168.1.10:502")
dstTr, _ := transport.DialTCP("broker.emqx.io:1883")
engine.Transform(ctx, srcTr, dstTr, &kernel.Request{
    Function: "read_holding_registers", Address: "40001", Count: 2,
})
```

### 自定义协议

实现一个新协议只需两个文件：

```go
// myproto/myproto.go
type MyProtocol struct{}

func (p *MyProtocol) Name() string              { return "myproto" }
func (p *MyProtocol) Variants() []string        { return []string{"tcp"} }
func (p *MyProtocol) DefaultPort() int          { return 9000 }
func (p *MyProtocol) NewCodec(v string) (kernel.Codec, error) {
    return &myCodec{}, nil
}

type myCodec struct{}

func (c *myCodec) Encode(req *kernel.Request) ([]byte, error) {
    return []byte(req.Function + ":" + req.Address), nil
}

func (c *myCodec) Decode(data []byte) (*kernel.Response, error) {
    return &kernel.Response{Data: data}, nil
}
```

---

## 项目结构

```
industrial-protocols-go/
├── go.work                       # workspace 聚合所有模块
├── kernel/                       # 核心模块
│   ├── connection/               # ConnectionManager + 连接池
│   ├── config/                   # ConfigRepository (YAML)
│   ├── gateway/                  # GatewayEngine 协议转换
│   ├── bridge/                   # Bridge 外部进程桥接
│   ├── vendor/                   # Vendor 厂商预设
│   ├── event/                    # 事件总线
│   ├── metrics/                  # 指标接口
│   ├── security/                 # TLS 安全
│   ├── pipeline/                 # 中间件链
│   └── transport/                # TCP/UDP/Pipe 传输
│
├── protocols/                    # 40 协议模块
│   ├── ethernet/                 # 5 工业以太网（全部完成）
│   ├── fieldbus/                 # 11 现场总线（4 完成）
│   ├── iot/                      # 2 IoT/消息（全部完成）
│   ├── automotive/               # 5 汽车总线（2 完成）
│   ├── building/                 # 2 楼宇/照明（1 完成）
│   ├── bridge/                   # 12 硬件桥接（stub）
│   └── system/                   # 3 系统总线（stub）
│
├── examples/modbus_basic/        # Modbus TCP 示例
├── _tools/                       # Makefile + 辅助脚本
└── docs/superpowers/             # 设计文档
```

## 测试

```bash
make test         # 全量测试
make test-unit    # 仅单元测试
make vet && make fmt
```

---

## License

MIT — Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz
