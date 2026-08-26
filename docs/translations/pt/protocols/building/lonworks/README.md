# SDK LonWorks GatewayBridge

Protocolo LonWorks (ANSI/CEA-709.1) via ponte de gateway TCP.

## Hardware

| Gateway | Interface | IP Padrão |
|---------|-----------|-------------|
| Echelon U60 | Interface de rede USB FT-10 | atribuído pelo host |
| Echelon U70 | Interface USB TP/XF-1250 | atribuído pelo host |
| Loytec L-IP | Roteador LonWorks/IP | 192.168.0.90 |

## Cabeamento

- FT-10 (Free Topology): par trançado insensível a polaridade, até 500 m de topologia livre
- TP/XF-1250: topologia em barramento com terminador de 105 ohms
- Ethernet no roteador L-IP para a ponte TCP

## Configuração do IP do Gateway

```go
driver, codec, err := lonworks.NewGatewayDriver("192.168.0.90:2009")
handler, err := lonworks.ReadyHandler("192.168.0.90:2009")
```
