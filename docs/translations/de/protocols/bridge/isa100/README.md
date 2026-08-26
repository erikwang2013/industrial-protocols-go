# ISA100.11a-CmdBridge-SDK

ISA100.11a ist ein drahtloser Industriestandard für die Prozessautomatisierung. Dieses Paket stellt einen ISA100.11a-Codec über das Kommandozeilenprogramm `yfgw410_cli` (Yokogawa-YFGW410-Field-Wireless-Gateway) bereit.

## CLI-Tool

Verwendet das CLI des Yokogawa-YFGW410-Field-Wireless-Gateways.

### Installation

```bash
# Install Yokogawa YFGW410 gateway software and tools
# Refer to Yokogawa documentation for Field Wireless Gateway setup
```

Verifizieren: `yfgw410_cli --help`

## Verwendung

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

## Treiber

```go
b, c, err := isa100.NewCmdDriver("/usr/bin/yfgw410_cli")
```

## Unterstützte Funktionen

| Funktion | Beschreibung                 |
|----------|-----------------------------|
| `read`   | Geräteattribut lesen       |
| `write`  | Geräteattribut schreiben      |
| `list`   | Bereitgestellte Geräte auflisten    |

## Testen

```bash
go test ./... -v
```
