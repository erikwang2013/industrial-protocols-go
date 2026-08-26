# DeviceNet-Protokoll-SDK

DeviceNet ist ein CAN-basiertes Industrienetzwerkprotokoll für die Fabrikautomatisierung. Dieses Paket stellt einen DeviceNet-Codec mit SocketCAN- und TCP-Gateway-Treibern bereit.

## Protokollübersicht

DeviceNet verwendet das Common Industrial Protocol (CIP) über CAN. Der vordefinierte Master/Slave-Verbindungssatz nutzt Group-2-Nachrichten (CAN-ID 0x400 + NodeID) für gepollte I/O und Explicit Messaging.

## Hardware-Anforderungen

### CAN-Modus (SocketCAN)
- **Linux**-System mit SocketCAN-Unterstützung (`CONFIG_CAN` aktiviert)
- Eine CAN-Schnittstelle (z. B. `can0`, `vcan0`)
- DeviceNet-fähige CAN-Hardware (z. B. Anybus Communicator, HMS IXXAT)

### Gateway-Modus
- TCP/IP-Konnektivität zum DeviceNet-Gateway
- Gateway mit Klartext-Kommando-Protokoll (z. B. HMS Anybus, Hilscher netX)

### Virtuelle CAN-Schnittstelle einrichten (für Tests)

```bash
sudo modprobe can
sudo modprobe can_raw
sudo modprobe vcan
sudo ip link add dev vcan0 type vcan
sudo ip link set up vcan0
```

## Verwendung

### CAN-Modus

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

### Gateway-Modus

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

## Unterstützte Funktionen

| Funktion | Beschreibung                    |
|----------|--------------------------------|
| `open`   | Explizite Verbindung öffnen       |
| `poll`   | I/O-Daten pollen (Group 2)        |

## Testen

```bash
go test ./... -v
```

Hinweis: SocketCAN-Treibertests erfordern Linux. Gateway-Treibertests laufen auf jedem Betriebssystem.
