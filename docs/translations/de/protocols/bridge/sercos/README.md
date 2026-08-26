# SERCOS-III-CmdBridge-SDK

SERCOS III (SErial Real-time COmmunication System) ist eine digitale Schnittstelle für die Bewegungssteuerung mit Ringtopologie über Glasfaser oder Kupfer. Dieses Paket stellt einen SERCOS-III-Codec über das Kommandozeilenprogramm `netx_cli` bereit.

## CLI-Tool

Verwendet das Hilscher-netX-SERCOS-III-CLI-Programm.

### Installation

```bash
# Install Hilscher netX driver and tools
# See https://www.hilscher.com/ for netX drivers
sudo apt-get install netx-driver
```

Verifizieren: `netx_cli --help`

## Verwendung

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

## Treiber

```go
b, c, err := sercos.NewCmdDriver("/usr/bin/netx_cli")
```

## Unterstützte Funktionen

| Funktion | Beschreibung                 |
|----------|-----------------------------|
| `read`   | SERCOS-IDN/S-Parameter lesen |
| `write`  | SERCOS-IDN/S-Parameter schreiben|
| `phase`  | Kommunikationsphase setzen     |

## Testen

```bash
go test ./... -v
```
