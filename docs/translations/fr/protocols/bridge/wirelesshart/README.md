# SDK CmdBridge WirelessHART

WirelessHART (IEC 62591) est une norme de réseau industriel sans fil basée sur le protocole HART. Ce package fournit un codec WirelessHART via l'utilitaire en ligne de commande `emerson_1410_cli` (passerelle sans fil Emerson 1410/1420).

## Outil CLI

Utilise l'utilitaire CLI de la passerelle sans fil Emerson 1410/1420.

### Installation

```bash
# Install Emerson Wireless Gateway software and tools
# Refer to Emerson documentation for 1410/1420 Gateway setup
```

Vérification : `emerson_1410_cli --help`

## Utilisation

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

## Pilote

```go
b, c, err := wirelesshart.NewCmdDriver("/usr/bin/emerson_1410_cli")
```

## Fonctions prises en charge

| Fonction | Description                |
|----------|----------------------------|
| `read`   | Lire un paramètre d'équipement |
| `write`  | Écrire un paramètre d'équipement |
| `scan`   | Rechercher les équipements sans fil |

## Tests

```bash
go test ./... -v
```
