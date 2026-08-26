# Lightbus GatewayBridge SDK

Beckhoff Lightbus fiber optic protocol via TCP gateway bridge.

## Hardware

| Gateway | Interface | Default IP |
|---------|-----------|-------------|
| Beckhoff FC2001 | Lightbus PCI card | host-assigned |
| Beckhoff BK2000 | Lightbus bus coupler | 192.168.0.80 |
| Beckhoff FC9001 | Lightbus Ethernet adapter | 192.168.0.81 |

## Wiring

- Plastic optical fiber (POF) ring topology
- FC2001/FC9001 connects ring to TCP bridge via Ethernet
- Each device has TX and RX fiber connectors

## Gateway IP Configuration

```go
driver, codec, err := lightbus.NewGatewayDriver("192.168.0.80:2008")
handler, err := lightbus.ReadyHandler("192.168.0.80:2008")
```
