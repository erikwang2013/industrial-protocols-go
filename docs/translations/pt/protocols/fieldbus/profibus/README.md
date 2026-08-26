# SDK PROFIBUS GatewayBridge

Protocolo PROFIBUS DP/PA via ponte de gateway TCP.

## Hardware

| Gateway | Interface | IP Padrão |
|---------|-----------|-------------|
| Anybus Communicator | Slave PROFIBUS DP-V1 | 192.168.0.50 |
| Proxy Siemens CP 5611 | Mestre PROFIBUS PCI/PCIe | atribuído pelo host |
| HMS Fieldbus Gateway | Anybus NP40 | 192.168.0.51 |

## Cabeamento

- DB9 fêmea no gateway para a rede PROFIBUS (linha A verde, linha B vermelha)
- Resistor de terminação ATIVADO em ambas as extremidades da rede
- Ethernet na porta de gerenciamento do gateway

## Configuração do IP do Gateway

```go
driver, codec, err := profibus.NewGatewayDriver("192.168.0.50:2000")
handler, err := profibus.ReadyHandler("192.168.0.50:2000")
```
