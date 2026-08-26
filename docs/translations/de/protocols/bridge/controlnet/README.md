# ControlNet-CmdBridge-SDK

ControlNet ist ein Echtzeit-Industrienetzwerkprotokoll von Allen-Bradley (Rockwell Automation) für schnellen, zeitkritischen Datenaustausch. Dieses Paket stellt einen ControlNet-Codec über das Kommandozeilenprogramm `1784-pcic-cli` bereit.

## CLI-Tool

Verwendet das CLI-Programm der 1784-PCIC-ControlNet-Schnittstellenkarte.

### Installation

```bash
# Install Rockwell 1784-PCIC driver and tools
# Refer to Rockwell Automation documentation for RSLinx Classic SDK
```

Verifizieren: `1784-pcic-cli --help`

## Verwendung

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

## Treiber

```go
b, c, err := controlnet.NewCmdDriver("/usr/bin/1784-pcic-cli")
```

## Unterstützte Funktionen

| Funktion | Beschreibung               |
|----------|---------------------------|
| `read`   | Von einem ControlNet-Knoten lesen |
| `write`  | An einen ControlNet-Knoten schreiben  |
| `status` | PCIC-Kartenstatus abfragen    |

## Testen

```bash
go test ./... -v
```
