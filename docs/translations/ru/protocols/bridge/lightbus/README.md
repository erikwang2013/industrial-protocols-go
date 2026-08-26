# SDK Lightbus GatewayBridge

Протокол Beckhoff Lightbus (оптоволокно) через TCP-мост-шлюз.

## Аппаратура

| Шлюз | Интерфейс | IP по умолчанию |
|---------|-----------|-------------|
| Beckhoff FC2001 | PCI-карта Lightbus | назначается хостом |
| Beckhoff BK2000 | Связующее устройство Lightbus | 192.168.0.80 |
| Beckhoff FC9001 | Ethernet-адаптер Lightbus | 192.168.0.81 |

## Подключение

- Кольцевая топология на пластиковом оптическом волокне (POF)
- FC2001/FC9001 подключают кольцо к TCP-мосту через Ethernet
- Каждое устройство имеет разъёмы оптоволокна TX и RX

## Настройка IP шлюза

```go
driver, codec, err := lightbus.NewGatewayDriver("192.168.0.80:2008")
handler, err := lightbus.ReadyHandler("192.168.0.80:2008")
```
