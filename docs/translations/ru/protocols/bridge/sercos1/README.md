# SDK SERCOS I/II CmdBridge

SERCOS I/II — устаревшая последовательная оптоволоконная версия интерфейса SERCOS для цифрового управления движением. Этот пакет предоставляет кодек SERCOS I/II через утилиту командной строки `sercos_cli`.

## CLI-инструмент

Использует CLI-утилиту оптоволоконного интерфейса SERCOS.

### Установка

```bash
# Драйвер и инструменты интерфейсной карты SERCOS
# Обратитесь к документации вендора по конкретной карте-мастеру SERCOS
```

Проверка: `sercos_cli --help`

## Использование

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/sercos1"
)

func main() {
    p := sercos1.New()
    codec, _ := p.NewCodec("cmd")

    req := &kernel.Request{
        Function: "read",
        Address:  "0x0100",
        Count:    4,
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## Драйвер

```go
b, c, err := sercos1.NewCmdDriver("/usr/bin/sercos_cli")
```

## Поддерживаемые функции

| Функция | Описание              |
|----------|--------------------------|
| `read`   | Чтение SERCOS IDN       |
| `write`  | Запись SERCOS IDN       |
| `status` | Запрос состояния привода |

## Тестирование

```bash
go test ./... -v
```
