# PROFIBUS GatewayBridge SDK

PROFIBUS DP/PA प्रोटोकॉल TCP गेटवे ब्रिज के माध्यम से।

## हार्डवेयर

| गेटवे | इंटरफ़ेस | डिफ़ॉल्ट IP |
|---------|-----------|-------------|
| Anybus Communicator | PROFIBUS DP-V1 स्लेव | 192.168.0.50 |
| Siemens CP 5611 प्रॉक्सी | PCI/PCIe PROFIBUS मास्टर | host-assigned |
| HMS Fieldbus Gateway | Anybus NP40 | 192.168.0.51 |

## वायरिंग

- गेटवे पर DB9 फीमेल से PROFIBUS नेटवर्क (A-लाइन हरी, B-लाइन लाल)
- नेटवर्क के दोनों सिरों पर टर्मिनेशन रेसिस्टर ON
- गेटवे प्रबंधन पोर्ट से ईथरनेट

## गेटवे IP कॉन्फ़िगरेशन

```go
driver, codec, err := profibus.NewGatewayDriver("192.168.0.50:2000")
handler, err := profibus.ReadyHandler("192.168.0.50:2000")
```
