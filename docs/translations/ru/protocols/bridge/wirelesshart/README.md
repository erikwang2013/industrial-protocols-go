# SDK WirelessHART CmdBridge

WirelessHART (IEC 62591) — беспроводной промышленный стандарт сети на основе протокола HART. Этот пакет предоставляет кодек WirelessHART через утилиту командной строки `emerson_1410_cli` (беспроводной шлюз Emerson 1410/1420).

## CLI-инструмент

Использует CLI-утилиту беспроводного шлюза Emerson 1410/1420.

### Установка

```bash
# Установите ПО и инструменты беспроводного шлюза Emerson
# Обратитесь к документации Emerson по настройке шлюза 1410/1420
```

Проверка: `emerson_1410_cli --help`

## Использование

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/wirelesshart"
)

func main() {
    p := wirelesshart.New()
    codec, _ := p.NewCodec("cmd")

    req := &kernel.Request{
        Function: "read",
        Address:  "TT101",
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## Драйвер

```go
b, c, err := wirelesshart.NewCmdDriver("/usr/bin/emerson_1410_cli")
```

## Поддерживаемые функции

| Функция | Описание                |
|----------|----------------------------|
| `read`   | Чтение параметра устройства |
| `write`  | Запись параметра устройства |
| `scan`   | Поиск беспроводных устройств |

## Тестирование

```bash
go test ./... -v
```
