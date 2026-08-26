# CC-Link IE Field GatewayBridge SDK

TCP গেটওয়ে ব্রিজের মাধ্যমে CC-Link IE Field প্রোটোকল।

## হার্ডওয়্যার

| গেটওয়ে | ইন্টারফেস | ডিফল্ট IP |
|---------|-----------|-------------|
| Mitsubishi MELSEC IQ-R | CC-Link IE Field মাস্টার | 192.168.3.40 |
| Mitsubishi RJ71GF11-T2 | CC-Link IE Field মডিউল | হোস্ট-নির্ধারিত |
| HMS Anybus CC-Link IE | এমবেডেড গেটওয়ে | 192.168.0.52 |

## ওয়্যারিং

- CC-Link IE Field-এর জন্য RJ45 ইথারনেট (1 Gbps রিং বা স্টার টপোলজি)
- ম্যানেজমেন্ট পোর্ট আলাদা নেটওয়ার্কে
- গেটওয়ে ফিল্ড নেটওয়ার্ককে TCP ব্রিজে সংযুক্ত করে

## গেটওয়ে IP কনফিগারেশন

```go
driver, codec, err := cclinkie.NewGatewayDriver("192.168.3.40:2005")
handler, err := cclinkie.ReadyHandler("192.168.3.40:2005")
```
