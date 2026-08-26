# SDK WorldFIP GatewayBridge

Protocolo WorldFIP via ponte de gateway TCP.

## Hardware

| Gateway | Interface | IP Padrão |
|---------|-----------|-------------|
| FIPIO Agent | Agente de fieldbus WorldFIP | 192.168.0.70 |
| FIP Gateway (Alstom) | WorldFIP para Ethernet | 192.168.0.71 |
| NI FIP-USB | Interface WorldFIP USB | atribuído pelo host |

## Cabeamento

- D-SUB de 9 pinos no gateway para o tronco WorldFIP (FIP1 = Data+, FIP2 = Data-)
- Terminador de linha (120 ohms) em ambas as extremidades do barramento
- Ethernet para a ponte TCP

## Configuração do IP do Gateway

```go
driver, codec, err := worldfip.NewGatewayDriver("192.168.0.70:2007")
handler, err := worldfip.ReadyHandler("192.168.0.70:2007")
```
