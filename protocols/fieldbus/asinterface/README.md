# AS-Interface GatewayBridge SDK

AS-Interface (ASi) protocol via TCP gateway bridge.

## Hardware

| Gateway | Interface | Default IP |
|---------|-----------|-------------|
| Bihl+Wiedemann BWU2044 | ASi Master (Ethernet) | 192.168.0.30 |
| P+F VBG-ENX-K20-DMD-EV | ASi Gateway | 192.168.0.31 |
| ifm AC1375 | ASi ControllerE | 192.168.0.32 |

## Wiring

- Yellow ASi cable (power + data) from gateway to slaves
- Black auxiliary power cable (24 VDC for actuators) optional
- Ethernet on gateway for TCP bridge

## Gateway IP Configuration

```go
driver, codec, err := asinterface.NewGatewayDriver("192.168.0.30:2003")
handler, err := asinterface.ReadyHandler("192.168.0.30:2003")
```
