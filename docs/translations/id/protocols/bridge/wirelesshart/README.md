# SDK WirelessHART CmdBridge

WirelessHART (IEC 62591) adalah standar jaringan industri nirkabel yang berbasis protokol HART. Paket ini menyediakan codec WirelessHART melalui utilitas baris perintah `emerson_1410_cli` (Emerson 1410/1420 Wireless Gateway).

## Alat CLI

Menggunakan utilitas CLI Emerson 1410/1420 Wireless Gateway.

### Instalasi

```bash
# Install perangkat lunak dan alat Emerson Wireless Gateway
# Lihat dokumentasi Emerson untuk penyiapan Gateway 1410/1420
```

Verifikasi: `emerson_1410_cli --help`

## Penggunaan

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/wirelesshart"
)

func main() {
    p := wirelesshart.New()
    codec, _ := p.NewCodec("cmd")

    req := &kernel.Request{
        Function: "read",
        Address:  "TT101",
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## Driver

```go
b, c, err := wirelesshart.NewCmdDriver("/usr/bin/emerson_1410_cli")
```

## Fungsi yang Didukung

| Fungsi | Deskripsi                |
|----------|----------------------------|
| `read`   | Membaca parameter perangkat      |
| `write`  | Menulis parameter perangkat     |
| `scan`   | Memindai perangkat nirkabel  |

## Pengujian

```bash
go test ./... -v
```
