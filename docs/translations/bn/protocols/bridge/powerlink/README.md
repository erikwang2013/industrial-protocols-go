# POWERLINK CmdBridge SDK

POWERLINK (Ethernet POWERLINK) শিল্প অটোমেশনের জন্য একটি রিয়েল-টাইম ইথারনেট প্রোটোকল। এই প্যাকেজটি `openPOWERLINK_demo` কমান্ড-লাইন ইউটিলিটির মাধ্যমে POWERLINK কোডেক সরবরাহ করে।

## CLI টুল

openPOWERLINK স্ট্যাকের ডেমো অ্যাপ্লিকেশন ব্যবহার করে।

### ইনস্টলেশন

```bash
git clone https://github.com/OpenAutomationTechnologies/openPOWERLINK.git
cd openPOWERLINK
cmake -B build && cmake --build build
sudo cmake --install build
```

যাচাই: `openPOWERLINK_demo --help`

## ব্যবহার

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

## ড্রাইভার

```go
b, c, err := powerlink.NewCmdDriver("/usr/bin/openPOWERLINK_demo")
```

## সমর্থিত ফাংশন

| ফাংশন | বিবরণ                   |
|----------|-------------------------------|
| `read`   | অবজেক্ট ডিকশনারি এন্ট্রি পড়া  |
| `write`  | অবজেক্ট ডিকশনারি এন্ট্রি লেখা |
| `status` | নোড/NMT স্টেট জিজ্ঞাসা          |

## টেস্টিং

```bash
go test ./... -v
```
