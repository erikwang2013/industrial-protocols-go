# SDK GatewayBridge AS-Interface

Protocolo AS-Interface (ASi) mediante puente de gateway TCP.

## Hardware

| Gateway | Interfaz | IP predeterminada |
|---------|-----------|-------------|
| Bihl+Wiedemann BWU2044 | Maestro ASi (Ethernet) | 192.168.0.30 |
| P+F VBG-ENX-K20-DMD-EV | Gateway ASi | 192.168.0.31 |
| ifm AC1375 | Controlador ASi E | 192.168.0.32 |

## Cableado

- Cable ASi amarillo (alimentación + datos) del gateway a los esclavos
- Cable negro de alimentación auxiliar (24 VDC para actuadores) opcional
- Ethernet en el gateway para el puente TCP

## Configuración de la IP del gateway

```go
driver, codec, err := asinterface.NewGatewayDriver("192.168.0.30:2003")
handler, err := asinterface.ReadyHandler("192.168.0.30:2003")
```
