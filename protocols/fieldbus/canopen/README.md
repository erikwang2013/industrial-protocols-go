# CANopen Protocol SDK

CANopen is a CAN-based higher-layer protocol for embedded control systems. This package provides a CANopen codec and SocketCAN driver.

## Protocol Overview

CANopen uses standard 11-bit CAN identifiers with the following predefined connection set:

| Function    | CAN ID              | Description                   |
|------------|---------------------|-------------------------------|
| NMT        | 0x000               | Network management            |
| SYNC       | 0x080               | Synchronization message       |
| SDO (tx)   | 0x580 + NodeID      | Service Data Object (server)  |
| SDO (rx)   | 0x600 + NodeID      | Service Data Object (client)  |
| PDO1 (tx)  | 0x180 + NodeID      | Process Data Object 1         |
| Heartbeat  | 0x700 + NodeID      | Heartbeat / Bootup            |

## Hardware Requirements

- **Linux** system with SocketCAN support (`CONFIG_CAN` enabled)
- A CAN interface (e.g. `can0`, `vcan0` for virtual CAN)
- CAN-capable transceiver hardware (e.g. MCP2515, SJA1000, or USB-CAN adapter)

### Setting up a virtual CAN interface (for testing)

```bash
sudo modprobe can
sudo modprobe can_raw
sudo modprobe vcan
sudo ip link add dev vcan0 type vcan
sudo ip link set up vcan0
```

## Usage

```go
package main

import (
    "fmt"
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/fieldbus/canopen"
)

func main() {
    p := canopen.New()
    codec, _ := p.NewCodec("can")

    // Read SDO object dictionary entry
    req := &kernel.Request{
        Function: "sdo_read",
        Metadata: map[string]any{
            "index": float64(0x1000), // Device type
            "sub":   float64(0),
        },
    }
    raw, _ := codec.Encode(req)
    fmt.Printf("SDO read frame: %X\n", raw)

    // NMT start remote node
    req2 := &kernel.Request{Function: "nmt_start"}
    raw2, _ := codec.Encode(req2)
    fmt.Printf("NMT start frame: %X\n", raw2)
}
```

## Supported Functions

| Function     | Description                          |
|-------------|--------------------------------------|
| `sdo_read`  | Read object dictionary entry         |
| `sdo_write` | Write object dictionary entry        |
| `nmt_start` | Start remote node (NMT)              |
| `nmt_stop`  | Stop remote node (NMT)               |
| `nmt_reset` | Reset remote node (NMT)              |
| `heartbeat` | Send heartbeat / bootup message      |

## Testing

```bash
go test ./... -v
```

Note: SocketCAN driver tests require a Linux system with CAN hardware or a virtual CAN interface.
