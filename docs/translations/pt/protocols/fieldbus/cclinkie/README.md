# SDK CC-Link IE Field GatewayBridge

Protocolo CC-Link IE Field via ponte de gateway TCP.

## Hardware

| Gateway | Interface | IP Padrão |
|---------|-----------|-------------|
| Mitsubishi MELSEC IQ-R | Mestre CC-Link IE Field | 192.168.3.40 |
| Mitsubishi RJ71GF11-T2 | Módulo CC-Link IE Field | atribuído pelo host |
| HMS Anybus CC-Link IE | Gateway embarcado | 192.168.0.52 |

## Cabeamento

- Ethernet RJ45 para CC-Link IE Field (topologia em anel ou estrela de 1 Gbps)
- Porta de gerenciamento em rede separada
- O gateway conecta a rede de campo à ponte TCP

## Configuração do IP do Gateway

```go
driver, codec, err := cclinkie.NewGatewayDriver("192.168.3.40:2005")
handler, err := cclinkie.ReadyHandler("192.168.3.40:2005")
```
