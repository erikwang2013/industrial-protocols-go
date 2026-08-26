# SDK протокола Modbus Plus

Modbus Plus (MB+) — высокоскоростная промышленная сеть с передачей маркера, разработанная Modicon (Schneider Electric).

## Варианты

| Вариант | Тип моста | Описание |
|---------|-------------|-------------|
| `gateway` | GatewayBridge | TCP-соединение с адаптером SA85/BM85 |
| `cmd` | CmdBridge | CLI-обёртка над утилитой `sa85_cli` |

## Драйвер шлюза

```go
driver, codec, err := modbusplus.NewGatewayDriver("192.168.0.160:2010")
handler, err := modbusplus.ReadyHandler("192.168.0.160:2010")
```

## Cmd-драйвер

```go
driver, codec, err := modbusplus.NewCmdDriver("/usr/bin/sa85_cli")
```

### Установка CLI-инструмента

```bash
# Установите драйвер и инструменты SA85 Modbus Plus
# Обратитесь к документации Schneider Electric по адаптеру SA85
```

Проверка: `sa85_cli --help`

## Аппаратура

| Шлюз | Интерфейс | IP по умолчанию |
|---------|-----------|-------------|
| Schneider SA85 | Адаптер ISA Modbus Plus | назначается хостом |
| Schneider BM85 | Мост/мультиплексор Modbus Plus | 192.168.0.A0 |
| ProSoft MVI56-MBP | Модуль ControlLogix MB+ | назначается хостом |

## Подключение

- Коаксиальный кабель (RG-62) с разъёмами BNC для магистрали MB+
- Терминаторы (78 Ом на каждом конце)
- Мост BM85 подключает MB+ к TCP через Ethernet

## Формат кадра протокола

- Магическое число (2 байта): `0x4D42`
- Назначение (1 байт): адрес узла
- Команда (1 байт): 0x01=чтение, 0x02=запись
- Длина (2 байта): размер полезной нагрузки (big-endian)
- Полезная нагрузка (N байт): данные
- CRC (2 байта): Modbus CRC-16 (little-endian)

## Тестирование

```bash
go test ./... -v
```
