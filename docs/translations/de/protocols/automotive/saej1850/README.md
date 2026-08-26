# SAE-J1850-CAN-Protokoll-SDK

SAE J1850 ist ein Fahrzeugkommunikationsstandard für die On-Board-Diagnose (OBD-II). Dieses Paket stellt einen J1850-Codec über den CAN-Bus (ISO 15765-4 / CAN TP) bereit.

## Protokollübersicht

SAE-J1850-OBD-II über CAN verwendet 29-Bit-erweiterte CAN-Identifier mit folgender Struktur:

### CAN-ID-Format (29-Bit)

| Bits       | Feld    | Beschreibung                        |
|-----------|----------|------------------------------------|
| 28-26     | Priorität | Nachrichtenpriorität (0-7, Standard 6)  |
| 25        | Ext-ID   | Bei erweiterten Frames immer 1       |
| 24-16     | PF       | Parameter Format (Header)          |
| 15-8      | PS       | Parameter Specific (Ziel/Quelle) |
| 7-0       | SA       | Quelladresse                     |

### Standard-OBD-II-CAN-IDs

| Typ                | CAN-ID (hex)    | Beschreibung                    |
|--------------------|-----------------|--------------------------------|
| Physical Request   | 0x18DAxxF1      | Anfrage an bestimmtes Steuergerät (xx=Adresse) |
| Physical Response  | 0x18DAF1xx      | Antwort vom Steuergerät (xx=Adresse)   |
| Functional Request | 0x18DB33F1      | Broadcast an alle Steuergeräte         |

### ISO-15765-2-Frame-Format

Einzelframe: hohes Nibble von Byte 0 = Datenlänge (0-7), niedriges Nibble + restliche Bytes = Diagnosedaten.

## Hardware-Anforderungen

- **Linux**-System mit SocketCAN-Unterstützung
- J1850-CAN-OBD-II-Adapter (z. B. ELM327-kompatibler USB-zu-CAN, OBDLink SX, Macchina M2)
- Fahrzeug mit OBD-II-Anschluss (die meisten Fahrzeuge ab Baujahr 1996)

### CAN-Schnittstelle einrichten

```bash
sudo modprobe can
sudo modprobe can_raw
sudo ip link set can0 type can bitrate 500000
sudo ip link set up can0
```

## Verwendung

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

## Unterstützte Funktionen

| Funktion        | OBD-II-Modus | Beschreibung                         |
|----------------|-------------|-------------------------------------|
| `mode01`       | $01         | Aktuelle Antriebsdaten abfragen     |
| `mode03`       | $03         | Abgasrelevante DTCs abfragen       |
| `mode0A`       | $0A         | Permanente DTCs abfragen              |
| `diag_request` | Benutzerdefiniert      | Allgemeine Diagnoseanfrage          |
| `diag_response`| -           | Diagnose-Antwortframe           |
| `broadcast`    | -           | Funktionaler Broadcast an alle Steuergeräte    |

## Testen

```bash
go test ./... -v
```
