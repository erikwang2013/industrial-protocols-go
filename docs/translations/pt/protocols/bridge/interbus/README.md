# SDK Interbus GatewayBridge

Protocolo Interbus via ponte de gateway TCP.

## Hardware

| Gateway | Interface | IP Padrão |
|---------|-----------|-------------|
| Phoenix Contact IBS S5 DSC/I-T | Controlador mestre Interbus | 192.168.0.60 |
| Phoenix Contact IBS PCI SC/I-T | Mestre Interbus PCI | atribuído pelo host |
| HMS Anybus Interbus | Gateway embarcado | 192.168.0.61 |

## Cabeamento

- D-SUB de 9 pinos no gateway para o barramento remoto Interbus (entrada/saída)
- Blindagem conectada ao FE em ambas as extremidades
- Ethernet para o gateway para a ponte TCP

## Configuração do IP do Gateway

```go
driver, codec, err := interbus.NewGatewayDriver("192.168.0.60:2006")
handler, err := interbus.ReadyHandler("192.168.0.60:2006")
```
