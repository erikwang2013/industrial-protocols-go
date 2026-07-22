# LonWorks GatewayBridge SDK

LonWorks (ANSI/CEA-709.1) protocol via TCP gateway bridge.

## Hardware

| Gateway | Interface | Default IP |
|---------|-----------|-------------|
| Echelon U60 | USB FT-10 network interface | host-assigned |
| Echelon U70 | USB TP/XF-1250 interface | host-assigned |
| Loytec L-IP | LonWorks/IP router | 192.168.0.90 |

## Wiring

- FT-10 (Free Topology): polarity-insensitive twisted pair, up to 500 m free topology
- TP/XF-1250: bus topology with 105 ohm terminator
- Ethernet on L-IP router for TCP bridge

## Gateway IP Configuration

```go
driver, codec, err := lonworks.NewGatewayDriver("192.168.0.90:2009")
handler, err := lonworks.ReadyHandler("192.168.0.90:2009")
```
