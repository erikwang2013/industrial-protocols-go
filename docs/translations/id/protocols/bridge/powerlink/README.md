# SDK POWERLINK CmdBridge

POWERLINK (Ethernet POWERLINK) adalah protokol Ethernet real-time untuk otomasi industri. Paket ini menyediakan codec POWERLINK melalui utilitas baris perintah `openPOWERLINK_demo`.

## Alat CLI

Menggunakan aplikasi demo stack openPOWERLINK.

### Instalasi

```bash
git clone https://github.com/OpenAutomationTechnologies/openPOWERLINK.git
cd openPOWERLINK
cmake -B build && cmake --build build
sudo cmake --install build
```

Verifikasi: `openPOWERLINK_demo --help`

## Penggunaan

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/powerlink"
)

func main() {
    p := powerlink.New()
    codec, _ := p.NewCodec("cmd")

    req := &kernel.Request{
        Function: "read",
        Address:  "0x2000",
        Count:    8,
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## Driver

```go
b, c, err := powerlink.NewCmdDriver("/usr/bin/openPOWERLINK_demo")
```

## Fungsi yang Didukung

| Fungsi | Deskripsi                   |
|----------|-------------------------------|
| `read`   | Membaca entri object dictionary  |
| `write`  | Menulis entri object dictionary |
| `status` | Menanyakan status node/NMT          |

## Pengujian

```bash
go test ./... -v
```
