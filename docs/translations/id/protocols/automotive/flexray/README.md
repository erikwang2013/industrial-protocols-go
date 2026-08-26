# SDK Protokol CAN FlexRay

FlexRay adalah protokol komunikasi otomotif deterministik berkecepatan tinggi. Paket ini menyediakan codec FlexRay di atas bus CAN, mendukung pembingkaian berbasis siklus dengan pemeriksaan integritas CRC-16/XMODEM.

## Ikhtisar Protokol

FlexRay menggunakan skema multipleks pembagian waktu (TDMA) dengan siklus komunikasi yang berulang. Setiap siklus terdiri dari segmen statis dan dinamis. Codec ini memetakan frame FlexRay ke frame CAN extended 29-bit.

### Format Wire

Payload siklus FlexRay:
- **Header** (2 byte): nomor siklus (little-endian)
- **Status** (1 byte): bit 7=PPI (Payload Preamble Indicator), bit 6=NFI, bit 5=SYF, bit 4=SUF
- **Data** (N byte): payload (maks 254 byte)
- **CRC** (2 byte): CRC-16/XMODEM atas header+status+data (little-endian)

Pengodean CAN ID (extended 29-bit):
- Bit 28-24: Tipe pesan (0x01=frame, 0x02=status)
- Bit 23-10: Cadangan
- Bit 15-10: Slot ID (6 bit)
- Bit 9-0: Nomor siklus (10 bit)

## Persyaratan Perangkat Keras

- Sistem **Linux** dengan dukungan SocketCAN
- Adaptor Vector VN7600/VN7640 atau Bosch FlexRay-CAN
- Jaringan FlexRay dengan terminasi yang benar (bias 2.5V)

### Menyiapkan antarmuka CAN

```bash
sudo modprobe can
sudo modprobe can_raw
sudo ip link set can0 type can bitrate 500000
sudo ip link set up can0
```

## Penggunaan

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/automotive/flexray"
)

func main() {
    p := flexray.New()
    codec, _ := p.NewCodec("can")

    // Kirim frame FlexRay di siklus 5 dengan PPI diset
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

## Fungsi yang Didukung

| Fungsi | Deskripsi                              |
|----------|------------------------------------------|
| `frame`  | Mengirim payload frame FlexRay               |
| `status` | Menanyakan konfigurasi slot (slot, cycle)   |

## Pengujian

```bash
go test ./... -v
```
