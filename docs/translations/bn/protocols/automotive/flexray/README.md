# FlexRay CAN প্রোটোকল SDK

FlexRay একটি উচ্চ-গতির ডিটারমিনিস্টিক অটোমোটিভ কমিউনিকেশন প্রোটোকল। এই প্যাকেজটি CAN বাসের উপর FlexRay কোডেক সরবরাহ করে, যা CRC-16/XMODEM ইন্টিগ্রিটি চেকসহ সাইকেল-ভিত্তিক ফ্রেমিং সমর্থন করে।

## প্রোটোকল ওভারভিউ

FlexRay পুনরাবৃত্ত কমিউনিকেশন সাইকেলসহ একটি টাইম-ডিভিশন মাল্টিপল অ্যাক্সেস (TDMA) স্কিম ব্যবহার করে। প্রতিটি সাইকেল স্ট্যাটিক ও ডাইনামিক সেগমেন্ট নিয়ে গঠিত। এই কোডেক FlexRay ফ্রেমগুলোকে এক্সটেন্ডেড ২৯-বিট CAN ফ্রেমে ম্যাপ করে।

### ওয়্যার ফরম্যাট

FlexRay সাইকেল পেলোড:
- **Header** (২ বাইট): সাইকেল নম্বর (লিটল-এন্ডিয়ান)
- **Status** (১ বাইট): বিট 7=PPI (Payload Preamble Indicator), বিট 6=NFI, বিট 5=SYF, বিট 4=SUF
- **Data** (N বাইট): পেলোড (সর্বোচ্চ ২৫৪ বাইট)
- **CRC** (২ বাইট): header+status+data এর উপর CRC-16/XMODEM (লিটল-এন্ডিয়ান)

CAN ID এনকোডিং (২৯-বিট এক্সটেন্ডেড):
- বিট 28-24: মেসেজ টাইপ (0x01=ফ্রেম, 0x02=স্ট্যাটাস)
- বিট 23-10: রিজার্ভড
- বিট 15-10: স্লট ID (৬ বিট)
- বিট 9-0: সাইকেল নম্বর (১০ বিট)

## হার্ডওয়্যার প্রয়োজনীয়তা

- SocketCAN সমর্থনসহ **Linux** সিস্টেম
- Vector VN7600/VN7640 বা Bosch FlexRay-CAN অ্যাডাপ্টার
- সঠিক টার্মিনেশনসহ FlexRay নেটওয়ার্ক (2.5V বায়াস)

### CAN ইন্টারফেস সেটআপ

```bash
sudo modprobe can
sudo modprobe can_raw
sudo ip link set can0 type can bitrate 500000
sudo ip link set up can0
```

## ব্যবহার

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/automotive/flexray"
)

func main() {
    p := flexray.New()
    codec, _ := p.NewCodec("can")

    // Send a FlexRay frame in cycle 5 with PPI set
    req := &kernel.Request{
        Function: "frame",
        Data:     []byte{0x42, 0x01},
        Metadata: map[string]any{
            "cycle": float64(5),
            "ppi":   true,
        },
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## সমর্থিত ফাংশন

| ফাংশন | বিবরণ                              |
|----------|------------------------------------------|
| `frame`  | FlexRay ফ্রেম পেলোড পাঠানো               |
| `status` | স্লট কনফিগারেশন (slot, cycle) জিজ্ঞাসা   |

## টেস্টিং

```bash
go test ./... -v
```
