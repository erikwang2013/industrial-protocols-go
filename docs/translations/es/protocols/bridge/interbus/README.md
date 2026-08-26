# SDK GatewayBridge Interbus

Protocolo Interbus mediante puente de gateway TCP.

## Hardware

| Gateway | Interfaz | IP predeterminada |
|---------|-----------|-------------|
| Phoenix Contact IBS S5 DSC/I-T | Controlador maestro Interbus | 192.168.0.60 |
| Phoenix Contact IBS PCI SC/I-T | Maestro PCI Interbus | asignada por el host |
| HMS Anybus Interbus | Gateway integrado | 192.168.0.61 |

## Cableado

- D-SUB de 9 pines en el gateway al bus remoto Interbus (entrada/salida)
- Blindaje conectado a FE en ambos extremos
- Ethernet al gateway para el puente TCP

## Configuración de la IP del gateway

```go
driver, codec, err := interbus.NewGatewayDriver("192.168.0.60:2006")
handler, err := interbus.ReadyHandler("192.168.0.60:2006")
```
