# WirelessHART CmdBridge SDK

WirelessHART (IEC 62591) HART প্রোটোকলের উপর ভিত্তি করে একটি ওয়্যারলেস শিল্প নেটওয়ার্কিং স্ট্যান্ডার্ড। এই প্যাকেজটি `emerson_1410_cli` (Emerson 1410/1420 ওয়্যারলেস গেটওয়ে) কমান্ড-লাইন ইউটিলিটির মাধ্যমে WirelessHART কোডেক সরবরাহ করে।

## CLI টুল

Emerson 1410/1420 ওয়্যারলেস গেটওয়ে CLI ইউটিলিটি ব্যবহার করে।

### ইনস্টলেশন

```bash
# Install Emerson Wireless Gateway software and tools
# Refer to Emerson documentation for 1410/1420 Gateway setup
```

যাচাই: `emerson_1410_cli --help`

## ব্যবহার

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

## ড্রাইভার

```go
b, c, err := wirelesshart.NewCmdDriver("/usr/bin/emerson_1410_cli")
```

## সমর্থিত ফাংশন

| ফাংশন | বিবরণ                |
|----------|----------------------------|
| `read`   | ডিভাইস প্যারামিটার পড়া      |
| `write`  | ডিভাইস প্যারামিটার লেখা     |
| `scan`   | ওয়্যারলেস ডিভাইস স্ক্যান  |

## টেস্টিং

```bash
go test ./... -v
```
