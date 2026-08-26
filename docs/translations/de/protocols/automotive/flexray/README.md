# FlexRay-CAN-Protokoll-SDK

FlexRay ist ein deterministisches Hochgeschwindigkeits-Kommunikationsprotokoll für die Automobilindustrie. Dieses Paket stellt einen FlexRay-Codec über den CAN-Bus bereit und unterstützt zyklusbasiertes Framing mit CRC-16/XMODEM-Integritätsprüfung.

## Protokollübersicht

FlexRay verwendet ein Zeitmultiplex-Verfahren (TDMA) mit sich wiederholenden Kommunikationszyklen. Jeder Zyklus besteht aus statischen und dynamischen Segmenten. Dieser Codec bildet FlexRay-Frames auf erweiterte 29-Bit-CAN-Frames ab.

### Wire-Format

FlexRay-Zyklus-Payload:
- **Header** (2 Bytes): Zyklusnummer (Little-Endian)
- **Status** (1 Byte): Bit 7=PPI (Payload Preamble Indicator), Bit 6=NFI, Bit 5=SYF, Bit 4=SUF
- **Daten** (N Bytes): Payload (max. 254 Bytes)
- **CRC** (2 Bytes): CRC-16/XMODEM über Header+Status+Daten (Little-Endian)

CAN-ID-Codierung (29-Bit erweitert):
- Bits 28-24: Nachrichtentyp (0x01=Frame, 0x02=Status)
- Bits 23-10: Reserviert
- Bits 15-10: Slot-ID (6 Bit)
- Bits 9-0: Zyklusnummer (10 Bit)

## Hardware-Anforderungen

- **Linux**-System mit SocketCAN-Unterstützung
- Vector VN7600/VN7640 oder Bosch FlexRay-CAN-Adapter
- FlexRay-Netzwerk mit korrektem Abschluss (2,5-V-Vorspannung)

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

## Unterstützte Funktionen

| Funktion | Beschreibung                              |
|----------|------------------------------------------|
| `frame`  | FlexRay-Frame-Payload senden               |
| `status` | Slot-Konfiguration abfragen (slot, cycle) |

## Testen

```bash
go test ./... -v
```
