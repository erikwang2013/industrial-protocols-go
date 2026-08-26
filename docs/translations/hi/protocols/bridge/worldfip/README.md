# WorldFIP GatewayBridge SDK

WorldFIP प्रोटोकॉल TCP गेटवे ब्रिज के माध्यम से।

## हार्डवेयर

| गेटवे | इंटरफ़ेस | डिफ़ॉल्ट IP |
|---------|-----------|-------------|
| FIPIO Agent | WorldFIP फील्डबस एजेंट | 192.168.0.70 |
| FIP Gateway (Alstom) | WorldFIP से ईथरनेट | 192.168.0.71 |
| NI FIP-USB | USB WorldFIP इंटरफ़ेस | host-assigned |

## वायरिंग

- गेटवे पर 9-पिन D-SUB से WorldFIP ट्रंक (FIP1 = Data+, FIP2 = Data-)
- बस के दोनों सिरों पर लाइन टर्मिनेटर (120 ओम)
- TCP ब्रिज के लिए ईथरनेट

## गेटवे IP कॉन्फ़िगरेशन

```go
driver, codec, err := worldfip.NewGatewayDriver("192.168.0.70:2007")
handler, err := worldfip.ReadyHandler("192.168.0.70:2007")
```
