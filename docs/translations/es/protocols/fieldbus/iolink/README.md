# SDK GatewayBridge IO-Link

Protocolo IO-Link mediante puente de gateway TCP.

## Hardware

| Gateway | Interfaz | IP predeterminada |
|---------|-----------|-------------|
| ifm AL1332 | Maestro IO-Link (EtherNet/IP) | 192.168.0.40 |
| Balluff BNI00AZ | Maestro IO-Link (PROFINET) | 192.168.0.41 |
| SICK SIG200 | Maestro IO-Link (Ethernet) | 192.168.0.42 |

## Cableado

- Conector M12 (4 pines) por cada puerto IO-Link: L+ (marrón), L- (azul), C/Q (negro), sin usar (blanco)
- Alimentación de 24 VDC para el maestro y los dispositivos
- Ethernet en el maestro para el puente TCP

## Configuración de la IP del gateway

```go
driver, codec, err := iolink.NewGatewayDriver("192.168.0.40:2004")
handler, err := iolink.ReadyHandler("192.168.0.40:2004")
```
