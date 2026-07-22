# SERCOS I/II CmdBridge SDK

SERCOS I/II is the legacy serial fiber optic version of the SERCOS interface for digital motion control. This package provides a SERCOS I/II codec via the `sercos_cli` command-line utility.

## CLI Tool

Uses a SERCOS fiber optic interface CLI utility.

### Installation

```bash
# SERCOS interface card driver and tools
# Refer to vendor documentation for the specific SERCOS master card
```

Verify: `sercos_cli --help`

## Usage

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/sercos1"
)

func main() {
    p := sercos1.New()
    codec, _ := p.NewCodec("cmd")

    req := &kernel.Request{
        Function: "read",
        Address:  "0x0100",
        Count:    4,
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## Driver

```go
b, c, err := sercos1.NewCmdDriver("/usr/bin/sercos_cli")
```

## Supported Functions

| Function | Description              |
|----------|--------------------------|
| `read`   | Read SERCOS IDN          |
| `write`  | Write SERCOS IDN         |
| `status` | Query drive status       |

## Testing

```bash
go test ./... -v
```
