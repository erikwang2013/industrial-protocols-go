[English](../../README.en.md) | [中文](../../README.md) | [한국어](../ko/README.md) | [Русский](../ru/README.md) | [Deutsch](../de/README.md) | [Français](../fr/README.md) | [Español](../es/README.md) | [Português](../pt/README.md) | [हिन्दी](../hi/README.md) | [العربية](../ar/README.md) | [বাংলা](../bn/README.md) | Bahasa Indonesia | [日本語](../ja/README.md)

# Industrial Protocols Go

Kumpulan protokol komunikasi jaringan industri dalam Go — arsitektur berlapis + middleware, mencakup 40 protokol industri, 14 implementasi murni perangkat lunak + 26 SDK perangkat keras (semuanya dengan driver).

> Referensi implementasi PHP: [github.com/erikwang2013/industrial-protocols](https://github.com/erikwang2013/industrial-protocols)

---

## Desain Proyek

### Lapisan Arsitektur

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
│  Modul terpisah per protokol       (TCP/UDP/Serial)│
│                                                │
│  14 Pure-Soft   +   26 Hardware SDKs (with drivers)│
└──────────────────────────────────────────────────┘
```

### Konsep Desain

**Mikrokernel + SDK protokol.** Kernel hanya mendefinisikan antarmuka dan perhatian lintas-potong (connection pool, retry, circuit breaker, event, metrik), tanpa menyertakan implementasi protokol apa pun. Setiap protokol adalah modul Go terpisah yang diimpor sesuai kebutuhan.

**Pemisahan lapisan.** Empat lapisan abstraksi:

| Lapisan | Tanggung Jawab | Tipe Inti |
|----|------|---------|
| **Transport** | Saluran komunikasi tingkat bawah | Antarmuka `Transport` — `TCPTransport`, `UDPTransport`, `PipeTransport` |
| **Codec** | Pengodean/dekode protokol | Antarmuka `Codec` — `Encode(*Request) ([]byte, error)` / `Decode([]byte) (*Response, error)` |
| **Pipeline** | Middleware lintas-potong | `Middleware` — `Timeout`, `Retry`, `CircuitBreaker`, `Logger` |
| **Kernel** | Siklus hidup perangkat | `ConnectionManager`, `ConfigRepository` |

**Penulis protokol hanya perlu mengimplementasikan dua antarmuka: `Protocol` + `Codec`.** Lapisan transport, connection pool, retry, timeout, dan circuit breaker semuanya dipakai ulang dari kernel.

### Modul Kernel

| Modul | Jalur | Deskripsi |
|------|------|------|
| ConnectionManager | `kernel/connection/` | Registrasi perangkat, connection pool, health check, mendukung tiga strategi: Lazy/Eager/Pooled |
| ConfigRepository | `kernel/config/` | Pemuatan konfigurasi perangkat YAML/JSON |
| GatewayEngine | `kernel/gateway/` | Mesin aturan konversi antar-protokol (Modbus→MQTT, dll.) |
| Bridge | `kernel/bridge/` | Jembatan proses eksternal (komunikasi stdin/stdout) |
| Vendor | `kernel/vendor/` | Prasetel parameter vendor (Siemens S7, Rockwell AB, dll.) |
| Event | `kernel/event/` | Bus peristiwa berbasis channel |
| Metrics | `kernel/metrics/` | Antarmuka pengumpulan metrik (Prometheus + noop) |
| Security | `kernel/security/` | Pembungkus keamanan lapisan transport TLS |

---

## Dukungan Protokol

### Ethernet Industri (5/5 selesai)

| Protokol | Transport | Port | Fungsi |
|------|------|------|------|
| **Modbus** | TCP, RTU | 502 | FC 01/02/03/04/05/06/16, CRC16, MBAP |
| **BACnet** | UDP | 47808 | Penemuan perangkat Who-Is/I-Am, ReadProperty |
| **EtherNet/IP** | TCP | 44818 | ENIP RegisterSession + CIP Read Tag |
| **OPC UA** | TCP | 4840 | Binary HEL + OpenSecureChannel |
| **PROFINET** | UDP | 34964 | DCP Identify/Set + Record Data Read/Write |

### Fieldbus (11/11 — 4 murni perangkat lunak + 7 SDK perangkat keras)

| Protokol | Transport | Port | Fungsi |
|------|------|------|------|
| **HART** | Serial FSK | — | Frame pendek/panjang, Command 0/3, checksum XOR |
| **CC-Link** | RS-485 | — | Polling master-slave, CRC-16/XMODEM |
| **DNP3** | TCP, Serial | 20000 | Segmentasi/reassembly lapisan transport, polling Class 0, CRC-16/DNP |
| **IEC 61850** | TCP | 102 | MMS Initiate/Conclude, Read/Write, BER-TLV |
| **PROFIBUS** | Serial | — | Implementasi murni perangkat lunak — memerlukan perangkat keras CP 5611 |
| **CANopen** | CAN, Gateway | — | Baca/tulis SDO, NMT start/stop/reset, Heartbeat |
| **DeviceNet** | CAN, Gateway | — | Explicit Messaging, Poll, koneksi I/O |
| **Foundation Fieldbus** | Serial | — | Implementasi murni perangkat lunak — memerlukan kartu antarmuka FF H1 |
| **AS-Interface** | Serial | — | Implementasi murni perangkat lunak — memerlukan gateway ASi |
| **IO-Link** | Serial | — | Implementasi murni perangkat lunak — memerlukan IO-Link Master |
| **CC-Link IE** | Ethernet | — | Implementasi murni perangkat lunak — memerlukan gateway CC-Link IE |

### IoT / Pesan (2/2 selesai)

| Protokol | Transport | Port | Fungsi |
|------|------|------|------|
| **MQTT** | TCP, WS | 1883 | 3.1.1 CONNECT/PUBLISH/SUBSCRIBE/PING |
| **HART-IP** | TCP, UDP | 5094 | HART over TCP, memakai ulang codec HART |

### Bus Otomotif (5/5 — 2 murni perangkat lunak + 3 SDK perangkat keras)

| Protokol | Transport | Baudrate | Fungsi |
|------|------|--------|------|
| **LIN** | UART | — | Frame master-slave, checksum PID, Classic/Enhanced Checksum |
| **K-Line** | Serial | 10400 | ISO 9141/14230, 5-baud Fast Init, OBD-II SID 01/03/09 |
| **FlexRay** | CAN, Serial | — | Encode/decode Slot/Frame, CRC-16/XMODEM, SocketCAN + SerialBridge |
| **SAE J1850** | CAN | — | Encode/decode PWM/VPW, Mode 01/03/0A, CRC-8 |
| **MOST** | Serial | — | Encode/decode Hex-frame, Read/Write/Status, SerialBridge |

### Gedung / Pencahayaan (2/2 — 1 murni perangkat lunak + 1 SDK perangkat keras)

| Protokol | Transport | Fungsi |
|------|------|------|
| **DALI** | Serial | Frame maju 16-bit, frame balik 8-bit, perintah standar (Off/Max/Dim) |
| **LonWorks** | Serial | Encode/decode murni perangkat lunak — memerlukan chip Neuron atau gateway |

### Jembatan Perangkat Keras (12/12 — semuanya memiliki driver CmdBridge/SerialBridge)

| Protokol | Tipe SDK | Fungsi |
|------|----------|------|
| **EtherCAT** | CmdBridge + hex-codec | Subperintah Upload/Download/Slaves, encode/decode heksadesimal CoE |
| **POWERLINK** | CmdBridge + hex-codec | Subperintah Read/Write/Status, frame SoC/Preq/Pres |
| **SERCOS III** | CmdBridge + hex-codec | Subperintah Read/Write/Phase, transisi fase (NRT→CP4) |
| **SERCOS I/II** | CmdBridge + hex-codec | Read/Write/Status, topologi cincin serat optik |
| **ControlNet** | CmdBridge + hex-codec | Read/Write/Status, penjadwalan CTDMA, produsen/konsumen |
| **Interbus** | CmdBridge | Read/Decode, topologi cincin, frame IBS CMD |
| **WorldFIP** | CmdBridge | Read/Write/Decode, produsen/konsumen, arbiter bus |
| **Lightbus** | CmdBridge | Read/Decode, interkoneksi serat optik, 32 node |
| **Modbus Plus** | CmdBridge + hex-codec | Read/Write/Decode, token passing, peer-to-peer |
| **ISA100.11a** | CmdBridge + hex-codec | Read/Write/List, nirkabel 6LoWPAN, mesh |
| **WirelessHART** | CmdBridge + hex-codec | Read/Write/Scan, TSMP MAC, self-organizing |
| **SAE J1850** | CmdBridge | Diagnostik otomotif dialihkan ke CmdBridge, memerlukan antarmuka J1850 |

### Bus Sistem (3/3 — semuanya memiliki driver sysfs/procfs)

| Protokol | Tipe SDK | Fungsi |
|------|----------|------|
| **PCI/PCIe** | sysfs — `/sys/bus/pci/devices/<BDF>/config` | Baca/tulis ruang konfigurasi, PipeTransport, memerlukan CAP_SYS_ADMIN |
| **VME/VPX** | procfs — `/proc/vme/<slot>` | Pengalamatan A16/A24/A32, PipeTransport, memerlukan modul vme_tsi148 |
| **CompactPCI** | sysfs — `/sys/bus/pci/devices/<BDF>/config` | PICMG 2.0, hot-swap, 3U/6U, memerlukan cpci_hotplug |

---

## Panduan Penggunaan

Referensi API dan contoh penggunaan (instalasi, baca/tulis Modbus, connection pool, MQTT, prasetel Vendor, konversi gateway, protokol kustom) telah dipindahkan ke dokumen terpisah:

> **[API.md](API.md) — Referensi API lengkap dan panduan penggunaan**
>
> Versi Inggris: [API.en.md](../../API.en.md)

---

## Struktur Proyek

```
industrial-protocols-go/
├── go.work                       # workspace yang menggabungkan semua modul
├── kernel/                       # modul inti
│   ├── connection/               # ConnectionManager + connection pool
│   ├── config/                   # ConfigRepository (YAML)
│   ├── gateway/                  # GatewayEngine konversi protokol
│   ├── bridge/                   # Bridge jembatan proses eksternal
│   ├── vendor/                   # Vendor prasetel vendor
│   ├── event/                    # bus peristiwa
│   ├── metrics/                  # antarmuka metrik
│   ├── security/                 # keamanan TLS
│   ├── pipeline/                 # rantai middleware
│   └── transport/                # transport TCP/UDP/Pipe
│
├── protocols/                    # 40 modul protokol
│   ├── ethernet/                 # 5 Ethernet industri (semua selesai)
│   ├── fieldbus/                 # 11 fieldbus (4 murni perangkat lunak + 7 SDK perangkat keras)
│   ├── iot/                      # 2 IoT/pesan (semua selesai)
│   ├── automotive/               # 5 bus otomotif (2 murni perangkat lunak + 3 SDK perangkat keras)
│   ├── building/                 # 2 gedung/pencahayaan (1 murni perangkat lunak + 1 SDK perangkat keras)
│   ├── bridge/                   # 12 jembatan perangkat keras (CmdBridge/SerialBridge)
│   └── system/                   # 3 bus sistem (driver sysfs/procfs)
│
├── examples/modbus_basic/        # contoh Modbus TCP
├── _tools/                       # Makefile + skrip bantu
└── docs/superpowers/             # dokumen desain
```

## Pengujian

```bash
make test         # pengujian penuh
make test-unit    # hanya unit test
make vet && make fmt
```

---

## Dukung Kami

Jika proyek ini bermanfaat bagi Anda, silakan beri donasi untuk mendukung kami terus memeliharanya.

### Alipay / WeChat

| Alipay | WeChat |
|--------|------|
| ![支付宝](../../alipay.png) | ![微信](../../weixinpay.png) |

### Transfer Global (Transfer Bank Internasional)

Untuk transfer bank dari pengguna di luar Tiongkok Daratan:

**Informasi Penerima:**

| Item | Isi |
|------|------|
| Nama Penerima | WANG KEXUN |
| Nomor Rekening Penerima | 881015918251 |

**Bank Penerima:**

| Item | Isi |
|------|------|
| Nama Bank | ZA Bank Limited |
| SWIFT Code | AABLHKHHXXX |
| Kode Bank | 387 |
| Alamat Bank | Core F, Cyberport 3, 100 Cyberport Road, Hong Kong |

**Bank Perantara Transfer Lintas Batas (jika diperlukan):**

> Perlu diperhatikan, ini adalah informasi bank perantara transfer lintas batas, bukan bank penerima. Silakan tanyakan ke bank pengirim apakah informasi bank perantara diperlukan.

- **Untuk transfer dalam HKD, CNY, dan USD** — bank perantara adalah Citibank:
  - Nama Bank: Citibank N.A. Hong Kong
  - SWIFT Code: CITIHKHXXXX
  - Kode Bank: 006
  - Nama Cabang: Hong Kong Branch
  - Kode Cabang: 391
  - Alamat Bank: Citibank Tower, Citibank Plaza, 3 Garden Road, Central, Hong Kong
- **Untuk transfer mata uang lainnya** — bank perantara adalah BNY Mellon:
  - Nama Bank: THE BANK OF NEW YORK MELLON
  - SWIFT Code: IRVTUS3NXXX
  - Alamat Bank: THE BANK OF NEW YORK MELLON, 240 GREENWICH STREET, NEW YORK, United States

---

## License

MIT — Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz
