# SDK du protocole FlexRay CAN

FlexRay est un protocole de communication automobile déterministe à haut débit. Ce package fournit un codec FlexRay sur bus CAN, avec un tramage basé sur les cycles et des contrôles d'intégrité CRC-16/XMODEM.

## Vue d'ensemble du protocole

FlexRay utilise un schéma d'accès multiple à répartition dans le temps (TDMA) avec des cycles de communication répétés. Chaque cycle se compose de segments statiques et dynamiques. Ce codec mappe les trames FlexRay sur des trames CAN étendues de 29 bits.

### Format sur le fil

Charge utile d'un cycle FlexRay :
- **En-tête** (2 octets) : numéro de cycle (little-endian)
- **Statut** (1 octet) : bit 7=PPI (Payload Preamble Indicator), bit 6=NFI, bit 5=SYF, bit 4=SUF
- **Données** (N octets) : charge utile (max 254 octets)
- **CRC** (2 octets) : CRC-16/XMODEM sur en-tête+statut+données (little-endian)

Encodage de l'ID CAN (étendu 29 bits) :
- Bits 28-24 : type de message (0x01=trame, 0x02=statut)
- Bits 23-10 : réservés
- Bits 15-10 : ID de slot (6 bits)
- Bits 9-0 : numéro de cycle (10 bits)

## Exigences matérielles

- Système **Linux** avec prise en charge SocketCAN
- Adaptateur Vector VN7600/VN7640 ou Bosch FlexRay-CAN
- Réseau FlexRay avec terminaison correcte (polarisation 2,5 V)

### Configuration de l'interface CAN

```bash
sudo modprobe can
sudo modprobe can_raw
sudo ip link set can0 type can bitrate 500000
sudo ip link set up can0
```

## Utilisation

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/automotive/flexray"
)

func main() {
    p := flexray.New()
    codec, _ := p.NewCodec("can")

    // Send a FlexRay frame in cycle 5 with PPI set
    req := &kernel.Request{
        Function: "frame",
        Data:     []byte{0x42, 0x01},
        Metadata: map[string]any{
            "cycle": float64(5),
            "ppi":   true,
        },
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## Fonctions prises en charge

| Fonction | Description |
|----------|------------------------------------------|
| `frame`  | Envoyer la charge utile d'une trame FlexRay |
| `status` | Interroger la configuration de slot (slot, cycle) |

## Tests

```bash
go test ./... -v
```
