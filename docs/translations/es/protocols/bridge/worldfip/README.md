# SDK GatewayBridge WorldFIP

Protocolo WorldFIP mediante puente de gateway TCP.

## Hardware

| Gateway | Interfaz | IP predeterminada |
|---------|-----------|-------------|
| Agente FIPIO | Agente de fieldbus WorldFIP | 192.168.0.70 |
| FIP Gateway (Alstom) | WorldFIP a Ethernet | 192.168.0.71 |
| NI FIP-USB | Interfaz USB WorldFIP | asignada por el host |

## Cableado

- D-SUB de 9 pines en el gateway al troncal WorldFIP (FIP1 = Data+, FIP2 = Data-)
- Terminador de línea (120 ohmios) en ambos extremos del bus
- Ethernet para el puente TCP

## Configuración de la IP del gateway

```go
driver, codec, err := worldfip.NewGatewayDriver("192.168.0.70:2007")
handler, err := worldfip.ReadyHandler("192.168.0.70:2007")
```
