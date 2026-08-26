# DeviceNet प्रोटोकॉल SDK

DeviceNet फैक्ट्री स्वचालन के लिए CAN-आधारित औद्योगिक नेटवर्क प्रोटोकॉल है। यह पैकेज SocketCAN और TCP गेटवे दोनों ड्राइवरों के साथ DeviceNet कोडेक प्रदान करता है।

## प्रोटोकॉल अवलोकन

DeviceNet CAN पर Common Industrial Protocol (CIP) का उपयोग करता है। पूर्वनिर्धारित मास्टर/स्लेव कनेक्शन सेट, पोल किए गए I/O और एक्सप्लिसिट मैसेजिंग के लिए Group 2 मैसेज (CAN ID 0x400 + NodeID) का उपयोग करता है।

## हार्डवेयर आवश्यकताएँ

### CAN मोड (SocketCAN)
- SocketCAN समर्थन वाला **Linux** सिस्टम (`CONFIG_CAN` सक्षम)
- एक CAN इंटरफ़ेस (जैसे `can0`, `vcan0`)
- DeviceNet-क्षम CAN हार्डवेयर (जैसे Anybus Communicator, HMS IXXAT)

### गेटवे मोड
- DeviceNet गेटवे के लिए TCP/IP कनेक्टिविटी
- प्लेन-टेक्स्ट कमांड प्रोटोकॉल समर्थन करने वाला गेटवे (जैसे HMS Anybus, Hilscher netX)

### वर्चुअल CAN इंटरफ़ेस सेट करना (परीक्षण के लिए)

```bash
sudo modprobe can
sudo modprobe can_raw
sudo modprobe vcan
sudo ip link add dev vcan0 type vcan
sudo ip link set up vcan0
```

## उपयोग

### CAN मोड

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/fieldbus/devicenet"
)

func main() {
    p := devicenet.New()
    codec, _ := p.NewCodec("can")

    // Poll request
    req := &kernel.Request{
        Function: "poll",
        Data:     []byte{0x01, 0x00},
    }
    raw, _ := codec.Encode(req)

    // Open connection
    req2 := &kernel.Request{Function: "open"}
    raw2, _ := codec.Encode(req2)
    _ = raw
    _ = raw2
}
```

### गेटवे मोड

```go
p := devicenet.New()
codec, _ := p.NewCodec("gateway")

req := &kernel.Request{
    Function: "poll",
    Data:     []byte{0xAB, 0xCD},
}
raw, _ := codec.Encode(req)
// raw will be: "poll abcd\n"
```

## समर्थित फ़ंक्शन

| फ़ंक्शन | विवरण                    |
|----------|--------------------------------|
| `open`   | एक्सप्लिसिट कनेक्शन खोलें       |
| `poll`   | I/O डेटा पोल करें (Group 2)        |

## परीक्षण

```bash
go test ./... -v
```

नोट: SocketCAN ड्राइवर परीक्षणों के लिए Linux आवश्यक है। गेटवे ड्राइवर परीक्षण किसी भी OS पर चलते हैं।
