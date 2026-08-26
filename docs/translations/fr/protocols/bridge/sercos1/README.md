# SDK CmdBridge SERCOS I/II

SERCOS I/II est la version série à fibre optique héritée de l'interface SERCOS pour la commande de mouvement numérique. Ce package fournit un codec SERCOS I/II via l'utilitaire en ligne de commande `sercos_cli`.

## Outil CLI

Utilise un utilitaire CLI d'interface à fibre optique SERCOS.

### Installation

```bash
# SERCOS interface card driver and tools
# Refer to vendor documentation for the specific SERCOS master card
```

Vérification : `sercos_cli --help`

## Utilisation

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

## Pilote

```go
b, c, err := sercos1.NewCmdDriver("/usr/bin/sercos_cli")
```

## Fonctions prises en charge

| Fonction | Description              |
|----------|--------------------------|
| `read`   | Lire un IDN SERCOS       |
| `write`  | Écrire un IDN SERCOS     |
| `status` | Interroger l'état du variateur |

## Tests

```bash
go test ./... -v
```
