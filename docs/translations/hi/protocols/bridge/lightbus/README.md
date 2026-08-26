# Lightbus GatewayBridge SDK

Beckhoff Lightbus फाइबर ऑप्टिक प्रोटोकॉल TCP गेटवे ब्रिज के माध्यम से।

## हार्डवेयर

| गेटवे | इंटरफ़ेस | डिफ़ॉल्ट IP |
|---------|-----------|-------------|
| Beckhoff FC2001 | Lightbus PCI कार्ड | host-assigned |
| Beckhoff BK2000 | Lightbus बस कपलर | 192.168.0.80 |
| Beckhoff FC9001 | Lightbus ईथरनेट एडाप्टर | 192.168.0.81 |

## वायरिंग

- प्लास्टिक ऑप्टिकल फाइबर (POF) रिंग टोपोलॉजी
- FC2001/FC9001 रिंग को ईथरनेट के माध्यम से TCP ब्रिज से जोड़ता है
- प्रत्येक डिवाइस में TX और RX फाइबर कनेक्टर होते हैं

## गेटवे IP कॉन्फ़िगरेशन

```go
driver, codec, err := lightbus.NewGatewayDriver("192.168.0.80:2008")
handler, err := lightbus.ReadyHandler("192.168.0.80:2008")
```
