# SDK CmdBridge ControlNet

ControlNet est un protocole de réseau industriel temps réel développé par Allen-Bradley (Rockwell Automation) pour l'échange de données à haut débit et à contrainte de temps. Ce package fournit un codec ControlNet via l'utilitaire en ligne de commande `1784-pcic-cli`.

## Outil CLI

Utilise l'utilitaire CLI de la carte d'interface ControlNet 1784-PCIC.

### Installation

```bash
# Install Rockwell 1784-PCIC driver and tools
# Refer to Rockwell Automation documentation for RSLinx Classic SDK
```

Vérification : `1784-pcic-cli --help`

## Utilisation

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

## Pilote

```go
b, c, err := controlnet.NewCmdDriver("/usr/bin/1784-pcic-cli")
```

## Fonctions prises en charge

| Fonction | Description               |
|----------|---------------------------|
| `read`   | Lire depuis un nœud ControlNet |
| `write`  | Écrire vers un nœud ControlNet |
| `status` | Interroger l'état de la carte PCIC |

## Tests

```bash
go test ./... -v
```
