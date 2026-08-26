# Interbus GatewayBridge SDK

Interbus प्रोटोकॉल TCP गेटवे ब्रिज के माध्यम से।

## हार्डवेयर

| गेटवे | इंटरफ़ेस | डिफ़ॉल्ट IP |
|---------|-----------|-------------|
| Phoenix Contact IBS S5 DSC/I-T | Interbus मास्टर कंट्रोलर | 192.168.0.60 |
| Phoenix Contact IBS PCI SC/I-T | PCI Interbus मास्टर | host-assigned |
| HMS Anybus Interbus | एम्बेडेड गेटवे | 192.168.0.61 |

## वायरिंग

- गेटवे पर 9-पिन D-SUB से Interbus रिमोट बस (इनकमिंग/आउटगोइंग)
- शील्ड दोनों सिरों पर FE से जुड़ी
- TCP ब्रिज के लिए गेटवे से ईथरनेट

## गेटवे IP कॉन्फ़िगरेशन

```go
driver, codec, err := interbus.NewGatewayDriver("192.168.0.60:2006")
handler, err := interbus.ReadyHandler("192.168.0.60:2006")
```
