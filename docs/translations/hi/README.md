[English](../../README.en.md) | [中文](../../README.md) | [한국어](../ko/README.md) | [Русский](../ru/README.md) | [Deutsch](../de/README.md) | [Français](../fr/README.md) | [Español](../es/README.md) | [Português](../pt/README.md) | हिन्दी | [العربية](../ar/README.md) | [বাংলা](../bn/README.md) | [Bahasa Indonesia](../id/README.md) | [日本語](../ja/README.md)

# Industrial Protocols Go

Go भाषा में औद्योगिक नेटवर्क संचार प्रोटोकॉल संग्रह — स्तरित + मिडलवेयर आर्किटेक्चर, 40 औद्योगिक प्रोटोकॉल को कवर करता है, 14 शुद्ध-सॉफ्टवेयर कार्यान्वयन + 26 हार्डवेयर SDK (सभी ड्राइवर के साथ)।

> PHP कार्यान्वयन देखें: [github.com/erikwang2013/industrial-protocols](https://github.com/erikwang2013/industrial-protocols)

---

## परियोजना डिज़ाइन

### आर्किटेक्चर परतें

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
│  प्रति प्रोटोकॉल अलग मॉड्यूल      (TCP/UDP/Serial)  │
│                                                │
│  14 Pure-Soft   +   26 Hardware SDKs (with drivers)│
└──────────────────────────────────────────────────┘
```

### डिज़ाइन दृष्टिकोण

**माइक्रो-कर्नेल + प्रोटोकॉल SDK।** कर्नेल केवल इंटरफ़ेस और क्रॉस-कटिंग चिंताओं (कनेक्शन पूल, पुनः प्रयास, सर्किट ब्रेकर, इवेंट, मेट्रिक्स) को परिभाषित करता है, इसमें कोई विशिष्ट प्रोटोकॉल कार्यान्वयन शामिल नहीं है। प्रत्येक प्रोटोकॉल एक स्वतंत्र Go मॉड्यूल है, जिसे आवश्यकता अनुसार जोड़ा जाता है।

**परतों में वियोजन।** चार परतों का अमूर्तन:

| परत | ज़िम्मेदारी | मुख्य प्रकार |
|----|------|---------|
| **Transport** | निचला संचार चैनल | `Transport` इंटरफ़ेस — `TCPTransport`, `UDPTransport`, `PipeTransport` |
| **Codec** | प्रोटोकॉल एनकोड/डिकोड | `Codec` इंटरफ़ेस — `Encode(*Request) ([]byte, error)` / `Decode([]byte) (*Response, error)` |
| **Pipeline** | क्रॉस-कटिंग मिडलवेयर | `Middleware` — `Timeout`, `Retry`, `CircuitBreaker`, `Logger` |
| **Kernel** | डिवाइस जीवनचक्र | `ConnectionManager`, `ConfigRepository` |

**प्रोटोकॉल लेखकों को केवल `Protocol` + `Codec` दो इंटरफ़ेस लागू करने होते हैं।** ट्रांसपोर्ट परत, कनेक्शन पूल, पुनः प्रयास, टाइमआउट, सर्किट ब्रेकर — सब कुछ kernel से पुनः उपयोग होता है।

### कर्नेल मॉड्यूल

| मॉड्यूल | पथ | विवरण |
|------|------|------|
| ConnectionManager | `kernel/connection/` | डिवाइस पंजीकरण, कनेक्शन पूल, स्वास्थ्य जाँच; Lazy/Eager/Pooled तीन रणनीतियों का समर्थन |
| ConfigRepository | `kernel/config/` | YAML/JSON डिवाइस कॉन्फ़िगरेशन लोडिंग |
| GatewayEngine | `kernel/gateway/` | क्रॉस-प्रोटोकॉल रूपांतरण नियम इंजन (Modbus→MQTT आदि) |
| Bridge | `kernel/bridge/` | बाहरी प्रक्रिया ब्रिजिंग (stdin/stdout संचार) |
| Vendor | `kernel/vendor/` | विक्रेता पैरामीटर प्रीसेट (Siemens S7, Rockwell AB आदि) |
| Event | `kernel/event/` | channel-आधारित इवेंट बस |
| Metrics | `kernel/metrics/` | मेट्रिक्स संग्रह इंटरफ़ेस (Prometheus + noop) |
| Security | `kernel/security/` | TLS ट्रांसपोर्ट परत सुरक्षा रैपर |

---

## प्रोटोकॉल समर्थन

### औद्योगिक ईथरनेट (5/5 पूर्ण)

| प्रोटोकॉल | ट्रांसपोर्ट | पोर्ट | कार्यक्षमता |
|------|------|------|------|
| **Modbus** | TCP, RTU | 502 | FC 01/02/03/04/05/06/16, CRC16, MBAP |
| **BACnet** | UDP | 47808 | Who-Is/I-Am डिवाइस खोज, ReadProperty |
| **EtherNet/IP** | TCP | 44818 | ENIP RegisterSession + CIP Read Tag |
| **OPC UA** | TCP | 4840 | Binary HEL + OpenSecureChannel |
| **PROFINET** | UDP | 34964 | DCP Identify/Set + Record Data Read/Write |

### फील्डबस (11/11 — 4 शुद्ध-सॉफ्टवेयर + 7 हार्डवेयर SDK)

| प्रोटोकॉल | ट्रांसपोर्ट | पोर्ट | कार्यक्षमता |
|------|------|------|------|
| **HART** | Serial FSK | — | छोटा फ्रेम/लंबा फ्रेम, Command 0/3, XOR चेकसम |
| **CC-Link** | RS-485 | — | मास्टर-स्लेव पोलिंग, CRC-16/XMODEM |
| **DNP3** | TCP, Serial | 20000 | ट्रांसपोर्ट परत विभाजन/पुनर्संयोजन, Class 0 पोलिंग, CRC-16/DNP |
| **IEC 61850** | TCP | 102 | MMS Initiate/Conclude, Read/Write, BER-TLV |
| **PROFIBUS** | Serial | — | शुद्ध-सॉफ्टवेयर कार्यान्वयन — CP 5611 हार्डवेयर की आवश्यकता |
| **CANopen** | CAN, Gateway | — | SDO पढ़ना/लिखना, NMT प्रारंभ/रोक/रीसेट, Heartbeat |
| **DeviceNet** | CAN, Gateway | — | Explicit Messaging, Poll, I/O कनेक्शन |
| **Foundation Fieldbus** | Serial | — | शुद्ध-सॉफ्टवेयर कार्यान्वयन — FF H1 इंटरफ़ेस कार्ड की आवश्यकता |
| **AS-Interface** | Serial | — | शुद्ध-सॉफ्टवेयर कार्यान्वयन — ASi गेटवे की आवश्यकता |
| **IO-Link** | Serial | — | शुद्ध-सॉफ्टवेयर कार्यान्वयन — IO-Link Master की आवश्यकता |
| **CC-Link IE** | Ethernet | — | शुद्ध-सॉफ्टवेयर कार्यान्वयन — CC-Link IE गेटवे की आवश्यकता |

### IoT / मैसेजिंग (2/2 पूर्ण)

| प्रोटोकॉल | ट्रांसपोर्ट | पोर्ट | कार्यक्षमता |
|------|------|------|------|
| **MQTT** | TCP, WS | 1883 | 3.1.1 CONNECT/PUBLISH/SUBSCRIBE/PING |
| **HART-IP** | TCP, UDP | 5094 | HART over TCP, HART कोडेक का पुनः उपयोग |

### ऑटोमोटिव बस (5/5 — 2 शुद्ध-सॉफ्टवेयर + 3 हार्डवेयर SDK)

| प्रोटोकॉल | ट्रांसपोर्ट | बॉड दर | कार्यक्षमता |
|------|------|--------|------|
| **LIN** | UART | — | मास्टर/स्लेव फ्रेम, PID चेकसम, Classic/Enhanced Checksum |
| **K-Line** | Serial | 10400 | ISO 9141/14230, 5-baud Fast Init, OBD-II SID 01/03/09 |
| **FlexRay** | CAN, Serial | — | Slot/Frame एनकोड/डिकोड, CRC-16/XMODEM, SocketCAN + SerialBridge |
| **SAE J1850** | CAN | — | PWM/VPW एनकोड/डिकोड, Mode 01/03/0A, CRC-8 |
| **MOST** | Serial | — | Hex-frame एनकोड/डिकोड, Read/Write/Status, SerialBridge |

### बिल्डिंग / लाइटिंग (2/2 — 1 शुद्ध-सॉफ्टवेयर + 1 हार्डवेयर SDK)

| प्रोटोकॉल | ट्रांसपोर्ट | कार्यक्षमता |
|------|------|------|
| **DALI** | Serial | 16-bit फॉरवर्ड फ्रेम, 8-bit बैकवर्ड फ्रेम, मानक कमांड (Off/Max/Dim) |
| **LonWorks** | Serial | शुद्ध-सॉफ्टवेयर कोडेक — Neuron चिप या गेटवे की आवश्यकता |

### हार्डवेयर ब्रिजिंग (12/12 — सभी में CmdBridge/SerialBridge ड्राइवर हैं)

| प्रोटोकॉल | SDK प्रकार | कार्यक्षमता |
|------|----------|------|
| **EtherCAT** | CmdBridge + hex-codec | Upload/Download/Slaves उप-कमांड, CoE हेक्साडेसिमल कोडेक |
| **POWERLINK** | CmdBridge + hex-codec | Read/Write/Status उप-कमांड, SoC/Preq/Pres फ्रेम |
| **SERCOS III** | CmdBridge + hex-codec | Read/Write/Phase उप-कमांड, फेज़ संक्रमण (NRT→CP4) |
| **SERCOS I/II** | CmdBridge + hex-codec | Read/Write/Status, फाइबर ऑप्टिक रिंग टोपोलॉजी |
| **ControlNet** | CmdBridge + hex-codec | Read/Write/Status, CTDMA शेड्यूलिंग, निर्माता/उपभोक्ता |
| **Interbus** | CmdBridge | Read/Decode, रिंग टोपोलॉजी, IBS CMD फ्रेम |
| **WorldFIP** | CmdBridge | Read/Write/Decode, निर्माता/उपभोक्ता, बस अर्बिटर |
| **Lightbus** | CmdBridge | Read/Decode, फाइबर ऑप्टिक इंटरकनेक्ट, 32 नोड |
| **Modbus Plus** | CmdBridge + hex-codec | Read/Write/Decode, टोकन पासिंग, पीयर-टू-पीयर |
| **ISA100.11a** | CmdBridge + hex-codec | Read/Write/List, 6LoWPAN वायरलेस, मेश |
| **WirelessHART** | CmdBridge + hex-codec | Read/Write/Scan, TSMP MAC, स्व-संगठित |
| **SAE J1850** | CmdBridge | वाहन-ग्रेड डायग्नोस्टिक्स से CmdBridge रूपांतरण, J1850 इंटरफ़ेस की आवश्यकता |

### सिस्टम बस (3/3 — सभी में sysfs/procfs ड्राइवर हैं)

| प्रोटोकॉल | SDK प्रकार | कार्यक्षमता |
|------|----------|------|
| **PCI/PCIe** | sysfs — `/sys/bus/pci/devices/<BDF>/config` | कॉन्फ़िगरेशन स्पेस पढ़ना/लिखना, PipeTransport, CAP_SYS_ADMIN की आवश्यकता |
| **VME/VPX** | procfs — `/proc/vme/<slot>` | A16/A24/A32 एड्रेसिंग, PipeTransport, vme_tsi148 मॉड्यूल की आवश्यकता |
| **CompactPCI** | sysfs — `/sys/bus/pci/devices/<BDF>/config` | PICMG 2.0, हॉट-प्लग, 3U/6U, cpci_hotplug की आवश्यकता |

---

## उपयोग गाइड

API संदर्भ और उपयोग उदाहरण (इंस्टॉलेशन, Modbus पढ़ना/लिखना, कनेक्शन पूल, MQTT, Vendor प्रीसेट, गेटवे रूपांतरण, कस्टम प्रोटोकॉल) एक अलग दस्तावेज़ में स्थानांतरित कर दिए गए हैं:

> **[API.md](API.md) — पूर्ण API संदर्भ और उपयोग गाइड**
>
> अंग्रेज़ी संस्करण: [API.en.md](../API.en.md)

---

## परियोजना संरचना

```
industrial-protocols-go/
├── go.work                       # workspace सभी मॉड्यूल को एकत्रित करता है
├── kernel/                       # कर्नेल मॉड्यूल
│   ├── connection/               # ConnectionManager + कनेक्शन पूल
│   ├── config/                   # ConfigRepository (YAML)
│   ├── gateway/                  # GatewayEngine प्रोटोकॉल रूपांतरण
│   ├── bridge/                   # Bridge बाहरी प्रक्रिया ब्रिजिंग
│   ├── vendor/                   # Vendor विक्रेता प्रीसेट
│   ├── event/                    # इवेंट बस
│   ├── metrics/                  # मेट्रिक्स इंटरफ़ेस
│   ├── security/                 # TLS सुरक्षा
│   ├── pipeline/                 # मिडलवेयर श्रृंखला
│   └── transport/                # TCP/UDP/Pipe ट्रांसपोर्ट
│
├── protocols/                    # 40 प्रोटोकॉल मॉड्यूल
│   ├── ethernet/                 # 5 औद्योगिक ईथरनेट (सभी पूर्ण)
│   ├── fieldbus/                 # 11 फील्डबस (4 शुद्ध-सॉफ्टवेयर + 7 हार्डवेयर SDK)
│   ├── iot/                      # 2 IoT/मैसेजिंग (सभी पूर्ण)
│   ├── automotive/               # 5 ऑटोमोटिव बस (2 शुद्ध-सॉफ्टवेयर + 3 हार्डवेयर SDK)
│   ├── building/                 # 2 बिल्डिंग/लाइटिंग (1 शुद्ध-सॉफ्टवेयर + 1 हार्डवेयर SDK)
│   ├── bridge/                   # 12 हार्डवेयर ब्रिजिंग (CmdBridge/SerialBridge)
│   └── system/                   # 3 सिस्टम बस (sysfs/procfs ड्राइवर)
│
├── examples/modbus_basic/        # Modbus TCP उदाहरण
├── _tools/                       # Makefile + सहायक स्क्रिप्ट
└── docs/superpowers/             # डिज़ाइन दस्तावेज़
```

## परीक्षण

```bash
make test         # पूर्ण परीक्षण
make test-unit    # केवल यूनिट परीक्षण
make vet && make fmt
```

---

## हमें समर्थन करें

यदि यह प्रोजेक्ट आपके काम आया है, तो कृपया दान करके हमारे निरंतर रखरखाव का समर्थन करें।

### अलीपे / वीचैट

| अलीपे | वीचैट |
|--------|------|
| ![अलीपे](../../alipay.png) | ![वीचैट](../../weixinpay.png) |

### वैश्विक स्थानांतरण (अंतर्राष्ट्रीय बैंक हस्तांतरण)

यह मुख्यभूमि चीन के बाहर के उपयोगकर्ताओं के लिए बैंक स्थानांतरण है:

**प्राप्तकर्ता की जानकारी:**

| मद | विवरण |
|------|------|
| प्राप्तकर्ता का नाम | WANG KEXUN |
| प्राप्तकर्ता खाता संख्या | 881015918251 |

**प्राप्तकर्ता बैंक:**

| मद | विवरण |
|------|------|
| बैंक का नाम | ZA Bank Limited |
| SWIFT Code | AABLHKHHXXX |
| बैंक कोड | 387 |
| बैंक का पता | Core F, Cyberport 3, 100 Cyberport Road, Hong Kong |

**क्रॉस-बॉर्डर रेमिटेंस एजेंट बैंक (मध्यवर्ती बैंक, यदि आवश्यक हो):**

> कृपया ध्यान दें, यह क्रॉस-बॉर्डर रेमिटेंस एजेंट बैंक (मध्यवर्ती बैंक) की जानकारी है, प्राप्तकर्ता बैंक की नहीं। कृपया अपने रेमिटेंस बैंक से पूछें कि क्या एजेंट बैंक की जानकारी देना आवश्यक है।

- **हांगकांग डॉलर, RMB और USD के लिए** — एजेंट बैंक Citibank है:
  - बैंक का नाम: Citibank N.A. Hong Kong
  - SWIFT Code: CITIHKHXXXX
  - बैंक कोड: 006
  - शाखा का नाम: Hong Kong Branch
  - शाखा कोड: 391
  - बैंक का पता: Citibank Tower, Citibank Plaza, 3 Garden Road, Central, Hong Kong
- **अन्य मुद्राओं के लिए** — एजेंट बैंक BNY Mellon है:
  - बैंक का नाम: THE BANK OF NEW YORK MELLON
  - SWIFT Code: IRVTUS3NXXX
  - बैंक का पता: THE BANK OF NEW YORK MELLON, 240 GREENWICH STREET, NEW YORK, United States

---

## License

MIT — Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz
