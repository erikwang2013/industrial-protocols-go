# SDK ISA100.11a CmdBridge

ISA100.11a adalah standar jaringan industri nirkabel untuk otomasi proses. Paket ini menyediakan codec ISA100.11a melalui utilitas baris perintah `yfgw410_cli` (gateway nirkabel lapangan Yokogawa YFGW410).

## Alat CLI

Menggunakan utilitas CLI Yokogawa YFGW410 Field Wireless Gateway.

### Instalasi

```bash
# Install perangkat lunak dan alat gateway Yokogawa YFGW410
# Lihat dokumentasi Yokogawa untuk penyiapan Field Wireless Gateway
```

Verifikasi: `yfgw410_cli --help`

## Penggunaan

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/isa100"
)

func main() {
    p := isa100.New()
    codec, _ := p.NewCodec("cmd")

    req := &kernel.Request{
        Function: "read",
        Address:  "DEV001",
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## Driver

```go
b, c, err := isa100.NewCmdDriver("/usr/bin/yfgw410_cli")
```

## Fungsi yang Didukung

| Fungsi | Deskripsi                 |
|----------|-----------------------------|
| `read`   | Membaca atribut perangkat       |
| `write`  | Menulis atribut perangkat      |
| `list`   | Menampilkan daftar perangkat ter-provisi    |

## Pengujian

```bash
go test ./... -v
```
