# SDK du protocole MOST série

MOST (Media Oriented Systems Transport) est une technologie de réseau multimédia à haut débit utilisée principalement dans les systèmes d'infodivertissement automobiles. Ce package fournit un codec MOST sur adaptateur série avec interface de commandes AT.

## Vue d'ensemble du protocole

MOST utilise une communication série synchrone sur couche physique à fibre optique. Cette implémentation se connecte via un adaptateur série qui expose une interface de commandes AT à 115200 bauds.

### Commandes AT

| Commande            | Description                  |
|--------------------|------------------------------|
| `AT+READ=<addr>,<len>`  | Lire des octets à une adresse |
| `AT+WRITE=<addr>,<hex>` | Écrire des octets hexadécimaux à une adresse |
| `AT+STATUS`             | Interroger l'état de l'anneau/réseau |

### Format de réponse

- `+OK:<hex_data>` -- réponse de succès
- `+ERR:<code>` -- réponse d'erreur

## Exigences matérielles

- Réseau à fibre optique MOST avec terminaison correcte
- Adaptateur MOST-série (p. ex. adaptateur USB MOST150)
- Port série à 115200 bauds, 8N1

## Utilisation

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/automotive/most"
)

func main() {
    p := most.New()
    codec, _ := p.NewCodec("serial")

    // Read 4 bytes from address 0x0100
    req := &kernel.Request{
        Function: "read",
        Address:  "0x0100",
        Count:    4,
    }
    raw, _ := codec.Encode(req)
    // raw = "AT+READ=0x0100,4\r\n"
    _ = raw

    // Write data to address 0x0200
    req2 := &kernel.Request{
        Function: "write",
        Address:  "0x0200",
        Data:     []byte{0xDE, 0xAD, 0xBE, 0xEF},
    }
    raw2, _ := codec.Encode(req2)
    // raw2 = "AT+WRITE=0x0200,DEADBEEF\r\n"
    _ = raw2
}
```

## Pilote

```go
b, c, err := most.NewSerialDriver("/dev/ttyUSB0")
```

## Fonctions prises en charge

| Fonction | Description                 |
|----------|-----------------------------|
| `read`   | Lire depuis une adresse MOST |
| `write`  | Écrire des données vers une adresse MOST |
| `status` | Interroger l'état de l'anneau/réseau |

## Tests

```bash
go test ./... -v
```
