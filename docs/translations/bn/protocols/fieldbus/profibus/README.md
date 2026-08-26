# PROFIBUS GatewayBridge SDK

TCP গেটওয়ে ব্রিজের মাধ্যমে PROFIBUS DP/PA প্রোটোকল।

## হার্ডওয়্যার

| গেটওয়ে | ইন্টারফেস | ডিফল্ট IP |
|---------|-----------|-------------|
| Anybus Communicator | PROFIBUS DP-V1 স্লেভ | 192.168.0.50 |
| Siemens CP 5611 proxy | PCI/PCIe PROFIBUS মাস্টার | হোস্ট-নির্ধারিত |
| HMS Fieldbus Gateway | Anybus NP40 | 192.168.0.51 |

## ওয়্যারিং

- গেটওয়ের DB9 ফিমেল থেকে PROFIBUS নেটওয়ার্কে (A-লাইন সবুজ, B-লাইন লাল)
- নেটওয়ার্কের দুই প্রান্তে টার্মিনেশন রেজিস্টর ON
- গেটওয়ে ম্যানেজমেন্ট পোর্টে ইথারনেট

## গেটওয়ে IP কনফিগারেশন

```go
driver, codec, err := profibus.NewGatewayDriver("192.168.0.50:2000")
handler, err := profibus.ReadyHandler("192.168.0.50:2000")
```
