# SDK CmdBridge ISA100.11a

ISA100.11a est une norme de réseau industriel sans fil pour l'automatisation des procédés. Ce package fournit un codec ISA100.11a via l'utilitaire en ligne de commande `yfgw410_cli` (passerelle sans fil de terrain Yokogawa YFGW410).

## Outil CLI

Utilise l'utilitaire CLI de la passerelle sans fil de terrain Yokogawa YFGW410.

### Installation

```bash
# Install Yokogawa YFGW410 gateway software and tools
# Refer to Yokogawa documentation for Field Wireless Gateway setup
```

Vérification : `yfgw410_cli --help`

## Utilisation

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

## Pilote

```go
b, c, err := isa100.NewCmdDriver("/usr/bin/yfgw410_cli")
```

## Fonctions prises en charge

| Fonction | Description                 |
|----------|-----------------------------|
| `read`   | Lire un attribut d'équipement |
| `write`  | Écrire un attribut d'équipement |
| `list`   | Lister les équipements provisionnés |

## Tests

```bash
go test ./... -v
```
