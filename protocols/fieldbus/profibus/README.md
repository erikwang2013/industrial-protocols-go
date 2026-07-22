# PROFIBUS GatewayBridge SDK

PROFIBUS DP/PA protocol via TCP gateway bridge.

## Hardware

| Gateway | Interface | Default IP |
|---------|-----------|-------------|
| Anybus Communicator | PROFIBUS DP-V1 slave | 192.168.0.50 |
| Siemens CP 5611 proxy | PCI/PCIe PROFIBUS master | host-assigned |
| HMS Fieldbus Gateway | Anybus NP40 | 192.168.0.51 |

## Wiring

- DB9 female on gateway to PROFIBUS network (A-line green, B-line red)
- Termination resistor ON at both network ends
- Ethernet to gateway management port

## Gateway IP Configuration

```go
driver, codec, err := profibus.NewGatewayDriver("192.168.0.50:2000")
handler, err := profibus.ReadyHandler("192.168.0.50:2000")
```
