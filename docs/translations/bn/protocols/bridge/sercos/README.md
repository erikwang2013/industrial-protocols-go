# SERCOS III CmdBridge SDK

SERCOS III (SErial Real-time COmmunication System) মোশন কন্ট্রোলের জন্য একটি ডিজিটাল ইন্টারফেস, যা ফাইবার-অপটিক বা কপার ক্যাবলের উপর রিং টপোলজি ব্যবহার করে। এই প্যাকেজটি `netx_cli` কমান্ড-লাইন ইউটিলিটির মাধ্যমে SERCOS III কোডেক সরবরাহ করে।

## CLI টুল

Hilscher netX SERCOS III CLI ইউটিলিটি ব্যবহার করে।

### ইনস্টলেশন

```bash
# Install Hilscher netX driver and tools
# See https://www.hilscher.com/ for netX drivers
sudo apt-get install netx-driver
```

যাচাই: `netx_cli --help`

## ব্যবহার

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

## ড্রাইভার

```go
b, c, err := sercos.NewCmdDriver("/usr/bin/netx_cli")
```

## সমর্থিত ফাংশন

| ফাংশন | বিবরণ                 |
|----------|-----------------------------|
| `read`   | SERCOS IDN/S প্যারামিটার পড়া |
| `write`  | SERCOS IDN/S প্যারামিটার লেখা|
| `phase`  | কমিউনিকেশন ফেজ সেট করা     |

## টেস্টিং

```bash
go test ./... -v
```
