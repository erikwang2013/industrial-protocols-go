# SDK Protokol CAN SAE J1850

SAE J1850 adalah standar komunikasi kendaraan yang digunakan untuk diagnostik on-board (OBD-II). Paket ini menyediakan codec J1850 di atas bus CAN (ISO 15765-4 / CAN TP).

## Ikhtisar Protokol

SAE J1850 OBD-II over CAN menggunakan identifier CAN extended 29-bit dengan struktur berikut:

### Format CAN ID (29-bit)

| Bit       | Field    | Deskripsi                        |
|-----------|----------|------------------------------------|
| 28-26     | Priority | Prioritas pesan (0-7, default 6)  |
| 25        | Ext ID   | Selalu 1 untuk frame extended       |
| 24-16     | PF       | Parameter Format (Header)          |
| 15-8      | PS       | Parameter Specific (target/sumber) |
| 7-0       | SA       | Alamat Sumber                     |

### CAN ID OBD-II Standar

| Tipe                | CAN ID (hex)    | Deskripsi                    |
|--------------------|-----------------|--------------------------------|
| Physical Request   | 0x18DAxxF1      | Permintaan ke ECU tertentu (xx=addr) |
| Physical Response  | 0x18DAF1xx      | Respons dari ECU (xx=addr)   |
| Functional Request | 0x18DB33F1      | Siaran ke semua ECU         |

### Format Frame ISO 15765-2

Single frame: nibble atas byte 0 = panjang data (0-7), nibble bawah + byte sisanya = data diagnostik.

## Persyaratan Perangkat Keras

- Sistem **Linux** dengan dukungan SocketCAN
- Adaptor OBD-II J1850-CAN (mis. USB-ke-CAN kompatibel ELM327, OBDLink SX, Macchina M2)
- Kendaraan dengan konektor OBD-II (sebagian besar kendaraan 1996+)

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
    "github.com/erikwang2013/industrial-protocols-go/protocols/automotive/saej1850"
)

func main() {
    p := saej1850.New()
    codec, _ := p.NewCodec("can")

    // Mode $01 PID $0C: RPM mesin
    req := &kernel.Request{
        Function: "mode01",
        Metadata: map[string]any{"pid": float64(0x0C)},
    }
    raw, _ := codec.Encode(req)
    _ = raw

    // Mode $03: Minta DTC terkait emisi
    req2 := &kernel.Request{Function: "mode03"}
    raw2, _ := codec.Encode(req2)
    _ = raw2
}
```

## Fungsi yang Didukung

| Fungsi        | Mode OBD-II | Deskripsi                         |
|----------------|-------------|-------------------------------------|
| `mode01`       | $01         | Minta data powertrain saat ini     |
| `mode03`       | $03         | Minta DTC terkait emisi       |
| `mode0A`       | $0A         | Minta DTC permanen              |
| `diag_request` | Kustom      | Permintaan diagnostik generik          |
| `diag_response`| -           | Frame respons diagnostik           |
| `broadcast`    | -           | Siaran fungsional ke semua ECU    |

## Pengujian

```bash
go test ./... -v
```
