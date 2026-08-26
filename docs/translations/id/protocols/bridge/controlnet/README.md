# SDK ControlNet CmdBridge

ControlNet adalah protokol jaringan industri real-time yang dikembangkan oleh Allen-Bradley (Rockwell Automation) untuk pertukaran data berkecepatan tinggi dan kritis terhadap waktu. Paket ini menyediakan codec ControlNet melalui utilitas baris perintah `1784-pcic-cli`.

## Alat CLI

Menggunakan utilitas CLI kartu antarmuka ControlNet 1784-PCIC.

### Instalasi

```bash
# Install driver dan alat Rockwell 1784-PCIC
# Lihat dokumentasi Rockwell Automation untuk SDK RSLinx Classic
```

Verifikasi: `1784-pcic-cli --help`

## Penggunaan

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/controlnet"
)

func main() {
    p := controlnet.New()
    codec, _ := p.NewCodec("cmd")

    req := &kernel.Request{
        Function: "read",
        Address:  "0x10",
        Count:    8,
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## Driver

```go
b, c, err := controlnet.NewCmdDriver("/usr/bin/1784-pcic-cli")
```

## Fungsi yang Didukung

| Fungsi | Deskripsi               |
|----------|---------------------------|
| `read`   | Membaca dari node ControlNet |
| `write`  | Menulis ke node ControlNet  |
| `status` | Menanyakan status kartu PCIC    |

## Pengujian

```bash
go test ./... -v
```
