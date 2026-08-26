# SERCOS-I/II-CmdBridge-SDK

SERCOS I/II ist die ältere serielle Glasfaserversion der SERCOS-Schnittstelle für digitale Bewegungssteuerung. Dieses Paket stellt einen SERCOS-I/II-Codec über das Kommandozeilenprogramm `sercos_cli` bereit.

## CLI-Tool

Verwendet ein CLI-Programm für eine SERCOS-Glasfaserschnittstelle.

### Installation

```bash
# SERCOS interface card driver and tools
# Refer to vendor documentation for the specific SERCOS master card
```

Verifizieren: `sercos_cli --help`

## Verwendung

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

## Treiber

```go
b, c, err := sercos1.NewCmdDriver("/usr/bin/sercos_cli")
```

## Unterstützte Funktionen

| Funktion | Beschreibung              |
|----------|--------------------------|
| `read`   | SERCOS-IDN lesen          |
| `write`  | SERCOS-IDN schreiben         |
| `status` | Antriebsstatus abfragen       |

## Testen

```bash
go test ./... -v
```
