# SDK GatewayBridge PROFIBUS

Protocolo PROFIBUS DP/PA mediante puente de gateway TCP.

## Hardware

| Gateway | Interfaz | IP predeterminada |
|---------|-----------|-------------|
| Anybus Communicator | Esclavo PROFIBUS DP-V1 | 192.168.0.50 |
| Proxy Siemens CP 5611 | Maestro PCI/PCIe PROFIBUS | asignada por el host |
| HMS Fieldbus Gateway | Anybus NP40 | 192.168.0.51 |

## Cableado

- DB9 hembra en el gateway a la red PROFIBUS (línea A verde, línea B roja)
- Resistencia de terminación ACTIVADA en ambos extremos de la red
- Ethernet al puerto de gestión del gateway

## Configuración de la IP del gateway

```go
driver, codec, err := profibus.NewGatewayDriver("192.168.0.50:2000")
handler, err := profibus.ReadyHandler("192.168.0.50:2000")
```
