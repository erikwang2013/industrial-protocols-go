# SDK GatewayBridge LonWorks

Protocolo LonWorks (ANSI/CEA-709.1) mediante puente de gateway TCP.

## Hardware

| Gateway | Interfaz | IP predeterminada |
|---------|-----------|-------------|
| Echelon U60 | Interfaz de red USB FT-10 | asignada por el host |
| Echelon U70 | Interfaz USB TP/XF-1250 | asignada por el host |
| Loytec L-IP | Enrutador LonWorks/IP | 192.168.0.90 |

## Cableado

- FT-10 (Free Topology): par trenzado insensible a la polaridad, hasta 500 m de topología libre
- TP/XF-1250: topología de bus con terminador de 105 ohmios
- Ethernet en el enrutador L-IP para el puente TCP

## Configuración de la IP del gateway

```go
driver, codec, err := lonworks.NewGatewayDriver("192.168.0.90:2009")
handler, err := lonworks.ReadyHandler("192.168.0.90:2009")
```
