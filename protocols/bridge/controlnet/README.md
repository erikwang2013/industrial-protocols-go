# ControlNet CmdBridge SDK

ControlNet is a real-time industrial network protocol developed by Allen-Bradley (Rockwell Automation) for high-speed, time-critical data exchange. This package provides a ControlNet codec via the `1784-pcic-cli` command-line utility.

## CLI Tool

Uses the 1784-PCIC ControlNet interface card CLI utility.

### Installation

```bash
# Install Rockwell 1784-PCIC driver and tools
# Refer to Rockwell Automation documentation for RSLinx Classic SDK
```

Verify: `1784-pcic-cli --help`

## Usage

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/controlnet"
)

func main() {
    p := controlnet.New()
    codec, _ := p.NewCodec("cmd")

    req := &kernel.Request{
        Function: "read",
        Address:  "0x10",
        Count:    8,
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## Driver

```go
b, c, err := controlnet.NewCmdDriver("/usr/bin/1784-pcic-cli")
```

## Supported Functions

| Function | Description               |
|----------|---------------------------|
| `read`   | Read from ControlNet node |
| `write`  | Write to ControlNet node  |
| `status` | Query PCIC card status    |

## Testing

```bash
go test ./... -v
```
