[English](README.en.md) | 中文 | [한국어](docs/translations/ko/README.md) | [Русский](docs/translations/ru/README.md) | [Deutsch](docs/translations/de/README.md) | [Français](docs/translations/fr/README.md) | [Español](docs/translations/es/README.md) | [Português](docs/translations/pt/README.md) | [हिन्दी](docs/translations/hi/README.md) | [العربية](docs/translations/ar/README.md) | [বাংলা](docs/translations/bn/README.md) | [Bahasa Indonesia](docs/translations/id/README.md) | [日本語](docs/translations/ja/README.md)

# Industrial Protocols Go

Go 语言工业网络通信协议集 —— 分层 + 中间件架构，覆盖 40 种工业协议，14 个纯软实现 + 26 个硬件 SDK（全部带驱动）。

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
│  14 Pure-Soft   +   26 Hardware SDKs (with drivers)│
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

### 现场总线（11/11 — 4 纯软 + 7 硬件 SDK）

| 协议 | 传输 | 端口 | 功能 |
|------|------|------|------|
| **HART** | Serial FSK | — | 短帧/长帧, Command 0/3, XOR 校验 |
| **CC-Link** | RS-485 | — | 主从轮询, CRC-16/XMODEM |
| **DNP3** | TCP, Serial | 20000 | 传输层分段/重组, Class 0 轮询, CRC-16/DNP |
| **IEC 61850** | TCP | 102 | MMS Initiate/Conclude, Read/Write, BER-TLV |
| **PROFIBUS** | Serial | — | 纯软实现 — 需 CP 5611 硬件 |
| **CANopen** | CAN, Gateway | — | SDO 读/写, NMT 启停/复位, Heartbeat |
| **DeviceNet** | CAN, Gateway | — | Explicit Messaging, Poll, I/O 连接 |
| **Foundation Fieldbus** | Serial | — | 纯软实现 — 需 FF H1 接口卡 |
| **AS-Interface** | Serial | — | 纯软实现 — 需 ASi 网关 |
| **IO-Link** | Serial | — | 纯软实现 — 需 IO-Link Master |
| **CC-Link IE** | Ethernet | — | 纯软实现 — 需 CC-Link IE 网关 |

### IoT / 消息（2/2 已完成）

| 协议 | 传输 | 端口 | 功能 |
|------|------|------|------|
| **MQTT** | TCP, WS | 1883 | 3.1.1 CONNECT/PUBLISH/SUBSCRIBE/PING |
| **HART-IP** | TCP, UDP | 5094 | HART over TCP, 复用 HART 编解码 |

### 汽车总线（5/5 — 2 纯软 + 3 硬件 SDK）

| 协议 | 传输 | 波特率 | 功能 |
|------|------|--------|------|
| **LIN** | UART | — | 主从帧, PID 校验, Classic/Enhanced Checksum |
| **K-Line** | Serial | 10400 | ISO 9141/14230, 5-baud Fast Init, OBD-II SID 01/03/09 |
| **FlexRay** | CAN, Serial | — | Slot/Frame 编解码, CRC-16/XMODEM, SocketCAN + SerialBridge |
| **SAE J1850** | CAN | — | PWM/VPW 编解码, Mode 01/03/0A, CRC-8 |
| **MOST** | Serial | — | Hex-frame 编解码, Read/Write/Status, SerialBridge |

### 楼宇 / 照明（2/2 — 1 纯软 + 1 硬件 SDK）

| 协议 | 传输 | 功能 |
|------|------|------|
| **DALI** | Serial | 16-bit 前向帧, 8-bit 后向帧, 标准命令 (Off/Max/Dim) |
| **LonWorks** | Serial | 纯软编解码 — 需 Neuron 芯片或网关 |

### 硬件桥接（12/12 — 全部有 CmdBridge/SerialBridge 驱动）

| 协议 | SDK 类型 | 功能 |
|------|----------|------|
| **EtherCAT** | CmdBridge + hex-codec | Upload/Download/Slaves 子命令, CoE 十六进制编解码 |
| **POWERLINK** | CmdBridge + hex-codec | Read/Write/Status 子命令, SoC/Preq/Pres 帧 |
| **SERCOS III** | CmdBridge + hex-codec | Read/Write/Phase 子命令, 阶段转换 (NRT→CP4) |
| **SERCOS I/II** | CmdBridge + hex-codec | Read/Write/Status, 光纤环拓扑 |
| **ControlNet** | CmdBridge + hex-codec | Read/Write/Status, CTDMA 调度, 生产者/消费者 |
| **Interbus** | CmdBridge | Read/Decode, 环形拓扑, IBS CMD 帧 |
| **WorldFIP** | CmdBridge | Read/Write/Decode, 生产者/消费者, 总线仲裁器 |
| **Lightbus** | CmdBridge | Read/Decode, 光纤互连, 32节点 |
| **Modbus Plus** | CmdBridge + hex-codec | Read/Write/Decode, 令牌传递, 对等 |
| **ISA100.11a** | CmdBridge + hex-codec | Read/Write/List, 6LoWPAN 无线, 网状 |
| **WirelessHART** | CmdBridge + hex-codec | Read/Write/Scan, TSMP MAC, 自组织 |
| **SAE J1850** | CmdBridge | 车规诊断转 CmdBridge, 需 J1850 接口 |

### 系统总线（3/3 — 全部有 sysfs/procfs 驱动）

| 协议 | SDK 类型 | 功能 |
|------|----------|------|
| **PCI/PCIe** | sysfs — `/sys/bus/pci/devices/<BDF>/config` | 配置空间读写, PipeTransport, 需 CAP_SYS_ADMIN |
| **VME/VPX** | procfs — `/proc/vme/<slot>` | A16/A24/A32 编址, PipeTransport, 需 vme_tsi148 模块 |
| **CompactPCI** | sysfs — `/sys/bus/pci/devices/<BDF>/config` | PICMG 2.0, 热插拔, 3U/6U, 需 cpci_hotplug |

---

## 使用指南

API 参考与使用示例（安装、Modbus 读写、连接池、MQTT、Vendor 预设、网关转换、自定义协议）已移至独立文档：

> **[docs/API.md](docs/API.md) — 完整 API 参考与使用指南**
>
> 英文版：[docs/API.en.md](docs/API.en.md)

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
│   ├── fieldbus/                 # 11 现场总线（4 纯软 + 7 硬件 SDK）
│   ├── iot/                      # 2 IoT/消息（全部完成）
│   ├── automotive/               # 5 汽车总线（2 纯软 + 3 硬件 SDK）
│   ├── building/                 # 2 楼宇/照明（1 纯软 + 1 硬件 SDK）
│   ├── bridge/                   # 12 硬件桥接（CmdBridge/SerialBridge）
│   └── system/                   # 3 系统总线（sysfs/procfs 驱动）
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

## 支持我们

如果这个项目帮助到了你，欢迎打赏支持我们继续维护。

### 支付宝 / 微信

| 支付宝 | 微信 |
|--------|------|
| ![支付宝](docs/alipay.png) | ![微信](docs/weixinpay.png) |

### 全球转账（国际银行汇款）

适用于中国大陆以外用户的银行转账汇款：

**收款人信息：**

| 项目 | 内容 |
|------|------|
| 收款人姓名 | WANG KEXUN |
| 收款账户号码 | 881015918251 |

**收款银行：**

| 项目 | 内容 |
|------|------|
| 银行名称 | ZA Bank Limited |
| SWIFT Code | AABLHKHHXXX |
| 银行编号 | 387 |
| 银行地址 | Core F, Cyberport 3, 100 Cyberport Road, Hong Kong |

**跨境汇款代理银行（中转银行，如需）：**

> 请留意，此为跨境汇款代理银行（中转银行）信息，非收款银行信息。请向汇款银行查询是否需要提供代理银行信息。

- **汇入港元、人民币及美元** — 代理银行为 Citibank：
  - 银行名称：Citibank N.A. Hong Kong
  - SWIFT Code：CITIHKHXXXX
  - 银行编号：006
  - 分行名称：Hong Kong Branch
  - 分行编号：391
  - 银行地址：Citibank Tower, Citibank Plaza, 3 Garden Road, Central, Hong Kong
- **汇入其他币种** — 代理银行为 BNY Mellon：
  - 银行名称：THE BANK OF NEW YORK MELLON
  - SWIFT Code：IRVTUS3NXXX
  - 银行地址：THE BANK OF NEW YORK MELLON, 240 GREENWICH STREET, NEW YORK, United States

---

## License

MIT — Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz
