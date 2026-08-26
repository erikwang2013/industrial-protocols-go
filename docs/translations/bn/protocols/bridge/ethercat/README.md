# EtherCAT CmdBridge SDK

EtherCAT (Ethernet for Control Automation Technology) একটি উচ্চ-পারফরম্যান্স শিল্প ইথারনেট ফিল্ডবাস। এই প্যাকেজটি `ethercat` কমান্ড-লাইন ইউটিলিটির মাধ্যমে EtherCAT কোডেক সরবরাহ করে।

## CLI টুল

IgH EtherCAT Master কমান্ড-লাইন টুল ব্যবহার করে।

### ইনস্টলেশন

```bash
# Debian/Ubuntu
sudo apt-get install ethercat-master

# From source
git clone https://gitlab.com/etherlab.org/ethercat.git
cd ethercat
make && sudo make install
```

যাচাই: `ethercat slaves`

## ব্যবহার

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/ethercat"
)

func main() {
    p := ethercat.New()
    codec, _ := p.NewCodec("cmd")

    // Upload SDO from address 0x1000
    req := &kernel.Request{
        Function: "upload",
        Address:  "0x1000",
        Count:    4,
    }
    raw, _ := codec.Encode(req)
    // raw = "upload 0x1000 4\n"
    _ = raw
}
```

## ড্রাইভার

```go
b, c, err := ethercat.NewCmdDriver("/usr/bin/ethercat")
```

## সমর্থিত ফাংশন

| ফাংশন   | বিবরণ            |
|------------|------------------------|
| `upload`   | ঠিকানা থেকে SDO পড়া  |
| `download` | ঠিকানায় SDO লেখা   |
| `slaves`   | EtherCAT স্লেভ তালিকা   |

## টেস্টিং

```bash
go test ./... -v
```
