# SDK IO-Link GatewayBridge

Protocolo IO-Link via ponte de gateway TCP.

## Hardware

| Gateway | Interface | IP Padrão |
|---------|-----------|-------------|
| ifm AL1332 | IO-Link Master (EtherNet/IP) | 192.168.0.40 |
| Balluff BNI00AZ | IO-Link Master (PROFINET) | 192.168.0.41 |
| SICK SIG200 | IO-Link Master (Ethernet) | 192.168.0.42 |

## Cabeamento

- Conector M12 (4 pinos) para cada porta IO-Link: L+ (marrom), L- (azul), C/Q (preto), sem uso (branco)
- Fonte de alimentação de 24 VCC para o master e os dispositivos
- Ethernet no master para a ponte TCP

## Configuração do IP do Gateway

```go
driver, codec, err := iolink.NewGatewayDriver("192.168.0.40:2004")
handler, err := iolink.ReadyHandler("192.168.0.40:2004")
```
