# SDK ControlNet CmdBridge

ControlNet — протокол промышленной сети реального времени, разработанный Allen-Bradley (Rockwell Automation) для высокоскоростного обмена данными, критичного ко времени. Этот пакет предоставляет кодек ControlNet через утилиту командной строки `1784-pcic-cli`.

## CLI-инструмент

Использует CLI-утилиту для интерфейсной карты ControlNet 1784-PCIC.

### Установка

```bash
# Установите драйвер и инструменты Rockwell 1784-PCIC
# Обратитесь к документации Rockwell Automation по SDK RSLinx Classic
```

Проверка: `1784-pcic-cli --help`

## Использование

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/controlnet"
)

func main() {
    p := controlnet.New()
    codec, _ := p.NewCodec("cmd")

    req := &kernel.Request{
        Function: "read",
        Address:  "0x10",
        Count:    8,
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## Драйвер

```go
b, c, err := controlnet.NewCmdDriver("/usr/bin/1784-pcic-cli")
```

## Поддерживаемые функции

| Функция | Описание               |
|----------|---------------------------|
| `read`   | Чтение из узла ControlNet |
| `write`  | Запись в узел ControlNet  |
| `status` | Запрос состояния карты PCIC |

## Тестирование

```bash
go test ./... -v
```
