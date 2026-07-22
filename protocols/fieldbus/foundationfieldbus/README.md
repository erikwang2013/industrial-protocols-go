# Foundation Fieldbus GatewayBridge SDK

Foundation Fieldbus H1/HSE via TCP gateway bridge.

## Hardware

| Gateway | Interface | Default IP |
|---------|-----------|-------------|
| NI USB-8486 | USB H1 interface | host-assigned |
| Softing FFusb | USB H1 interface | host-assigned |
| P+F HD2-GTR-4PA | H1 to Ethernet gateway | 192.168.0.20 |

## Wiring

- H1 trunk (twisted-pair, shielded) with terminator at both ends
- Fieldbus power conditioner for H1 (24 VDC, 350-500 mA per segment)
- Gateway connects H1 segment to TCP bridge via Ethernet

## Gateway IP Configuration

```go
driver, codec, err := foundationfieldbus.NewGatewayDriver("192.168.0.20:2002")
handler, err := foundationfieldbus.ReadyHandler("192.168.0.20:2002")
```
