# SDK CC-Link IE Field GatewayBridge

Протокол CC-Link IE Field через TCP-мост-шлюз.

## Аппаратура

| Шлюз | Интерфейс | IP по умолчанию |
|---------|-----------|-------------|
| Mitsubishi MELSEC IQ-R | Мастер CC-Link IE Field | 192.168.3.40 |
| Mitsubishi RJ71GF11-T2 | Модуль CC-Link IE Field | назначается хостом |
| HMS Anybus CC-Link IE | Встраиваемый шлюз | 192.168.0.52 |

## Подключение

- RJ45 Ethernet для CC-Link IE Field (кольцевая или звёздная топология 1 Гбит/с)
- Порт управления — в отдельной сети
- Шлюз соединяет полевую сеть с TCP-мостом

## Настройка IP шлюза

```go
driver, codec, err := cclinkie.NewGatewayDriver("192.168.3.40:2005")
handler, err := cclinkie.ReadyHandler("192.168.3.40:2005")
```
