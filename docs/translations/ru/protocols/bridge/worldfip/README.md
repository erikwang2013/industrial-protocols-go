# SDK WorldFIP GatewayBridge

Протокол WorldFIP через TCP-мост-шлюз.

## Аппаратура

| Шлюз | Интерфейс | IP по умолчанию |
|---------|-----------|-------------|
| FIPIO Agent | Агент полевой шины WorldFIP | 192.168.0.70 |
| FIP Gateway (Alstom) | WorldFIP в Ethernet | 192.168.0.71 |
| NI FIP-USB | USB-интерфейс WorldFIP | назначается хостом |

## Подключение

- 9-контактный D-SUB на шлюзе к магистрали WorldFIP (FIP1 = Data+, FIP2 = Data-)
- Терминатор (120 Ом) на обоих концах шины
- Ethernet для TCP-моста

## Настройка IP шлюза

```go
driver, codec, err := worldfip.NewGatewayDriver("192.168.0.70:2007")
handler, err := worldfip.ReadyHandler("192.168.0.70:2007")
```
