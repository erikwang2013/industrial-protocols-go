# WirelessHART CmdBridge SDK

WirelessHART (IEC 62591) is a wireless industrial networking standard based on the HART protocol. This package provides a WirelessHART codec via the `emerson_1410_cli` (Emerson 1410/1420 Wireless Gateway) command-line utility.

## CLI Tool

Uses the Emerson 1410/1420 Wireless Gateway CLI utility.

### Installation

```bash
# Install Emerson Wireless Gateway software and tools
# Refer to Emerson documentation for 1410/1420 Gateway setup
```

Verify: `emerson_1410_cli --help`

## Usage

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/wirelesshart"
)

func main() {
    p := wirelesshart.New()
    codec, _ := p.NewCodec("cmd")

    req := &kernel.Request{
        Function: "read",
        Address:  "TT101",
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## Driver

```go
b, c, err := wirelesshart.NewCmdDriver("/usr/bin/emerson_1410_cli")
```

## Supported Functions

| Function | Description                |
|----------|----------------------------|
| `read`   | Read device parameter      |
| `write`  | Write device parameter     |
| `scan`   | Scan for wireless devices  |

## Testing

```bash
go test ./... -v
```
