# SERCOS III CmdBridge SDK

SERCOS III (SErial Real-time COmmunication System) is a digital interface for motion control, using a ring topology over fiber optic or copper. This package provides a SERCOS III codec via the `netx_cli` command-line utility.

## CLI Tool

Uses the Hilscher netX SERCOS III CLI utility.

### Installation

```bash
# Install Hilscher netX driver and tools
# See https://www.hilscher.com/ for netX drivers
sudo apt-get install netx-driver
```

Verify: `netx_cli --help`

## Usage

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/sercos"
)

func main() {
    p := sercos.New()
    codec, _ := p.NewCodec("cmd")

    req := &kernel.Request{
        Function: "read",
        Address:  "S-0-51",
        Count:    2,
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## Driver

```go
b, c, err := sercos.NewCmdDriver("/usr/bin/netx_cli")
```

## Supported Functions

| Function | Description                 |
|----------|-----------------------------|
| `read`   | Read SERCOS IDN/S parameter |
| `write`  | Write SERCOS IDN/S parameter|
| `phase`  | Set communication phase     |

## Testing

```bash
go test ./... -v
```
