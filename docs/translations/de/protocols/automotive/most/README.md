# MOST-Serial-Protokoll-SDK

MOST (Media Oriented Systems Transport) ist eine Hochgeschwindigkeits-Multimedia-Netzwerktechnologie, die vor allem in Kfz-Infotainmentsystemen eingesetzt wird. Dieses Paket stellt einen MOST-Codec über einen seriellen Adapter mit AT-Befehlsschnittstelle bereit.

## Protokollübersicht

MOST nutzt synchrone serielle Kommunikation über eine Glasfaser-Physik. Diese Implementierung verbindet sich über einen seriellen Adapter, der eine AT-Befehlsschnittstelle mit 115200 Baud bereitstellt.

### AT-Befehle

| Befehl            | Beschreibung                  |
|--------------------|------------------------------|
| `AT+READ=<addr>,<len>`  | Bytes von einer Adresse lesen |
| `AT+WRITE=<addr>,<hex>` | Hex-Bytes an eine Adresse schreiben |
| `AT+STATUS`             | Ring-/Netzwerkstatus abfragen |

### Antwortformat

- `+OK:<hex_data>` — erfolgreiche Antwort
- `+ERR:<code>` — Fehlerantwort

## Hardware-Anforderungen

- MOST-Glasfasernetzwerk mit korrektem Abschluss
- MOST-zu-Serial-Adapter (z. B. MOST150-USB-Adapter)
- Serielle Schnittstelle mit 115200 Baud, 8N1

## Verwendung

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

## Treiber

```go
b, c, err := most.NewSerialDriver("/dev/ttyUSB0")
```

## Unterstützte Funktionen

| Funktion | Beschreibung                 |
|----------|-----------------------------|
| `read`   | Von einer MOST-Adresse lesen |
| `write`  | Daten an eine MOST-Adresse schreiben |
| `status` | Ring-/Netzwerkstatus abfragen |

## Testen

```bash
go test ./... -v
```
