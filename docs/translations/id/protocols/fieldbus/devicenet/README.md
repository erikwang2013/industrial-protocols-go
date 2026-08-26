# SDK Protokol DeviceNet

DeviceNet adalah protokol jaringan industri berbasis CAN untuk otomasi pabrik. Paket ini menyediakan codec DeviceNet dengan driver SocketCAN dan gateway TCP.

## Ikhtisar Protokol

DeviceNet menggunakan Common Industrial Protocol (CIP) di atas CAN. Set Koneksi Master/Slave yang telah ditentukan menggunakan pesan Group 2 (CAN ID 0x400 + NodeID) untuk I/O polling dan explicit messaging.

## Persyaratan Perangkat Keras

### Mode CAN (SocketCAN)
- Sistem **Linux** dengan dukungan SocketCAN (`CONFIG_CAN` diaktifkan)
- Antarmuka CAN (mis. `can0`, `vcan0`)
- Perangkat keras CAN berkemampuan DeviceNet (mis. Anybus Communicator, HMS IXXAT)

### Mode Gateway
- Konektivitas TCP/IP ke gateway DeviceNet
- Gateway yang mendukung protokol perintah teks polos (mis. HMS Anybus, Hilscher netX)

### Menyiapkan antarmuka CAN virtual (untuk pengujian)

```bash
sudo modprobe can
sudo modprobe can_raw
sudo modprobe vcan
sudo ip link add dev vcan0 type vcan
sudo ip link set up vcan0
```

## Penggunaan

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

    // Permintaan poll
    req := &kernel.Request{
        Function: "poll",
        Data:     []byte{0x01, 0x00},
    }
    raw, _ := codec.Encode(req)

    // Buka koneksi
    req2 := &kernel.Request{Function: "open"}
    raw2, _ := codec.Encode(req2)
    _ = raw
    _ = raw2
}
```

### Mode Gateway

```go
p := devicenet.New()
codec, _ := p.NewCodec("gateway")

req := &kernel.Request{
    Function: "poll",
    Data:     []byte{0xAB, 0xCD},
}
raw, _ := codec.Encode(req)
// raw akan berupa: "poll abcd\n"
```

## Fungsi yang Didukung

| Fungsi | Deskripsi                    |
|----------|--------------------------------|
| `open`   | Membuka koneksi eksplisit       |
| `poll`   | Poll data I/O (Group 2)        |

## Pengujian

```bash
go test ./... -v
```

Catatan: pengujian driver SocketCAN memerlukan Linux. Pengujian driver gateway berjalan di OS apa pun.
