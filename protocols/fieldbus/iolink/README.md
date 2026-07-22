# IO-Link GatewayBridge SDK

IO-Link protocol via TCP gateway bridge.

## Hardware

| Gateway | Interface | Default IP |
|---------|-----------|-------------|
| ifm AL1332 | IO-Link Master (EtherNet/IP) | 192.168.0.40 |
| Balluff BNI00AZ | IO-Link Master (PROFINET) | 192.168.0.41 |
| SICK SIG200 | IO-Link Master (Ethernet) | 192.168.0.42 |

## Wiring

- M12 connector (4-pin) for each IO-Link port: L+ (brown), L- (blue), C/Q (black), unused (white)
- 24 VDC power supply for master and devices
- Ethernet on master for TCP bridge

## Gateway IP Configuration

```go
driver, codec, err := iolink.NewGatewayDriver("192.168.0.40:2004")
handler, err := iolink.ReadyHandler("192.168.0.40:2004")
```
