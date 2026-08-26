[English](../../README.en.md) | [中文](../../README.md) | [한국어](../ko/README.md) | [Русский](../ru/README.md) | [Deutsch](../de/README.md) | [Français](../fr/README.md) | Español | [Português](../pt/README.md) | [हिन्दी](../hi/README.md) | [العربية](../ar/README.md) | [বাংলা](../bn/README.md) | [Bahasa Indonesia](../id/README.md) | [日本語](../ja/README.md)

# Industrial Protocols Go

Conjunto de protocolos de comunicación industrial en Go —— arquitectura por capas + middleware, cubre 40 protocolos industriales, 14 implementaciones 100% software + 26 SDKs de hardware (todos con drivers).

> Referencia a la implementación en PHP: [github.com/erikwang2013/industrial-protocols](https://github.com/erikwang2013/industrial-protocols)

---

## Diseño del proyecto

### Capas de la arquitectura

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
│  Módulo independiente por protocolo  (TCP/UDP/Serial)  │
│                                                │
│  14 Pure-Soft   +   26 Hardware SDKs (with drivers)│
└──────────────────────────────────────────────────┘
```

### Enfoque de diseño

**Microkernel + SDK de protocolo.** El kernel solo define interfaces y preocupaciones transversales (pool de conexiones, reintentos, interruptor de circuito, eventos, métricas); no contiene ninguna implementación concreta de protocolo. Cada protocolo es un módulo Go independiente que se importa según necesidad.

**Desacoplamiento por capas.** Cuatro niveles de abstracción:

| Capa | Responsabilidad | Tipos centrales |
|----|------|---------|
| **Transport** | Canal de comunicación de bajo nivel | Interfaz `Transport` — `TCPTransport`, `UDPTransport`, `PipeTransport` |
| **Codec** | Codificación/decodificación del protocolo | Interfaz `Codec` — `Encode(*Request) ([]byte, error)` / `Decode([]byte) (*Response, error)` |
| **Pipeline** | Middleware transversal | `Middleware` — `Timeout`, `Retry`, `CircuitBreaker`, `Logger` |
| **Kernel** | Ciclo de vida de los dispositivos | `ConnectionManager`, `ConfigRepository` |

**El autor de un protocolo solo debe implementar las interfaces `Protocol` + `Codec`.** La capa de transporte, el pool de conexiones, los reintentos, los timeouts y el interruptor de circuito se reutilizan del kernel.

### Módulos del kernel

| Módulo | Ruta | Descripción |
|------|------|------|
| ConnectionManager | `kernel/connection/` | Registro de dispositivos, pool de conexiones, comprobación de estado; admite las tres estrategias Lazy/Eager/Pooled |
| ConfigRepository | `kernel/config/` | Carga de configuración de dispositivos YAML/JSON |
| GatewayEngine | `kernel/gateway/` | Motor de reglas de conversión entre protocolos (Modbus→MQTT, etc.) |
| Bridge | `kernel/bridge/` | Puente con procesos externos (comunicación stdin/stdout) |
| Vendor | `kernel/vendor/` | Preajustes de parámetros de fabricantes (Siemens S7, Rockwell AB, etc.) |
| Event | `kernel/event/` | Bus de eventos basado en canales |
| Metrics | `kernel/metrics/` | Interfaz de recolección de métricas (Prometheus + noop) |
| Security | `kernel/security/` | Envoltorio de seguridad TLS en la capa de transporte |

---

## Protocolos soportados

### Ethernet industrial (5/5 completados)

| Protocolo | Transporte | Puerto | Funciones |
|------|------|------|------|
| **Modbus** | TCP, RTU | 502 | FC 01/02/03/04/05/06/16, CRC16, MBAP |
| **BACnet** | UDP | 47808 | Descubrimiento de dispositivos Who-Is/I-Am, ReadProperty |
| **EtherNet/IP** | TCP | 44818 | ENIP RegisterSession + CIP Read Tag |
| **OPC UA** | TCP | 4840 | Binary HEL + OpenSecureChannel |
| **PROFINET** | UDP | 34964 | DCP Identify/Set + Lectura/Escritura de Record Data |

### Fieldbus (11/11 — 4 software puro + 7 SDK de hardware)

| Protocolo | Transporte | Puerto | Funciones |
|------|------|------|------|
| **HART** | Serial FSK | — | Trama corta/larga, Command 0/3, checksum XOR |
| **CC-Link** | RS-485 | — | Polling maestro-esclavo, CRC-16/XMODEM |
| **DNP3** | TCP, Serial | 20000 | Segmentación/reensamblado de capa de transporte, polling Clase 0, CRC-16/DNP |
| **IEC 61850** | TCP | 102 | MMS Initiate/Conclude, Read/Write, BER-TLV |
| **PROFIBUS** | Serial | — | Implementación pura en software — requiere hardware CP 5611 |
| **CANopen** | CAN, Gateway | — | Lectura/escritura SDO, NMT arranque/parada/reset, Heartbeat |
| **DeviceNet** | CAN, Gateway | — | Explicit Messaging, Poll, conexiones I/O |
| **Foundation Fieldbus** | Serial | — | Implementación pura en software — requiere tarjeta de interfaz FF H1 |
| **AS-Interface** | Serial | — | Implementación pura en software — requiere gateway ASi |
| **IO-Link** | Serial | — | Implementación pura en software — requiere maestro IO-Link |
| **CC-Link IE** | Ethernet | — | Implementación pura en software — requiere gateway CC-Link IE |

### IoT / Mensajería (2/2 completados)

| Protocolo | Transporte | Puerto | Funciones |
|------|------|------|------|
| **MQTT** | TCP, WS | 1883 | 3.1.1 CONNECT/PUBLISH/SUBSCRIBE/PING |
| **HART-IP** | TCP, UDP | 5094 | HART sobre TCP, reutiliza la codificación HART |

### Bus de automoción (5/5 — 2 software puro + 3 SDK de hardware)

| Protocolo | Transporte | Baudios | Funciones |
|------|------|--------|------|
| **LIN** | UART | — | Tramas maestro-esclavo, checksum PID, Classic/Enhanced Checksum |
| **K-Line** | Serial | 10400 | ISO 9141/14230, 5-baud Fast Init, OBD-II SID 01/03/09 |
| **FlexRay** | CAN, Serial | — | Codificación/decodificación Slot/Frame, CRC-16/XMODEM, SocketCAN + SerialBridge |
| **SAE J1850** | CAN | — | Codificación/decodificación PWM/VPW, Mode 01/03/0A, CRC-8 |
| **MOST** | Serial | — | Codificación/decodificación de tramas hex, Read/Write/Status, SerialBridge |

### Edificios / Iluminación (2/2 — 1 software puro + 1 SDK de hardware)

| Protocolo | Transporte | Funciones |
|------|------|------|
| **DALI** | Serial | Trama directa de 16 bits, trama de retorno de 8 bits, comandos estándar (Off/Max/Dim) |
| **LonWorks** | Serial | Codificación/decodificación pura en software — requiere chip Neuron o gateway |

### Puentes de hardware (12/12 — todos con drivers CmdBridge/SerialBridge)

| Protocolo | Tipo de SDK | Funciones |
|------|----------|------|
| **EtherCAT** | CmdBridge + hex-codec | Subcomandos Upload/Download/Slaves, codificación hexadecimal CoE |
| **POWERLINK** | CmdBridge + hex-codec | Subcomandos Read/Write/Status, tramas SoC/Preq/Pres |
| **SERCOS III** | CmdBridge + hex-codec | Subcomandos Read/Write/Phase, transición de fases (NRT→CP4) |
| **SERCOS I/II** | CmdBridge + hex-codec | Read/Write/Status, topología de anillo de fibra óptica |
| **ControlNet** | CmdBridge + hex-codec | Read/Write/Status, planificación CTDMA, productor/consumidor |
| **Interbus** | CmdBridge | Read/Decode, topología de anillo, tramas IBS CMD |
| **WorldFIP** | CmdBridge | Read/Write/Decode, productor/consumidor, árbitro de bus |
| **Lightbus** | CmdBridge | Read/Decode, interconexión de fibra óptica, 32 nodos |
| **Modbus Plus** | CmdBridge + hex-codec | Read/Write/Decode, paso de testigo, par a par |
| **ISA100.11a** | CmdBridge + hex-codec | Read/Write/List, inalámbrico 6LoWPAN, malla |
| **WirelessHART** | CmdBridge + hex-codec | Read/Write/Scan, MAC TSMP, autoorganizado |
| **SAE J1850** | CmdBridge | Diagnóstico automotriz a CmdBridge, requiere interfaz J1850 |

### Buses de sistema (3/3 — todos con drivers sysfs/procfs)

| Protocolo | Tipo de SDK | Funciones |
|------|----------|------|
| **PCI/PCIe** | sysfs — `/sys/bus/pci/devices/<BDF>/config` | Lectura/escritura del espacio de configuración, PipeTransport, requiere CAP_SYS_ADMIN |
| **VME/VPX** | procfs — `/proc/vme/<slot>` | Direccionamiento A16/A24/A32, PipeTransport, requiere módulo vme_tsi148 |
| **CompactPCI** | sysfs — `/sys/bus/pci/devices/<BDF>/config` | PICMG 2.0, conexión en caliente, 3U/6U, requiere cpci_hotplug |

---

## Guía de uso

La referencia de API y los ejemplos de uso (instalación, lectura/escritura Modbus, pool de conexiones, MQTT, preajustes Vendor, conversión por gateway, protocolos personalizados) se han trasladado a un documento independiente:

> **[API.md](API.md) — Referencia completa de API y guía de uso**
>
> Versión en inglés: [docs/API.en.md](../../API.en.md)

---

## Estructura del proyecto

```
industrial-protocols-go/
├── go.work                       # workspace que agrega todos los módulos
├── kernel/                       # módulos centrales
│   ├── connection/               # ConnectionManager + pool de conexiones
│   ├── config/                   # ConfigRepository (YAML)
│   ├── gateway/                  # GatewayEngine conversión de protocolos
│   ├── bridge/                   # Bridge puente con procesos externos
│   ├── vendor/                   # Vendor preajustes de fabricantes
│   ├── event/                    # bus de eventos
│   ├── metrics/                  # interfaz de métricas
│   ├── security/                 # seguridad TLS
│   ├── pipeline/                 # cadena de middleware
│   └── transport/                # transporte TCP/UDP/Pipe
│
├── protocols/                    # 40 módulos de protocolo
│   ├── ethernet/                 # 5 Ethernet industrial (todos completados)
│   ├── fieldbus/                 # 11 fieldbus (4 software puro + 7 SDK de hardware)
│   ├── iot/                      # 2 IoT/mensajería (todos completados)
│   ├── automotive/               # 5 buses de automoción (2 software puro + 3 SDK de hardware)
│   ├── building/                 # 2 edificios/iluminación (1 software puro + 1 SDK de hardware)
│   ├── bridge/                   # 12 puentes de hardware (CmdBridge/SerialBridge)
│   └── system/                   # 3 buses de sistema (drivers sysfs/procfs)
│
├── examples/modbus_basic/        # ejemplo Modbus TCP
├── _tools/                       # Makefile + scripts auxiliares
└── docs/superpowers/             # documentos de diseño
```

## Pruebas

```bash
make test         # prueba completa
make test-unit    # solo pruebas unitarias
make vet && make fmt
```

---

## Apóyanos

Si este proyecto te ha ayudado, eres bienvenido a donar para apoyarnos en su mantenimiento.

### Alipay / WeChat

| Alipay | WeChat |
|--------|------|
| ![Alipay](../../alipay.png) | ![WeChat](../../weixinpay.png) |

### Transferencia internacional (transferencia bancaria internacional)

Para transferencias bancarias de usuarios fuera de la China continental:

**Información del beneficiario:**

| Elemento | Contenido |
|------|------|
| Nombre del beneficiario | WANG KEXUN |
| Número de cuenta del beneficiario | 881015918251 |

**Banco beneficiario:**

| Elemento | Contenido |
|------|------|
| Nombre del banco | ZA Bank Limited |
| SWIFT Code | AABLHKHHXXX |
| Número de banco | 387 |
| Dirección del banco | Core F, Cyberport 3, 100 Cyberport Road, Hong Kong |

**Banco corresponsal de remesas transfronterizas (banco intermediario, si es necesario):**

> Tenga en cuenta que esta es la información del banco corresponsal (banco intermediario) para remesas transfronterizas, no la del banco beneficiario. Consulte con su banco emisor si necesita proporcionar la información del banco intermediario.

- **Remesas en dólares de Hong Kong, RMB y dólares estadounidenses** — el banco intermediario es Citibank:
  - Nombre del banco: Citibank N.A. Hong Kong
  - SWIFT Code: CITIHKHXXXX
  - Número de banco: 006
  - Nombre de la sucursal: Hong Kong Branch
  - Número de sucursal: 391
  - Dirección del banco: Citibank Tower, Citibank Plaza, 3 Garden Road, Central, Hong Kong
- **Remesas en otras divisas** — el banco intermediario es BNY Mellon:
  - Nombre del banco: THE BANK OF NEW YORK MELLON
  - SWIFT Code: IRVTUS3NXXX
  - Dirección del banco: THE BANK OF NEW YORK MELLON, 240 GREENWICH STREET, NEW YORK, United States

---

## License

MIT — Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz
