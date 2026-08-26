[English](../../README.en.md) | [中文](../../README.md) | 한국어 | [Русский](../ru/README.md) | [Deutsch](../de/README.md) | [Français](../fr/README.md) | [Español](../es/README.md) | [Português](../pt/README.md) | [हिन्दी](../hi/README.md) | [العربية](../ar/README.md) | [বাংলা](../bn/README.md) | [Bahasa Indonesia](../id/README.md) | [日本語](../ja/README.md)

# Industrial Protocols Go

Go 언어 산업 네트워크 통신 프로토콜 모음 —— 계층 + 미들웨어 아키텍처, 40가지 산업 프로토콜, 14개 순수 소프트웨어 구현 + 26개 하드웨어 SDK(전부 드라이버 포함).

> PHP 구현 참고: [github.com/erikwang2013/industrial-protocols](https://github.com/erikwang2013/industrial-protocols)

---

## 프로젝트 설계

### 아키텍처 계층

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
│  프로토콜별 독립 모듈               (TCP/UDP/Serial)  │
│                                                │
│  14 Pure-Soft   +   26 Hardware SDKs (with drivers)│
└──────────────────────────────────────────────────┘
```

### 설계 철학

**마이크로 커널 + 프로토콜 SDK.** 커널은 인터페이스와 횡단 관심사(연결 풀, 재시도, 차단기, 이벤트, 메트릭)만 정의하며 구체적인 프로토콜 구현을 포함하지 않습니다. 각 프로토콜은 독립적인 Go 모듈로, 필요에 따라 import 합니다.

**계층 분리.** 4계층 추상화:

| 계층 | 역할 | 핵심 타입 |
|----|------|---------|
| **Transport** | 하위 통신 채널 | `Transport` 인터페이스 — `TCPTransport`, `UDPTransport`, `PipeTransport` |
| **Codec** | 프로토콜 인코딩/디코딩 | `Codec` 인터페이스 — `Encode(*Request) ([]byte, error)` / `Decode([]byte) (*Response, error)` |
| **Pipeline** | 횡단 미들웨어 | `Middleware` — `Timeout`, `Retry`, `CircuitBreaker`, `Logger` |
| **Kernel** | 장치 수명주기 | `ConnectionManager`, `ConfigRepository` |

**프로토콜 작성자는 `Protocol` + `Codec` 두 인터페이스만 구현하면 됩니다.** 전송 계층, 연결 풀, 재시도, 타임아웃, 차단기는 전부 kernel에서 재사용합니다.

### 커널 모듈

| 모듈 | 경로 | 설명 |
|------|------|------|
| ConnectionManager | `kernel/connection/` | 장치 등록, 연결 풀, 헬스 체크, Lazy/Eager/Pooled 3가지 전략 지원 |
| ConfigRepository | `kernel/config/` | YAML/JSON 장치 설정 로드 |
| GatewayEngine | `kernel/gateway/` | 크로스 프로토콜 변환 규칙 엔진(Modbus→MQTT 등) |
| Bridge | `kernel/bridge/` | 외부 프로세스 브리지(stdin/stdout 통신) |
| Vendor | `kernel/vendor/` | 벤더 파라미터 프리셋(Siemens S7, Rockwell AB 등) |
| Event | `kernel/event/` | channel 기반 이벤트 버스 |
| Metrics | `kernel/metrics/` | 메트릭 수집 인터페이스(Prometheus + noop) |
| Security | `kernel/security/` | TLS 전송 계층 보안 래퍼 |

---

## 프로토콜 지원

### 산업용 이더넷(5/5 완료)

| 프로토콜 | 전송 | 포트 | 기능 |
|------|------|------|------|
| **Modbus** | TCP, RTU | 502 | FC 01/02/03/04/05/06/16, CRC16, MBAP |
| **BACnet** | UDP | 47808 | Who-Is/I-Am 장치 발견, ReadProperty |
| **EtherNet/IP** | TCP | 44818 | ENIP RegisterSession + CIP Read Tag |
| **OPC UA** | TCP | 4840 | Binary HEL + OpenSecureChannel |
| **PROFINET** | UDP | 34964 | DCP Identify/Set + Record Data Read/Write |

### 필드버스(11/11 — 4 순수 소프트웨어 + 7 하드웨어 SDK)

| 프로토콜 | 전송 | 포트 | 기능 |
|------|------|------|------|
| **HART** | Serial FSK | — | 단프레임/장프레임, Command 0/3, XOR 검증 |
| **CC-Link** | RS-485 | — | 마스터-슬레이브 폴링, CRC-16/XMODEM |
| **DNP3** | TCP, Serial | 20000 | 전송 계층 분할/재조립, Class 0 폴링, CRC-16/DNP |
| **IEC 61850** | TCP | 102 | MMS Initiate/Conclude, Read/Write, BER-TLV |
| **PROFIBUS** | Serial | — | 순수 소프트웨어 구현 — CP 5611 하드웨어 필요 |
| **CANopen** | CAN, Gateway | — | SDO 읽기/쓰기, NMT 시작/정지/리셋, Heartbeat |
| **DeviceNet** | CAN, Gateway | — | Explicit Messaging, Poll, I/O 연결 |
| **Foundation Fieldbus** | Serial | — | 순수 소프트웨어 구현 — FF H1 인터페이스 카드 필요 |
| **AS-Interface** | Serial | — | 순수 소프트웨어 구현 — ASi 게이트웨이 필요 |
| **IO-Link** | Serial | — | 순수 소프트웨어 구현 — IO-Link Master 필요 |
| **CC-Link IE** | Ethernet | — | 순수 소프트웨어 구현 — CC-Link IE 게이트웨이 필요 |

### IoT / 메시지(2/2 완료)

| 프로토콜 | 전송 | 포트 | 기능 |
|------|------|------|------|
| **MQTT** | TCP, WS | 1883 | 3.1.1 CONNECT/PUBLISH/SUBSCRIBE/PING |
| **HART-IP** | TCP, UDP | 5094 | TCP 기반 HART, HART 인코딩/디코딩 재사용 |

### 자동차 버스(5/5 — 2 순수 소프트웨어 + 3 하드웨어 SDK)

| 프로토콜 | 전송 | 보드레이트 | 기능 |
|------|------|--------|------|
| **LIN** | UART | — | 마스터/슬레이브 프레임, PID 검증, Classic/Enhanced Checksum |
| **K-Line** | Serial | 10400 | ISO 9141/14230, 5-baud Fast Init, OBD-II SID 01/03/09 |
| **FlexRay** | CAN, Serial | — | Slot/Frame 인코딩/디코딩, CRC-16/XMODEM, SocketCAN + SerialBridge |
| **SAE J1850** | CAN | — | PWM/VPW 인코딩/디코딩, Mode 01/03/0A, CRC-8 |
| **MOST** | Serial | — | Hex-frame 인코딩/디코딩, Read/Write/Status, SerialBridge |

### 빌딩 / 조명(2/2 — 1 순수 소프트웨어 + 1 하드웨어 SDK)

| 프로토콜 | 전송 | 기능 |
|------|------|------|
| **DALI** | Serial | 16-bit 순방향 프레임, 8-bit 역방향 프레임, 표준 명령(Off/Max/Dim) |
| **LonWorks** | Serial | 순수 소프트웨어 인코딩/디코딩 — Neuron 칩 또는 게이트웨이 필요 |

### 하드웨어 브리지(12/12 — 전부 CmdBridge/SerialBridge 드라이버 포함)

| 프로토콜 | SDK 유형 | 기능 |
|------|----------|------|
| **EtherCAT** | CmdBridge + hex-codec | Upload/Download/Slaves 하위 명령, CoE 16진수 인코딩/디코딩 |
| **POWERLINK** | CmdBridge + hex-codec | Read/Write/Status 하위 명령, SoC/Preq/Pres 프레임 |
| **SERCOS III** | CmdBridge + hex-codec | Read/Write/Phase 하위 명령, 단계 전환(NRT→CP4) |
| **SERCOS I/II** | CmdBridge + hex-codec | Read/Write/Status, 광섬유 링 토폴로지 |
| **ControlNet** | CmdBridge + hex-codec | Read/Write/Status, CTDMA 스케줄링, 프로듀서/컨슈머 |
| **Interbus** | CmdBridge | Read/Decode, 링 토폴로지, IBS CMD 프레임 |
| **WorldFIP** | CmdBridge | Read/Write/Decode, 프로듀서/컨슈머, 버스 중재기 |
| **Lightbus** | CmdBridge | Read/Decode, 광섬유 상호 연결, 32개 노드 |
| **Modbus Plus** | CmdBridge + hex-codec | Read/Write/Decode, 토큰 전달, 피어 투 피어 |
| **ISA100.11a** | CmdBridge + hex-codec | Read/Write/List, 6LoWPAN 무선, 메시 토폴로지 |
| **WirelessHART** | CmdBridge + hex-codec | Read/Write/Scan, TSMP MAC, 자가 구성 |
| **SAE J1850** | CmdBridge | 차량 진단을 CmdBridge로 변환, J1850 인터페이스 필요 |

### 시스템 버스(3/3 — 전부 sysfs/procfs 드라이버 포함)

| 프로토콜 | SDK 유형 | 기능 |
|------|----------|------|
| **PCI/PCIe** | sysfs — `/sys/bus/pci/devices/<BDF>/config` | 구성 공간 읽기/쓰기, PipeTransport, CAP_SYS_ADMIN 필요 |
| **VME/VPX** | procfs — `/proc/vme/<slot>` | A16/A24/A32 어드레싱, PipeTransport, vme_tsi148 모듈 필요 |
| **CompactPCI** | sysfs — `/sys/bus/pci/devices/<BDF>/config` | PICMG 2.0, 핫 플러그, 3U/6U, cpci_hotplug 필요 |

---

## 사용 가이드

API 레퍼런스와 사용 예제(설치, Modbus 읽기/쓰기, 연결 풀, MQTT, Vendor 프리셋, 게이트웨이 변환, 커스텀 프로토콜)는 별도 문서로 이동했습니다:

> **[API.md](API.md) — 전체 API 레퍼런스 및 사용 가이드**
>
> 영문 버전: [API.en.md](../../API.en.md)

---

## 프로젝트 구조

```
industrial-protocols-go/
├── go.work                       # workspace, 모든 모듈 집계
├── kernel/                       # 핵심 모듈
│   ├── connection/               # ConnectionManager + 연결 풀
│   ├── config/                   # ConfigRepository (YAML)
│   ├── gateway/                  # GatewayEngine 프로토콜 변환
│   ├── bridge/                   # Bridge 외부 프로세스 브리지
│   ├── vendor/                   # Vendor 벤더 프리셋
│   ├── event/                    # 이벤트 버스
│   ├── metrics/                  # 메트릭 인터페이스
│   ├── security/                 # TLS 보안
│   ├── pipeline/                 # 미들웨어 체인
│   └── transport/                # TCP/UDP/Pipe 전송
│
├── protocols/                    # 40개 프로토콜 모듈
│   ├── ethernet/                 # 5 산업용 이더넷(전부 완료)
│   ├── fieldbus/                 # 11 필드버스(4 순수 소프트웨어 + 7 하드웨어 SDK)
│   ├── iot/                      # 2 IoT/메시지(전부 완료)
│   ├── automotive/               # 5 자동차 버스(2 순수 소프트웨어 + 3 하드웨어 SDK)
│   ├── building/                 # 2 빌딩/조명(1 순수 소프트웨어 + 1 하드웨어 SDK)
│   ├── bridge/                   # 12 하드웨어 브리지(CmdBridge/SerialBridge)
│   └── system/                   # 3 시스템 버스(sysfs/procfs 드라이버)
│
├── examples/modbus_basic/        # Modbus TCP 예제
├── _tools/                       # Makefile + 보조 스크립트
└── docs/superpowers/             # 설계 문서
```

## 테스트

```bash
make test         # 전체 테스트
make test-unit    # 단위 테스트만
make vet && make fmt
```

---

## 지원해 주세요

이 프로젝트가 도움이 되었다면, 후원으로 유지보수를 계속할 수 있게 지원해 주시면 감사하겠습니다.

### 알리페이 / 위챗

| 알리페이 | 위챗 |
|--------|------|
| ![알리페이](../../alipay.png) | ![위챗](../../weixinpay.png) |

### 해외 송금(국제 은행 송금)

중국 본토 외 사용자를 위한 은행 송금:

**수취인 정보:**

| 항목 | 내용 |
|------|------|
| 수취인 이름 | WANG KEXUN |
| 수취 계좌 번호 | 881015918251 |

**수취 은행:**

| 항목 | 내용 |
|------|------|
| 은행 이름 | ZA Bank Limited |
| SWIFT Code | AABLHKHHXXX |
| 은행 번호 | 387 |
| 은행 주소 | Core F, Cyberport 3, 100 Cyberport Road, Hong Kong |

**해외 송금 중계 은행(필요한 경우):**

> 참고: 이는 해외 송금 중계 은행(중개 은행) 정보이며, 수취 은행 정보가 아닙니다. 송금 은행에 중계 은행 정보가 필요한지 문의하시기 바랍니다.

- **홍콩 달러, 위안화, 미 달러 송금** — 중계 은행은 Citibank:
  - 은행 이름: Citibank N.A. Hong Kong
  - SWIFT Code: CITIHKHXXXX
  - 은행 번호: 006
  - 지점 이름: Hong Kong Branch
  - 지점 번호: 391
  - 은행 주소: Citibank Tower, Citibank Plaza, 3 Garden Road, Central, Hong Kong
- **기타 통화 송금** — 중계 은행은 BNY Mellon:
  - 은행 이름: THE BANK OF NEW YORK MELLON
  - SWIFT Code: IRVTUS3NXXX
  - 은행 주소: THE BANK OF NEW YORK MELLON, 240 GREENWICH STREET, NEW YORK, United States

---

## License

MIT — Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz
