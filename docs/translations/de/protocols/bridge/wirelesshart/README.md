# WirelessHART-CmdBridge-SDK

WirelessHART (IEC 62591) ist ein drahtloser Industriestandard auf Basis des HART-Protokolls. Dieses Paket stellt einen WirelessHART-Codec über das Kommandozeilenprogramm `emerson_1410_cli` (Emerson-1410/1420-Wireless-Gateway) bereit.

## CLI-Tool

Verwendet das CLI-Programm des Emerson-1410/1420-Wireless-Gateways.

### Installation

```bash
# Install Emerson Wireless Gateway software and tools
# Refer to Emerson documentation for 1410/1420 Gateway setup
```

Verifizieren: `emerson_1410_cli --help`

## Verwendung

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

## Treiber

```go
b, c, err := wirelesshart.NewCmdDriver("/usr/bin/emerson_1410_cli")
```

## Unterstützte Funktionen

| Funktion | Beschreibung                |
|----------|----------------------------|
| `read`   | Geräteparameter lesen      |
| `write`  | Geräteparameter schreiben     |
| `scan`   | Nach drahtlosen Geräten suchen  |

## Testen

```bash
go test ./... -v
```
