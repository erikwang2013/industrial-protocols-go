# MOST Serial Protocol SDK

MOST (Media Oriented Systems Transport) is a high-speed multimedia network technology used primarily in automotive infotainment systems. This package provides a MOST codec over a serial adapter using AT-command interface.

## Protocol Overview

MOST uses synchronous serial communication over fiber optic physical layer. This implementation connects via a serial adapter that exposes an AT-command interface at 115200 baud.

### AT Commands

| Command            | Description                  |
|--------------------|------------------------------|
| `AT+READ=<addr>,<len>`  | Read bytes from address |
| `AT+WRITE=<addr>,<hex>` | Write hex bytes to address |
| `AT+STATUS`             | Query ring/network status |

### Response Format

- `+OK:<hex_data>` -- successful response
- `+ERR:<code>` -- error response

## Hardware Requirements

- MOST fiber optic network with proper termination
- MOST-to-serial adapter (e.g. MOST150 USB adapter)
- Serial port at 115200 baud, 8N1

## Usage

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/automotive/most"
)

func main() {
    p := most.New()
    codec, _ := p.NewCodec("serial")

    // Read 4 bytes from address 0x0100
    req := &kernel.Request{
        Function: "read",
        Address:  "0x0100",
        Count:    4,
    }
    raw, _ := codec.Encode(req)
    // raw = "AT+READ=0x0100,4\r\n"
    _ = raw

    // Write data to address 0x0200
    req2 := &kernel.Request{
        Function: "write",
        Address:  "0x0200",
        Data:     []byte{0xDE, 0xAD, 0xBE, 0xEF},
    }
    raw2, _ := codec.Encode(req2)
    // raw2 = "AT+WRITE=0x0200,DEADBEEF\r\n"
    _ = raw2
}
```

## Driver

```go
b, c, err := most.NewSerialDriver("/dev/ttyUSB0")
```

## Supported Functions

| Function | Description                 |
|----------|-----------------------------|
| `read`   | Read from a MOST address    |
| `write`  | Write data to a MOST address |
| `status` | Query ring/network status   |

## Testing

```bash
go test ./... -v
```
