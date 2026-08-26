# SDK ISA100.11a CmdBridge

ISA100.11a — беспроводной промышленный стандарт сети для автоматизации процессов. Этот пакет предоставляет кодек ISA100.11a через утилиту командной строки `yfgw410_cli` (беспроводной полевой шлюз Yokogawa YFGW410).

## CLI-инструмент

Использует CLI беспроводного полевого шлюза Yokogawa YFGW410.

### Установка

```bash
# Установите ПО и инструменты шлюза Yokogawa YFGW410
# Обратитесь к документации Yokogawa по настройке Field Wireless Gateway
```

Проверка: `yfgw410_cli --help`

## Использование

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/isa100"
)

func main() {
    p := isa100.New()
    codec, _ := p.NewCodec("cmd")

    req := &kernel.Request{
        Function: "read",
        Address:  "DEV001",
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## Драйвер

```go
b, c, err := isa100.NewCmdDriver("/usr/bin/yfgw410_cli")
```

## Поддерживаемые функции

| Функция | Описание                 |
|----------|-----------------------------|
| `read`   | Чтение атрибута устройства  |
| `write`  | Запись атрибута устройства  |
| `list`   | Список зарегистрированных устройств |

## Тестирование

```bash
go test ./... -v
```
