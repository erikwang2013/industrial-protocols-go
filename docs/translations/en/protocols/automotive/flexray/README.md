# FlexRay CAN Protocol SDK

FlexRay is a high-speed deterministic automotive communication protocol. This package provides a FlexRay codec over CAN bus, supporting cycle-based framing with CRC-16/XMODEM integrity checks.

## Protocol Overview

FlexRay uses a time-division multiple access (TDMA) scheme with repeated communication cycles. Each cycle consists of static and dynamic segments. This codec maps FlexRay frames onto extended 29-bit CAN frames.

### Wire Format

FlexRay cycle payload:
- **Header** (2 bytes): cycle number (little-endian)
- **Status** (1 byte): bit 7=PPI (Payload Preamble Indicator), bit 6=NFI, bit 5=SYF, bit 4=SUF
- **Data** (N bytes): payload (max 254 bytes)
- **CRC** (2 bytes): CRC-16/XMODEM over header+status+data (little-endian)

CAN ID encoding (29-bit extended):
- Bits 28-24: Message type (0x01=frame, 0x02=status)
- Bits 23-10: Reserved
- Bits 15-10: Slot ID (6 bits)
- Bits 9-0: Cycle number (10 bits)

## Hardware Requirements

- **Linux** system with SocketCAN support
- Vector VN7600/VN7640 or Bosch FlexRay-CAN adapter
- FlexRay network with proper termination (2.5V bias)

### Setting up CAN interface

```bash
sudo modprobe can
sudo modprobe can_raw
sudo ip link set can0 type can bitrate 500000
sudo ip link set up can0
```

## Usage

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/automotive/flexray"
)

func main() {
    p := flexray.New()
    codec, _ := p.NewCodec("can")

    // Send a FlexRay frame in cycle 5 with PPI set
    req := &kernel.Request{
        Function: "frame",
        Data:     []byte{0x42, 0x01},
        Metadata: map[string]any{
            "cycle": float64(5),
            "ppi":   true,
        },
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## Supported Functions

| Function | Description                              |
|----------|------------------------------------------|
| `frame`  | Send FlexRay frame payload               |
| `status` | Query slot configuration (slot, cycle)   |

## Testing

```bash
go test ./... -v
```
