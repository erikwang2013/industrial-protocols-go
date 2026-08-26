# Lightbus GatewayBridge SDK

TCP গেটওয়ে ব্রিজের মাধ্যমে Beckhoff Lightbus ফাইবার-অপটিক প্রোটোকল।

## হার্ডওয়্যার

| গেটওয়ে | ইন্টারফেস | ডিফল্ট IP |
|---------|-----------|-------------|
| Beckhoff FC2001 | Lightbus PCI কার্ড | হোস্ট-নির্ধারিত |
| Beckhoff BK2000 | Lightbus বাস কাপলার | 192.168.0.80 |
| Beckhoff FC9001 | Lightbus ইথারনেট অ্যাডাপ্টার | 192.168.0.81 |

## ওয়্যারিং

- প্লাস্টিক অপটিক্যাল ফাইবার (POF) রিং টপোলজি
- FC2001/FC9001 ইথারনেটের মাধ্যমে রিংকে TCP ব্রিজে সংযুক্ত করে
- প্রতিটি ডিভাইসে TX ও RX ফাইবার কানেক্টর থাকে

## গেটওয়ে IP কনফিগারেশন

```go
driver, codec, err := lightbus.NewGatewayDriver("192.168.0.80:2008")
handler, err := lightbus.ReadyHandler("192.168.0.80:2008")
```
