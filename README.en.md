# Industrial Protocols Go

Go industrial network communication protocol suite — layered + middleware architecture, covering 40 industrial protocols: 14 pure-software implementations + 26 hardware SDKs (all with drivers).

> Reference PHP implementation: [github.com/erikwang2013/industrial-protocols](https://github.com/erikwang2013/industrial-protocols)

---

## Project Design

### Architecture Layers

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
│  Per-protocol modules            (TCP/UDP/Serial)  │
│                                                │
│  14 Pure-Soft   +   26 Hardware SDKs (with drivers)│
└──────────────────────────────────────────────────┘
```

### Design Philosophy

**Micro-kernel + Protocol SDK.** The kernel defines only interfaces and cross-cutting concerns (connection pooling, retry, circuit breaker, events, metrics), and contains no concrete protocol implementations. Each protocol is an independent Go module, imported on demand.

**Layered decoupling.** Four abstraction layers:

| Layer | Responsibility | Core Types |
|----|------|---------|
| **Transport** | Low-level communication channel | `Transport` interface — `TCPTransport`, `UDPTransport`, `PipeTransport` |
| **Codec** | Protocol encode/decode | `Codec` interface — `Encode(*Request) ([]byte, error)` / `Decode([]byte) (*Response, error)` |
| **Pipeline** | Cross-cutting middleware | `Middleware` — `Timeout`, `Retry`, `CircuitBreaker`, `Logger` |
| **Kernel** | Device lifecycle | `ConnectionManager`, `ConfigRepository` |

**Protocol authors only need to implement `Protocol` + `Codec`.** Transport layer, connection pool, retry, timeout, and circuit breaker are all reused from the kernel.

### Core Interfaces

```go
// kernel/protocol.go — every protocol module must implement
type Protocol interface {
    Name()        string              // "modbus"
    Variants()    []string            // ["tcp", "rtu", "ascii"]
    DefaultPort() int                 // 502
    NewCodec(variant string) (Codec, error)
}

// kernel/codec.go — encode/decode
type Codec interface {
    Encode(req *Request) ([]byte, error)
    Decode(data []byte) (*Response, error)
}

// kernel/transport/transport.go — transport channel
type Transport interface {
    io.ReadWriter
    io.Closer
    Addr()  string
    Alive() bool
}

// kernel/pipeline/pipeline.go — middleware
type Handler    func(ctx context.Context, req *Request) (*Response, error)
type Middleware func(next Handler) Handler
```

### Middleware Chain

```
Request → [Timeout] → [Retry] → [CircuitBreaker] → [Logger] → Codec.Encode → Transport.Write
Response ← Codec.Decode ← Transport.Read
```

Built-in middleware:

| Middleware | Constructor | Description |
|--------|------|------|
| `Timeout` | `Timeout(d time.Duration)` | Context-based timeout control |
| `Retry` | `Retry(max int, backoff BackoffFunc)` | Linear / exponential / jitter backoff strategies |
| `CircuitBreaker` | `NewCircuitBreaker(threshold int, cooldown time.Duration)` | closed → open → half-open three-state breaker |
| `Logger` | `Logger(logger *log.Logger)` | Structured request logging |

### Kernel Modules

| Module | Path | Description |
|------|------|------|
| ConnectionManager | `kernel/connection/` | Device registration, connection pool, health checks; supports Lazy/Eager/Pooled strategies |
| ConfigRepository | `kernel/config/` | YAML/JSON device configuration loading |
| GatewayEngine | `kernel/gateway/` | Cross-protocol transformation rule engine (Modbus→MQTT, etc.) |
| Bridge | `kernel/bridge/` | External process bridging (stdin/stdout communication) |
| Vendor | `kernel/vendor/` | Vendor parameter presets (Siemens S7, Rockwell AB, etc.) |
| Event | `kernel/event/` | Channel-based event bus |
| Metrics | `kernel/metrics/` | Metrics collection interface (Prometheus + noop) |
| Security | `kernel/security/` | TLS transport-layer security wrapper |

---

## Protocol Support

### Industrial Ethernet (5/5 Complete)

| Protocol | Transport | Port | Features |
|------|------|------|------|
| **Modbus** | TCP, RTU | 502 | FC 01/02/03/04/05/06/16, CRC16, MBAP |
| **BACnet** | UDP | 47808 | Who-Is/I-Am device discovery, ReadProperty |
| **EtherNet/IP** | TCP | 44818 | ENIP RegisterSession + CIP Read Tag |
| **OPC UA** | TCP | 4840 | Binary HEL + OpenSecureChannel |
| **PROFINET** | UDP | 34964 | DCP Identify/Set + Record Data Read/Write |

### Fieldbus (11/11 — 4 Pure-Soft + 7 Hardware SDK)

| Protocol | Transport | Port | Features |
|------|------|------|------|
| **HART** | Serial FSK | — | Short/long frames, Command 0/3, XOR checksum |
| **CC-Link** | RS-485 | — | Master-slave polling, CRC-16/XMODEM |
| **DNP3** | TCP, Serial | 20000 | Transport-layer segmentation/reassembly, Class 0 polling, CRC-16/DNP |
| **IEC 61850** | TCP | 102 | MMS Initiate/Conclude, Read/Write, BER-TLV |
| **PROFIBUS** | Serial | — | Pure-soft implementation — requires CP 5611 hardware |
| **CANopen** | CAN, Gateway | — | SDO read/write, NMT start/stop/reset, Heartbeat |
| **DeviceNet** | CAN, Gateway | — | Explicit Messaging, Poll, I/O connections |
| **Foundation Fieldbus** | Serial | — | Pure-soft implementation — requires FF H1 interface card |
| **AS-Interface** | Serial | — | Pure-soft implementation — requires ASi gateway |
| **IO-Link** | Serial | — | Pure-soft implementation — requires IO-Link Master |
| **CC-Link IE** | Ethernet | — | Pure-soft implementation — requires CC-Link IE gateway |

### IoT / Messaging (2/2 Complete)

| Protocol | Transport | Port | Features |
|------|------|------|------|
| **MQTT** | TCP, WS | 1883 | 3.1.1 CONNECT/PUBLISH/SUBSCRIBE/PING |
| **HART-IP** | TCP, UDP | 5094 | HART over TCP, reuses HART codec |

### Automotive Bus (5/5 — 2 Pure-Soft + 3 Hardware SDK)

| Protocol | Transport | Baud Rate | Features |
|------|------|--------|------|
| **LIN** | UART | — | Master-slave frames, PID verification, Classic/Enhanced Checksum |
| **K-Line** | Serial | 10400 | ISO 9141/14230, 5-baud Fast Init, OBD-II SID 01/03/09 |
| **FlexRay** | CAN, Serial | — | Slot/Frame encode/decode, CRC-16/XMODEM, SocketCAN + SerialBridge |
| **SAE J1850** | CAN | — | PWM/VPW encode/decode, Mode 01/03/0A, CRC-8 |
| **MOST** | Serial | — | Hex-frame encode/decode, Read/Write/Status, SerialBridge |

### Building / Lighting (2/2 — 1 Pure-Soft + 1 Hardware SDK)

| Protocol | Transport | Features |
|------|------|------|
| **DALI** | Serial | 16-bit forward frames, 8-bit backward frames, standard commands (Off/Max/Dim) |
| **LonWorks** | Serial | Pure-soft codec — requires Neuron chip or gateway |

### Hardware Bridge (12/12 — All with CmdBridge/SerialBridge drivers)

| Protocol | SDK Type | Features |
|------|----------|------|
| **EtherCAT** | CmdBridge + hex-codec | Upload/Download/Slaves subcommands, CoE hex encoding/decoding |
| **POWERLINK** | CmdBridge + hex-codec | Read/Write/Status subcommands, SoC/Preq/Pres frames |
| **SERCOS III** | CmdBridge + hex-codec | Read/Write/Phase subcommands, phase transitions (NRT→CP4) |
| **SERCOS I/II** | CmdBridge + hex-codec | Read/Write/Status, fiber ring topology |
| **ControlNet** | CmdBridge + hex-codec | Read/Write/Status, CTDMA scheduling, producer/consumer |
| **Interbus** | CmdBridge | Read/Decode, ring topology, IBS CMD frames |
| **WorldFIP** | CmdBridge | Read/Write/Decode, producer/consumer, bus arbitrator |
| **Lightbus** | CmdBridge | Read/Decode, fiber interconnection, 32 nodes |
| **Modbus Plus** | CmdBridge + hex-codec | Read/Write/Decode, token passing, peer-to-peer |
| **ISA100.11a** | CmdBridge + hex-codec | Read/Write/List, 6LoWPAN wireless, mesh |
| **WirelessHART** | CmdBridge + hex-codec | Read/Write/Scan, TSMP MAC, self-organizing |
| **SAE J1850** | CmdBridge | Automotive diagnostics via CmdBridge, requires J1850 interface |

### System Bus (3/3 — All with sysfs/procfs drivers)

| Protocol | SDK Type | Features |
|------|----------|------|
| **PCI/PCIe** | sysfs — `/sys/bus/pci/devices/<BDF>/config` | Config space read/write, PipeTransport, requires CAP_SYS_ADMIN |
| **VME/VPX** | procfs — `/proc/vme/<slot>` | A16/A24/A32 addressing, PipeTransport, requires vme_tsi148 module |
| **CompactPCI** | sysfs — `/sys/bus/pci/devices/<BDF>/config` | PICMG 2.0, hot-plug, 3U/6U, requires cpci_hotplug |

---

## Usage Guide

### Installation

The kernel library is required. Protocol modules are imported on demand. All modules share version `v1.1.1`.

```bash
# Kernel (required)
go get github.com/erikwang2013/industrial-protocols-go/kernel@v1.1.1
```

**Industrial Ethernet (5):**

```bash
go get github.com/erikwang2013/industrial-protocols-go/protocols/ethernet/modbus@v1.1.1
go get github.com/erikwang2013/industrial-protocols-go/protocols/ethernet/bacnet@v1.1.1
go get github.com/erikwang2013/industrial-protocols-go/protocols/ethernet/ethernetip@v1.1.1
go get github.com/erikwang2013/industrial-protocols-go/protocols/ethernet/opcua@v1.1.1
go get github.com/erikwang2013/industrial-protocols-go/protocols/ethernet/profinet@v1.1.1
```

**Fieldbus (11):**

```bash
go get github.com/erikwang2013/industrial-protocols-go/protocols/fieldbus/hart@v1.1.1
go get github.com/erikwang2013/industrial-protocols-go/protocols/fieldbus/cclink@v1.1.1
go get github.com/erikwang2013/industrial-protocols-go/protocols/fieldbus/dnp3@v1.1.1
go get github.com/erikwang2013/industrial-protocols-go/protocols/fieldbus/iec61850@v1.1.1
go get github.com/erikwang2013/industrial-protocols-go/protocols/fieldbus/profibus@v1.1.1
go get github.com/erikwang2013/industrial-protocols-go/protocols/fieldbus/canopen@v1.1.1
go get github.com/erikwang2013/industrial-protocols-go/protocols/fieldbus/devicenet@v1.1.1
go get github.com/erikwang2013/industrial-protocols-go/protocols/fieldbus/foundationfieldbus@v1.1.1
go get github.com/erikwang2013/industrial-protocols-go/protocols/fieldbus/asinterface@v1.1.1
go get github.com/erikwang2013/industrial-protocols-go/protocols/fieldbus/iolink@v1.1.1
go get github.com/erikwang2013/industrial-protocols-go/protocols/fieldbus/cclinkie@v1.1.1
```

**IoT / Messaging (2):**

```bash
go get github.com/erikwang2013/industrial-protocols-go/protocols/iot/mqtt@v1.1.1
go get github.com/erikwang2013/industrial-protocols-go/protocols/iot/hartip@v1.1.1
```

**Automotive (5):**

```bash
go get github.com/erikwang2013/industrial-protocols-go/protocols/automotive/lin@v1.1.1
go get github.com/erikwang2013/industrial-protocols-go/protocols/automotive/kline@v1.1.1
go get github.com/erikwang2013/industrial-protocols-go/protocols/automotive/flexray@v1.1.1
go get github.com/erikwang2013/industrial-protocols-go/protocols/automotive/saej1850@v1.1.1
go get github.com/erikwang2013/industrial-protocols-go/protocols/automotive/most@v1.1.1
```

**Building / Lighting (2):**

```bash
go get github.com/erikwang2013/industrial-protocols-go/protocols/building/dali@v1.1.1
go get github.com/erikwang2013/industrial-protocols-go/protocols/building/lonworks@v1.1.1
```

**Hardware Bridges (12):**

```bash
go get github.com/erikwang2013/industrial-protocols-go/protocols/bridge/ethercat@v1.1.1
go get github.com/erikwang2013/industrial-protocols-go/protocols/bridge/powerlink@v1.1.1
go get github.com/erikwang2013/industrial-protocols-go/protocols/bridge/sercos@v1.1.1
go get github.com/erikwang2013/industrial-protocols-go/protocols/bridge/sercos1@v1.1.1
go get github.com/erikwang2013/industrial-protocols-go/protocols/bridge/controlnet@v1.1.1
go get github.com/erikwang2013/industrial-protocols-go/protocols/bridge/interbus@v1.1.1
go get github.com/erikwang2013/industrial-protocols-go/protocols/bridge/worldfip@v1.1.1
go get github.com/erikwang2013/industrial-protocols-go/protocols/bridge/lightbus@v1.1.1
go get github.com/erikwang2013/industrial-protocols-go/protocols/bridge/modbusplus@v1.1.1
go get github.com/erikwang2013/industrial-protocols-go/protocols/bridge/isa100@v1.1.1
go get github.com/erikwang2013/industrial-protocols-go/protocols/bridge/wirelesshart@v1.1.1
```

**System Bus (3):**

```bash
go get github.com/erikwang2013/industrial-protocols-go/protocols/system/pci@v1.1.1
go get github.com/erikwang2013/industrial-protocols-go/protocols/system/vme@v1.1.1
go get github.com/erikwang2013/industrial-protocols-go/protocols/system/cpci@v1.1.1
```

> **pkg.go.dev URL:** `https://pkg.go.dev/github.com/erikwang2013/industrial-protocols-go/` + submodule path (e.g. `/kernel`, `/protocols/ethernet/modbus`)

### Modbus TCP Read/Write

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

### Using ConnectionManager + Connection Pool

```go
mgr := connection.NewManager(&connection.PooledStrategy{})

// Register device (from code)
mgr.Register(&connection.DeviceConfig{
    Name: "plc-001", Network: "tcp", Addr: "192.168.1.10:502",
    Timeout: 3 * time.Second,
    Pool:    &connection.PoolConfig{MaxSize: 5},
})

// Or load from YAML config
repo, _ := config.Load("devices.yaml")
for _, cfg := range repo.All() {
    mgr.Register(cfg)
}

// Acquire connection → use → release
tr, _ := mgr.Acquire("plc-001")
// ... protocol read/write ...
mgr.Release(tr)

// Health check
status := mgr.Health("plc-001")
```

### MQTT Publish/Subscribe

```go
tr, _ := transport.DialTCP("broker.emqx.io:1883")
codec, _ := mqtt.New().NewCodec("tcp")

// Send CONNECT
connect, _ := codec.Encode(&kernel.Request{
    Function: "connect",
    Metadata: map[string]any{"client_id": "goclient"},
})
tr.Write(connect)
buf := make([]byte, 256)
n, _ := tr.Read(buf)
codec.Decode(buf[:n]) // CONNACK

// Publish message
pub, _ := codec.Encode(&kernel.Request{
    Function: "publish",
    Address:  "sensor/temperature",
    Data:     []byte("25.5"),
})
tr.Write(pub)

// Subscribe to topic
sub, _ := codec.Encode(&kernel.Request{
    Function: "subscribe",
    Address:  "sensor/#",
})
tr.Write(sub)
```

### Using Vendor Presets

```go
reg := vendor.NewRegistry()
vendor.RegisterSiemens(reg)
vendor.RegisterRockwell(reg)

profile, _ := reg.Find("Siemens", "S7-1200")
// profile.Defaults["protocol"] = "modbus"
// profile.Defaults["port"]     = 502
// profile.Defaults["endian"]   = "big"
```

### Cross-Protocol Gateway Transformation

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

### Custom Protocol

Implementing a new protocol requires only two files:

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

## Project Structure

```
industrial-protocols-go/
├── go.work                       # workspace aggregates all modules
├── kernel/                       # core module
│   ├── connection/               # ConnectionManager + connection pool
│   ├── config/                   # ConfigRepository (YAML)
│   ├── gateway/                  # GatewayEngine protocol transformation
│   ├── bridge/                   # Bridge external process bridging
│   ├── vendor/                   # Vendor presets
│   ├── event/                    # Event bus
│   ├── metrics/                  # Metrics interface
│   ├── security/                 # TLS security
│   ├── pipeline/                 # Middleware chain
│   └── transport/                # TCP/UDP/Pipe transport
│
├── protocols/                    # 40 protocol modules
│   ├── ethernet/                 # 5 industrial Ethernet (all complete)
│   ├── fieldbus/                 # 11 fieldbus (4 pure-soft + 7 hardware SDK)
│   ├── iot/                      # 2 IoT/Messaging (all complete)
│   ├── automotive/               # 5 automotive bus (2 pure-soft + 3 hardware SDK)
│   ├── building/                 # 2 building/lighting (1 pure-soft + 1 hardware SDK)
│   ├── bridge/                   # 12 hardware bridges (CmdBridge/SerialBridge)
│   └── system/                   # 3 system buses (sysfs/procfs drivers)
│
├── examples/modbus_basic/        # Modbus TCP example
├── _tools/                       # Makefile + helper scripts
└── docs/superpowers/             # Design documentation
```

## Testing

```bash
make test         # Full test suite
make test-unit    # Unit tests only
make vet && make fmt
```

---

## License

MIT — Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz
