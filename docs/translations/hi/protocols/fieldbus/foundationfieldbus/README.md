# Foundation Fieldbus GatewayBridge SDK

Foundation Fieldbus H1/HSE TCP गेटवे ब्रिज के माध्यम से।

## हार्डवेयर

| गेटवे | इंटरफ़ेस | डिफ़ॉल्ट IP |
|---------|-----------|-------------|
| NI USB-8486 | USB H1 इंटरफ़ेस | host-assigned |
| Softing FFusb | USB H1 इंटरफ़ेस | host-assigned |
| P+F HD2-GTR-4PA | H1 से ईथरनेट गेटवे | 192.168.0.20 |

## वायरिंग

- H1 ट्रंक (ट्विस्टेड-पेयर, शील्डेड) दोनों सिरों पर टर्मिनेटर के साथ
- H1 के लिए फील्डबस पावर कंडीशनर (24 VDC, प्रति सेगमेंट 350-500 mA)
- गेटवे H1 सेगमेंट को ईथरनेट के माध्यम से TCP ब्रिज से जोड़ता है

## गेटवे IP कॉन्फ़िगरेशन

```go
driver, codec, err := foundationfieldbus.NewGatewayDriver("192.168.0.20:2002")
handler, err := foundationfieldbus.ReadyHandler("192.168.0.20:2002")
```
