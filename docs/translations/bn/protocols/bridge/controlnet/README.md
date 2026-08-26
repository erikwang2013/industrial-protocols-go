# ControlNet CmdBridge SDK

ControlNet হল Allen-Bradley (Rockwell Automation) কর্তৃক উন্নত একটি রিয়েল-টাইম শিল্প নেটওয়ার্ক প্রোটোকল, যা উচ্চ-গতির, টাইম-ক্রিটিকাল ডেটা এক্সচেঞ্জের জন্য ব্যবহৃত হয়। এই প্যাকেজটি `1784-pcic-cli` কমান্ড-লাইন ইউটিলিটির মাধ্যমে ControlNet কোডেক সরবরাহ করে।

## CLI টুল

1784-PCIC ControlNet ইন্টারফেস কার্ডের CLI ইউটিলিটি ব্যবহার করে।

### ইনস্টলেশন

```bash
# Install Rockwell 1784-PCIC driver and tools
# Refer to Rockwell Automation documentation for RSLinx Classic SDK
```

যাচাই: `1784-pcic-cli --help`

## ব্যবহার

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

## ড্রাইভার

```go
b, c, err := controlnet.NewCmdDriver("/usr/bin/1784-pcic-cli")
```

## সমর্থিত ফাংশন

| ফাংশন | বিবরণ               |
|----------|---------------------------|
| `read`   | ControlNet নোড থেকে পড়া |
| `write`  | ControlNet নোডে লেখা  |
| `status` | PCIC কার্ড স্ট্যাটাস জিজ্ঞাসা    |

## টেস্টিং

```bash
go test ./... -v
```
