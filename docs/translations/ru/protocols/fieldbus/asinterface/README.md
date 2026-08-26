# SDK AS-Interface GatewayBridge

Протокол AS-Interface (ASi) через TCP-мост-шлюз.

## Аппаратура

| Шлюз | Интерфейс | IP по умолчанию |
|---------|-----------|-------------|
| Bihl+Wiedemann BWU2044 | Мастер ASi (Ethernet) | 192.168.0.30 |
| P+F VBG-ENX-K20-DMD-EV | Шлюз ASi | 192.168.0.31 |
| ifm AC1375 | Контроллер ASi ControllerE | 192.168.0.32 |

## Подключение

- Жёлтый кабель ASi (питание + данные) от шлюза к ведомым устройствам
- Дополнительно: чёрный кабель вспомогательного питания (24 В DC для приводов)
- Ethernet на шлюзе для TCP-моста

## Настройка IP шлюза

```go
driver, codec, err := asinterface.NewGatewayDriver("192.168.0.30:2003")
handler, err := asinterface.ReadyHandler("192.168.0.30:2003")
```
