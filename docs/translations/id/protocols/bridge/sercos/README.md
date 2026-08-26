# SDK SERCOS III CmdBridge

SERCOS III (SErial Real-time COmmunication System) adalah antarmuka digital untuk kontrol gerak, menggunakan topologi cincin di atas serat optik atau tembaga. Paket ini menyediakan codec SERCOS III melalui utilitas baris perintah `netx_cli`.

## Alat CLI

Menggunakan utilitas CLI Hilscher netX SERCOS III.

### Instalasi

```bash
# Install driver dan alat Hilscher netX
# Lihat https://www.hilscher.com/ untuk driver netX
sudo apt-get install netx-driver
```

Verifikasi: `netx_cli --help`

## Penggunaan

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/sercos"
)

func main() {
    p := sercos.New()
    codec, _ := p.NewCodec("cmd")

    req := &kernel.Request{
        Function: "read",
        Address:  "S-0-51",
        Count:    2,
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## Driver

```go
b, c, err := sercos.NewCmdDriver("/usr/bin/netx_cli")
```

## Fungsi yang Didukung

| Fungsi | Deskripsi                 |
|----------|-----------------------------|
| `read`   | Membaca parameter SERCOS IDN/S |
| `write`  | Menulis parameter SERCOS IDN/S|
| `phase`  | Mengatur fase komunikasi     |

## Pengujian

```bash
go test ./... -v
```
