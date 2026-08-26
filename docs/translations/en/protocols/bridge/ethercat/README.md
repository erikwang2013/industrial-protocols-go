# EtherCAT CmdBridge SDK

EtherCAT (Ethernet for Control Automation Technology) is a high-performance industrial Ethernet fieldbus. This package provides an EtherCAT codec via the `ethercat` command-line utility.

## CLI Tool

Uses the IgH EtherCAT Master command-line tool.

### Installation

```bash
# Debian/Ubuntu
sudo apt-get install ethercat-master

# From source
git clone https://gitlab.com/etherlab.org/ethercat.git
cd ethercat
make && sudo make install
```

Verify: `ethercat slaves`

## Usage

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/ethercat"
)

func main() {
    p := ethercat.New()
    codec, _ := p.NewCodec("cmd")

    // Upload SDO from address 0x1000
    req := &kernel.Request{
        Function: "upload",
        Address:  "0x1000",
        Count:    4,
    }
    raw, _ := codec.Encode(req)
    // raw = "upload 0x1000 4\n"
    _ = raw
}
```

## Driver

```go
b, c, err := ethercat.NewCmdDriver("/usr/bin/ethercat")
```

## Supported Functions

| Function   | Description            |
|------------|------------------------|
| `upload`   | Read SDO from address  |
| `download` | Write SDO to address   |
| `slaves`   | List EtherCAT slaves   |

## Testing

```bash
go test ./... -v
```
