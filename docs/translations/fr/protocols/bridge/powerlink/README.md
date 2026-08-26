# SDK CmdBridge POWERLINK

POWERLINK (Ethernet POWERLINK) est un protocole Ethernet temps réel pour l'automatisation industrielle. Ce package fournit un codec POWERLINK via l'utilitaire en ligne de commande `openPOWERLINK_demo`.

## Outil CLI

Utilise l'application de démonstration de la pile openPOWERLINK.

### Installation

```bash
git clone https://github.com/OpenAutomationTechnologies/openPOWERLINK.git
cd openPOWERLINK
cmake -B build && cmake --build build
sudo cmake --install build
```

Vérification : `openPOWERLINK_demo --help`

## Utilisation

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

## Pilote

```go
b, c, err := powerlink.NewCmdDriver("/usr/bin/openPOWERLINK_demo")
```

## Fonctions prises en charge

| Fonction | Description                   |
|----------|-------------------------------|
| `read`   | Lire une entrée du dictionnaire d'objets |
| `write`  | Écrire une entrée du dictionnaire d'objets |
| `status` | Interroger l'état du nœud/NMT |

## Tests

```bash
go test ./... -v
```
