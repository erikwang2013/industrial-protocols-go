# SDK EtherCAT CmdBridge

EtherCAT (Ethernet for Control Automation Technology) adalah fieldbus Ethernet industri berkinerja tinggi. Paket ini menyediakan codec EtherCAT melalui utilitas baris perintah `ethercat`.

## Alat CLI

Menggunakan alat baris perintah IgH EtherCAT Master.

### Instalasi

```bash
# Debian/Ubuntu
sudo apt-get install ethercat-master

# Dari sumber
git clone https://gitlab.com/etherlab.org/ethercat.git
cd ethercat
make && sudo make install
```

Verifikasi: `ethercat slaves`

## Penggunaan

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/ethercat"
)

func main() {
    p := ethercat.New()
    codec, _ := p.NewCodec("cmd")

    // Upload SDO dari alamat 0x1000
    req := &kernel.Request{
        Function: "upload",
        Address:  "0x1000",
        Count:    4,
    }
    raw, _ := codec.Encode(req)
    // raw = "upload 0x1000 4\n"
    _ = raw
}
```

## Driver

```go
b, c, err := ethercat.NewCmdDriver("/usr/bin/ethercat")
```

## Fungsi yang Didukung

| Fungsi   | Deskripsi            |
|------------|------------------------|
| `upload`   | Membaca SDO dari alamat  |
| `download` | Menulis SDO ke alamat   |
| `slaves`   | Menampilkan daftar slave EtherCAT   |

## Pengujian

```bash
go test ./... -v
```
