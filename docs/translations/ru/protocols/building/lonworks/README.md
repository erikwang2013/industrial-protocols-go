# SDK LonWorks GatewayBridge

Протокол LonWorks (ANSI/CEA-709.1) через TCP-мост-шлюз.

## Аппаратура

| Шлюз | Интерфейс | IP по умолчанию |
|---------|-----------|-------------|
| Echelon U60 | Сетевой интерфейс USB FT-10 | назначается хостом |
| Echelon U70 | Интерфейс USB TP/XF-1250 | назначается хостом |
| Loytec L-IP | Маршрутизатор LonWorks/IP | 192.168.0.90 |

## Подключение

- FT-10 (Free Topology): витая пара без учёта полярности, свободная топология до 500 м
- TP/XF-1250: шинная топология с терминатором 105 Ом
- Ethernet на маршрутизаторе L-IP для TCP-моста

## Настройка IP шлюза

```go
driver, codec, err := lonworks.NewGatewayDriver("192.168.0.90:2009")
handler, err := lonworks.ReadyHandler("192.168.0.90:2009")
```
