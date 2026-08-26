# SDK протокола FlexRay CAN

FlexRay — высокоскоростной детерминированный протокол автомобильной связи. Этот пакет предоставляет кодек FlexRay поверх шины CAN с кадрированием на основе циклов и проверкой целостности CRC-16/XMODEM.

## Обзор протокола

FlexRay использует схему множественного доступа с разделением времени (TDMA) с повторяющимися циклами связи. Каждый цикл состоит из статического и динамического сегментов. Этот кодек отображает кадры FlexRay на расширенные 29-битные CAN-кадры.

### Линейный формат

Полезная нагрузка цикла FlexRay:
- **Заголовок** (2 байта): номер цикла (little-endian)
- **Статус** (1 байт): бит 7=PPI (Payload Preamble Indicator), бит 6=NFI, бит 5=SYF, бит 4=SUF
- **Данные** (N байт): полезная нагрузка (макс. 254 байта)
- **CRC** (2 байта): CRC-16/XMODEM по заголовок+статус+данные (little-endian)

Кодирование CAN ID (29-битный расширенный):
- Биты 28–24: тип сообщения (0x01=кадр, 0x02=статус)
- Биты 23–10: зарезервировано
- Биты 15–10: идентификатор слота (6 бит)
- Биты 9–0: номер цикла (10 бит)

## Требования к аппаратуре

- **Linux** с поддержкой SocketCAN
- Адаптер Vector VN7600/VN7640 или Bosch FlexRay-CAN
- Сеть FlexRay с правильной терминацией (смещение 2,5 В)

### Настройка интерфейса CAN

```bash
sudo modprobe can
sudo modprobe can_raw
sudo ip link set can0 type can bitrate 500000
sudo ip link set up can0
```

## Использование

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/automotive/flexray"
)

func main() {
    p := flexray.New()
    codec, _ := p.NewCodec("can")

    // Отправить кадр FlexRay в цикле 5 с установленным PPI
    req := &kernel.Request{
        Function: "frame",
        Data:     []byte{0x42, 0x01},
        Metadata: map[string]any{
            "cycle": float64(5),
            "ppi":   true,
        },
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## Поддерживаемые функции

| Функция | Описание                              |
|----------|------------------------------------------|
| `frame`  | Отправить полезную нагрузку кадра FlexRay |
| `status` | Запросить конфигурацию слота (слот, цикл) |

## Тестирование

```bash
go test ./... -v
```
