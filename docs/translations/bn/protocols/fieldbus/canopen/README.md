# CANopen প্রোটোকল SDK

CANopen এমবেডেড কন্ট্রোল সিস্টেমের জন্য CAN-ভিত্তিক একটি উচ্চ-স্তরের প্রোটোকল। এই প্যাকেজটি একটি CANopen কোডেক এবং SocketCAN ড্রাইভার সরবরাহ করে।

## প্রোটোকল ওভারভিউ

CANopen স্ট্যান্ডার্ড ১১-বিট CAN আইডেন্টিফায়ার ব্যবহার করে, নিম্নলিখিত প্রি-ডিফাইনড কানেকশন সেটসহ:

| ফাংশন    | CAN ID              | বিবরণ                   |
|------------|---------------------|-------------------------------|
| NMT        | 0x000               | নেটওয়ার্ক ম্যানেজমেন্ট            |
| SYNC       | 0x080               | সিংক্রোনাইজেশন মেসেজ       |
| SDO (tx)   | 0x580 + NodeID      | Service Data Object (সার্ভার)  |
| SDO (rx)   | 0x600 + NodeID      | Service Data Object (ক্লায়েন্ট)  |
| PDO1 (tx)  | 0x180 + NodeID      | Process Data Object 1         |
| Heartbeat  | 0x700 + NodeID      | Heartbeat / Bootup            |

## হার্ডওয়্যার প্রয়োজনীয়তা

- SocketCAN সমর্থনসহ **Linux** সিস্টেম (`CONFIG_CAN` সক্রিয়)
- একটি CAN ইন্টারফেস (যেমন `can0`, ভার্চুয়াল CAN-এর জন্য `vcan0`)
- CAN-সক্ষম ট্রান্সসিভার হার্ডওয়্যার (যেমন MCP2515, SJA1000, বা USB-CAN অ্যাডাপ্টার)

### ভার্চুয়াল CAN ইন্টারফেস সেটআপ (টেস্টিংয়ের জন্য)

```bash
sudo modprobe can
sudo modprobe can_raw
sudo modprobe vcan
sudo ip link add dev vcan0 type vcan
sudo ip link set up vcan0
```

## ব্যবহার

```go
package main

import (
    "fmt"
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/fieldbus/canopen"
)

func main() {
    p := canopen.New()
    codec, _ := p.NewCodec("can")

    // Read SDO object dictionary entry
    req := &kernel.Request{
        Function: "sdo_read",
        Metadata: map[string]any{
            "index": float64(0x1000), // Device type
            "sub":   float64(0),
        },
    }
    raw, _ := codec.Encode(req)
    fmt.Printf("SDO read frame: %X\n", raw)

    // NMT start remote node
    req2 := &kernel.Request{Function: "nmt_start"}
    raw2, _ := codec.Encode(req2)
    fmt.Printf("NMT start frame: %X\n", raw2)
}
```

## সমর্থিত ফাংশন

| ফাংশন     | বিবরণ                          |
|-------------|--------------------------------------|
| `sdo_read`  | অবজেক্ট ডিকশনারি এন্ট্রি পড়া         |
| `sdo_write` | অবজেক্ট ডিকশনারি এন্ট্রি লেখা        |
| `nmt_start` | রিমোট নোড চালু করা (NMT)              |
| `nmt_stop`  | রিমোট নোড বন্ধ করা (NMT)               |
| `nmt_reset` | রিমোট নোড রিসেট করা (NMT)              |
| `heartbeat` | heartbeat / bootup মেসেজ পাঠানো      |

## টেস্টিং

```bash
go test ./... -v
```

দ্রষ্টব্য: SocketCAN ড্রাইভার টেস্টের জন্য CAN হার্ডওয়্যার বা ভার্চুয়াল CAN ইন্টারফেসসহ একটি Linux সিস্টেম প্রয়োজন।
