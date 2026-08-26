# SDK EtherCAT CmdBridge

EtherCAT (Ethernet for Control Automation Technology) — высокопроизводительная промышленная шина Ethernet. Этот пакет предоставляет кодек EtherCAT через утилиту командной строки `ethercat`.

## CLI-инструмент

Использует инструмент командной строки IgH EtherCAT Master.

### Установка

```bash
# Debian/Ubuntu
sudo apt-get install ethercat-master

# Из исходников
git clone https://gitlab.com/etherlab.org/ethercat.git
cd ethercat
make && sudo make install
```

Проверка: `ethercat slaves`

## Использование

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/ethercat"
)

func main() {
    p := ethercat.New()
    codec, _ := p.NewCodec("cmd")

    // Загрузить SDO с адреса 0x1000
    req := &kernel.Request{
        Function: "upload",
        Address:  "0x1000",
        Count:    4,
    }
    raw, _ := codec.Encode(req)
    // raw = "upload 0x1000 4\n"
    _ = raw
}
```

## Драйвер

```go
b, c, err := ethercat.NewCmdDriver("/usr/bin/ethercat")
```

## Поддерживаемые функции

| Функция   | Описание            |
|------------|------------------------|
| `upload`   | Чтение SDO по адресу  |
| `download` | Запись SDO по адресу  |
| `slaves`   | Список ведомых устройств EtherCAT |

## Тестирование

```bash
go test ./... -v
```
