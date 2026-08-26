[English](../../README.en.md) | [中文](../../README.md) | [한국어](../ko/README.md) | Русский | [Deutsch](../de/README.md) | [Français](../fr/README.md) | [Español](../es/README.md) | [Português](../pt/README.md) | [हिन्दी](../hi/README.md) | [العربية](../ar/README.md) | [বাংলা](../bn/README.md) | [Bahasa Indonesia](../id/README.md) | [日本語](../ja/README.md)

# Industrial Protocols Go

Набор протоколов промышленной сетевой связи на Go — слоистая архитектура + middleware, охватывает 40 промышленных протоколов: 14 чисто программных реализаций + 26 аппаратных SDK (все с драйверами).

> См. реализацию на PHP: [github.com/erikwang2013/industrial-protocols](https://github.com/erikwang2013/industrial-protocols)

---

## Проектные решения

### Слои архитектуры

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
│  Отдельный модуль на каждый протокол  (TCP/UDP/Serial) │
│                                                │
│  14 Pure-Soft   +   26 Hardware SDKs (with drivers)│
└──────────────────────────────────────────────────┘
```

### Идея дизайна

**Микроядро + протокольные SDK.** Ядро определяет только интерфейсы и сквозные задачи (пул соединений, повторы, размыкание цепи, события, метрики) и не содержит реализаций конкретных протоколов. Каждый протокол — это отдельный Go-модуль, подключаемый по необходимости.

**Разделение на слои.** Четыре уровня абстракции:

| Слой | Обязанности | Ключевые типы |
|----|------|---------|
| **Transport** | Нижний канал связи | интерфейс `Transport` — `TCPTransport`, `UDPTransport`, `PipeTransport` |
| **Codec** | Кодирование/декодирование протокола | интерфейс `Codec` — `Encode(*Request) ([]byte, error)` / `Decode([]byte) (*Response, error)` |
| **Pipeline** | Сквозные middleware | `Middleware` — `Timeout`, `Retry`, `CircuitBreaker`, `Logger` |
| **Kernel** | Жизненный цикл устройств | `ConnectionManager`, `ConfigRepository` |

**Автору протокола достаточно реализовать два интерфейса: `Protocol` + `Codec`.** Транспортный уровень, пул соединений, повторы, таймауты и размыкание цепи — всё переиспользуется из kernel.

### Модули ядра

| Модуль | Путь | Описание |
|------|------|------|
| ConnectionManager | `kernel/connection/` | Регистрация устройств, пул соединений, проверка здоровья; стратегии Lazy/Eager/Pooled |
| ConfigRepository | `kernel/config/` | Загрузка конфигурации устройств из YAML/JSON |
| GatewayEngine | `kernel/gateway/` | Движок правил преобразования между протоколами (Modbus→MQTT и т. п.) |
| Bridge | `kernel/bridge/` | Мост к внешним процессам (связь через stdin/stdout) |
| Vendor | `kernel/vendor/` | Предустановленные параметры вендоров (Siemens S7, Rockwell AB и др.) |
| Event | `kernel/event/` | Шина событий на основе каналов |
| Metrics | `kernel/metrics/` | Интерфейс сбора метрик (Prometheus + noop) |
| Security | `kernel/security/` | Обёртка безопасности транспортного уровня TLS |

---

## Поддержка протоколов

### Промышленный Ethernet (5/5 готово)

| Протокол | Транспорт | Порт | Функции |
|------|------|------|------|
| **Modbus** | TCP, RTU | 502 | FC 01/02/03/04/05/06/16, CRC16, MBAP |
| **BACnet** | UDP | 47808 | Обнаружение устройств Who-Is/I-Am, ReadProperty |
| **EtherNet/IP** | TCP | 44818 | ENIP RegisterSession + чтение тегов CIP |
| **OPC UA** | TCP | 4840 | Binary HEL + OpenSecureChannel |
| **PROFINET** | UDP | 34964 | DCP Identify/Set + Record Data Read/Write |

### Полевые шины (11/11 — 4 чисто программных + 7 аппаратных SDK)

| Протокол | Транспорт | Порт | Функции |
|------|------|------|------|
| **HART** | Serial FSK | — | Короткие/длинные кадры, Command 0/3, контроль XOR |
| **CC-Link** | RS-485 | — | Опрос «ведущий-ведомый», CRC-16/XMODEM |
| **DNP3** | TCP, Serial | 20000 | Сегментация/сборка на транспортном уровне, опрос Class 0, CRC-16/DNP |
| **IEC 61850** | TCP | 102 | MMS Initiate/Conclude, Read/Write, BER-TLV |
| **PROFIBUS** | Serial | — | Чисто программная реализация — требуется аппаратура CP 5611 |
| **CANopen** | CAN, Gateway | — | SDO чтение/запись, NMT пуск/стоп/сброс, Heartbeat |
| **DeviceNet** | CAN, Gateway | — | Explicit Messaging, Poll, I/O-соединения |
| **Foundation Fieldbus** | Serial | — | Чисто программная реализация — требуется интерфейсная карта FF H1 |
| **AS-Interface** | Serial | — | Чисто программная реализация — требуется шлюз ASi |
| **IO-Link** | Serial | — | Чисто программная реализация — требуется IO-Link Master |
| **CC-Link IE** | Ethernet | — | Чисто программная реализация — требуется шлюз CC-Link IE |

### IoT / сообщения (2/2 готово)

| Протокол | Транспорт | Порт | Функции |
|------|------|------|------|
| **MQTT** | TCP, WS | 1883 | 3.1.1 CONNECT/PUBLISH/SUBSCRIBE/PING |
| **HART-IP** | TCP, UDP | 5094 | HART поверх TCP, переиспользует кодек HART |

### Автомобильные шины (5/5 — 2 чисто программных + 3 аппаратных SDK)

| Протокол | Транспорт | Скорость | Функции |
|------|------|--------|------|
| **LIN** | UART | — | Кадры «ведущий-ведомый», контроль PID, Classic/Enhanced Checksum |
| **K-Line** | Serial | 10400 | ISO 9141/14230, 5-baud Fast Init, OBD-II SID 01/03/09 |
| **FlexRay** | CAN, Serial | — | Кодирование Slot/Frame, CRC-16/XMODEM, SocketCAN + SerialBridge |
| **SAE J1850** | CAN | — | Кодирование PWM/VPW, Mode 01/03/0A, CRC-8 |
| **MOST** | Serial | — | Кодирование Hex-frame, Read/Write/Status, SerialBridge |

### Здания / освещение (2/2 — 1 чисто программный + 1 аппаратный SDK)

| Протокол | Транспорт | Функции |
|------|------|------|
| **DALI** | Serial | 16-битные прямые кадры, 8-битные обратные, стандартные команды (Off/Max/Dim) |
| **LonWorks** | Serial | Чисто программное кодирование — требуется чип Neuron или шлюз |

### Аппаратные мосты (12/12 — у всех есть драйверы CmdBridge/SerialBridge)

| Протокол | Тип SDK | Функции |
|------|----------|------|
| **EtherCAT** | CmdBridge + hex-codec | Подкоманды Upload/Download/Slaves, шестнадцатеричное кодирование CoE |
| **POWERLINK** | CmdBridge + hex-codec | Подкоманды Read/Write/Status, кадры SoC/Preq/Pres |
| **SERCOS III** | CmdBridge + hex-codec | Подкоманды Read/Write/Phase, переключение фаз (NRT→CP4) |
| **SERCOS I/II** | CmdBridge + hex-codec | Read/Write/Status, оптоволоконная кольцевая топология |
| **ControlNet** | CmdBridge + hex-codec | Read/Write/Status, планировщик CTDMA, модель «производитель/потребитель» |
| **Interbus** | CmdBridge | Read/Decode, кольцевая топология, кадры IBS CMD |
| **WorldFIP** | CmdBridge | Read/Write/Decode, «производитель/потребитель», арбитр шины |
| **Lightbus** | CmdBridge | Read/Decode, оптоволоконное соединение, 32 узла |
| **Modbus Plus** | CmdBridge + hex-codec | Read/Write/Decode, передача маркера, одноранговая сеть |
| **ISA100.11a** | CmdBridge + hex-codec | Read/Write/List, беспроводная 6LoWPAN, mesh-сеть |
| **WirelessHART** | CmdBridge + hex-codec | Read/Write/Scan, TSMP MAC, самоорганизация |
| **SAE J1850** | CmdBridge | Автомобильная диагностика через CmdBridge, требуется интерфейс J1850 |

### Системные шины (3/3 — у всех есть драйверы sysfs/procfs)

| Протокол | Тип SDK | Функции |
|------|----------|------|
| **PCI/PCIe** | sysfs — `/sys/bus/pci/devices/<BDF>/config` | Чтение/запись конфигурационного пространства, PipeTransport, требуется CAP_SYS_ADMIN |
| **VME/VPX** | procfs — `/proc/vme/<slot>` | Адресация A16/A24/A32, PipeTransport, требуется модуль vme_tsi148 |
| **CompactPCI** | sysfs — `/sys/bus/pci/devices/<BDF>/config` | PICMG 2.0, горячая замена, 3U/6U, требуется cpci_hotplug |

---

## Руководство по использованию

Справочник по API и примеры использования (установка, чтение/запись Modbus, пул соединений, MQTT, предустановки Vendor, преобразование шлюза, пользовательские протоколы) перенесены в отдельный документ:

> **[API.md](API.md) — полный справочник по API и руководство по использованию**
>
> Английская версия: [API.en.md](../../API.en.md)

---

## Структура проекта

```
industrial-protocols-go/
├── go.work                       # workspace объединяет все модули
├── kernel/                       # модули ядра
│   ├── connection/               # ConnectionManager + пул соединений
│   ├── config/                   # ConfigRepository (YAML)
│   ├── gateway/                  # GatewayEngine преобразование протоколов
│   ├── bridge/                   # Bridge мост к внешним процессам
│   ├── vendor/                   # Vendor предустановки вендоров
│   ├── event/                    # шина событий
│   ├── metrics/                  # интерфейс метрик
│   ├── security/                 # безопасность TLS
│   ├── pipeline/                 # цепочка middleware
│   └── transport/                # транспорт TCP/UDP/Pipe
│
├── protocols/                    # 40 протокольных модулей
│   ├── ethernet/                 # 5 промышленных Ethernet (все готовы)
│   ├── fieldbus/                 # 11 полевых шин (4 чисто программных + 7 аппаратных SDK)
│   ├── iot/                      # 2 IoT/сообщения (все готовы)
│   ├── automotive/               # 5 автомобильных шин (2 чисто программных + 3 аппаратных SDK)
│   ├── building/                 # 2 здания/освещение (1 чисто программный + 1 аппаратный SDK)
│   ├── bridge/                   # 12 аппаратных мостов (CmdBridge/SerialBridge)
│   └── system/                   # 3 системные шины (драйверы sysfs/procfs)
│
├── examples/modbus_basic/        # пример Modbus TCP
├── _tools/                       # Makefile + вспомогательные скрипты
└── docs/superpowers/             # проектная документация
```

## Тестирование

```bash
make test         # полный прогон тестов
make test-unit    # только модульные тесты
make vet && make fmt
```

---

## Поддержка проекта

Если этот проект был вам полезен, вы можете поддержать нас, чтобы мы продолжали его развивать.

### Alipay / WeChat

| Alipay | WeChat |
|--------|------|
| ![Alipay](../../alipay.png) | ![WeChat](../../weixinpay.png) |

### Международный перевод (банковский)

Для банковских переводов от пользователей за пределами материкового Китая:

**Информация о получателе:**

| Поле | Значение |
|------|------|
| Имя получателя | WANG KEXUN |
| Номер счёта получателя | 881015918251 |

**Банк получателя:**

| Поле | Значение |
|------|------|
| Название банка | ZA Bank Limited |
| SWIFT Code | AABLHKHHXXX |
| Банковский код | 387 |
| Адрес банка | Core F, Cyberport 3, 100 Cyberport Road, Hong Kong |

**Банк-корреспондент для трансграничных переводов (промежуточный банк, при необходимости):**

> Обратите внимание: это информация о банке-корреспонденте (промежуточном банке) для трансграничных переводов, а не о банке получателя. Уточните в своём банке, требуется ли предоставлять информацию о банке-корреспонденте.

- **Переводы в гонконгских долларах, юанях и долларах США** — банк-корреспондент Citibank:
  - Название банка: Citibank N.A. Hong Kong
  - SWIFT Code: CITIHKHXXXX
  - Банковский код: 006
  - Отделение: Hong Kong Branch
  - Код отделения: 391
  - Адрес банка: Citibank Tower, Citibank Plaza, 3 Garden Road, Central, Hong Kong
- **Переводы в других валютах** — банк-корреспондент BNY Mellon:
  - Название банка: THE BANK OF NEW YORK MELLON
  - SWIFT Code: IRVTUS3NXXX
  - Адрес банка: THE BANK OF NEW YORK MELLON, 240 GREENWICH STREET, NEW YORK, United States

---

## Лицензия

MIT — Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz
