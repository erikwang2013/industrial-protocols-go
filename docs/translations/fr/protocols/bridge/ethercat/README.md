# SDK CmdBridge EtherCAT

EtherCAT (Ethernet for Control Automation Technology) est un bus de terrain Ethernet industriel hautes performances. Ce package fournit un codec EtherCAT via l'utilitaire en ligne de commande `ethercat`.

## Outil CLI

Utilise l'outil en ligne de commande IgH EtherCAT Master.

### Installation

```bash
# Debian/Ubuntu
sudo apt-get install ethercat-master

# From source
git clone https://gitlab.com/etherlab.org/ethercat.git
cd ethercat
make && sudo make install
```

Vérification : `ethercat slaves`

## Utilisation

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/ethercat"
)

func main() {
    p := ethercat.New()
    codec, _ := p.NewCodec("cmd")

    // Upload SDO from address 0x1000
    req := &kernel.Request{
        Function: "upload",
        Address:  "0x1000",
        Count:    4,
    }
    raw, _ := codec.Encode(req)
    // raw = "upload 0x1000 4\n"
    _ = raw
}
```

## Pilote

```go
b, c, err := ethercat.NewCmdDriver("/usr/bin/ethercat")
```

## Fonctions prises en charge

| Fonction   | Description            |
|------------|------------------------|
| `upload`   | Lire un SDO à une adresse |
| `download` | Écrire un SDO à une adresse |
| `slaves`   | Lister les esclaves EtherCAT |

## Tests

```bash
go test ./... -v
```
