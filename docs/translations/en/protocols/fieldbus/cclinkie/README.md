# CC-Link IE Field GatewayBridge SDK

CC-Link IE Field protocol via TCP gateway bridge.

## Hardware

| Gateway | Interface | Default IP |
|---------|-----------|-------------|
| Mitsubishi MELSEC IQ-R | CC-Link IE Field master | 192.168.3.40 |
| Mitsubishi RJ71GF11-T2 | CC-Link IE Field module | host-assigned |
| HMS Anybus CC-Link IE | Embedded gateway | 192.168.0.52 |

## Wiring

- RJ45 Ethernet for CC-Link IE Field (1 Gbps ring or star topology)
- Management port on separate network
- Gateway connects field network to TCP bridge

## Gateway IP Configuration

```go
driver, codec, err := cclinkie.NewGatewayDriver("192.168.3.40:2005")
handler, err := cclinkie.ReadyHandler("192.168.3.40:2005")
```
