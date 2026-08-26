# CC-Link IE Field GatewayBridge SDK

CC-Link IE Field प्रोटोकॉल TCP गेटवे ब्रिज के माध्यम से।

## हार्डवेयर

| गेटवे | इंटरफ़ेस | डिफ़ॉल्ट IP |
|---------|-----------|-------------|
| Mitsubishi MELSEC IQ-R | CC-Link IE Field मास्टर | 192.168.3.40 |
| Mitsubishi RJ71GF11-T2 | CC-Link IE Field मॉड्यूल | host-assigned |
| HMS Anybus CC-Link IE | एम्बेडेड गेटवे | 192.168.0.52 |

## वायरिंग

- CC-Link IE Field के लिए RJ45 ईथरनेट (1 Gbps रिंग या स्टार टोपोलॉजी)
- प्रबंधन पोर्ट अलग नेटवर्क पर
- गेटवे फील्ड नेटवर्क को TCP ब्रिज से जोड़ता है

## गेटवे IP कॉन्फ़िगरेशन

```go
driver, codec, err := cclinkie.NewGatewayDriver("192.168.3.40:2005")
handler, err := cclinkie.ReadyHandler("192.168.3.40:2005")
```
