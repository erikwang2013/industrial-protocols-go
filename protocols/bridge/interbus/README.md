# Interbus GatewayBridge SDK

Interbus protocol via TCP gateway bridge.

## Hardware

| Gateway | Interface | Default IP |
|---------|-----------|-------------|
| Phoenix Contact IBS S5 DSC/I-T | Interbus master controller | 192.168.0.60 |
| Phoenix Contact IBS PCI SC/I-T | PCI Interbus master | host-assigned |
| HMS Anybus Interbus | Embedded gateway | 192.168.0.61 |

## Wiring

- 9-pin D-SUB on gateway to Interbus remote bus (incoming/outgoing)
- Shield connected to FE at both ends
- Ethernet to gateway for TCP bridge

## Gateway IP Configuration

```go
driver, codec, err := interbus.NewGatewayDriver("192.168.0.60:2006")
handler, err := interbus.ReadyHandler("192.168.0.60:2006")
```
