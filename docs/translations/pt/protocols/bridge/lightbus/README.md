# SDK Lightbus GatewayBridge

Protocolo de fibra óptica Beckhoff Lightbus via ponte de gateway TCP.

## Hardware

| Gateway | Interface | IP Padrão |
|---------|-----------|-------------|
| Beckhoff FC2001 | Placa PCI Lightbus | atribuído pelo host |
| Beckhoff BK2000 | Acoplador de barramento Lightbus | 192.168.0.80 |
| Beckhoff FC9001 | Adaptador Ethernet Lightbus | 192.168.0.81 |

## Cabeamento

- Topologia em anel de fibra óptica plástica (POF)
- FC2001/FC9001 conecta o anel à ponte TCP via Ethernet
- Cada dispositivo tem conectores de fibra TX e RX

## Configuração do IP do Gateway

```go
driver, codec, err := lightbus.NewGatewayDriver("192.168.0.80:2008")
handler, err := lightbus.ReadyHandler("192.168.0.80:2008")
```
