# POWERLINK-CmdBridge-SDK

POWERLINK (Ethernet POWERLINK) ist ein Echtzeit-Ethernet-Protokoll für die industrielle Automatisierung. Dieses Paket stellt einen POWERLINK-Codec über das Kommandozeilenprogramm `openPOWERLINK_demo` bereit.

## CLI-Tool

Verwendet die Demo-Anwendung des openPOWERLINK-Stacks.

### Installation

```bash
git clone https://github.com/OpenAutomationTechnologies/openPOWERLINK.git
cd openPOWERLINK
cmake -B build && cmake --build build
sudo cmake --install build
```

Verifizieren: `openPOWERLINK_demo --help`

## Verwendung

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

## Treiber

```go
b, c, err := powerlink.NewCmdDriver("/usr/bin/openPOWERLINK_demo")
```

## Unterstützte Funktionen

| Funktion | Beschreibung                   |
|----------|-------------------------------|
| `read`   | Objektverzeichniseintrag lesen  |
| `write`  | Objektverzeichniseintrag schreiben |
| `status` | Knoten-/NMT-Zustand abfragen          |

## Testen

```bash
go test ./... -v
```
