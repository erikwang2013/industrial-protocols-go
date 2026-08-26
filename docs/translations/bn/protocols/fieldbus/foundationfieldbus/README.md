# Foundation Fieldbus GatewayBridge SDK

TCP গেটওয়ে ব্রিজের মাধ্যমে Foundation Fieldbus H1/HSE।

## হার্ডওয়্যার

| গেটওয়ে | ইন্টারফেস | ডিফল্ট IP |
|---------|-----------|-------------|
| NI USB-8486 | USB H1 ইন্টারফেস | হোস্ট-নির্ধারিত |
| Softing FFusb | USB H1 ইন্টারফেস | হোস্ট-নির্ধারিত |
| P+F HD2-GTR-4PA | H1 থেকে ইথারনেট গেটওয়ে | 192.168.0.20 |

## ওয়্যারিং

- H1 ট্রাঙ্ক (টুইস্টেড-পেয়ার, শিল্ডেড), দুই প্রান্তে টার্মিনেটরসহ
- H1-এর জন্য ফিল্ডবাস পাওয়ার কন্ডিশনার (24 VDC, প্রতি সেগমেন্টে 350-500 mA)
- গেটওয়ে ইথারনেটের মাধ্যমে H1 সেগমেন্টকে TCP ব্রিজে সংযুক্ত করে

## গেটওয়ে IP কনফিগারেশন

```go
driver, codec, err := foundationfieldbus.NewGatewayDriver("192.168.0.20:2002")
handler, err := foundationfieldbus.ReadyHandler("192.168.0.20:2002")
```
