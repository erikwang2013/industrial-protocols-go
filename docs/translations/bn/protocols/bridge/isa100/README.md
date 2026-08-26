# ISA100.11a CmdBridge SDK

ISA100.11a প্রসেস অটোমেশনের জন্য একটি ওয়্যারলেস শিল্প নেটওয়ার্কিং স্ট্যান্ডার্ড। এই প্যাকেজটি `yfgw410_cli` (Yokogawa YFGW410 ফিল্ড ওয়্যারলেস গেটওয়ে) কমান্ড-লাইন ইউটিলিটির মাধ্যমে ISA100.11a কোডেক সরবরাহ করে।

## CLI টুল

Yokogawa YFGW410 ফিল্ড ওয়্যারলেস গেটওয়ে CLI ব্যবহার করে।

### ইনস্টলেশন

```bash
# Install Yokogawa YFGW410 gateway software and tools
# Refer to Yokogawa documentation for Field Wireless Gateway setup
```

যাচাই: `yfgw410_cli --help`

## ব্যবহার

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

## ড্রাইভার

```go
b, c, err := isa100.NewCmdDriver("/usr/bin/yfgw410_cli")
```

## সমর্থিত ফাংশন

| ফাংশন | বিবরণ                 |
|----------|-----------------------------|
| `read`   | ডিভাইস অ্যাট্রিবিউট পড়া       |
| `write`  | ডিভাইস অ্যাট্রিবিউট লেখা      |
| `list`   | প্রভিশন্ড ডিভাইসের তালিকা    |

## টেস্টিং

```bash
go test ./... -v
```
