# ISA100.11a CmdBridge SDK

ISA100.11a is a wireless industrial networking standard for process automation. This package provides an ISA100.11a codec via the `yfgw410_cli` (Yokogawa YFGW410 field wireless gateway) command-line utility.

## CLI Tool

Uses the Yokogawa YFGW410 Field Wireless Gateway CLI.

### Installation

```bash
# Install Yokogawa YFGW410 gateway software and tools
# Refer to Yokogawa documentation for Field Wireless Gateway setup
```

Verify: `yfgw410_cli --help`

## Usage

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/isa100"
)

func main() {
    p := isa100.New()
    codec, _ := p.NewCodec("cmd")

    req := &kernel.Request{
        Function: "read",
        Address:  "DEV001",
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## Driver

```go
b, c, err := isa100.NewCmdDriver("/usr/bin/yfgw410_cli")
```

## Supported Functions

| Function | Description                 |
|----------|-----------------------------|
| `read`   | Read device attribute       |
| `write`  | Write device attribute      |
| `list`   | List provisioned devices    |

## Testing

```bash
go test ./... -v
```
