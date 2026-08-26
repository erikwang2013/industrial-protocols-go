# SAE J1850 CAN প্রোটোকল SDK

SAE J1850 একটি ভেহিকেল কমিউনিকেশন স্ট্যান্ডার্ড, যা অন-বোর্ড ডায়াগনস্টিকস (OBD-II)-এ ব্যবহৃত হয়। এই প্যাকেজটি CAN বাসের (ISO 15765-4 / CAN TP) উপর J1850 কোডেক সরবরাহ করে।

## প্রোটোকল ওভারভিউ

CAN-এর উপর SAE J1850 OBD-II ২৯-বিট এক্সটেন্ডেড CAN আইডেন্টিফায়ার ব্যবহার করে, নিম্নলিখিত কাঠামোতে:

### CAN ID ফরম্যাট (২৯-বিট)

| বিট       | ফিল্ড    | বিবরণ                        |
|-----------|----------|------------------------------------|
| 28-26     | Priority | মেসেজ প্রায়োরিটি (0-7, ডিফল্ট 6)  |
| 25        | Ext ID   | এক্সটেন্ডেড ফ্রেমের জন্য সর্বদা 1       |
| 24-16     | PF       | Parameter Format (Header)          |
| 15-8      | PS       | Parameter Specific (টার্গেট/সোর্স) |
| 7-0       | SA       | সোর্স অ্যাড্রেস                     |

### স্ট্যান্ডার্ড OBD-II CAN ID

| টাইপ                | CAN ID (হেক্স)    | বিবরণ                    |
|--------------------|-----------------|--------------------------------|
| Physical Request   | 0x18DAxxF1      | নির্দিষ্ট ECU-তে রিকোয়েস্ট (xx=addr) |
| Physical Response  | 0x18DAF1xx      | ECU থেকে রেসপন্স (xx=addr)   |
| Functional Request | 0x18DB33F1      | সব ECU-তে ব্রডকাস্ট         |

### ISO 15765-2 ফ্রেম ফরম্যাট

সিঙ্গেল ফ্রেম: বাইট 0-এর উচ্চ নিবল = ডেটা দৈর্ঘ্য (0-7), নিম্ন নিবল + বাকি বাইট = ডায়াগনস্টিক ডেটা।

## হার্ডওয়্যার প্রয়োজনীয়তা

- SocketCAN সমর্থনসহ **Linux** সিস্টেম
- J1850-CAN OBD-II অ্যাডাপ্টার (যেমন ELM327-সামঞ্জস্যপূর্ণ USB-টু-CAN, OBDLink SX, Macchina M2)
- OBD-II কানেক্টরসহ গাড়ি (১৯৯৬+ অধিকাংশ গাড়ি)

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
    "github.com/erikwang2013/industrial-protocols-go/protocols/automotive/saej1850"
)

func main() {
    p := saej1850.New()
    codec, _ := p.NewCodec("can")

    // Mode $01 PID $0C: Engine RPM
    req := &kernel.Request{
        Function: "mode01",
        Metadata: map[string]any{"pid": float64(0x0C)},
    }
    raw, _ := codec.Encode(req)
    _ = raw

    // Mode $03: Request emission-related DTCs
    req2 := &kernel.Request{Function: "mode03"}
    raw2, _ := codec.Encode(req2)
    _ = raw2
}
```

## সমর্থিত ফাংশন

| ফাংশন        | OBD-II মোড | বিবরণ                         |
|----------------|-------------|-------------------------------------|
| `mode01`       | $01         | বর্তমান পাওয়ারট্রেন ডেটা রিকোয়েস্ট     |
| `mode03`       | $03         | নির্গমন-সম্পর্কিত DTC রিকোয়েস্ট       |
| `mode0A`       | $0A         | স্থায়ী DTC রিকোয়েস্ট              |
| `diag_request` | কাস্টম      | জেনেরিক ডায়াগনস্টিক রিকোয়েস্ট          |
| `diag_response`| -           | ডায়াগনস্টিক রেসপন্স ফ্রেম           |
| `broadcast`    | -           | সব ECU-তে ফাংশনাল ব্রডকাস্ট    |

## টেস্টিং

```bash
go test ./... -v
```
