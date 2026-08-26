# Interbus GatewayBridge SDK

TCP গেটওয়ে ব্রিজের মাধ্যমে Interbus প্রোটোকল।

## হার্ডওয়্যার

| গেটওয়ে | ইন্টারফেস | ডিফল্ট IP |
|---------|-----------|-------------|
| Phoenix Contact IBS S5 DSC/I-T | Interbus মাস্টার কন্ট্রোলার | 192.168.0.60 |
| Phoenix Contact IBS PCI SC/I-T | PCI Interbus মাস্টার | হোস্ট-নির্ধারিত |
| HMS Anybus Interbus | এমবেডেড গেটওয়ে | 192.168.0.61 |

## ওয়্যারিং

- গেটওয়ের 9-পিন D-SUB থেকে Interbus রিমোট বাসে (ইনকামিং/আউটগোয়িং)
- শিল্ড দুই প্রান্তে FE-তে সংযুক্ত
- TCP ব্রিজের জন্য গেটওয়েতে ইথারনেট

## গেটওয়ে IP কনফিগারেশন

```go
driver, codec, err := interbus.NewGatewayDriver("192.168.0.60:2006")
handler, err := interbus.ReadyHandler("192.168.0.60:2006")
```
