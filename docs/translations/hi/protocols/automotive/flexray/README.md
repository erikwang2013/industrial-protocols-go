# FlexRay CAN प्रोटोकॉल SDK

FlexRay एक उच्च-गति, नियतात्मक ऑटोमोटिव संचार प्रोटोकॉल है। यह पैकेज CAN बस पर FlexRay कोडेक प्रदान करता है, जो CRC-16/XMODEM अखंडता जाँच के साथ चक्र-आधारित फ्रेमिंग का समर्थन करता है।

## प्रोटोकॉल अवलोकन

FlexRay बार-बार होने वाले संचार चक्रों के साथ टाइम-डिवीज़न मल्टीपल एक्सेस (TDMA) योजना का उपयोग करता है। प्रत्येक चक्र में स्थिर और गतिशील खंड (static and dynamic segments) होते हैं। यह कोडेक FlexRay फ्रेम को विस्तारित 29-बिट CAN फ्रेम पर मैप करता है।

### वायर फॉर्मेट

FlexRay चक्र पेलोड:
- **हेडर** (2 बाइट): चक्र संख्या (little-endian)
- **स्थिति** (1 बाइट): बिट 7=PPI (Payload Preamble Indicator), बिट 6=NFI, बिट 5=SYF, बिट 4=SUF
- **डेटा** (N बाइट): पेलोड (अधिकतम 254 बाइट)
- **CRC** (2 बाइट): हेडर+स्थिति+डेटा पर CRC-16/XMODEM (little-endian)

CAN ID एन्कोडिंग (29-बिट विस्तारित):
- बिट्स 28-24: मैसेज प्रकार (0x01=फ्रेम, 0x02=स्थिति)
- बिट्स 23-10: आरक्षित
- बिट्स 15-10: स्लॉट ID (6 बिट)
- बिट्स 9-0: चक्र संख्या (10 बिट)

## हार्डवेयर आवश्यकताएँ

- SocketCAN समर्थन वाला **Linux** सिस्टम
- Vector VN7600/VN7640 या Bosch FlexRay-CAN एडाप्टर
- उचित टर्मिनेशन वाला FlexRay नेटवर्क (2.5V बायस)

### CAN इंटरफ़ेस सेट करना

```bash
sudo modprobe can
sudo modprobe can_raw
sudo ip link set can0 type can bitrate 500000
sudo ip link set up can0
```

## उपयोग

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/automotive/flexray"
)

func main() {
    p := flexray.New()
    codec, _ := p.NewCodec("can")

    // Send a FlexRay frame in cycle 5 with PPI set
    req := &kernel.Request{
        Function: "frame",
        Data:     []byte{0x42, 0x01},
        Metadata: map[string]any{
            "cycle": float64(5),
            "ppi":   true,
        },
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## समर्थित फ़ंक्शन

| फ़ंक्शन | विवरण |
|----------|------------------------------------------|
| `frame`  | FlexRay फ्रेम पेलोड भेजें |
| `status` | स्लॉट कॉन्फ़िगरेशन पूछें (slot, cycle) |

## परीक्षण

```bash
go test ./... -v
```
