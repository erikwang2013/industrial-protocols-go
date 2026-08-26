# SDK PROFIBUS GatewayBridge

Протокол PROFIBUS DP/PA через TCP-мост-шлюз.

## Аппаратура

| Шлюз | Интерфейс | IP по умолчанию |
|---------|-----------|-------------|
| Anybus Communicator | Ведомое устройство PROFIBUS DP-V1 | 192.168.0.50 |
| Прокси Siemens CP 5611 | Мастер PCI/PCIe PROFIBUS | назначается хостом |
| HMS Fieldbus Gateway | Anybus NP40 | 192.168.0.51 |

## Подключение

- DB9 (розетка) на шлюзе к сети PROFIBUS (линия A — зелёный, линия B — красный)
- Терминатор включён на обоих концах сети
- Ethernet к порту управления шлюза

## Настройка IP шлюза

```go
driver, codec, err := profibus.NewGatewayDriver("192.168.0.50:2000")
handler, err := profibus.ReadyHandler("192.168.0.50:2000")
```
