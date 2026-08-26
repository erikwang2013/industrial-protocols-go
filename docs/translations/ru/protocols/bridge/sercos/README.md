# SDK SERCOS III CmdBridge

SERCOS III (SErial Real-time COmmunication System) — цифровой интерфейс для управления движением с кольцевой топологией по оптоволокну или меди. Этот пакет предоставляет кодек SERCOS III через утилиту командной строки `netx_cli`.

## CLI-инструмент

Использует CLI-утилиту Hilscher netX SERCOS III.

### Установка

```bash
# Установите драйвер и инструменты Hilscher netX
# См. https://www.hilscher.com/ — драйверы netX
sudo apt-get install netx-driver
```

Проверка: `netx_cli --help`

## Использование

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/sercos"
)

func main() {
    p := sercos.New()
    codec, _ := p.NewCodec("cmd")

    req := &kernel.Request{
        Function: "read",
        Address:  "S-0-51",
        Count:    2,
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## Драйвер

```go
b, c, err := sercos.NewCmdDriver("/usr/bin/netx_cli")
```

## Поддерживаемые функции

| Функция | Описание                 |
|----------|-----------------------------|
| `read`   | Чтение параметра SERCOS IDN/S |
| `write`  | Запись параметра SERCOS IDN/S |
| `phase`  | Установка фазы связи        |

## Тестирование

```bash
go test ./... -v
```
