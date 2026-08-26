# SDK Foundation Fieldbus GatewayBridge

Foundation Fieldbus H1/HSE via ponte de gateway TCP.

## Hardware

| Gateway | Interface | IP Padrão |
|---------|-----------|-------------|
| NI USB-8486 | Interface H1 USB | atribuído pelo host |
| Softing FFusb | Interface H1 USB | atribuído pelo host |
| P+F HD2-GTR-4PA | Gateway H1 para Ethernet | 192.168.0.20 |

## Cabeamento

- Tronco H1 (par trançado, blindado) com terminador em ambas as extremidades
- Condicionador de alimentação do fieldbus para H1 (24 VCC, 350-500 mA por segmento)
- O gateway conecta o segmento H1 à ponte TCP via Ethernet

## Configuração do IP do Gateway

```go
driver, codec, err := foundationfieldbus.NewGatewayDriver("192.168.0.20:2002")
handler, err := foundationfieldbus.ReadyHandler("192.168.0.20:2002")
```
