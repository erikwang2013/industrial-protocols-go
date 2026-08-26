# EtherCAT-CmdBridge-SDK

EtherCAT (Ethernet for Control Automation Technology) ist ein leistungsfähiger industrieller Ethernet-Feldbus. Dieses Paket stellt einen EtherCAT-Codec über das Kommandozeilenprogramm `ethercat` bereit.

## CLI-Tool

Verwendet das Kommandozeilentool des IgH-EtherCAT-Masters.

### Installation

```bash
# Debian/Ubuntu
sudo apt-get install ethercat-master

# From source
git clone https://gitlab.com/etherlab.org/ethercat.git
cd ethercat
make && sudo make install
```

Verifizieren: `ethercat slaves`

## Verwendung

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

## Treiber

```go
b, c, err := ethercat.NewCmdDriver("/usr/bin/ethercat")
```

## Unterstützte Funktionen

| Funktion   | Beschreibung            |
|------------|------------------------|
| `upload`   | SDO von einer Adresse lesen  |
| `download` | SDO an eine Adresse schreiben   |
| `slaves`   | EtherCAT-Slaves auflisten   |

## Testen

```bash
go test ./... -v
```
