# SDK du protocole SAE J1850 CAN

SAE J1850 est une norme de communication véhicule utilisée pour le diagnostic embarqué (OBD-II). Ce package fournit un codec J1850 sur bus CAN (ISO 15765-4 / CAN TP).

## Vue d'ensemble du protocole

L'OBD-II SAE J1850 sur CAN utilise des identifiants CAN étendus de 29 bits avec la structure suivante :

### Format d'ID CAN (29 bits)

| Bits       | Champ    | Description                        |
|-----------|----------|------------------------------------|
| 28-26     | Priority | Priorité du message (0-7, défaut 6) |
| 25        | Ext ID   | Toujours 1 pour les trames étendues |
| 24-16     | PF       | Format du paramètre (en-tête)       |
| 15-8      | PS       | Spécifique au paramètre (cible/source) |
| 7-0       | SA       | Adresse source                     |

### ID CAN OBD-II standard

| Type                | ID CAN (hex)    | Description                    |
|--------------------|-----------------|--------------------------------|
| Requête physique   | 0x18DAxxF1      | Requête vers un ECU spécifique (xx=addr) |
| Réponse physique   | 0x18DAF1xx      | Réponse d'un ECU (xx=addr)     |
| Requête fonctionnelle | 0x18DB33F1    | Diffusion vers tous les ECU    |

### Format de trame ISO 15765-2

Trame unique : le demi-octet supérieur de l'octet 0 = longueur des données (0-7), le demi-octet inférieur + les octets restants = données de diagnostic.

## Exigences matérielles

- Système **Linux** avec prise en charge SocketCAN
- Adaptateur OBD-II J1850-CAN (p. ex. USB-to-CAN compatible ELM327, OBDLink SX, Macchina M2)
- Véhicule avec connecteur OBD-II (la plupart des véhicules à partir de 1996)

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
    "github.com/erikwang2013/industrial-protocols-go/protocols/automotive/saej1850"
)

func main() {
    p := saej1850.New()
    codec, _ := p.NewCodec("can")

    // Mode $01 PID $0C: Engine RPM
    req := &kernel.Request{
        Function: "mode01",
        Metadata: map[string]any{"pid": float64(0x0C)},
    }
    raw, _ := codec.Encode(req)
    _ = raw

    // Mode $03: Request emission-related DTCs
    req2 := &kernel.Request{Function: "mode03"}
    raw2, _ := codec.Encode(req2)
    _ = raw2
}
```

## Fonctions prises en charge

| Fonction        | Mode OBD-II | Description                         |
|----------------|-------------|-------------------------------------|
| `mode01`       | $01         | Demander les données actuelles du groupe motopropulseur |
| `mode03`       | $03         | Demander les DTC liés aux émissions |
| `mode0A`       | $0A         | Demander les DTC permanents         |
| `diag_request` | Personnalisé | Requête de diagnostic générique     |
| `diag_response`| -           | Trame de réponse de diagnostic      |
| `broadcast`    | -           | Diffusion fonctionnelle vers tous les ECU |

## Tests

```bash
go test ./... -v
```
