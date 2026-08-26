# POWERLINK CmdBridge SDK

POWERLINK (Ethernet POWERLINK) is a real-time Ethernet protocol for industrial automation. This package provides a POWERLINK codec via the `openPOWERLINK_demo` command-line utility.

## CLI Tool

Uses the openPOWERLINK stack demo application.

### Installation

```bash
git clone https://github.com/OpenAutomationTechnologies/openPOWERLINK.git
cd openPOWERLINK
cmake -B build && cmake --build build
sudo cmake --install build
```

Verify: `openPOWERLINK_demo --help`

## Usage

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/powerlink"
)

func main() {
    p := powerlink.New()
    codec, _ := p.NewCodec("cmd")

    req := &kernel.Request{
        Function: "read",
        Address:  "0x2000",
        Count:    8,
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## Driver

```go
b, c, err := powerlink.NewCmdDriver("/usr/bin/openPOWERLINK_demo")
```

## Supported Functions

| Function | Description                   |
|----------|-------------------------------|
| `read`   | Read object dictionary entry  |
| `write`  | Write object dictionary entry |
| `status` | Query node/NMT state          |

## Testing

```bash
go test ./... -v
```
