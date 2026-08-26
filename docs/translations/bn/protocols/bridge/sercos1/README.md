# SERCOS I/II CmdBridge SDK

SERCOS I/II ডিজিটাল মোশন কন্ট্রোলের জন্য SERCOS ইন্টারফেসের লিগ্যাসি সিরিয়াল ফাইবার-অপটিক সংস্করণ। এই প্যাকেজটি `sercos_cli` কমান্ড-লাইন ইউটিলিটির মাধ্যমে SERCOS I/II কোডেক সরবরাহ করে।

## CLI টুল

একটি SERCOS ফাইবার-অপটিক ইন্টারফেস CLI ইউটিলিটি ব্যবহার করে।

### ইনস্টলেশন

```bash
# SERCOS interface card driver and tools
# Refer to vendor documentation for the specific SERCOS master card
```

যাচাই: `sercos_cli --help`

## ব্যবহার

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

## ড্রাইভার

```go
b, c, err := sercos1.NewCmdDriver("/usr/bin/sercos_cli")
```

## সমর্থিত ফাংশন

| ফাংশন | বিবরণ              |
|----------|--------------------------|
| `read`   | SERCOS IDN পড়া          |
| `write`  | SERCOS IDN লেখা         |
| `status` | ড্রাইভ স্ট্যাটাস জিজ্ঞাসা       |

## টেস্টিং

```bash
go test ./... -v
```
