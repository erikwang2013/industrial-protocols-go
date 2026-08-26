# SDK SERCOS I/II CmdBridge

SERCOS I/II adalah versi serial serat optik lama dari antarmuka SERCOS untuk kontrol gerak digital. Paket ini menyediakan codec SERCOS I/II melalui utilitas baris perintah `sercos_cli`.

## Alat CLI

Menggunakan utilitas CLI antarmuka serat optik SERCOS.

### Instalasi

```bash
# Driver dan alat kartu antarmuka SERCOS
# Lihat dokumentasi vendor untuk kartu master SERCOS tertentu
```

Verifikasi: `sercos_cli --help`

## Penggunaan

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/sercos1"
)

func main() {
    p := sercos1.New()
    codec, _ := p.NewCodec("cmd")

    req := &kernel.Request{
        Function: "read",
        Address:  "0x0100",
        Count:    4,
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## Driver

```go
b, c, err := sercos1.NewCmdDriver("/usr/bin/sercos_cli")
```

## Fungsi yang Didukung

| Fungsi | Deskripsi              |
|----------|--------------------------|
| `read`   | Membaca SERCOS IDN          |
| `write`  | Menulis SERCOS IDN         |
| `status` | Menanyakan status drive       |

## Pengujian

```bash
go test ./... -v
```
