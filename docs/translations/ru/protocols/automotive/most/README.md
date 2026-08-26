# SDK протокола MOST Serial

MOST (Media Oriented Systems Transport) — высокоскоростная мультимедийная сетевая технология, используемая в основном в автомобильных информационно-развлекательных системах. Этот пакет предоставляет кодек MOST поверх последовательного адаптера с AT-командным интерфейсом.

## Обзор протокола

MOST использует синхронную последовательную связь по оптоволоконному физическому уровню. Данная реализация подключается через последовательный адаптер, предоставляющий AT-командный интерфейс со скоростью 115200 бод.

### AT-команды

| Команда            | Описание                  |
|--------------------|------------------------------|
| `AT+READ=<addr>,<len>`  | Чтение байтов по адресу |
| `AT+WRITE=<addr>,<hex>` | Запись hex-байтов по адресу |
| `AT+STATUS`             | Запрос состояния кольца/сети |

### Формат ответа

- `+OK:<hex_data>` — успешный ответ
- `+ERR:<code>` — ответ об ошибке

## Требования к аппаратуре

- Оптоволоконная сеть MOST с правильной терминацией
- Адаптер MOST-to-serial (например, USB-адаптер MOST150)
- Последовательный порт: 115200 бод, 8N1

## Использование

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/automotive/most"
)

func main() {
    p := most.New()
    codec, _ := p.NewCodec("serial")

    // Чтение 4 байт по адресу 0x0100
    req := &kernel.Request{
        Function: "read",
        Address:  "0x0100",
        Count:    4,
    }
    raw, _ := codec.Encode(req)
    // raw = "AT+READ=0x0100,4\r\n"
    _ = raw

    // Запись данных по адресу 0x0200
    req2 := &kernel.Request{
        Function: "write",
        Address:  "0x0200",
        Data:     []byte{0xDE, 0xAD, 0xBE, 0xEF},
    }
    raw2, _ := codec.Encode(req2)
    // raw2 = "AT+WRITE=0x0200,DEADBEEF\r\n"
    _ = raw2
}
```

## Драйвер

```go
b, c, err := most.NewSerialDriver("/dev/ttyUSB0")
```

## Поддерживаемые функции

| Функция | Описание                 |
|----------|-----------------------------|
| `read`   | Чтение из адреса MOST      |
| `write`  | Запись данных в адрес MOST |
| `status` | Запрос состояния кольца/сети |

## Тестирование

```bash
go test ./... -v
```
