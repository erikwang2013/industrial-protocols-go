# SDK GatewayBridge Foundation Fieldbus

Foundation Fieldbus H1/HSE mediante puente de gateway TCP.

## Hardware

| Gateway | Interfaz | IP predeterminada |
|---------|-----------|-------------|
| NI USB-8486 | Interfaz USB H1 | asignada por el host |
| Softing FFusb | Interfaz USB H1 | asignada por el host |
| P+F HD2-GTR-4PA | Gateway H1 a Ethernet | 192.168.0.20 |

## Cableado

- Troncal H1 (par trenzado, blindado) con terminador en ambos extremos
- Acondicionador de alimentación de fieldbus para H1 (24 VDC, 350-500 mA por segmento)
- El gateway conecta el segmento H1 al puente TCP mediante Ethernet

## Configuración de la IP del gateway

```go
driver, codec, err := foundationfieldbus.NewGatewayDriver("192.168.0.20:2002")
handler, err := foundationfieldbus.ReadyHandler("192.168.0.20:2002")
```
