# SDK du protocole DeviceNet

DeviceNet est un protocole de réseau industriel basé sur CAN pour l'automatisation d'usine. Ce package fournit un codec DeviceNet avec des pilotes SocketCAN et passerelle TCP.

## Vue d'ensemble du protocole

DeviceNet utilise le protocole industriel commun (CIP) sur CAN. L'ensemble de connexions Maître/Esclave prédéfini utilise les messages du groupe 2 (CAN ID 0x400 + NodeID) pour les E/S interrogées (polled) et la messagerie explicite.

## Exigences matérielles

### Mode CAN (SocketCAN)
- Système **Linux** avec prise en charge SocketCAN (`CONFIG_CAN` activé)
- Une interface CAN (p. ex. `can0`, `vcan0`)
- Matériel CAN compatible DeviceNet (p. ex. Anybus Communicator, HMS IXXAT)

### Mode passerelle
- Connectivité TCP/IP vers la passerelle DeviceNet
- Passerelle prenant en charge un protocole de commande en texte brut (p. ex. HMS Anybus, Hilscher netX)

### Configuration d'une interface CAN virtuelle (pour les tests)

```bash
sudo modprobe can
sudo modprobe can_raw
sudo modprobe vcan
sudo ip link add dev vcan0 type vcan
sudo ip link set up vcan0
```

## Utilisation

### Mode CAN

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/fieldbus/devicenet"
)

func main() {
    p := devicenet.New()
    codec, _ := p.NewCodec("can")

    // Poll request
    req := &kernel.Request{
        Function: "poll",
        Data:     []byte{0x01, 0x00},
    }
    raw, _ := codec.Encode(req)

    // Open connection
    req2 := &kernel.Request{Function: "open"}
    raw2, _ := codec.Encode(req2)
    _ = raw
    _ = raw2
}
```

### Mode passerelle

```go
p := devicenet.New()
codec, _ := p.NewCodec("gateway")

req := &kernel.Request{
    Function: "poll",
    Data:     []byte{0xAB, 0xCD},
}
raw, _ := codec.Encode(req)
// raw will be: "poll abcd\n"
```

## Fonctions prises en charge

| Fonction | Description                    |
|----------|--------------------------------|
| `open`   | Ouvrir une connexion explicite |
| `poll`   | Interroger les données E/S (groupe 2) |

## Tests

```bash
go test ./... -v
```

Note : les tests du pilote SocketCAN nécessitent Linux. Les tests du pilote de passerelle fonctionnent sur tous les systèmes d'exploitation.
