[English](../../README.en.md) | [中文](../../README.md) | [한국어](../ko/README.md) | [Русский](../ru/README.md) | [Deutsch](../de/README.md) | [Français](../fr/README.md) | [Español](../es/README.md) | [Português](../pt/README.md) | [हिन्दी](../hi/README.md) | [العربية](../ar/README.md) | বাংলা | [Bahasa Indonesia](../id/README.md) | [日本語](../ja/README.md)

# Industrial Protocols Go

Go ভাষায় শিল্প নেটওয়ার্ক কমিউনিকেশন প্রোটোকল সংকলন —— স্তরভিত্তিক + মিডলওয়্যার আর্কিটেকচার, ৪০টি শিল্প প্রোটোকল কভার করে, ১৪টি সম্পূর্ণ সফটওয়্যার বাস্তবায়ন + ২৬টি হার্ডওয়্যার SDK (সবগুলোই ড্রাইভারসহ)।

> PHP বাস্তবায়ন দেখুন: [github.com/erikwang2013/industrial-protocols](https://github.com/erikwang2013/industrial-protocols)

---

## প্রকল্প নকশা

### আর্কিটেকচারের স্তর

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

### নকশা দর্শন

**মাইক্রো-কার্নেল + প্রোটোকল SDK।** কার্নেল শুধুমাত্র ইন্টারফেস এবং ক্রস-কাটিং কনসার্ন (কানেকশন পুল, রিট্রাই, সার্কিট ব্রেকার, ইভেন্ট, মেট্রিক্স) সংজ্ঞায়িত করে, কোনো নির্দিষ্ট প্রোটোকল বাস্তবায়ন ধারণ করে না। প্রতিটি প্রোটোকল একটি স্বাধীন Go মডিউল, প্রয়োজন অনুযায়ী ইমপোর্ট করা হয়।

**স্তরভিত্তিক ডিকপলিং।** চার স্তরের অ্যাবস্ট্রাকশন:

| স্তর | দায়িত্ব | মূল টাইপ |
|----|------|---------|
| **Transport** | অন্তর্নিহিত যোগাযোগ চ্যানেল | `Transport` ইন্টারফেস — `TCPTransport`, `UDPTransport`, `PipeTransport` |
| **Codec** | প্রোটোকল এনকোড/ডিকোড | `Codec` ইন্টারফেস — `Encode(*Request) ([]byte, error)` / `Decode([]byte) (*Response, error)` |
| **Pipeline** | ক্রস-কাটিং মিডলওয়্যার | `Middleware` — `Timeout`, `Retry`, `CircuitBreaker`, `Logger` |
| **Kernel** | ডিভাইস লাইফসাইকেল | `ConnectionManager`, `ConfigRepository` |

**প্রোটোকল লেখককে শুধুমাত্র `Protocol` + `Codec` দুটি ইন্টারফেস বাস্তবায়ন করতে হয়।** ট্রান্সপোর্ট স্তর, কানেকশন পুল, রিট্রাই, টাইমআউট, সার্কিট ব্রেকার — সবগুলোই kernel থেকে পুনর্ব্যবহৃত হয়।

### কার্নেল মডিউল

| মডিউল | পাথ | বিবরণ |
|------|------|------|
| ConnectionManager | `kernel/connection/` | ডিভাইস নিবন্ধন, কানেকশন পুল, হেলথ চেক; Lazy/Eager/Pooled তিনটি কৌশল সমর্থন করে |
| ConfigRepository | `kernel/config/` | YAML/JSON ডিভাইস কনফিগ লোডিং |
| GatewayEngine | `kernel/gateway/` | ক্রস-প্রোটোকল রূপান্তর নিয়ম ইঞ্জিন (Modbus→MQTT ইত্যাদি) |
| Bridge | `kernel/bridge/` | বাহ্যিক প্রসেস ব্রিজিং (stdin/stdout কমিউনিকেশন) |
| Vendor | `kernel/vendor/` | ভেন্ডর প্যারামিটার প্রিসেট (Siemens S7, Rockwell AB ইত্যাদি) |
| Event | `kernel/event/` | channel-ভিত্তিক ইভেন্ট বাস |
| Metrics | `kernel/metrics/` | মেট্রিক্স সংগ্রহ ইন্টারফেস (Prometheus + noop) |
| Security | `kernel/security/` | TLS ট্রান্সপোর্ট লেয়ার সিকিউরিটি র‍্যাপার |

---

## প্রোটোকল সমর্থন

### শিল্প ইথারনেট (৫/৫ সম্পন্ন)

| প্রোটোকল | ট্রান্সপোর্ট | পোর্ট | কার্যকারিতা |
|------|------|------|------|
| **Modbus** | TCP, RTU | 502 | FC 01/02/03/04/05/06/16, CRC16, MBAP |
| **BACnet** | UDP | 47808 | Who-Is/I-Am ডিভাইস ডিসকভারি, ReadProperty |
| **EtherNet/IP** | TCP | 44818 | ENIP RegisterSession + CIP Read Tag |
| **OPC UA** | TCP | 4840 | Binary HEL + OpenSecureChannel |
| **PROFINET** | UDP | 34964 | DCP Identify/Set + Record Data Read/Write |

### ফিল্ডবাস (১১/১১ — ৪টি পিওর-সফট + ৭টি হার্ডওয়্যার SDK)

| প্রোটোকল | ট্রান্সপোর্ট | পোর্ট | কার্যকারিতা |
|------|------|------|------|
| **HART** | Serial FSK | — | শর্ট/লং ফ্রেম, Command 0/3, XOR চেকসাম |
| **CC-Link** | RS-485 | — | মাস্টার-স্লেভ পোলিং, CRC-16/XMODEM |
| **DNP3** | TCP, Serial | 20000 | ট্রান্সপোর্ট লেয়ার সেগমেন্টেশন/রিঅ্যাসেম্বলি, Class 0 পোলিং, CRC-16/DNP |
| **IEC 61850** | TCP | 102 | MMS Initiate/Conclude, Read/Write, BER-TLV |
| **PROFIBUS** | Serial | — | পিওর-সফট বাস্তবায়ন — CP 5611 হার্ডওয়্যার প্রয়োজন |
| **CANopen** | CAN, Gateway | — | SDO পড়া/লেখা, NMT চালু/বন্ধ/রিসেট, Heartbeat |
| **DeviceNet** | CAN, Gateway | — | Explicit Messaging, Poll, I/O কানেকশন |
| **Foundation Fieldbus** | Serial | — | পিওর-সফট বাস্তবায়ন — FF H1 ইন্টারফেস কার্ড প্রয়োজন |
| **AS-Interface** | Serial | — | পিওর-সফট বাস্তবায়ন — ASi গেটওয়ে প্রয়োজন |
| **IO-Link** | Serial | — | পিওর-সফট বাস্তবায়ন — IO-Link Master প্রয়োজন |
| **CC-Link IE** | Ethernet | — | পিওর-সফট বাস্তবায়ন — CC-Link IE গেটওয়ে প্রয়োজন |

### IoT / মেসেজিং (২/২ সম্পন্ন)

| প্রোটোকল | ট্রান্সপোর্ট | পোর্ট | কার্যকারিতা |
|------|------|------|------|
| **MQTT** | TCP, WS | 1883 | 3.1.1 CONNECT/PUBLISH/SUBSCRIBE/PING |
| **HART-IP** | TCP, UDP | 5094 | TCP-র ওপর HART, HART এনকোড/ডিকোড পুনর্ব্যবহার |

### অটোমোটিভ বাস (৫/৫ — ২টি পিওর-সফট + ৩টি হার্ডওয়্যার SDK)

| প্রোটোকল | ট্রান্সপোর্ট | বড রেট | কার্যকারিতা |
|------|------|--------|------|
| **LIN** | UART | — | মাস্টার/স্লেভ ফ্রেম, PID চেকসাম, Classic/Enhanced Checksum |
| **K-Line** | Serial | 10400 | ISO 9141/14230, 5-baud Fast Init, OBD-II SID 01/03/09 |
| **FlexRay** | CAN, Serial | — | Slot/Frame এনকোড/ডিকোড, CRC-16/XMODEM, SocketCAN + SerialBridge |
| **SAE J1850** | CAN | — | PWM/VPW এনকোড/ডিকোড, Mode 01/03/0A, CRC-8 |
| **MOST** | Serial | — | Hex-frame এনকোড/ডিকোড, Read/Write/Status, SerialBridge |

### বিল্ডিং / লাইটিং (২/২ — ১টি পিওর-সফট + ১টি হার্ডওয়্যার SDK)

| প্রোটোকল | ট্রান্সপোর্ট | কার্যকারিতা |
|------|------|------|
| **DALI** | Serial | 16-bit ফরোয়ার্ড ফ্রেম, 8-bit ব্যাকওয়ার্ড ফ্রেম, স্ট্যান্ডার্ড কমান্ড (Off/Max/Dim) |
| **LonWorks** | Serial | পিওর-সফট এনকোড/ডিকোড — Neuron চিপ বা গেটওয়ে প্রয়োজন |

### হার্ডওয়্যার ব্রিজিং (১২/১২ — সবগুলোতেই CmdBridge/SerialBridge ড্রাইভার)

| প্রোটোকল | SDK ধরন | কার্যকারিতা |
|------|----------|------|
| **EtherCAT** | CmdBridge + hex-codec | Upload/Download/Slaves সাবকমান্ড, CoE হেক্সাডেসিমেল এনকোড/ডিকোড |
| **POWERLINK** | CmdBridge + hex-codec | Read/Write/Status সাবকমান্ড, SoC/Preq/Pres ফ্রেম |
| **SERCOS III** | CmdBridge + hex-codec | Read/Write/Phase সাবকমান্ড, ফেজ ট্রানজিশন (NRT→CP4) |
| **SERCOS I/II** | CmdBridge + hex-codec | Read/Write/Status, ফাইবার-অপটিক রিং টপোলজি |
| **ControlNet** | CmdBridge + hex-codec | Read/Write/Status, CTDMA শিডিউলিং, প্রোডিউসার/কনজিউমার |
| **Interbus** | CmdBridge | Read/Decode, রিং টপোলজি, IBS CMD ফ্রেম |
| **WorldFIP** | CmdBridge | Read/Write/Decode, প্রোডিউসার/কনজিউমার, বাস অরবিট্রেটর |
| **Lightbus** | CmdBridge | Read/Decode, ফাইবার-অপটিক ইন্টারকানেক্ট, ৩২ নোড |
| **Modbus Plus** | CmdBridge + hex-codec | Read/Write/Decode, টোকেন পাসিং, পিয়ার-টু-পিয়ার |
| **ISA100.11a** | CmdBridge + hex-codec | Read/Write/List, 6LoWPAN ওয়্যারলেস, মেশ টপোলজি |
| **WirelessHART** | CmdBridge + hex-codec | Read/Write/Scan, TSMP MAC, সেলফ-অর্গানাইজিং |
| **SAE J1850** | CmdBridge | ভেহিকেল ডায়াগনস্টিককে CmdBridge-তে রূপান্তর, J1850 ইন্টারফেস প্রয়োজন |

### সিস্টেম বাস (৩/৩ — সবগুলোতেই sysfs/procfs ড্রাইভার)

| প্রোটোকল | SDK ধরন | কার্যকারিতা |
|------|----------|------|
| **PCI/PCIe** | sysfs — `/sys/bus/pci/devices/<BDF>/config` | কনফিগ স্পেস পড়া/লেখা, PipeTransport, CAP_SYS_ADMIN প্রয়োজন |
| **VME/VPX** | procfs — `/proc/vme/<slot>` | A16/A24/A32 অ্যাড্রেসিং, PipeTransport, vme_tsi148 মডিউল প্রয়োজন |
| **CompactPCI** | sysfs — `/sys/bus/pci/devices/<BDF>/config` | PICMG 2.0, হট-প্লাগ, 3U/6U, cpci_hotplug প্রয়োজন |

---

## ব্যবহার নির্দেশিকা

API রেফারেন্স ও ব্যবহারের উদাহরণ (ইনস্টলেশন, Modbus পড়া/লেখা, কানেকশন পুল, MQTT, Vendor প্রিসেট, গেটওয়ে রূপান্তর, কাস্টম প্রোটোকল) আলাদা ডকুমেন্টে স্থানান্তরিত হয়েছে:

> **[API.md](API.md) — সম্পূর্ণ API রেফারেন্স ও ব্যবহার নির্দেশিকা**
>
> ইংরেজি সংস্করণ: [API.en.md](../../API.en.md)

---

## প্রকল্প কাঠামো

```
industrial-protocols-go/
├── go.work                       # ওয়ার্কস্পেস — সব মডিউল একত্রিত করে
├── kernel/                       # কোর মডিউল
│   ├── connection/               # ConnectionManager + কানেকশন পুল
│   ├── config/                   # ConfigRepository (YAML)
│   ├── gateway/                  # GatewayEngine প্রোটোকল রূপান্তর
│   ├── bridge/                   # Bridge বাহ্যিক প্রসেস ব্রিজিং
│   ├── vendor/                   # Vendor ভেন্ডর প্রিসেট
│   ├── event/                    # ইভেন্ট বাস
│   ├── metrics/                  # মেট্রিক্স ইন্টারফেস
│   ├── security/                 # TLS নিরাপত্তা
│   ├── pipeline/                 # মিডলওয়্যার চেইন
│   └── transport/                # TCP/UDP/Pipe ট্রান্সপোর্ট
│
├── protocols/                    # 40টি প্রোটোকল মডিউল
│   ├── ethernet/                 # 5টি শিল্প ইথারনেট (সব সম্পন্ন)
│   ├── fieldbus/                 # 11টি ফিল্ডবাস (4 পিওর-সফট + 7 হার্ডওয়্যার SDK)
│   ├── iot/                      # 2টি IoT/মেসেজিং (সব সম্পন্ন)
│   ├── automotive/               # 5টি অটোমোটিভ বাস (2 পিওর-সফট + 3 হার্ডওয়্যার SDK)
│   ├── building/                 # 2টি বিল্ডিং/লাইটিং (1 পিওর-সফট + 1 হার্ডওয়্যার SDK)
│   ├── bridge/                   # 12টি হার্ডওয়্যার ব্রিজ (CmdBridge/SerialBridge)
│   └── system/                   # 3টি সিস্টেম বাস (sysfs/procfs ড্রাইভার)
│
├── examples/modbus_basic/        # Modbus TCP উদাহরণ
├── _tools/                       # Makefile + সহায়ক স্ক্রিপ্ট
└── docs/superpowers/             # ডিজাইন ডকুমেন্ট
```

## পরীক্ষা

```bash
make test         # সম্পূর্ণ টেস্ট
make test-unit    # শুধুমাত্র ইউনিট টেস্ট
make vet && make fmt
```
---

## আমাদের সমর্থন করুন

যদি এই প্রকল্পটি আপনার কাজে লেগে থাকে, তাহলে আমাদের রক্ষণাবেক্ষণ চালিয়ে যেতে দান করে সাহায্য করতে পারেন।

### Alipay / WeChat

| আলিপে | উইচ্যাট |
|--------|------|
| ![支付宝](../../alipay.png) | ![微信](../../weixinpay.png) |

### আন্তর্জাতিক ব্যাংক ট্রান্সফার

চীনের মূল ভূখণ্ডের বাইরের ব্যবহারকারীদের জন্য ব্যাংক ট্রান্সফার:

**প্রাপকের তথ্য:**

| আইটেম | বিবরণ |
|------|------|
| প্রাপকের নাম | WANG KEXUN |
| অ্যাকাউন্ট নম্বর | 881015918251 |

**প্রাপক ব্যাংক:**

| আইটেম | বিবরণ |
|------|------|
| ব্যাংকের নাম | ZA Bank Limited |
| SWIFT Code | AABLHKHHXXX |
| ব্যাংক কোড | 387 |
| ব্যাংকের ঠিকানা | Core F, Cyberport 3, 100 Cyberport Road, Hong Kong |

**ক্রস-বর্ডার রেমিট্যান্স করেসপন্ডেন্ট ব্যাংক (মধ্যস্থ ব্যাংক, প্রয়োজনে):**

> মনে রাখবেন, এটি ক্রস-বর্ডার রেমিট্যান্স করেসপন্ডেন্ট ব্যাংক (মধ্যস্থ ব্যাংক) এর তথ্য, প্রাপক ব্যাংকের তথ্য নয়। করেসপন্ডেন্ট ব্যাংকের তথ্য প্রদানের প্রয়োজন আছে কিনা তা প্রেরক ব্যাংককে জিজ্ঞাসা করুন।

- **হংকং ডলার, চীনা ইউয়ান ও মার্কিন ডলার জমা** — করেসপন্ডেন্ট ব্যাংক Citibank:
  - ব্যাংকের নাম: Citibank N.A. Hong Kong
  - SWIFT Code: CITIHKHXXXX
  - ব্যাংক কোড: 006
  - শাখার নাম: Hong Kong Branch
  - শাখা কোড: 391
  - ব্যাংকের ঠিকানা: Citibank Tower, Citibank Plaza, 3 Garden Road, Central, Hong Kong
- **অন্যান্য মুদ্রা জমা** — করেসপন্ডেন্ট ব্যাংক BNY Mellon:
  - ব্যাংকের নাম: THE BANK OF NEW YORK MELLON
  - SWIFT Code: IRVTUS3NXXX
  - ব্যাংকের ঠিকানা: THE BANK OF NEW YORK MELLON, 240 GREENWICH STREET, NEW YORK, United States

---

## লাইসেন্স

MIT — Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz
