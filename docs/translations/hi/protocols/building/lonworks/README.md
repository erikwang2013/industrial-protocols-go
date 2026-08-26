# LonWorks GatewayBridge SDK

LonWorks (ANSI/CEA-709.1) प्रोटोकॉल TCP गेटवे ब्रिज के माध्यम से।

## हार्डवेयर

| गेटवे | इंटरफ़ेस | डिफ़ॉल्ट IP |
|---------|-----------|-------------|
| Echelon U60 | USB FT-10 नेटवर्क इंटरफ़ेस | host-assigned |
| Echelon U70 | USB TP/XF-1250 इंटरफ़ेस | host-assigned |
| Loytec L-IP | LonWorks/IP राउटर | 192.168.0.90 |

## वायरिंग

- FT-10 (Free Topology): पोलैरिटी-असंवेदनशील ट्विस्टेड पेयर, फ्री टोपोलॉजी में 500 मीटर तक
- TP/XF-1250: 105 ओम टर्मिनेटर के साथ बस टोपोलॉजी
- TCP ब्रिज के लिए L-IP राउटर पर ईथरनेट

## गेटवे IP कॉन्फ़िगरेशन

```go
driver, codec, err := lonworks.NewGatewayDriver("192.168.0.90:2009")
handler, err := lonworks.ReadyHandler("192.168.0.90:2009")
```
