# SDK Protokol CANopen

CANopen adalah protokol lapisan atas berbasis CAN untuk sistem kontrol tertanam. Paket ini menyediakan codec CANopen dan driver SocketCAN.

## Ikhtisar Protokol

CANopen menggunakan identifier CAN standar 11-bit dengan set koneksi bawaan berikut:

| Fungsi    | CAN ID              | Deskripsi                   |
|------------|---------------------|-------------------------------|
| NMT        | 0x000               | Manajemen jaringan            |
| SYNC       | 0x080               | Pesan sinkronisasi       |
| SDO (tx)   | 0x580 + NodeID      | Service Data Object (server)  |
| SDO (rx)   | 0x600 + NodeID      | Service Data Object (client)  |
| PDO1 (tx)  | 0x180 + NodeID      | Process Data Object 1         |
| Heartbeat  | 0x700 + NodeID      | Heartbeat / Bootup            |

## Persyaratan Perangkat Keras

- Sistem **Linux** dengan dukungan SocketCAN (`CONFIG_CAN` diaktifkan)
- Antarmuka CAN (mis. `can0`, `vcan0` untuk CAN virtual)
- Perangkat keras transceiver berkemampuan CAN (mis. MCP2515, SJA1000, atau adaptor USB-CAN)

### Menyiapkan antarmuka CAN virtual (untuk pengujian)

```bash
sudo modprobe can
sudo modprobe can_raw
sudo modprobe vcan
sudo ip link add dev vcan0 type vcan
sudo ip link set up vcan0
```

## Penggunaan

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

    // Baca entri object dictionary SDO
    req := &kernel.Request{
        Function: "sdo_read",
        Metadata: map[string]any{
            "index": float64(0x1000), // Tipe perangkat
            "sub":   float64(0),
        },
    }
    raw, _ := codec.Encode(req)
    fmt.Printf("Frame baca SDO: %X\n", raw)

    // NMT start node remote
    req2 := &kernel.Request{Function: "nmt_start"}
    raw2, _ := codec.Encode(req2)
    fmt.Printf("Frame start NMT: %X\n", raw2)
}
```

## Fungsi yang Didukung

| Fungsi     | Deskripsi                          |
|-------------|--------------------------------------|
| `sdo_read`  | Membaca entri object dictionary         |
| `sdo_write` | Menulis entri object dictionary        |
| `nmt_start` | Memulai node remote (NMT)              |
| `nmt_stop`  | Menghentikan node remote (NMT)               |
| `nmt_reset` | Mereset node remote (NMT)              |
| `heartbeat` | Mengirim pesan heartbeat / bootup      |

## Pengujian

```bash
go test ./... -v
```

Catatan: pengujian driver SocketCAN memerlukan sistem Linux dengan perangkat keras CAN atau antarmuka CAN virtual.
