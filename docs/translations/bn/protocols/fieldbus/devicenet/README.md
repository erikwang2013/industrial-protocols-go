# DeviceNet প্রোটোকল SDK

DeviceNet ফ্যাক্টরি অটোমেশনের জন্য একটি CAN-ভিত্তিক শিল্প নেটওয়ার্ক প্রোটোকল। এই প্যাকেজটি SocketCAN ও TCP গেটওয়ে উভয় ড্রাইভারসহ একটি DeviceNet কোডেক সরবরাহ করে।

## প্রোটোকল ওভারভিউ

DeviceNet CAN-এর উপর Common Industrial Protocol (CIP) ব্যবহার করে। প্রি-ডিফাইনড মাস্টার/স্লেভ কানেকশন সেট পোলড I/O ও এক্সপ্লিসিট মেসেজিংয়ের জন্য Group 2 মেসেজ (CAN ID 0x400 + NodeID) ব্যবহার করে।

## হার্ডওয়্যার প্রয়োজনীয়তা

### CAN মোড (SocketCAN)
- SocketCAN সমর্থনসহ **Linux** সিস্টেম (`CONFIG_CAN` সক্রিয়)
- একটি CAN ইন্টারফেস (যেমন `can0`, `vcan0`)
- DeviceNet-সক্ষম CAN হার্ডওয়্যার (যেমন Anybus Communicator, HMS IXXAT)

### গেটওয়ে মোড
- DeviceNet গেটওয়েতে TCP/IP কানেক্টিভিটি
- প্লেইন-টেক্সট কমান্ড প্রোটোকল সমর্থনকারী গেটওয়ে (যেমন HMS Anybus, Hilscher netX)

### ভার্চুয়াল CAN ইন্টারফেস সেটআপ (টেস্টিংয়ের জন্য)

```bash
sudo modprobe can
sudo modprobe can_raw
sudo modprobe vcan
sudo ip link add dev vcan0 type vcan
sudo ip link set up vcan0
```

## ব্যবহার

### CAN মোড

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/fieldbus/devicenet"
)

func main() {
    p := devicenet.New()
    codec, _ := p.NewCodec("can")

    // Poll request
    req := &kernel.Request{
        Function: "poll",
        Data:     []byte{0x01, 0x00},
    }
    raw, _ := codec.Encode(req)

    // Open connection
    req2 := &kernel.Request{Function: "open"}
    raw2, _ := codec.Encode(req2)
    _ = raw
    _ = raw2
}
```

### গেটওয়ে মোড

```go
p := devicenet.New()
codec, _ := p.NewCodec("gateway")

req := &kernel.Request{
    Function: "poll",
    Data:     []byte{0xAB, 0xCD},
}
raw, _ := codec.Encode(req)
// raw will be: "poll abcd\n"
```

## সমর্থিত ফাংশন

| ফাংশন | বিবরণ                    |
|----------|--------------------------------|
| `open`   | এক্সপ্লিসিট কানেকশন খোলা       |
| `poll`   | I/O ডেটা পোল করা (Group 2)        |

## টেস্টিং

```bash
go test ./... -v
```

দ্রষ্টব্য: SocketCAN ড্রাইভার টেস্টের জন্য Linux প্রয়োজন। গেটওয়ে ড্রাইভার টেস্ট যেকোনো OS-এ চলে।
