# IO-Link GatewayBridge SDK

IO-Link प्रोटोकॉल TCP गेटवे ब्रिज के माध्यम से।

## हार्डवेयर

| गेटवे | इंटरफ़ेस | डिफ़ॉल्ट IP |
|---------|-----------|-------------|
| ifm AL1332 | IO-Link मास्टर (EtherNet/IP) | 192.168.0.40 |
| Balluff BNI00AZ | IO-Link मास्टर (PROFINET) | 192.168.0.41 |
| SICK SIG200 | IO-Link मास्टर (Ethernet) | 192.168.0.42 |

## वायरिंग

- प्रत्येक IO-Link पोर्ट के लिए M12 कनेक्टर (4-पिन): L+ (भूरा), L- (नीला), C/Q (काला), अनुपयोगी (सफेद)
- मास्टर और डिवाइस के लिए 24 VDC पावर सप्लाई
- TCP ब्रिज के लिए मास्टर पर ईथरनेट

## गेटवे IP कॉन्फ़िगरेशन

```go
driver, codec, err := iolink.NewGatewayDriver("192.168.0.40:2004")
handler, err := iolink.ReadyHandler("192.168.0.40:2004")
```
