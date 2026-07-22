# Industrial Protocols Go — Phase 3 Design Spec

Date: 2026-07-22
Phases 1-2: kernel + 14 soft protocol implementations + 26 stubs

## Overview

Phase 3 将 26 个 stub 协议升级为硬件 SDK — 每个协议定义硬件接口抽象、桥接方式、协议数据结构，通过 kernel/bridge 模块的 4 种桥接类型对接真实硬件。

## Bridge Types

```
kernel/bridge/
├── bridge.go       # Bridge 接口 (已有)
├── cmd.go          # CmdBridge — CLI 子进程 stdin/stdout (已有)
├── gateway.go      # GatewayBridge — TCP 网关 (新增)
├── can.go          # CANBridge — Linux SocketCAN, build tag linux (新增)
└── serial.go       # SerialBridge — 串口 AT 指令 (新增)
```

### Bridge 接口

```go
type Bridge interface {
    Start() error
    Stop() error
    Transport() transport.Transport
}
```

### CmdBridge (已有，不变)

```go
type CmdBridge struct { ... }
func NewCmdBridge(name string, args ...string) *CmdBridge
```

### GatewayBridge (新增)

```go
type GatewayBridge struct { conn net.Conn; addr string; running bool }
func NewGatewayBridge(address string) *GatewayBridge
```

### CANBridge (新增，linux only)

```go
type CANFrame struct { ID uint32; Data []byte; Ext bool }
type CANBridge struct { iface string; socket int }
func NewCANBridge(iface string) *CANBridge
func (b *CANBridge) Send(frame CANFrame) error
func (b *CANBridge) Recv() (CANFrame, error)
```

### SerialBridge (新增)

```go
type SerialBridge struct { port io.ReadWriteCloser; dev string; baud int }
func NewSerialBridge(dev string, baud int) *SerialBridge
```

## Hardware SDK Protocol Structure

每个硬件协议模块：

```
protocols/<category>/<name>/
├── go.mod
├── <name>.go          # Protocol + Codec 实现
├── driver.go          # NewXxxDriver() → Bridge + Codec
├── <name>_test.go     # Codec 单元测试（mock transport）
└── README.md          # 硬件要求 + 接线说明
```

### driver.go 统一模式

```go
func ReadyHandler(driver func() (bridge.Bridge, kernel.Codec, error)) (pipeline.Handler, error) {
    b, codec, err := driver()
    if err != nil { return nil, err }
    if err := b.Start(); err != nil { return nil, err }
    tr := b.Transport()
    return pipeline.Chain(
        pipeline.Timeout(5 * time.Second),
        pipeline.Retry(3, pipeline.LinearBackoff(100*time.Millisecond)),
    )(func(ctx context.Context, req *kernel.Request) (*kernel.Response, error) {
        raw, err := codec.Encode(req)
        if err != nil { return nil, err }
        tr.Write(raw)
        buf := make([]byte, 4096)
        n, _ := tr.Read(buf)
        return codec.Decode(buf[:n])
    }), nil
}
```

## Protocol Allocation

### CmdBridge (8 protocols)

| Protocol | Driver | Hardware | CLI Tool |
|----------|--------|----------|----------|
| EtherCAT | NewTwinCATDriver / NewSOEMDriver | Beckhoff NIC / PC NIC | `TwinCAT.Cmd.exe` / `soem_master` |
| POWERLINK | NewOpenPOWERLINKDriver | Standard NIC | `openPOWERLINK` daemon |
| SERCOS III | NewSercosIIIDriver | Hilscher netX | `netx_cli` |
| SERCOS I/II | NewSercosDriver | SERCOS fiber card | `sercos_cli` |
| ControlNet | NewControlNetDriver | AB 1784-PCIC/S | `1784-pcic-cli` |
| Modbus Plus | NewModbusPlusDriver | Schneider SA85/BM85 | `sa85_cli` |
| ISA100 | NewISA100Driver | Yokogawa YFGW410 | `yfgw410_cli` |
| WirelessHART | NewWirelessHARTDriver | Emerson 1410/1420 | `emerson_1410_cli` |

### GatewayBridge (10 protocols)

| Protocol | Driver | Hardware | Default Port |
|----------|--------|----------|-------------|
| PROFIBUS | NewAnybusPBDriver | Anybus/Siemens CP 5611 TCP proxy | 2000 |
| DeviceNet | NewAnybusDNDriver | Anybus DeviceNet Scanner | 2001 |
| Foundation Fieldbus | NewNIFFDriver / NewSoftingFFDriver | NI USB-8486 / Softing FFusb | 2002 |
| AS-Interface | NewBihlWiedemannDriver | Bihl+Wiedemann / P+F ASi | 2003 |
| IO-Link | NewIfmDriver / NewBalluffDriver | ifm / Balluff IO-Link Master | 2004 |
| CC-Link IE | NewMitsubishiCLIEDriver | Mitsubishi CC-Link IE gateway | 2005 |
| Interbus | NewPhoenixIBSDriver | Phoenix Contact IBS | 2006 |
| WorldFIP | NewWorldFIPDriver | WorldFIP/Fipio gateway | 2007 |
| Lightbus | NewBeckhoffLightbusDriver | Beckhoff Lightbus gateway | 2008 |
| LonWorks | NewEchelonDriver | Echelon U60/U70 gateway | 2009 |

### CANBridge (3 protocols)

| Protocol | Driver | Hardware | CAN ID Base |
|----------|--------|----------|------------|
| CANopen | NewSocketCANDriver("can0") | PCAN-USB / vcan0 | 0x600+NodeID |
| FlexRay (CAN) | NewSocketCANDriver("can0") | Vector/Bosch FlexRay-CAN adapter | per ECU |
| SAE J1850 | NewSocketCANDriver("can0") | J1850-CAN OBD-II adapter | per OBD-II |

### SerialBridge (2 protocols)

| Protocol | Driver | Hardware | Baud |
|----------|--------|----------|------|
| FlexRay (serial) | NewFlexRaySerialDriver | FlexRay serial debugger | 115200 |
| MOST | NewMOSTSerialDriver | MOST fiber serial adapter | 115200 |

### SystemBus (3 protocols)

| Protocol | Driver | Probe Method | Path |
|----------|--------|-------------|------|
| PCI/PCIe | NewPCIDriver | sysfs scan | `/sys/bus/pci/devices/` |
| VME/VPX | NewVMEDriver | procfs | `/proc/vme` |
| CompactPCI | NewCPCIDriver | PCI scan | `/sys/bus/pci/devices/` |

SystemBus 协议的 `Transport` 返回 `io.ReadWriter` 包装的 `os.File`（读写 `/dev/pciX` 或 sysfs 节点）。

## Implementation Plan

### Priority (by bridge availability for testing)

| Priority | Bridge Type | Count | Testable without hardware |
|----------|-----------|-------|--------------------------|
| 1 | CANBridge + CAN protocols | 3 | vcan0 (virtual CAN) |
| 2 | GatewayBridge protocols | 10 | Mock TCP server |
| 3 | SerialBridge protocols | 2 | Mock reader |
| 4 | CmdBridge protocols | 8 | Needs CLI tools |
| 5 | SystemBus | 3 | Needs Linux kernel |

### Files Created/Modified

| Category | Files |
|----------|-------|
| kernel/bridge/ (3 new bridges) | gateway.go, can.go, serial.go + tests |
| protocols/bridge/ (12 SDK upgrades) | driver.go + README.md per protocol |
| protocols/fieldbus/ (6 SDK upgrades) | driver.go + README.md |
| protocols/automotive/ (3 SDK upgrades) | driver.go + README.md |
| protocols/building/ (1 SDK upgrade) | driver.go + README.md |
| protocols/system/ (3 SDK upgrades) | driver.go + README.md |
| **Total** | **~80 files** |

## Decisions Summary

| Dimension | Decision |
|-----------|----------|
| Bridge architecture | 4 types: Cmd + Gateway + CAN + Serial + SystemBus (sysfs) |
| Driver pattern | Each protocol provides NewXxxDriver() → Bridge+Codec |
| CAN | Linux SocketCAN, build tag `linux` |
| Gateway | TCP connection to hardware gateway |
| Testing | Codec unit tests (mock transport) for all; CAN tested with vcan0 |
