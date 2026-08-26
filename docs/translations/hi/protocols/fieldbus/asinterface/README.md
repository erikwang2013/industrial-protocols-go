# AS-Interface GatewayBridge SDK

AS-Interface (ASi) प्रोटोकॉल TCP गेटवे ब्रिज के माध्यम से।

## हार्डवेयर

| गेटवे | इंटरफ़ेस | डिफ़ॉल्ट IP |
|---------|-----------|-------------|
| Bihl+Wiedemann BWU2044 | ASi मास्टर (Ethernet) | 192.168.0.30 |
| P+F VBG-ENX-K20-DMD-EV | ASi गेटवे | 192.168.0.31 |
| ifm AC1375 | ASi ControllerE | 192.168.0.32 |

## वायरिंग

- गेटवे से स्लेव तक पीली ASi केबल (पावर + डेटा)
- काली सहायक पावर केबल (एक्चुएटर्स के लिए 24 VDC) वैकल्पिक
- TCP ब्रिज के लिए गेटवे पर ईथरनेट

## गेटवे IP कॉन्फ़िगरेशन

```go
driver, codec, err := asinterface.NewGatewayDriver("192.168.0.30:2003")
handler, err := asinterface.ReadyHandler("192.168.0.30:2003")
```
