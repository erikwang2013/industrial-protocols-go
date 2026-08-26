# SDK Protokol Serial MOST

MOST (Media Oriented Systems Transport) adalah teknologi jaringan multimedia berkecepatan tinggi yang terutama digunakan pada sistem infotainment otomotif. Paket ini menyediakan codec MOST melalui adaptor serial dengan antarmuka perintah AT.

## Ikhtisar Protokol

MOST menggunakan komunikasi serial sinkron di atas lapisan fisik serat optik. Implementasi ini terhubung melalui adaptor serial yang mengekspos antarmuka perintah AT pada baudrate 115200.

### Perintah AT

| Perintah            | Deskripsi                  |
|--------------------|------------------------------|
| `AT+READ=<addr>,<len>`  | Membaca byte dari alamat |
| `AT+WRITE=<addr>,<hex>` | Menulis byte heksadesimal ke alamat |
| `AT+STATUS`             | Menanyakan status ring/jaringan |

### Format Respons

- `+OK:<hex_data>` -- respons sukses
- `+ERR:<code>` -- respons error

## Persyaratan Perangkat Keras

- Jaringan serat optik MOST dengan terminasi yang benar
- Adaptor MOST-ke-serial (mis. adaptor USB MOST150)
- Port serial pada baudrate 115200, 8N1

## Penggunaan

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/automotive/most"
)

func main() {
    p := most.New()
    codec, _ := p.NewCodec("serial")

    // Baca 4 byte dari alamat 0x0100
    req := &kernel.Request{
        Function: "read",
        Address:  "0x0100",
        Count:    4,
    }
    raw, _ := codec.Encode(req)
    // raw = "AT+READ=0x0100,4\r\n"
    _ = raw

    // Tulis data ke alamat 0x0200
    req2 := &kernel.Request{
        Function: "write",
        Address:  "0x0200",
        Data:     []byte{0xDE, 0xAD, 0xBE, 0xEF},
    }
    raw2, _ := codec.Encode(req2)
    // raw2 = "AT+WRITE=0x0200,DEADBEEF\r\n"
    _ = raw2
}
```

## Driver

```go
b, c, err := most.NewSerialDriver("/dev/ttyUSB0")
```

## Fungsi yang Didukung

| Fungsi | Deskripsi                 |
|----------|-----------------------------|
| `read`   | Membaca dari alamat MOST    |
| `write`  | Menulis data ke alamat MOST |
| `status` | Menanyakan status ring/jaringan   |

## Pengujian

```bash
go test ./... -v
```
