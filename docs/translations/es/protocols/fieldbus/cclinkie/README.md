# SDK GatewayBridge CC-Link IE Field

Protocolo CC-Link IE Field mediante puente de gateway TCP.

## Hardware

| Gateway | Interfaz | IP predeterminada |
|---------|-----------|-------------|
| Mitsubishi MELSEC IQ-R | Maestro CC-Link IE Field | 192.168.3.40 |
| Mitsubishi RJ71GF11-T2 | Módulo CC-Link IE Field | asignada por el host |
| HMS Anybus CC-Link IE | Gateway integrado | 192.168.0.52 |

## Cableado

- Ethernet RJ45 para CC-Link IE Field (anillo de 1 Gbps o topología de estrella)
- Puerto de gestión en una red separada
- El gateway conecta la red de campo al puente TCP

## Configuración de la IP del gateway

```go
driver, codec, err := cclinkie.NewGatewayDriver("192.168.3.40:2005")
handler, err := cclinkie.ReadyHandler("192.168.3.40:2005")
```
