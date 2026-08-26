# SDK Foundation Fieldbus GatewayBridge

Foundation Fieldbus H1/HSE через TCP-мост-шлюз.

## Аппаратура

| Шлюз | Интерфейс | IP по умолчанию |
|---------|-----------|-------------|
| NI USB-8486 | USB-интерфейс H1 | назначается хостом |
| Softing FFusb | USB-интерфейс H1 | назначается хостом |
| P+F HD2-GTR-4PA | Шлюз H1 в Ethernet | 192.168.0.20 |

## Подключение

- Магистраль H1 (экранированная витая пара) с терминаторами на обоих концах
- Кондиционер питания Fieldbus для H1 (24 В DC, 350–500 мА на сегмент)
- Шлюз соединяет сегмент H1 с TCP-мостом через Ethernet

## Настройка IP шлюза

```go
driver, codec, err := foundationfieldbus.NewGatewayDriver("192.168.0.20:2002")
handler, err := foundationfieldbus.ReadyHandler("192.168.0.20:2002")
```
