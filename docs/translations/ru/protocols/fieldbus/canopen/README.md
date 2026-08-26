# SDK протокола CANopen

CANopen — протокол более высокого уровня на основе CAN для встраиваемых систем управления. Этот пакет предоставляет кодек CANopen и драйвер SocketCAN.

## Обзор протокола

CANopen использует стандартные 11-битные идентификаторы CAN со следующим предопределённым набором соединений:

| Функция    | CAN ID              | Описание                   |
|------------|---------------------|-------------------------------|
| NMT        | 0x000               | Управление сетью             |
| SYNC       | 0x080               | Сообщение синхронизации      |
| SDO (tx)   | 0x580 + NodeID      | Service Data Object (сервер) |
| SDO (rx)   | 0x600 + NodeID      | Service Data Object (клиент) |
| PDO1 (tx)  | 0x180 + NodeID      | Process Data Object 1        |
| Heartbeat  | 0x700 + NodeID      | Heartbeat / Bootup           |

## Требования к аппаратуре

- **Linux** с поддержкой SocketCAN (включён `CONFIG_CAN`)
- CAN-интерфейс (например, `can0`, `vcan0` для виртуальной CAN)
- Аппаратный трансивер CAN (например, MCP2515, SJA1000 или USB-CAN адаптер)

### Настройка виртуального CAN-интерфейса (для тестирования)

```bash
sudo modprobe can
sudo modprobe can_raw
sudo modprobe vcan
sudo ip link add dev vcan0 type vcan
sudo ip link set up vcan0
```

## Использование

```go
package main

import (
    "fmt"
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/fieldbus/canopen"
)

func main() {
    p := canopen.New()
    codec, _ := p.NewCodec("can")

    // Чтение записи словаря объектов SDO
    req := &kernel.Request{
        Function: "sdo_read",
        Metadata: map[string]any{
            "index": float64(0x1000), // Тип устройства
            "sub":   float64(0),
        },
    }
    raw, _ := codec.Encode(req)
    fmt.Printf("SDO read frame: %X\n", raw)

    // NMT: запуск удалённого узла
    req2 := &kernel.Request{Function: "nmt_start"}
    raw2, _ := codec.Encode(req2)
    fmt.Printf("NMT start frame: %X\n", raw2)
}
```

## Поддерживаемые функции

| Функция     | Описание                          |
|-------------|--------------------------------------|
| `sdo_read`  | Чтение записи словаря объектов       |
| `sdo_write` | Запись записи словаря объектов       |
| `nmt_start` | Запуск удалённого узла (NMT)         |
| `nmt_stop`  | Остановка удалённого узла (NMT)      |
| `nmt_reset` | Сброс удалённого узла (NMT)          |
| `heartbeat` | Отправка сообщения heartbeat/bootup  |

## Тестирование

```bash
go test ./... -v
```

Примечание: тесты драйвера SocketCAN требуют Linux с CAN-аппаратурой или виртуальным CAN-интерфейсом.
