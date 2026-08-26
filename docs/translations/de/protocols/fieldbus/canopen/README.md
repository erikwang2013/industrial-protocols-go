# CANopen-Protokoll-SDK

CANopen ist ein CAN-basiertes Higher-Layer-Protokoll für eingebettete Steuerungssysteme. Dieses Paket stellt einen CANopen-Codec und einen SocketCAN-Treiber bereit.

## Protokollübersicht

CANopen verwendet standardmäßige 11-Bit-CAN-Identifier mit folgendem vordefinierten Verbindungssatz:

| Funktion    | CAN-ID              | Beschreibung                   |
|------------|---------------------|-------------------------------|
| NMT        | 0x000               | Netzwerkverwaltung            |
| SYNC       | 0x080               | Synchronisationsnachricht       |
| SDO (tx)   | 0x580 + NodeID      | Service Data Object (Server)  |
| SDO (rx)   | 0x600 + NodeID      | Service Data Object (Client)  |
| PDO1 (tx)  | 0x180 + NodeID      | Process Data Object 1         |
| Heartbeat  | 0x700 + NodeID      | Heartbeat / Bootup            |

## Hardware-Anforderungen

- **Linux**-System mit SocketCAN-Unterstützung (`CONFIG_CAN` aktiviert)
- Eine CAN-Schnittstelle (z. B. `can0`, `vcan0` für virtuelles CAN)
- CAN-fähige Transceiver-Hardware (z. B. MCP2515, SJA1000 oder USB-CAN-Adapter)

### Virtuelle CAN-Schnittstelle einrichten (für Tests)

```bash
sudo modprobe can
sudo modprobe can_raw
sudo modprobe vcan
sudo ip link add dev vcan0 type vcan
sudo ip link set up vcan0
```

## Verwendung

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

## Unterstützte Funktionen

| Funktion     | Beschreibung                          |
|-------------|--------------------------------------|
| `sdo_read`  | Objektverzeichniseintrag lesen         |
| `sdo_write` | Objektverzeichniseintrag schreiben        |
| `nmt_start` | Remote-Knoten starten (NMT)              |
| `nmt_stop`  | Remote-Knoten stoppen (NMT)               |
| `nmt_reset` | Remote-Knoten zurücksetzen (NMT)              |
| `heartbeat` | Heartbeat-/Bootup-Nachricht senden      |

## Testen

```bash
go test ./... -v
```

Hinweis: Die SocketCAN-Treibertests erfordern ein Linux-System mit CAN-Hardware oder einer virtuellen CAN-Schnittstelle.
