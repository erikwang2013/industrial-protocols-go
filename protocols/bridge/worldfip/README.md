# WorldFIP GatewayBridge SDK

WorldFIP protocol via TCP gateway bridge.

## Hardware

| Gateway | Interface | Default IP |
|---------|-----------|-------------|
| FIPIO Agent | WorldFIP fieldbus agent | 192.168.0.70 |
| FIP Gateway (Alstom) | WorldFIP to Ethernet | 192.168.0.71 |
| NI FIP-USB | USB WorldFIP interface | host-assigned |

## Wiring

- 9-pin D-SUB on gateway to WorldFIP trunk (FIP1 = Data+, FIP2 = Data-)
- Line terminator (120 ohm) at both bus ends
- Ethernet for TCP bridge

## Gateway IP Configuration

```go
driver, codec, err := worldfip.NewGatewayDriver("192.168.0.70:2007")
handler, err := worldfip.ReadyHandler("192.168.0.70:2007")
```
