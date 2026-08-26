# SDK Interbus GatewayBridge

Протокол Interbus через TCP-мост-шлюз.

## Аппаратура

| Шлюз | Интерфейс | IP по умолчанию |
|---------|-----------|-------------|
| Phoenix Contact IBS S5 DSC/I-T | Контроллер-мастер Interbus | 192.168.0.60 |
| Phoenix Contact IBS PCI SC/I-T | PCI-мастер Interbus | назначается хостом |
| HMS Anybus Interbus | Встраиваемый шлюз | 192.168.0.61 |

## Подключение

- 9-контактный D-SUB на шлюзе к удалённой шине Interbus (вход/выход)
- Экран подключён к FE (защитному заземлению) с обоих концов
- Ethernet на шлюзе для TCP-моста

## Настройка IP шлюза

```go
driver, codec, err := interbus.NewGatewayDriver("192.168.0.60:2006")
handler, err := interbus.ReadyHandler("192.168.0.60:2006")
```
