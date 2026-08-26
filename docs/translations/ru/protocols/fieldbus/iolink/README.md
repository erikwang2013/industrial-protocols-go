# SDK IO-Link GatewayBridge

Протокол IO-Link через TCP-мост-шлюз.

## Аппаратура

| Шлюз | Интерфейс | IP по умолчанию |
|---------|-----------|-------------|
| ifm AL1332 | Мастер IO-Link (EtherNet/IP) | 192.168.0.40 |
| Balluff BNI00AZ | Мастер IO-Link (PROFINET) | 192.168.0.41 |
| SICK SIG200 | Мастер IO-Link (Ethernet) | 192.168.0.42 |

## Подключение

- Разъём M12 (4-контактный) для каждого порта IO-Link: L+ (коричневый), L- (синий), C/Q (чёрный), не используется (белый)
- Источник питания 24 В DC для мастера и устройств
- Ethernet на мастере для TCP-моста

## Настройка IP шлюза

```go
driver, codec, err := iolink.NewGatewayDriver("192.168.0.40:2004")
handler, err := iolink.ReadyHandler("192.168.0.40:2004")
```
