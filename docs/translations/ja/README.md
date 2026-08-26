[English](../../README.en.md) | [中文](../../README.md) | [한국어](../ko/README.md) | [Русский](../ru/README.md) | [Deutsch](../de/README.md) | [Français](../fr/README.md) | [Español](../es/README.md) | [Português](../pt/README.md) | [हिन्दी](../hi/README.md) | [العربية](../ar/README.md) | [বাংলা](../bn/README.md) | [Bahasa Indonesia](../id/README.md) | 日本語

# Industrial Protocols Go

Go言語による工業ネットワーク通信プロトコル集 —— 階層 + ミドルウェアアーキテクチャで、40種類の工業プロトコル、14の純ソフトウェア実装 + 26のハードウェアSDK（すべてドライバ付き）をカバー。

> PHP実装の参考: [github.com/erikwang2013/industrial-protocols](https://github.com/erikwang2013/industrial-protocols)

---

## プロジェクト設計

### アーキテクチャ階層

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
│  プロトコル別独立モジュール           (TCP/UDP/Serial)  │
│                                                │
│  14 Pure-Soft   +   26 Hardware SDKs (with drivers)│
└──────────────────────────────────────────────────┘
```

### 設計思想

**マイクロカーネル + プロトコルSDK。** カーネルはインターフェースと横断的関心事（コネクションプール、リトライ、サーキットブレーカー、イベント、メトリクス）のみを定義し、具体的なプロトコル実装は一切含みません。各プロトコルは独立したGoモジュールとして、必要に応じて導入します。

**階層による分離。** 4層の抽象化：

| 層 | 役割 | コア型 |
|----|------|---------|
| **Transport** | 最下層の通信チャネル | `Transport` インターフェース — `TCPTransport`, `UDPTransport`, `PipeTransport` |
| **Codec** | プロトコルエンコード/デコード | `Codec` インターフェース — `Encode(*Request) ([]byte, error)` / `Decode([]byte) (*Response, error)` |
| **Pipeline** | 横断的ミドルウェア | `Middleware` — `Timeout`, `Retry`, `CircuitBreaker`, `Logger` |
| **Kernel** | デバイスのライフサイクル | `ConnectionManager`, `ConfigRepository` |

**プロトコル作者は `Protocol` + `Codec` の2つのインターフェースを実装するだけです。** トランスポート層、コネクションプール、リトライ、タイムアウト、サーキットブレーカーはすべてkernelから再利用します。

### カーネルモジュール

| モジュール | パス | 説明 |
|------|------|------|
| ConnectionManager | `kernel/connection/` | デバイス登録、コネクションプール、ヘルスチェック。Lazy/Eager/Pooled の3つの戦略をサポート |
| ConfigRepository | `kernel/config/` | YAML/JSONデバイス設定のロード |
| GatewayEngine | `kernel/gateway/` | クロスプロトコル変換ルールエンジン（Modbus→MQTT など） |
| Bridge | `kernel/bridge/` | 外部プロセスとのブリッジ（stdin/stdout通信） |
| Vendor | `kernel/vendor/` | ベンダーパラメータプリセット（Siemens S7, Rockwell AB など） |
| Event | `kernel/event/` | channelベースのイベントバス |
| Metrics | `kernel/metrics/` | メトリクス収集インターフェース（Prometheus + noop） |
| Security | `kernel/security/` | TLSトランスポート層セキュリティラッパー |

---

## プロトコル対応

### 産業イーサネット（5/5 完了）

| プロトコル | トランスポート | ポート | 機能 |
|------|------|------|------|
| **Modbus** | TCP, RTU | 502 | FC 01/02/03/04/05/06/16, CRC16, MBAP |
| **BACnet** | UDP | 47808 | Who-Is/I-Am デバイス発見, ReadProperty |
| **EtherNet/IP** | TCP | 44818 | ENIP RegisterSession + CIP Read Tag |
| **OPC UA** | TCP | 4840 | Binary HEL + OpenSecureChannel |
| **PROFINET** | UDP | 34964 | DCP Identify/Set + Record Data Read/Write |

### フィールドバス（11/11 — 4純ソフト + 7ハードウェアSDK）

| プロトコル | トランスポート | ポート | 機能 |
|------|------|------|------|
| **HART** | Serial FSK | — | 短フレーム/長フレーム, Command 0/3, XORチェック |
| **CC-Link** | RS-485 | — | マスター/スレーブポーリング, CRC-16/XMODEM |
| **DNP3** | TCP, Serial | 20000 | トランスポート層のセグメンテーション/再構築, Class 0 ポーリング, CRC-16/DNP |
| **IEC 61850** | TCP | 102 | MMS Initiate/Conclude, Read/Write, BER-TLV |
| **PROFIBUS** | Serial | — | 純ソフト実装 — CP 5611 ハードウェアが必要 |
| **CANopen** | CAN, Gateway | — | SDO 読み/書き, NMT 起動/停止/リセット, Heartbeat |
| **DeviceNet** | CAN, Gateway | — | Explicit Messaging, Poll, I/O 接続 |
| **Foundation Fieldbus** | Serial | — | 純ソフト実装 — FF H1 インターフェースカードが必要 |
| **AS-Interface** | Serial | — | 純ソフト実装 — ASi ゲートウェイが必要 |
| **IO-Link** | Serial | — | 純ソフト実装 — IO-Link Master が必要 |
| **CC-Link IE** | Ethernet | — | 純ソフト実装 — CC-Link IE ゲートウェイが必要 |

### IoT / メッセージング（2/2 完了）

| プロトコル | トランスポート | ポート | 機能 |
|------|------|------|------|
| **MQTT** | TCP, WS | 1883 | 3.1.1 CONNECT/PUBLISH/SUBSCRIBE/PING |
| **HART-IP** | TCP, UDP | 5094 | HART over TCP, HART エンコード/デコードを再利用 |

### 自動車バス（5/5 — 2純ソフト + 3ハードウェアSDK）

| プロトコル | トランスポート | ボーレート | 機能 |
|------|------|--------|------|
| **LIN** | UART | — | マスター/スレーブフレーム, PIDチェック, Classic/Enhanced Checksum |
| **K-Line** | Serial | 10400 | ISO 9141/14230, 5-baud Fast Init, OBD-II SID 01/03/09 |
| **FlexRay** | CAN, Serial | — | Slot/Frame エンコード/デコード, CRC-16/XMODEM, SocketCAN + SerialBridge |
| **SAE J1850** | CAN | — | PWM/VPW エンコード/デコード, Mode 01/03/0A, CRC-8 |
| **MOST** | Serial | — | Hex-frame エンコード/デコード, Read/Write/Status, SerialBridge |

### ビル / 照明（2/2 — 1純ソフト + 1ハードウェアSDK）

| プロトコル | トランスポート | 機能 |
|------|------|------|
| **DALI** | Serial | 16-bit 前方向フレーム, 8-bit 後方向フレーム, 標準コマンド (Off/Max/Dim) |
| **LonWorks** | Serial | 純ソフトエンコード/デコード — Neuron チップまたはゲートウェイが必要 |

### ハードウェアブリッジ（12/12 — すべて CmdBridge/SerialBridge ドライバ付き）

| プロトコル | SDKタイプ | 機能 |
|------|----------|------|
| **EtherCAT** | CmdBridge + hex-codec | Upload/Download/Slaves サブコマンド, CoE 16進エンコード/デコード |
| **POWERLINK** | CmdBridge + hex-codec | Read/Write/Status サブコマンド, SoC/Preq/Pres フレーム |
| **SERCOS III** | CmdBridge + hex-codec | Read/Write/Phase サブコマンド, フェーズ遷移 (NRT→CP4) |
| **SERCOS I/II** | CmdBridge + hex-codec | Read/Write/Status, 光ファイバリングトポロジ |
| **ControlNet** | CmdBridge + hex-codec | Read/Write/Status, CTDMA スケジューリング, プロデューサー/コンシューマー |
| **Interbus** | CmdBridge | Read/Decode, リングトポロジ, IBS CMD フレーム |
| **WorldFIP** | CmdBridge | Read/Write/Decode, プロデューサー/コンシューマー, バスアービタ |
| **Lightbus** | CmdBridge | Read/Decode, 光ファイバ相互接続, 32ノード |
| **Modbus Plus** | CmdBridge + hex-codec | Read/Write/Decode, トークンパッシング, ピアツーピア |
| **ISA100.11a** | CmdBridge + hex-codec | Read/Write/List, 6LoWPAN ワイヤレス, メッシュ |
| **WirelessHART** | CmdBridge + hex-codec | Read/Write/Scan, TSMP MAC, 自己組織化 |
| **SAE J1850** | CmdBridge | 車載診断をCmdBridgeに変換, J1850 インターフェースが必要 |

### システムバス（3/3 — すべて sysfs/procfs ドライバ付き）

| プロトコル | SDKタイプ | 機能 |
|------|----------|------|
| **PCI/PCIe** | sysfs — `/sys/bus/pci/devices/<BDF>/config` | コンフィグ空間の読み書き, PipeTransport, CAP_SYS_ADMIN が必要 |
| **VME/VPX** | procfs — `/proc/vme/<slot>` | A16/A24/A32 アドレッシング, PipeTransport, vme_tsi148 モジュールが必要 |
| **CompactPCI** | sysfs — `/sys/bus/pci/devices/<BDF>/config` | PICMG 2.0, ホットプラグ, 3U/6U, cpci_hotplug が必要 |

---

## 使用ガイド

APIリファレンスと使用例（インストール、Modbus読み書き、コネクションプール、MQTT、Vendorプリセット、ゲートウェイ変換、カスタムプロトコル）は独立したドキュメントに移動しました：

> **[API.md](API.md) — 完全なAPIリファレンスと使用ガイド**
>
> 英語版：[API.en.md](../../API.en.md)

---

## プロジェクト構成

```
industrial-protocols-go/
├── go.work                       # 全モジュールを集約するworkspace
├── kernel/                       # コアモジュール
│   ├── connection/               # ConnectionManager + コネクションプール
│   ├── config/                   # ConfigRepository (YAML)
│   ├── gateway/                  # GatewayEngine プロトコル変換
│   ├── bridge/                   # Bridge 外部プロセスブリッジ
│   ├── vendor/                   # Vendor ベンダープリセット
│   ├── event/                    # イベントバス
│   ├── metrics/                  # メトリクスインターフェース
│   ├── security/                 # TLS セキュリティ
│   ├── pipeline/                 # ミドルウェアチェーン
│   └── transport/                # TCP/UDP/Pipe トランスポート
│
├── protocols/                    # 40プロトコルモジュール
│   ├── ethernet/                 # 5産業イーサネット（すべて完了）
│   ├── fieldbus/                 # 11フィールドバス（4純ソフト + 7ハードウェアSDK）
│   ├── iot/                      # 2 IoT/メッセージング（すべて完了）
│   ├── automotive/               # 5自動車バス（2純ソフト + 3ハードウェアSDK）
│   ├── building/                 # 2ビル/照明（1純ソフト + 1ハードウェアSDK）
│   ├── bridge/                   # 12ハードウェアブリッジ（CmdBridge/SerialBridge）
│   └── system/                   # 3システムバス（sysfs/procfs ドライバ）
│
├── examples/modbus_basic/        # Modbus TCP サンプル
├── _tools/                       # Makefile + 補助スクリプト
└── docs/superpowers/             # 設計ドキュメント
```

## テスト

```bash
make test         # 全量テスト
make test-unit    # ユニットテストのみ
make vet && make fmt
```

---

## サポート

このプロジェクトがお役に立ったなら、寄付で開発の継続をサポートしていただけると幸いです。

### Alipay / WeChat

| Alipay | WeChat |
|--------|------|
| ![支付宝](../../alipay.png) | ![微信](../../weixinpay.png) |

### 海外送金（国際銀行送金）

中国本土以外のユーザー向けの銀行送金：

**受取人情報：**

| 項目 | 内容 |
|------|------|
| 受取人名義 | WANG KEXUN |
| 受取口座番号 | 881015918251 |

**受取銀行：**

| 項目 | 内容 |
|------|------|
| 銀行名 | ZA Bank Limited |
| SWIFT Code | AABLHKHHXXX |
| 銀行コード | 387 |
| 銀行住所 | Core F, Cyberport 3, 100 Cyberport Road, Hong Kong |

**クロスボーダー送金代理銀行（中継銀行、必要な場合）：**

> ご注意：これはクロスボーダー送金の代理銀行（中継銀行）の情報であり、受取銀行の情報ではありません。代理銀行の情報が必要かどうかは、送金銀行にご確認ください。

- **香港ドル・人民元・米ドルの送金** — 代理銀行は Citibank：
  - 銀行名：Citibank N.A. Hong Kong
  - SWIFT Code：CITIHKHXXXX
  - 銀行コード：006
  - 支店名：Hong Kong Branch
  - 支店コード：391
  - 銀行住所：Citibank Tower, Citibank Plaza, 3 Garden Road, Central, Hong Kong
- **その他の通貨の送金** — 代理銀行は BNY Mellon：
  - 銀行名：THE BANK OF NEW YORK MELLON
  - SWIFT Code：IRVTUS3NXXX
  - 銀行住所：THE BANK OF NEW YORK MELLON, 240 GREENWICH STREET, NEW YORK, United States

---

## License

MIT — Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz
