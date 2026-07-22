# SAE J1850 CAN Protocol SDK

SAE J1850 is a vehicle communication standard used for on-board diagnostics (OBD-II). This package provides a J1850 codec over CAN bus (ISO 15765-4 / CAN TP).

## Protocol Overview

SAE J1850 OBD-II over CAN uses 29-bit extended CAN identifiers with the following structure:

### CAN ID Format (29-bit)

| Bits       | Field    | Description                        |
|-----------|----------|------------------------------------|
| 28-26     | Priority | Message priority (0-7, default 6)  |
| 25        | Ext ID   | Always 1 for extended frames       |
| 24-16     | PF       | Parameter Format (Header)          |
| 15-8      | PS       | Parameter Specific (target/source) |
| 7-0       | SA       | Source Address                     |

### Standard OBD-II CAN IDs

| Type                | CAN ID (hex)    | Description                    |
|--------------------|-----------------|--------------------------------|
| Physical Request   | 0x18DAxxF1      | Request to specific ECU (xx=addr) |
| Physical Response  | 0x18DAF1xx      | Response from ECU (xx=addr)   |
| Functional Request | 0x18DB33F1      | Broadcast to all ECUs         |

### ISO 15765-2 Frame Format

Single frame: byte 0 upper nibble = data length (0-7), lower nibble + remaining bytes = diagnostic data.

## Hardware Requirements

- **Linux** system with SocketCAN support
- J1850-CAN OBD-II adapter (e.g. ELM327-compatible USB-to-CAN, OBDLink SX, Macchina M2)
- Vehicle with OBD-II connector (1996+ most vehicles)

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
    "github.com/erikwang2013/industrial-protocols-go/protocols/automotive/saej1850"
)

func main() {
    p := saej1850.New()
    codec, _ := p.NewCodec("can")

    // Mode $01 PID $0C: Engine RPM
    req := &kernel.Request{
        Function: "mode01",
        Metadata: map[string]any{"pid": float64(0x0C)},
    }
    raw, _ := codec.Encode(req)
    _ = raw

    // Mode $03: Request emission-related DTCs
    req2 := &kernel.Request{Function: "mode03"}
    raw2, _ := codec.Encode(req2)
    _ = raw2
}
```

## Supported Functions

| Function        | OBD-II Mode | Description                         |
|----------------|-------------|-------------------------------------|
| `mode01`       | $01         | Request current powertrain data     |
| `mode03`       | $03         | Request emission-related DTCs       |
| `mode0A`       | $0A         | Request permanent DTCs              |
| `diag_request` | Custom      | Generic diagnostic request          |
| `diag_response`| -           | Diagnostic response frame           |
| `broadcast`    | -           | Functional broadcast to all ECUs    |

## Testing

```bash
go test ./... -v
```
