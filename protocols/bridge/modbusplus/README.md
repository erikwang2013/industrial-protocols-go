# Modbus Plus Protocol SDK

Modbus Plus (MB+) is a high-speed token-passing industrial network developed by Modicon (Schneider Electric).

## Variants

| Variant | Bridge Type | Description |
|---------|-------------|-------------|
| `gateway` | GatewayBridge | TCP connection to SA85/BM85 adapter |
| `cmd` | CmdBridge | CLI wrapper for `sa85_cli` utility |

## Gateway Driver

```go
driver, codec, err := modbusplus.NewGatewayDriver("192.168.0.160:2010")
handler, err := modbusplus.ReadyHandler("192.168.0.160:2010")
```

## Cmd Driver

```go
driver, codec, err := modbusplus.NewCmdDriver("/usr/bin/sa85_cli")
```

### CLI Tool Installation

```bash
# Install SA85 Modbus Plus driver and tools
# Refer to Schneider Electric documentation for SA85 adapter
```

Verify: `sa85_cli --help`

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

## Protocol Frame Format

- Magic (2 bytes): `0x4D42`
- Destination (1 byte): Node address
- Command (1 byte): 0x01=read, 0x02=write
- Length (2 bytes): Payload size (big-endian)
- Payload (N bytes): Data
- CRC (2 bytes): Modbus CRC-16 (little-endian)

## Testing

```bash
go test ./... -v
```
