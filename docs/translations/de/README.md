[English](../../README.en.md) | [中文](../../README.md) | [한국어](../ko/README.md) | [Русский](../ru/README.md) | Deutsch | [Français](../fr/README.md) | [Español](../es/README.md) | [Português](../pt/README.md) | [हिन्दी](../hi/README.md) | [العربية](../ar/README.md) | [বাংলা](../bn/README.md) | [Bahasa Indonesia](../id/README.md) | [日本語](../ja/README.md)

# Industrial Protocols Go

Go-Protokollsatz für industrielle Netzwerkkommunikation — Schichten- und Middleware-Architektur, abdeckend 40 industrielle Protokolle, 14 reine Software-Implementierungen + 26 Hardware-SDKs (alle mit Treibern).

> Referenz-PHP-Implementierung: [github.com/erikwang2013/industrial-protocols](https://github.com/erikwang2013/industrial-protocols)

---

## Projektdesign

### Architektur-Schichten

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
│  Eigenständiges Modul pro Protokoll (TCP/UDP/Serial)│
│                                                │
│  14 Pure-Soft   +   26 Hardware SDKs (with drivers)│
└──────────────────────────────────────────────────┘
```

### Designphilosophie

**Mikrokernel + Protokoll-SDKs.** Der Kernel definiert nur Schnittstellen und Querschnittsbelange (Verbindungspool, Retry, Circuit Breaker, Events, Metriken) und enthält keine konkrete Protokollimplementierung. Jedes Protokoll ist ein eigenständiges Go-Modul, das bei Bedarf eingebunden wird.

**Schichtenentkopplung.** Vier Abstraktionsebenen:

| Ebene | Verantwortung | Kerntypen |
|----|------|---------|
| **Transport** | Kommunikationskanäle der unteren Ebene | `Transport`-Schnittstelle — `TCPTransport`, `UDPTransport`, `PipeTransport` |
| **Codec** | Protokoll-Codierung/-Decodierung | `Codec`-Schnittstelle — `Encode(*Request) ([]byte, error)` / `Decode([]byte) (*Response, error)` |
| **Pipeline** | Querschnitts-Middleware | `Middleware` — `Timeout`, `Retry`, `CircuitBreaker`, `Logger` |
| **Kernel** | Gerätelebenszyklus | `ConnectionManager`, `ConfigRepository` |

**Protokollautoren implementieren nur die beiden Schnittstellen `Protocol` + `Codec`.** Transportschicht, Verbindungspool, Retry, Timeout und Circuit Breaker werden aus dem Kernel wiederverwendet.

### Kernel-Module

| Modul | Pfad | Beschreibung |
|------|------|------|
| ConnectionManager | `kernel/connection/` | Geräteregistrierung, Verbindungspool, Health Checks; unterstützt die drei Strategien Lazy/Eager/Pooled |
| ConfigRepository | `kernel/config/` | Laden von YAML/JSON-Gerätekonfigurationen |
| GatewayEngine | `kernel/gateway/` | Regel-Engine für protokollübergreifende Konvertierung (Modbus→MQTT usw.) |
| Bridge | `kernel/bridge/` | Brücke zu externen Prozessen (stdin/stdout-Kommunikation) |
| Vendor | `kernel/vendor/` | Herstellerspezifische Parameter-Presets (Siemens S7, Rockwell AB usw.) |
| Event | `kernel/event/` | Channel-basierter Event-Bus |
| Metrics | `kernel/metrics/` | Metrikerfassungs-Schnittstelle (Prometheus + noop) |
| Security | `kernel/security/` | TLS-Sicherheitskapselung der Transportschicht |

---

## Protokollunterstützung

### Industrielles Ethernet (5/5 abgeschlossen)

| Protokoll | Transport | Port | Funktionen |
|------|------|------|------|
| **Modbus** | TCP, RTU | 502 | FC 01/02/03/04/05/06/16, CRC16, MBAP |
| **BACnet** | UDP | 47808 | Who-Is/I-Am-Geräteerkennung, ReadProperty |
| **EtherNet/IP** | TCP | 44818 | ENIP RegisterSession + CIP Read Tag |
| **OPC UA** | TCP | 4840 | Binary HEL + OpenSecureChannel |
| **PROFINET** | UDP | 34964 | DCP Identify/Set + Record Data Read/Write |

### Feldbusse (11/11 — 4 rein softwarebasiert + 7 Hardware-SDKs)

| Protokoll | Transport | Port | Funktionen |
|------|------|------|------|
| **HART** | Serial FSK | — | Kurz-/Langrahmen, Command 0/3, XOR-Prüfung |
| **CC-Link** | RS-485 | — | Master-Slave-Polling, CRC-16/XMODEM |
| **DNP3** | TCP, Serial | 20000 | Transportlayer-Segmentierung/-Rekombination, Class-0-Polling, CRC-16/DNP |
| **IEC 61850** | TCP | 102 | MMS Initiate/Conclude, Read/Write, BER-TLV |
| **PROFIBUS** | Serial | — | Reine Software-Implementierung — erfordert CP-5611-Hardware |
| **CANopen** | CAN, Gateway | — | SDO Lesen/Schreiben, NMT Start/Stopp/Reset, Heartbeat |
| **DeviceNet** | CAN, Gateway | — | Explicit Messaging, Poll, I/O-Verbindungen |
| **Foundation Fieldbus** | Serial | — | Reine Software-Implementierung — erfordert FF-H1-Interfacekarte |
| **AS-Interface** | Serial | — | Reine Software-Implementierung — erfordert ASi-Gateway |
| **IO-Link** | Serial | — | Reine Software-Implementierung — erfordert IO-Link-Master |
| **CC-Link IE** | Ethernet | — | Reine Software-Implementierung — erfordert CC-Link-IE-Gateway |

### IoT / Messaging (2/2 abgeschlossen)

| Protokoll | Transport | Port | Funktionen |
|------|------|------|------|
| **MQTT** | TCP, WS | 1883 | 3.1.1 CONNECT/PUBLISH/SUBSCRIBE/PING |
| **HART-IP** | TCP, UDP | 5094 | HART über TCP, nutzt die HART-Codierung/-Decodierung |

### Automobilbusse (5/5 — 2 rein softwarebasiert + 3 Hardware-SDKs)

| Protokoll | Transport | Baudrate | Funktionen |
|------|------|--------|------|
| **LIN** | UART | — | Master-/Slave-Frames, PID-Prüfung, Classic/Enhanced Checksum |
| **K-Line** | Serial | 10400 | ISO 9141/14230, 5-Baud-Fast-Init, OBD-II SID 01/03/09 |
| **FlexRay** | CAN, Serial | — | Slot-/Frame-Codierung, CRC-16/XMODEM, SocketCAN + SerialBridge |
| **SAE J1850** | CAN | — | PWM/VPW-Codierung, Mode 01/03/0A, CRC-8 |
| **MOST** | Serial | — | Hex-Frame-Codierung, Read/Write/Status, SerialBridge |

### Gebäude / Beleuchtung (2/2 — 1 rein softwarebasiert + 1 Hardware-SDK)

| Protokoll | Transport | Funktionen |
|------|------|------|
| **DALI** | Serial | 16-Bit-Vorwärtsframes, 8-Bit-Rückwärtsframes, Standardbefehle (Off/Max/Dim) |
| **LonWorks** | Serial | Reine Software-Codierung — erfordert Neuron-Chip oder Gateway |

### Hardware-Brücken (12/12 — alle mit CmdBridge/SerialBridge-Treibern)

| Protokoll | SDK-Typ | Funktionen |
|------|----------|------|
| **EtherCAT** | CmdBridge + Hex-Codec | Upload/Download/Slaves-Unterbefehle, CoE-Hex-Codierung |
| **POWERLINK** | CmdBridge + Hex-Codec | Read/Write/Status-Unterbefehle, SoC/Preq/Pres-Frames |
| **SERCOS III** | CmdBridge + Hex-Codec | Read/Write/Phase-Unterbefehle, Phasenübergänge (NRT→CP4) |
| **SERCOS I/II** | CmdBridge + Hex-Codec | Read/Write/Status, Glasfaser-Ringtopologie |
| **ControlNet** | CmdBridge + Hex-Codec | Read/Write/Status, CTDMA-Scheduling, Produzent/Konsument |
| **Interbus** | CmdBridge | Read/Decode, Ringtopologie, IBS-CMD-Frames |
| **WorldFIP** | CmdBridge | Read/Write/Decode, Produzent/Konsument, Bus-Arbiter |
| **Lightbus** | CmdBridge | Read/Decode, Glasfaser-Verbindung, 32 Knoten |
| **Modbus Plus** | CmdBridge + Hex-Codec | Read/Write/Decode, Token-Passing, Peer-to-Peer |
| **ISA100.11a** | CmdBridge + Hex-Codec | Read/Write/List, 6LoWPAN-Funk, Mesh |
| **WirelessHART** | CmdBridge + Hex-Codec | Read/Write/Scan, TSMP-MAC, selbstorganisierend |
| **SAE J1850** | CmdBridge | Fahrzeugdiagnose auf CmdBridge, erfordert J1850-Interface |

### Systembusse (3/3 — alle mit sysfs/procfs-Treibern)

| Protokoll | SDK-Typ | Funktionen |
|------|----------|------|
| **PCI/PCIe** | sysfs — `/sys/bus/pci/devices/<BDF>/config` | Konfigurationsraum Lesen/Schreiben, PipeTransport, erfordert CAP_SYS_ADMIN |
| **VME/VPX** | procfs — `/proc/vme/<slot>` | A16/A24/A32-Adressierung, PipeTransport, erfordert vme_tsi148-Modul |
| **CompactPCI** | sysfs — `/sys/bus/pci/devices/<BDF>/config` | PICMG 2.0, Hot-Swap, 3U/6U, erfordert cpci_hotplug |

---

## Verwendung

Die API-Referenz und Verwendungsbeispiele (Installation, Modbus Lesen/Schreiben, Verbindungspool, MQTT, Vendor-Presets, Gateway-Konvertierung, benutzerdefinierte Protokolle) wurden in ein separates Dokument verschoben:

> **[API.md](API.md) — vollständige API-Referenz und Verwendungsanleitung**
>
> Englische Version: [API.en.md](../../docs/API.en.md)

---

## Projektstruktur

```
industrial-protocols-go/
├── go.work                       # Workspace aggregiert alle Module
├── kernel/                       # Kernmodul
│   ├── connection/               # ConnectionManager + Verbindungspool
│   ├── config/                   # ConfigRepository (YAML)
│   ├── gateway/                  # GatewayEngine Protokollkonvertierung
│   ├── bridge/                   # Bridge externe Prozessbrücke
│   ├── vendor/                   # Vendor Hersteller-Presets
│   ├── event/                    # Event-Bus
│   ├── metrics/                  # Metrik-Schnittstellen
│   ├── security/                 # TLS-Sicherheit
│   ├── pipeline/                 # Middleware-Kette
│   └── transport/                # TCP/UDP/Pipe-Transport
│
├── protocols/                    # 40 Protokollmodule
│   ├── ethernet/                 # 5 Industrielles Ethernet (alle abgeschlossen)
│   ├── fieldbus/                 # 11 Feldbusse (4 rein softwarebasiert + 7 Hardware-SDKs)
│   ├── iot/                      # 2 IoT/Messaging (alle abgeschlossen)
│   ├── automotive/               # 5 Automobilbusse (2 rein softwarebasiert + 3 Hardware-SDKs)
│   ├── building/                 # 2 Gebäude/Beleuchtung (1 rein softwarebasiert + 1 Hardware-SDK)
│   ├── bridge/                   # 12 Hardware-Brücken (CmdBridge/SerialBridge)
│   └── system/                   # 3 Systembusse (sysfs/procfs-Treiber)
│
├── examples/modbus_basic/        # Modbus-TCP-Beispiel
├── _tools/                       # Makefile + Hilfsskripte
└── docs/superpowers/             # Designdokumente
```

## Tests

```bash
make test         # vollständiger Testlauf
make test-unit    # nur Unit-Tests
make vet && make fmt
```

---

## Unterstützung

Wenn dieses Projekt dir geholfen hat, freuen wir uns über eine Spende zur Unterstützung der weiteren Pflege.

### Alipay / WeChat

| Alipay | WeChat |
|--------|------|
| ![Alipay](../../alipay.png) | ![WeChat](../../weixinpay.png) |

### Internationale Überweisung (SWIFT)

Banküberweisung für Nutzer außerhalb Chinas:

**Empfänger:**

| Posten | Inhalt |
|------|------|
| Empfängername | WANG KEXUN |
| Kontonummer | 881015918251 |

**Empfängerbank:**

| Posten | Inhalt |
|------|------|
| Bankname | ZA Bank Limited |
| SWIFT-Code | AABLHKHHXXX |
| Bankleitzahl | 387 |
| Bankadresse | Core F, Cyberport 3, 100 Cyberport Road, Hong Kong |

**Korrespondenzbank für grenzüberschreitende Überweisungen (Zwischenbank, falls erforderlich):**

> Bitte beachten: Dies sind die Angaben der Korrespondenzbank (Zwischenbank) für grenzüberschreitende Überweisungen, nicht die der Empfängerbank. Erkundige dich bei deiner überweisenden Bank, ob die Angaben der Korrespondenzbank benötigt werden.

- **Überweisungen in Hongkong-Dollar, Renminbi und US-Dollar** — Korrespondenzbank ist Citibank:
  - Bankname: Citibank N.A. Hong Kong
  - SWIFT-Code: CITIHKHXXXX
  - Bankleitzahl: 006
  - Filiale: Hong Kong Branch
  - Filialnummer: 391
  - Bankadresse: Citibank Tower, Citibank Plaza, 3 Garden Road, Central, Hong Kong
- **Überweisungen in anderen Währungen** — Korrespondenzbank ist BNY Mellon:
  - Bankname: THE BANK OF NEW YORK MELLON
  - SWIFT-Code: IRVTUS3NXXX
  - Bankadresse: THE BANK OF NEW YORK MELLON, 240 GREENWICH STREET, NEW YORK, United States

---

## License

MIT — Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz
