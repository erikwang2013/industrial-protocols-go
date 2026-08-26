# MOST Serial प्रोटोकॉल SDK

MOST (Media Oriented Systems Transport) एक उच्च-गति मल्टीमीडिया नेटवर्क तकनीक है जिसका उपयोग मुख्य रूप से ऑटोमोटिव इंफोटेनमेंट सिस्टम में होता है। यह पैकेज AT-कमांड इंटरफ़ेस वाले सीरियल एडाप्टर पर MOST कोडेक प्रदान करता है।

## प्रोटोकॉल अवलोकन

MOST फाइबर ऑप्टिक फिजिकल लेयर पर सिंक्रोनस सीरियल संचार का उपयोग करता है। यह कार्यान्वयन एक सीरियल एडाप्टर के माध्यम से जुड़ता है जो 115200 बॉड पर AT-कमांड इंटरफ़ेस उजागर करता है।

### AT कमांड

| कमांड            | विवरण                  |
|--------------------|------------------------------|
| `AT+READ=<addr>,<len>`  | पते से बाइट्स पढ़ें |
| `AT+WRITE=<addr>,<hex>` | पते पर हेक्स बाइट्स लिखें |
| `AT+STATUS`             | रिंग/नेटवर्क स्थिति पूछें |

### प्रतिक्रिया फॉर्मेट

- `+OK:<hex_data>` -- सफल प्रतिक्रिया
- `+ERR:<code>` -- त्रुटि प्रतिक्रिया

## हार्डवेयर आवश्यकताएँ

- उचित टर्मिनेशन वाला MOST फाइबर ऑप्टिक नेटवर्क
- MOST-से-सीरियल एडाप्टर (जैसे MOST150 USB एडाप्टर)
- सीरियल पोर्ट 115200 बॉड, 8N1

## उपयोग

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/automotive/most"
)

func main() {
    p := most.New()
    codec, _ := p.NewCodec("serial")

    // Read 4 bytes from address 0x0100
    req := &kernel.Request{
        Function: "read",
        Address:  "0x0100",
        Count:    4,
    }
    raw, _ := codec.Encode(req)
    // raw = "AT+READ=0x0100,4\r\n"
    _ = raw

    // Write data to address 0x0200
    req2 := &kernel.Request{
        Function: "write",
        Address:  "0x0200",
        Data:     []byte{0xDE, 0xAD, 0xBE, 0xEF},
    }
    raw2, _ := codec.Encode(req2)
    // raw2 = "AT+WRITE=0x0200,DEADBEEF\r\n"
    _ = raw2
}
```

## ड्राइवर

```go
b, c, err := most.NewSerialDriver("/dev/ttyUSB0")
```

## समर्थित फ़ंक्शन

| फ़ंक्शन | विवरण |
|----------|-----------------------------|
| `read`   | MOST पते से पढ़ें |
| `write`  | MOST पते पर डेटा लिखें |
| `status` | रिंग/नेटवर्क स्थिति पूछें |

## परीक्षण

```bash
go test ./... -v
```
