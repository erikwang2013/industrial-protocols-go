# WorldFIP GatewayBridge SDK

TCP গেটওয়ে ব্রিজের মাধ্যমে WorldFIP প্রোটোকল।

## হার্ডওয়্যার

| গেটওয়ে | ইন্টারফেস | ডিফল্ট IP |
|---------|-----------|-------------|
| FIPIO Agent | WorldFIP ফিল্ডবাস এজেন্ট | 192.168.0.70 |
| FIP Gateway (Alstom) | WorldFIP থেকে ইথারনেট | 192.168.0.71 |
| NI FIP-USB | USB WorldFIP ইন্টারফেস | হোস্ট-নির্ধারিত |

## ওয়্যারিং

- গেটওয়ের 9-পিন D-SUB থেকে WorldFIP ট্রাঙ্কে (FIP1 = Data+, FIP2 = Data-)
- বাসের দুই প্রান্তে লাইন টার্মিনেটর (১২০ ওহম)
- TCP ব্রিজের জন্য ইথারনেট

## গেটওয়ে IP কনফিগারেশন

```go
driver, codec, err := worldfip.NewGatewayDriver("192.168.0.70:2007")
handler, err := worldfip.ReadyHandler("192.168.0.70:2007")
```
