# SDK POWERLINK CmdBridge

POWERLINK (Ethernet POWERLINK) — протокол Ethernet реального времени для промышленной автоматизации. Этот пакет предоставляет кодек POWERLINK через утилиту командной строки `openPOWERLINK_demo`.

## CLI-инструмент

Использует демонстрационное приложение стека openPOWERLINK.

### Установка

```bash
git clone https://github.com/OpenAutomationTechnologies/openPOWERLINK.git
cd openPOWERLINK
cmake -B build && cmake --build build
sudo cmake --install build
```

Проверка: `openPOWERLINK_demo --help`

## Использование

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/powerlink"
)

func main() {
    p := powerlink.New()
    codec, _ := p.NewCodec("cmd")

    req := &kernel.Request{
        Function: "read",
        Address:  "0x2000",
        Count:    8,
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## Драйвер

```go
b, c, err := powerlink.NewCmdDriver("/usr/bin/openPOWERLINK_demo")
```

## Поддерживаемые функции

| Функция | Описание                   |
|----------|-------------------------------|
| `read`   | Чтение записи словаря объектов  |
| `write`  | Запись записи словаря объектов |
| `status` | Запрос состояния узла/NMT      |

## Тестирование

```bash
go test ./... -v
```
