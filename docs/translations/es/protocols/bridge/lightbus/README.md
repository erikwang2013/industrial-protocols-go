# SDK GatewayBridge Lightbus

Protocolo de fibra óptica Beckhoff Lightbus mediante puente de gateway TCP.

## Hardware

| Gateway | Interfaz | IP predeterminada |
|---------|-----------|-------------|
| Beckhoff FC2001 | Tarjeta PCI Lightbus | asignada por el host |
| Beckhoff BK2000 | Acoplador de bus Lightbus | 192.168.0.80 |
| Beckhoff FC9001 | Adaptador Ethernet Lightbus | 192.168.0.81 |

## Cableado

- Topología de anillo de fibra óptica plástica (POF)
- FC2001/FC9001 conecta el anillo al puente TCP mediante Ethernet
- Cada dispositivo tiene conectores de fibra TX y RX

## Configuración de la IP del gateway

```go
driver, codec, err := lightbus.NewGatewayDriver("192.168.0.80:2008")
handler, err := lightbus.ReadyHandler("192.168.0.80:2008")
```
