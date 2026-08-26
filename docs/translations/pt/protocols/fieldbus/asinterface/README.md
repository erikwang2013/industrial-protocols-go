# SDK AS-Interface GatewayBridge

Protocolo AS-Interface (ASi) via ponte de gateway TCP.

## Hardware

| Gateway | Interface | IP Padrão |
|---------|-----------|-------------|
| Bihl+Wiedemann BWU2044 | ASi Master (Ethernet) | 192.168.0.30 |
| P+F VBG-ENX-K20-DMD-EV | Gateway ASi | 192.168.0.31 |
| ifm AC1375 | ASi ControllerE | 192.168.0.32 |

## Cabeamento

- Cabo ASi amarelo (alimentação + dados) do gateway para os slaves
- Cabo de alimentação auxiliar preto (24 VCC para atuadores) opcional
- Ethernet no gateway para a ponte TCP

## Configuração do IP do Gateway

```go
driver, codec, err := asinterface.NewGatewayDriver("192.168.0.30:2003")
handler, err := asinterface.ReadyHandler("192.168.0.30:2003")
```
