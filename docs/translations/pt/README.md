[English](../../README.en.md) | [中文](../../README.md) | [한국어](../ko/README.md) | [Русский](../ru/README.md) | [Deutsch](../de/README.md) | [Français](../fr/README.md) | [Español](../es/README.md) | Português | [हिन्दी](../hi/README.md) | [العربية](../ar/README.md) | [বাংলা](../bn/README.md) | [Bahasa Indonesia](../id/README.md) | [日本語](../ja/README.md)

# Industrial Protocols Go

Conjunto de protocolos de comunicação industrial em Go — arquitetura em camadas + middleware, cobrindo 40 protocolos industriais, 14 implementações 100% software + 26 SDKs de hardware (todos com drivers).

> Referência da implementação PHP: [github.com/erikwang2013/industrial-protocols](https://github.com/erikwang2013/industrial-protocols)

---

## Design do Projeto

### Camadas da Arquitetura

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
│  Módulo indep. por protocolo     (TCP/UDP/Serial) │
│                                                │
│  14 Pure-Soft   +   26 Hardware SDKs (with drivers)│
└──────────────────────────────────────────────────┘
```

### Filosofia de Design

**Microkernel + SDK de protocolos.** O kernel define apenas as interfaces e as preocupações transversais (pool de conexões, retry, circuit breaker, eventos, métricas) e não contém nenhuma implementação concreta de protocolo. Cada protocolo é um módulo Go independente, importado conforme a necessidade.

**Desacoplamento em camadas.** Quatro camadas de abstração:

| Camada | Responsabilidade | Tipos principais |
|----|------|---------|
| **Transport** | Canal de comunicação de baixo nível | Interface `Transport` — `TCPTransport`, `UDPTransport`, `PipeTransport` |
| **Codec** | Codificação/decodificação do protocolo | Interface `Codec` — `Encode(*Request) ([]byte, error)` / `Decode([]byte) (*Response, error)` |
| **Pipeline** | Middlewares transversais | `Middleware` — `Timeout`, `Retry`, `CircuitBreaker`, `Logger` |
| **Kernel** | Ciclo de vida dos dispositivos | `ConnectionManager`, `ConfigRepository` |

**Os autores de protocolos precisam apenas implementar as duas interfaces `Protocol` + `Codec`.** Transporte, pool de conexões, retry, timeout e circuit breaker são todos reutilizados do kernel.

### Módulos do Kernel

| Módulo | Caminho | Descrição |
|------|------|------|
| ConnectionManager | `kernel/connection/` | Registro de dispositivos, pool de conexões, health check; suporta as três estratégias Lazy/Eager/Pooled |
| ConfigRepository | `kernel/config/` | Carregamento de configuração de dispositivos em YAML/JSON |
| GatewayEngine | `kernel/gateway/` | Motor de regras de conversão entre protocolos (Modbus→MQTT etc.) |
| Bridge | `kernel/bridge/` | Ponte para processos externos (comunicação via stdin/stdout) |
| Vendor | `kernel/vendor/` | Presets de parâmetros de fabricantes (Siemens S7, Rockwell AB etc.) |
| Event | `kernel/event/` | Barramento de eventos baseado em channels |
| Metrics | `kernel/metrics/` | Interface de coleta de métricas (Prometheus + noop) |
| Security | `kernel/security/` | Wrapper de segurança da camada de transporte TLS |

---

## Protocolos Suportados

### Ethernet Industrial (5/5 concluídos)

| Protocolo | Transporte | Porta | Função |
|------|------|------|------|
| **Modbus** | TCP, RTU | 502 | FC 01/02/03/04/05/06/16, CRC16, MBAP |
| **BACnet** | UDP | 47808 | Descoberta de dispositivos Who-Is/I-Am, ReadProperty |
| **EtherNet/IP** | TCP | 44818 | ENIP RegisterSession + CIP Read Tag |
| **OPC UA** | TCP | 4840 | HEL binário + OpenSecureChannel |
| **PROFINET** | UDP | 34964 | DCP Identify/Set + Record Data Read/Write |

### Fieldbus (11/11 — 4 100% software + 7 SDKs de hardware)

| Protocolo | Transporte | Porta | Função |
|------|------|------|------|
| **HART** | Serial FSK | — | Quadros curtos/longos, Command 0/3, checksum XOR |
| **CC-Link** | RS-485 | — | Polling mestre-escravo, CRC-16/XMODEM |
| **DNP3** | TCP, Serial | 20000 | Segmentação/remontagem na camada de transporte, polling Classe 0, CRC-16/DNP |
| **IEC 61850** | TCP | 102 | MMS Initiate/Conclude, Read/Write, BER-TLV |
| **PROFIBUS** | Serial | — | Implementação 100% software — requer hardware CP 5611 |
| **CANopen** | CAN, Gateway | — | Leitura/escrita SDO, NMT start/stop/reset, Heartbeat |
| **DeviceNet** | CAN, Gateway | — | Explicit Messaging, Poll, conexões I/O |
| **Foundation Fieldbus** | Serial | — | Implementação 100% software — requer placa de interface FF H1 |
| **AS-Interface** | Serial | — | Implementação 100% software — requer gateway ASi |
| **IO-Link** | Serial | — | Implementação 100% software — requer IO-Link Master |
| **CC-Link IE** | Ethernet | — | Implementação 100% software — requer gateway CC-Link IE |

### IoT / Mensagens (2/2 concluídos)

| Protocolo | Transporte | Porta | Função |
|------|------|------|------|
| **MQTT** | TCP, WS | 1883 | 3.1.1 CONNECT/PUBLISH/SUBSCRIBE/PING |
| **HART-IP** | TCP, UDP | 5094 | HART sobre TCP, reutiliza a codificação/decodificação HART |

### Barramento Automotivo (5/5 — 2 100% software + 3 SDKs de hardware)

| Protocolo | Transporte | Baud rate | Função |
|------|------|--------|------|
| **LIN** | UART | — | Quadros mestre/escravo, checksum PID, Classic/Enhanced Checksum |
| **K-Line** | Serial | 10400 | ISO 9141/14230, Fast Init de 5 baud, OBD-II SID 01/03/09 |
| **FlexRay** | CAN, Serial | — | Codificação/decodificação Slot/Frame, CRC-16/XMODEM, SocketCAN + SerialBridge |
| **SAE J1850** | CAN | — | Codificação/decodificação PWM/VPW, Mode 01/03/0A, CRC-8 |
| **MOST** | Serial | — | Codificação/decodificação de quadros hex, Read/Write/Status, SerialBridge |

### Edifícios / Iluminação (2/2 — 1 100% software + 1 SDK de hardware)

| Protocolo | Transporte | Função |
|------|------|------|
| **DALI** | Serial | Quadros diretos de 16 bits, quadros de retorno de 8 bits, comandos padrão (Off/Max/Dim) |
| **LonWorks** | Serial | Codificação/decodificação 100% software — requer chip Neuron ou gateway |

### Ponte de Hardware (12/12 — todos com drivers CmdBridge/SerialBridge)

| Protocolo | Tipo de SDK | Função |
|------|----------|------|
| **EtherCAT** | CmdBridge + hex-codec | Subcomandos Upload/Download/Slaves, codificação/decodificação hexadecimal CoE |
| **POWERLINK** | CmdBridge + hex-codec | Subcomandos Read/Write/Status, quadros SoC/Preq/Pres |
| **SERCOS III** | CmdBridge + hex-codec | Subcomandos Read/Write/Phase, transição de fases (NRT→CP4) |
| **SERCOS I/II** | CmdBridge + hex-codec | Read/Write/Status, topologia em anel de fibra óptica |
| **ControlNet** | CmdBridge + hex-codec | Read/Write/Status, agendamento CTDMA, produtor/consumidor |
| **Interbus** | CmdBridge | Read/Decode, topologia em anel, quadros IBS CMD |
| **WorldFIP** | CmdBridge | Read/Write/Decode, produtor/consumidor, árbitro de barramento |
| **Lightbus** | CmdBridge | Read/Decode, interconexão por fibra óptica, 32 nós |
| **Modbus Plus** | CmdBridge + hex-codec | Read/Write/Decode, passagem de token, peer-to-peer |
| **ISA100.11a** | CmdBridge + hex-codec | Read/Write/List, sem fio 6LoWPAN, malha |
| **WirelessHART** | CmdBridge + hex-codec | Read/Write/Scan, MAC TSMP, auto-organização |
| **SAE J1850** | CmdBridge | Diagnóstico automotivo convertido para CmdBridge, requer interface J1850 |

### Barramento de Sistema (3/3 — todos com drivers sysfs/procfs)

| Protocolo | Tipo de SDK | Função |
|------|----------|------|
| **PCI/PCIe** | sysfs — `/sys/bus/pci/devices/<BDF>/config` | Leitura/escrita do espaço de configuração, PipeTransport, requer CAP_SYS_ADMIN |
| **VME/VPX** | procfs — `/proc/vme/<slot>` | Endereçamento A16/A24/A32, PipeTransport, requer módulo vme_tsi148 |
| **CompactPCI** | sysfs — `/sys/bus/pci/devices/<BDF>/config` | PICMG 2.0, hot plug, 3U/6U, requer cpci_hotplug |

---

## Guia de Uso

A referência da API e os exemplos de uso (instalação, leitura/escrita Modbus, pool de conexões, MQTT, presets Vendor, conversão por gateway, protocolos personalizados) foram movidos para um documento separado:

> **[API.md](API.md) — Referência completa da API e guia de uso**
>
> Versão em inglês: [API.en.md](../../API.en.md)

---

## Estrutura do Projeto

```
industrial-protocols-go/
├── go.work                       # workspace que agrega todos os módulos
├── kernel/                       # módulos principais
│   ├── connection/               # ConnectionManager + pool de conexões
│   ├── config/                   # ConfigRepository (YAML)
│   ├── gateway/                  # GatewayEngine conversão de protocolos
│   ├── bridge/                   # Bridge ponte para processos externos
│   ├── vendor/                   # Vendor presets de fabricantes
│   ├── event/                    # barramento de eventos
│   ├── metrics/                  # interface de métricas
│   ├── security/                 # segurança TLS
│   ├── pipeline/                 # cadeia de middlewares
│   └── transport/                # transportes TCP/UDP/Pipe
│
├── protocols/                    # 40 módulos de protocolos
│   ├── ethernet/                 # 5 Ethernet industrial (todos concluídos)
│   ├── fieldbus/                 # 11 fieldbus (4 100% software + 7 SDKs de hardware)
│   ├── iot/                      # 2 IoT/mensagens (todos concluídos)
│   ├── automotive/               # 5 barramento automotivo (2 100% software + 3 SDKs de hardware)
│   ├── building/                 # 2 edifícios/iluminação (1 100% software + 1 SDK de hardware)
│   ├── bridge/                   # 12 pontes de hardware (CmdBridge/SerialBridge)
│   └── system/                   # 3 barramentos de sistema (drivers sysfs/procfs)
│
├── examples/modbus_basic/        # exemplo Modbus TCP
├── _tools/                       # Makefile + scripts auxiliares
└── docs/superpowers/             # documentos de design
```

## Testes

```bash
make test         # teste completo
make test-unit    # apenas testes unitários
make vet && make fmt
```

---

## Apoie-nos

Se este projeto ajudou você, sinta-se à vontade para fazer uma doação e apoiar nossa manutenção contínua.

### Alipay / WeChat

| Alipay | WeChat |
|--------|------|
| ![Alipay](../../alipay.png) | ![WeChat](../../weixinpay.png) |

### Transferências Globais (remessa bancária internacional)

Remessa bancária para usuários fora da China continental:

**Informações do beneficiário:**

| Campo | Conteúdo |
|------|------|
| Nome do beneficiário | WANG KEXUN |
| Número da conta do beneficiário | 881015918251 |

**Banco beneficiário:**

| Campo | Conteúdo |
|------|------|
| Nome do banco | ZA Bank Limited |
| SWIFT Code | AABLHKHHXXX |
| Código bancário | 387 |
| Endereço do banco | Core F, Cyberport 3, 100 Cyberport Road, Hong Kong |

**Banco intermediário para remessas transfronteiriças (banco correspondente, se necessário):**

> Observe: estas são informações do banco intermediário (correspondente) para remessas transfronteiriças, não do banco beneficiário. Consulte o seu banco emissor para saber se é necessário fornecer informações do banco intermediário.

- **Para depósitos em dólares de Hong Kong, renminbi e dólares americanos** — o banco intermediário é o Citibank:
  - Nome do banco: Citibank N.A. Hong Kong
  - SWIFT Code: CITIHKHXXXX
  - Código bancário: 006
  - Nome da agência: Hong Kong Branch
  - Código da agência: 391
  - Endereço do banco: Citibank Tower, Citibank Plaza, 3 Garden Road, Central, Hong Kong
- **Para outras moedas** — o banco intermediário é o BNY Mellon:
  - Nome do banco: THE BANK OF NEW YORK MELLON
  - SWIFT Code: IRVTUS3NXXX
  - Endereço do banco: THE BANK OF NEW YORK MELLON, 240 GREENWICH STREET, NEW YORK, United States

---

## License

MIT — Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz
