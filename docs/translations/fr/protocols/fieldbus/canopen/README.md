# SDK du protocole CANopen

CANopen est un protocole de couche supérieure basé sur CAN pour les systèmes de contrôle embarqués. Ce package fournit un codec CANopen et un pilote SocketCAN.

## Vue d'ensemble du protocole

CANopen utilise des identifiants CAN standard de 11 bits avec l'ensemble de connexions prédéfini suivant :

| Fonction    | ID CAN              | Description                   |
|------------|---------------------|-------------------------------|
| NMT        | 0x000               | Gestion du réseau            |
| SYNC       | 0x080               | Message de synchronisation   |
| SDO (tx)   | 0x580 + NodeID      | Objet de données de service (serveur) |
| SDO (rx)   | 0x600 + NodeID      | Objet de données de service (client) |
| PDO1 (tx)  | 0x180 + NodeID      | Objet de données de processus 1 |
| Heartbeat  | 0x700 + NodeID      | Heartbeat / Bootup            |

## Exigences matérielles

- Système **Linux** avec prise en charge SocketCAN (`CONFIG_CAN` activé)
- Une interface CAN (p. ex. `can0`, `vcan0` pour CAN virtuel)
- Matériel transcepteur compatible CAN (p. ex. MCP2515, SJA1000, ou adaptateur USB-CAN)

### Configuration d'une interface CAN virtuelle (pour les tests)

```bash
sudo modprobe can
sudo modprobe can_raw
sudo modprobe vcan
sudo ip link add dev vcan0 type vcan
sudo ip link set up vcan0
```

## Utilisation

```go
package main

import (
    "fmt"
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/fieldbus/canopen"
)

func main() {
    p := canopen.New()
    codec, _ := p.NewCodec("can")

    // Read SDO object dictionary entry
    req := &kernel.Request{
        Function: "sdo_read",
        Metadata: map[string]any{
            "index": float64(0x1000), // Device type
            "sub":   float64(0),
        },
    }
    raw, _ := codec.Encode(req)
    fmt.Printf("SDO read frame: %X\n", raw)

    // NMT start remote node
    req2 := &kernel.Request{Function: "nmt_start"}
    raw2, _ := codec.Encode(req2)
    fmt.Printf("NMT start frame: %X\n", raw2)
}
```

## Fonctions prises en charge

| Fonction     | Description                          |
|-------------|--------------------------------------|
| `sdo_read`  | Lire une entrée du dictionnaire d'objets |
| `sdo_write` | Écrire une entrée du dictionnaire d'objets |
| `nmt_start` | Démarrer un nœud distant (NMT)       |
| `nmt_stop`  | Arrêter un nœud distant (NMT)        |
| `nmt_reset` | Réinitialiser un nœud distant (NMT)  |
| `heartbeat` | Envoyer un message heartbeat/bootup  |

## Tests

```bash
go test ./... -v
```

Note : les tests du pilote SocketCAN nécessitent un système Linux avec du matériel CAN ou une interface CAN virtuelle.
