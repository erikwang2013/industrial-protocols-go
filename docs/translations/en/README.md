English | [中文](../../README.md) | [한국어](../ko/README.md) | [Русский](../ru/README.md) | [Deutsch](../de/README.md) | [Français](../fr/README.md) | [Español](../es/README.md) | [Português](../pt/README.md) | [हिन्दी](../hi/README.md) | [العربية](../ar/README.md) | [বাংলা](../bn/README.md) | [Bahasa Indonesia](../id/README.md) | [日本語](../ja/README.md)

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

The full API reference and usage examples (installation, Modbus read/write, connection pool, MQTT, Vendor presets, gateway transformation, custom protocols) have moved to a standalone document:

> **[API.md](API.md) — Complete API Reference & Usage Guide**
>
> Chinese version: [API reference](../API.md)

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

## Support Us

If this project has helped you, feel free to tip us to support continued maintenance.

### Alipay / WeChat Pay

| Alipay | WeChat Pay |
|--------|------|
| ![Alipay](../../alipay.png) | ![WeChat Pay](../../weixinpay.png) |

### Global Transfer (International Bank Transfer)

Bank transfer for users outside mainland China:

**Recipient Information:**

| Item | Value |
|------|------|
| Recipient Name | WANG KEXUN |
| Recipient Account Number | 881015918251 |

**Receiving Bank:**

| Item | Value |
|------|------|
| Bank Name | ZA Bank Limited |
| SWIFT Code | AABLHKHHXXX |
| Bank Code | 387 |
| Bank Address | Core F, Cyberport 3, 100 Cyberport Road, Hong Kong |

**Cross-Border Correspondent Bank (if required):**

> Please note: this is the cross-border correspondent (intermediary) bank information, not the receiving bank. Please ask your remitting bank whether correspondent bank details are required.

- **For HKD, CNY and USD remittances** — correspondent bank is Citibank:
  - Bank Name: Citibank N.A. Hong Kong
  - SWIFT Code: CITIHKHXXXX
  - Bank Code: 006
  - Branch Name: Hong Kong Branch
  - Branch Code: 391
  - Bank Address: Citibank Tower, Citibank Plaza, 3 Garden Road, Central, Hong Kong
- **For other currencies** — correspondent bank is BNY Mellon:
  - Bank Name: THE BANK OF NEW YORK MELLON
  - SWIFT Code: IRVTUS3NXXX
  - Bank Address: THE BANK OF NEW YORK MELLON, 240 GREENWICH STREET, NEW YORK, United States

---

## License

MIT — Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz
