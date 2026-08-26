[English](../../README.en.md) | [中文](../../README.md) | [한국어](../ko/README.md) | [Русский](../ru/README.md) | [Deutsch](../de/README.md) | [Français](../fr/README.md) | [Español](../es/README.md) | [Português](../pt/README.md) | [हिन्दी](../hi/README.md) | العربية | [বাংলা](../bn/README.md) | [Bahasa Indonesia](../id/README.md) | [日本語](../ja/README.md)

# Industrial Protocols Go

حزمة بروتوكولات الاتصال الصناعي الشبكية بلغة Go — بنية طبقية + وسيطات (Middleware)، تغطي 40 بروتوكولًا صناعيًا: 14 تنفيذًا برمجيًا خالصًا + 26 SDK للأجهزة (جميعها مع برامج تشغيل).

> المرجع في PHP: [github.com/erikwang2013/industrial-protocols](https://github.com/erikwang2013/industrial-protocols)

---

## تصميم المشروع

### الطبقات المعمارية

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
│  وحدة مستقلة لكل بروتوكول            (TCP/UDP/Serial)  │
│                                                │
│  14 Pure-Soft   +   26 Hardware SDKs (with drivers)│
└──────────────────────────────────────────────────┘
```

### فلسفة التصميم

**نواة مصغّرة + SDK للبروتوكولات.** تُعرّف النواة الواجهات والاهتمامات العرضية فقط (تجمع الاتصالات، إعادة المحاولة، قاطع الدائرة، الأحداث، المقاييس)، ولا تتضمن أي تنفيذ بروتوكول محدد. كل بروتوكول هو وحدة Go مستقلة، تُستورد عند الحاجة.

**فصل الطبقات.** أربع طبقات من التجريد:

| الطبقة | المسؤولية | الأنواع الأساسية |
|----|------|---------|
| **Transport** | قناة الاتصال السفلى | واجهة `Transport` — `TCPTransport`, `UDPTransport`, `PipeTransport` |
| **Codec** | ترميز/فك ترميز البروتوكول | واجهة `Codec` — `Encode(*Request) ([]byte, error)` / `Decode([]byte) (*Response, error)` |
| **Pipeline** | الوسيطات العرضية | `Middleware` — `Timeout`, `Retry`, `CircuitBreaker`, `Logger` |
| **Kernel** | دورة حياة الأجهزة | `ConnectionManager`, `ConfigRepository` |

**يكفي أن ينفّذ مؤلف البروتوكول الواجهتين `Protocol` + `Codec` فقط.** طبقة النقل، تجمع الاتصالات، إعادة المحاولة، المهلة، وقاطع الدائرة تُعاد جميعها من kernel.

### وحدات النواة

| الوحدة | المسار | الوصف |
|------|------|------|
| ConnectionManager | `kernel/connection/` | تسجيل الأجهزة، تجمع الاتصالات، فحص الصحة، يدعم ثلاث استراتيجيات Lazy/Eager/Pooled |
| ConfigRepository | `kernel/config/` | تحميل إعدادات الأجهزة من YAML/JSON |
| GatewayEngine | `kernel/gateway/` | محرك قواعد التحويل بين البروتوكولات (Modbus→MQTT وغيرها) |
| Bridge | `kernel/bridge/` | جسر العمليات الخارجية (اتصال عبر stdin/stdout) |
| Vendor | `kernel/vendor/` | إعدادات مُسبقة لمعاملات المصنّعين (Siemens S7, Rockwell AB وغيرها) |
| Event | `kernel/event/` | ناقل أحداث قائم على القنوات (channel-based) |
| Metrics | `kernel/metrics/` | واجهة جمع المقاييس (Prometheus + noop) |
| Security | `kernel/security/` | تغليف أمان طبقة النقل TLS |

---

## دعم البروتوكولات

### الإيثرنت الصناعي (5/5 مكتمل)

| البروتوكول | النقل | المنفذ | الوظائف |
|------|------|------|------|
| **Modbus** | TCP, RTU | 502 | FC 01/02/03/04/05/06/16, CRC16, MBAP |
| **BACnet** | UDP | 47808 | اكتشاف الأجهزة Who-Is/I-Am, ReadProperty |
| **EtherNet/IP** | TCP | 44818 | ENIP RegisterSession + CIP Read Tag |
| **OPC UA** | TCP | 4840 | Binary HEL + OpenSecureChannel |
| **PROFINET** | UDP | 34964 | DCP Identify/Set + Record Data Read/Write |

### ناقل الميدان (11/11 — 4 برمجية خالصة + 7 SDK أجهزة)

| البروتوكول | النقل | المنفذ | الوظائف |
|------|------|------|------|
| **HART** | Serial FSK | — | إطار قصير/طويل, Command 0/3, فحص XOR |
| **CC-Link** | RS-485 | — | استقصاء رئيسي/تابع, CRC-16/XMODEM |
| **DNP3** | TCP, Serial | 20000 | تجزئة/إعادة تجميع طبقة النقل, استقصاء Class 0, CRC-16/DNP |
| **IEC 61850** | TCP | 102 | MMS Initiate/Conclude, Read/Write, BER-TLV |
| **PROFIBUS** | Serial | — | تنفيذ برمجي خالص — يتطلب أجهزة CP 5611 |
| **CANopen** | CAN, Gateway | — | قراءة/كتابة SDO, NMT تشغيل/إيقاف/إعادة ضبط, Heartbeat |
| **DeviceNet** | CAN, Gateway | — | Explicit Messaging, Poll, اتصالات I/O |
| **Foundation Fieldbus** | Serial | — | تنفيذ برمجي خالص — يتطلب بطاقة واجهة FF H1 |
| **AS-Interface** | Serial | — | تنفيذ برمجي خالص — يتطلب بوابة ASi |
| **IO-Link** | Serial | — | تنفيذ برمجي خالص — يتطلب IO-Link Master |
| **CC-Link IE** | Ethernet | — | تنفيذ برمجي خالص — يتطلب بوابة CC-Link IE |

### IoT / الرسائل (2/2 مكتمل)

| البروتوكول | النقل | المنفذ | الوظائف |
|------|------|------|------|
| **MQTT** | TCP, WS | 1883 | 3.1.1 CONNECT/PUBLISH/SUBSCRIBE/PING |
| **HART-IP** | TCP, UDP | 5094 | HART عبر TCP, إعادة استخدام ترميز HART |

### ناقل السيارات (5/5 — 2 برمجية خالصة + 3 SDK أجهزة)

| البروتوكول | النقل | معدل البود | الوظائف |
|------|------|--------|------|
| **LIN** | UART | — | إطارات رئيسية/تابعة, فحص PID, Classic/Enhanced Checksum |
| **K-Line** | Serial | 10400 | ISO 9141/14230, تهيئة سريعة 5-baud, OBD-II SID 01/03/09 |
| **FlexRay** | CAN, Serial | — | ترميز/فك Slot/Frame, CRC-16/XMODEM, SocketCAN + SerialBridge |
| **SAE J1850** | CAN | — | ترميز/فك PWM/VPW, Mode 01/03/0A, CRC-8 |
| **MOST** | Serial | — | ترميز/فك Hex-frame, Read/Write/Status, SerialBridge |

### المباني / الإضاءة (2/2 — 1 برمجية خالصة + 1 SDK أجهزة)

| البروتوكول | النقل | الوظائف |
|------|------|------|
| **DALI** | Serial | إطار أمامي 16-bit, إطار خلفي 8-bit, أوامر قياسية (Off/Max/Dim) |
| **LonWorks** | Serial | ترميز/فك برمجي خالص — يتطلب شريحة Neuron أو بوابة |

### الجسور الإلكترونية (12/12 — جميعها مع برامج تشغيل CmdBridge/SerialBridge)

| البروتوكول | نوع SDK | الوظائف |
|------|----------|------|
| **EtherCAT** | CmdBridge + hex-codec | أوامر فرعية Upload/Download/Slaves, ترميز/فك CoE بالنظام الست عشري |
| **POWERLINK** | CmdBridge + hex-codec | أوامر فرعية Read/Write/Status, إطارات SoC/Preq/Pres |
| **SERCOS III** | CmdBridge + hex-codec | أوامر فرعية Read/Write/Phase, انتقال المراحل (NRT→CP4) |
| **SERCOS I/II** | CmdBridge + hex-codec | Read/Write/Status, طوبولوجيا حلقة الألياف الضوئية |
| **ControlNet** | CmdBridge + hex-codec | Read/Write/Status, جدولة CTDMA, منتج/مستهلك |
| **Interbus** | CmdBridge | Read/Decode, طوبولوجيا حلقية, إطارات IBS CMD |
| **WorldFIP** | CmdBridge | Read/Write/Decode, منتج/مستهلك, محكم ناقل |
| **Lightbus** | CmdBridge | Read/Decode, وصل ألياف ضوئية, 32 عقدة |
| **Modbus Plus** | CmdBridge + hex-codec | Read/Write/Decode, تمرير الرمز (token), نظير-لنظير |
| **ISA100.11a** | CmdBridge + hex-codec | Read/Write/List, لاسلكي 6LoWPAN, شبكي |
| **WirelessHART** | CmdBridge + hex-codec | Read/Write/Scan, TSMP MAC, ذاتي التنظيم |
| **SAE J1850** | CmdBridge | تشخيص المركبات إلى CmdBridge, يتطلب واجهة J1850 |

### ناقل النظام (3/3 — جميعها مع برامج تشغيل sysfs/procfs)

| البروتوكول | نوع SDK | الوظائف |
|------|----------|------|
| **PCI/PCIe** | sysfs — `/sys/bus/pci/devices/<BDF>/config` | قراءة/كتابة مساحة الإعداد, PipeTransport, يتطلب CAP_SYS_ADMIN |
| **VME/VPX** | procfs — `/proc/vme/<slot>` | عنونة A16/A24/A32, PipeTransport, يتطلب وحدة vme_tsi148 |
| **CompactPCI** | sysfs — `/sys/bus/pci/devices/<BDF>/config` | PICMG 2.0, إدخال/إخراج ساخن, 3U/6U, يتطلب cpci_hotplug |

---

## دليل الاستخدام

تم نقل مرجع API وأمثلة الاستخدام (التثبيت، قراءة/كتابة Modbus، تجمع الاتصالات، MQTT، إعدادات Vendor المسبقة، تحويل البوابات، البروتوكولات المخصصة) إلى مستند مستقل:

> **[API.md](API.md) — مرجع API الكامل ودليل الاستخدام**
>
> النسخة الإنجليزية: [API.en.md](../../API.en.md)

---

## هيكل المشروع

```
industrial-protocols-go/
├── go.work                       # workspace لتجميع جميع الوحدات
├── kernel/                       # الوحدات الأساسية
│   ├── connection/               # ConnectionManager + تجمع الاتصالات
│   ├── config/                   # ConfigRepository (YAML)
│   ├── gateway/                  # GatewayEngine تحويل البروتوكولات
│   ├── bridge/                   # Bridge جسر العمليات الخارجية
│   ├── vendor/                   # Vendor الإعدادات المسبقة للمصنّعين
│   ├── event/                    # ناقل الأحداث
│   ├── metrics/                  # واجهة المقاييس
│   ├── security/                 # أمان TLS
│   ├── pipeline/                 # سلسلة الوسيطات
│   └── transport/                # نقل TCP/UDP/Pipe
│
├── protocols/                    # 40 وحدة بروتوكول
│   ├── ethernet/                 # 5 إيثرنت صناعي (مكتملة جميعها)
│   ├── fieldbus/                 # 11 ناقل ميدان (4 برمجية خالصة + 7 SDK أجهزة)
│   ├── iot/                      # 2 IoT/رسائل (مكتملة جميعها)
│   ├── automotive/               # 5 ناقل سيارات (2 برمجية خالصة + 3 SDK أجهزة)
│   ├── building/                 # 2 مباني/إضاءة (1 برمجية خالصة + 1 SDK أجهزة)
│   ├── bridge/                   # 12 جسر أجهزة (CmdBridge/SerialBridge)
│   └── system/                   # 3 ناقل نظام (برامج تشغيل sysfs/procfs)
│
├── examples/modbus_basic/        # مثال Modbus TCP
├── _tools/                       # Makefile + نصوص مساعدة
└── docs/superpowers/             # مستندات التصميم
```

## الاختبار

```bash
make test         # الاختبار الشامل
make test-unit    # اختبارات الوحدات فقط
make vet && make fmt
```

---

## ادعمنا

إذا كان هذا المشروع قد أفادك، فمرحبًا بك في دعمنا بالتبرع لمواصلة الصيانة.

### Alipay / WeChat

| Alipay | WeChat |
|--------|------|
| ![Alipay](../../alipay.png) | ![WeChat](../../weixinpay.png) |

### التحويلات العالمية (تحويل مصرفي دولي)

للمستخدمين خارج الصين القارية الذين يرغبون في التحويل المصرفي:

**معلومات المستلم:**

| البند | المحتوى |
|------|------|
| اسم المستلم | WANG KEXUN |
| رقم حساب المستلم | 881015918251 |

**البنك المستلم:**

| البند | المحتوى |
|------|------|
| اسم البنك | ZA Bank Limited |
| رمز SWIFT | AABLHKHHXXX |
| رقم البنك | 387 |
| عنوان البنك | Core F, Cyberport 3, 100 Cyberport Road, Hong Kong |

**البنك الوسيط للتحويلات العابرة للحدود (بنك وسيط، عند الحاجة):**

> يرجى الانتباه: هذه معلومات البنك الوسيط للتحويلات العابرة للحدود، وليست معلومات البنك المستلم. يُرجى الاستفسار من البنك المُرسِل عما إذا كان يتطلب معلومات البنك الوسيط.

- **التحويل بدولار هونغ كونغ واليوان الرينمينبي والدولار الأمريكي** — البنك الوسيط هو Citibank:
  - اسم البنك: Citibank N.A. Hong Kong
  - رمز SWIFT: CITIHKHXXXX
  - رقم البنك: 006
  - اسم الفرع: Hong Kong Branch
  - رقم الفرع: 391
  - عنوان البنك: Citibank Tower, Citibank Plaza, 3 Garden Road, Central, Hong Kong
- **التحويل بعملات أخرى** — البنك الوسيط هو BNY Mellon:
  - اسم البنك: THE BANK OF NEW YORK MELLON
  - رمز SWIFT: IRVTUS3NXXX
  - عنوان البنك: THE BANK OF NEW YORK MELLON, 240 GREENWICH STREET, NEW YORK, United States

---

## الترخيص

MIT — حقوق النشر (c) 2026 erik <erik@erik.xyz> — https://erik.xyz
