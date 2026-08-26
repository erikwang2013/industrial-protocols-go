# AS-Interface GatewayBridge SDK

TCP গেটওয়ে ব্রিজের মাধ্যমে AS-Interface (ASi) প্রোটোকল।

## হার্ডওয়্যার

| গেটওয়ে | ইন্টারফেস | ডিফল্ট IP |
|---------|-----------|-------------|
| Bihl+Wiedemann BWU2044 | ASi মাস্টার (ইথারনেট) | 192.168.0.30 |
| P+F VBG-ENX-K20-DMD-EV | ASi গেটওয়ে | 192.168.0.31 |
| ifm AC1375 | ASi ControllerE | 192.168.0.32 |

## ওয়্যারিং

- গেটওয়ে থেকে স্লেভ পর্যন্ত হলুদ ASi ক্যাবল (পাওয়ার + ডেটা)
- ঐচ্ছিক কালো অক্সিলিয়ারি পাওয়ার ক্যাবল (অ্যাকচুয়েটরের জন্য 24 VDC)
- TCP ব্রিজের জন্য গেটওয়েতে ইথারনেট

## গেটওয়ে IP কনফিগারেশন

```go
driver, codec, err := asinterface.NewGatewayDriver("192.168.0.30:2003")
handler, err := asinterface.ReadyHandler("192.168.0.30:2003")
```
