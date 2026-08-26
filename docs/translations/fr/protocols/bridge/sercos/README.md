# SDK CmdBridge SERCOS III

SERCOS III (SErial Real-time COmmunication System) est une interface numérique pour la commande de mouvement, utilisant une topologie en anneau sur fibre optique ou cuivre. Ce package fournit un codec SERCOS III via l'utilitaire en ligne de commande `netx_cli`.

## Outil CLI

Utilise l'utilitaire CLI Hilscher netX SERCOS III.

### Installation

```bash
# Install Hilscher netX driver and tools
# See https://www.hilscher.com/ for netX drivers
sudo apt-get install netx-driver
```

Vérification : `netx_cli --help`

## Utilisation

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

## Pilote

```go
b, c, err := sercos.NewCmdDriver("/usr/bin/netx_cli")
```

## Fonctions prises en charge

| Fonction | Description                 |
|----------|-----------------------------|
| `read`   | Lire un paramètre IDN/S SERCOS |
| `write`  | Écrire un paramètre IDN/S SERCOS |
| `phase`  | Régler la phase de communication |

## Tests

```bash
go test ./... -v
```
