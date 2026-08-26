# SAE J1850 CAN प्रोटोकॉल SDK

SAE J1850 एक वाहन संचार मानक है जिसका उपयोग ऑन-बोर्ड डायग्नोस्टिक्स (OBD-II) के लिए किया जाता है। यह पैकेज CAN बस (ISO 15765-4 / CAN TP) पर J1850 कोडेक प्रदान करता है।

## प्रोटोकॉल अवलोकन

CAN पर SAE J1850 OBD-II निम्नलिखित संरचना के साथ 29-बिट विस्तारित CAN पहचानकर्ताओं का उपयोग करता है:

### CAN ID फॉर्मेट (29-बिट)

| बिट्स       | फ़ील्ड    | विवरण                        |
|-----------|----------|------------------------------------|
| 28-26     | Priority | मैसेज प्राथमिकता (0-7, डिफ़ॉल्ट 6)  |
| 25        | Ext ID   | विस्तारित फ्रेम के लिए हमेशा 1       |
| 24-16     | PF       | Parameter Format (हेडर)          |
| 15-8      | PS       | Parameter Specific (लक्ष्य/स्रोत) |
| 7-0       | SA       | स्रोत पता                     |

### मानक OBD-II CAN ID

| प्रकार                | CAN ID (hex)    | विवरण                    |
|--------------------|-----------------|--------------------------------|
| Physical Request   | 0x18DAxxF1      | विशिष्ट ECU को अनुरोध (xx=addr) |
| Physical Response  | 0x18DAF1xx      | ECU से प्रतिक्रिया (xx=addr)   |
| Functional Request | 0x18DB33F1      | सभी ECU को प्रसारण         |

### ISO 15765-2 फ्रेम फॉर्मेट

सिंगल फ्रेम: बाइट 0 का ऊपरी निबल = डेटा लंबाई (0-7), निचला निबल + शेष बाइट = डायग्नोस्टिक डेटा।

## हार्डवेयर आवश्यकताएँ

- SocketCAN समर्थन वाला **Linux** सिस्टम
- J1850-CAN OBD-II एडाप्टर (जैसे ELM327-संगत USB-to-CAN, OBDLink SX, Macchina M2)
- OBD-II कनेक्टर वाला वाहन (1996+ अधिकांश वाहन)

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
    "github.com/erikwang2013/industrial-protocols-go/protocols/automotive/saej1850"
)

func main() {
    p := saej1850.New()
    codec, _ := p.NewCodec("can")

    // Mode $01 PID $0C: Engine RPM
    req := &kernel.Request{
        Function: "mode01",
        Metadata: map[string]any{"pid": float64(0x0C)},
    }
    raw, _ := codec.Encode(req)
    _ = raw

    // Mode $03: Request emission-related DTCs
    req2 := &kernel.Request{Function: "mode03"}
    raw2, _ := codec.Encode(req2)
    _ = raw2
}
```

## समर्थित फ़ंक्शन

| फ़ंक्शन        | OBD-II मोड | विवरण                         |
|----------------|-------------|-------------------------------------|
| `mode01`       | $01         | वर्तमान पॉवरट्रेन डेटा का अनुरोध करें |
| `mode03`       | $03         | उत्सर्जन-संबंधित DTC का अनुरोध करें |
| `mode0A`       | $0A         | स्थायी DTC का अनुरोध करें          |
| `diag_request` | कस्टम      | सामान्य डायग्नोस्टिक अनुरोध          |
| `diag_response`| -           | डायग्नोस्टिक प्रतिक्रिया फ्रेम           |
| `broadcast`    | -           | सभी ECU को फंक्शनल प्रसारण    |

## परीक्षण

```bash
go test ./... -v
```
