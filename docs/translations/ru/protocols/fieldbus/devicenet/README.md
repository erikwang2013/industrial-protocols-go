# SDK протокола DeviceNet

DeviceNet — промышленный сетевой протокол на основе CAN для автоматизации производства. Этот пакет предоставляет кодек DeviceNet с драйверами SocketCAN и TCP-шлюза.

## Обзор протокола

DeviceNet использует Common Industrial Protocol (CIP) поверх CAN. Предопределённый набор соединений «Мастер/Ведомый» использует сообщения группы 2 (CAN ID 0x400 + NodeID) для опрашиваемого ввода/вывода и явных сообщений (explicit messaging).

## Требования к аппаратуре

### Режим CAN (SocketCAN)
- **Linux** с поддержкой SocketCAN (включён `CONFIG_CAN`)
- CAN-интерфейс (например, `can0`, `vcan0`)
- CAN-аппаратура с поддержкой DeviceNet (например, Anybus Communicator, HMS IXXAT)

### Режим шлюза
- Подключение TCP/IP к шлюзу DeviceNet
- Шлюз с поддержкой протокола текстовых команд (например, HMS Anybus, Hilscher netX)

### Настройка виртуального CAN-интерфейса (для тестирования)

```bash
sudo modprobe can
sudo modprobe can_raw
sudo modprobe vcan
sudo ip link add dev vcan0 type vcan
sudo ip link set up vcan0
```

## Использование

### Режим CAN

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/fieldbus/devicenet"
)

func main() {
    p := devicenet.New()
    codec, _ := p.NewCodec("can")

    // Запрос опроса
    req := &kernel.Request{
        Function: "poll",
        Data:     []byte{0x01, 0x00},
    }
    raw, _ := codec.Encode(req)

    // Открыть соединение
    req2 := &kernel.Request{Function: "open"}
    raw2, _ := codec.Encode(req2)
    _ = raw
    _ = raw2
}
```

### Режим шлюза

```go
p := devicenet.New()
codec, _ := p.NewCodec("gateway")

req := &kernel.Request{
    Function: "poll",
    Data:     []byte{0xAB, 0xCD},
}
raw, _ := codec.Encode(req)
// raw будет: "poll abcd\n"
```

## Поддерживаемые функции

| Функция | Описание                    |
|----------|--------------------------------|
| `open`   | Открыть явное соединение       |
| `poll`   | Опрос данных ввода/вывода (группа 2) |

## Тестирование

```bash
go test ./... -v
```

Примечание: тесты драйвера SocketCAN требуют Linux. Тесты драйвера шлюза работают в любой ОС.
