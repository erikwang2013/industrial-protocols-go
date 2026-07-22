# Modbus Plus GatewayBridge SDK

Modbus Plus (MB+) protocol via TCP gateway bridge.

## Hardware

| Gateway | Interface | Default IP |
|---------|-----------|-------------|
| Schneider SA85 | ISA Modbus Plus adapter | host-assigned |
| Schneider BM85 | Modbus Plus bridge/multiplexer | 192.168.0.A0 |
| ProSoft MVI56-MBP | ControlLogix MB+ module | host-assigned |

## Wiring

- Twinaxial cable (RG-62) with BNC connectors for MB+ trunk
- Terminating resistor (78 ohm at each end)
- BM85 bridge connects MB+ to TCP via Ethernet

## Gateway IP Configuration

```go
driver, codec, err := modbusplus.NewGatewayDriver("192.168.0.160:2010")
handler, err := modbusplus.ReadyHandler("192.168.0.160:2010")
```
