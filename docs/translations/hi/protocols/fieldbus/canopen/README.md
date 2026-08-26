# CANopen प्रोटोकॉल SDK

CANopen एम्बेडेड कंट्रोल सिस्टम के लिए CAN-आधारित उच्च-स्तरीय प्रोटोकॉल है। यह पैकेज CANopen कोडेक और SocketCAN ड्राइवर प्रदान करता है।

## प्रोटोकॉल अवलोकन

CANopen निम्नलिखित पूर्वनिर्धारित कनेक्शन सेट के साथ मानक 11-बिट CAN पहचानकर्ताओं का उपयोग करता है:

| फ़ंक्शन    | CAN ID              | विवरण                   |
|------------|---------------------|-------------------------------|
| NMT        | 0x000               | नेटवर्क प्रबंधन            |
| SYNC       | 0x080               | सिंक्रोनाइज़ेशन मैसेज       |
| SDO (tx)   | 0x580 + NodeID      | Service Data Object (सर्वर)  |
| SDO (rx)   | 0x600 + NodeID      | Service Data Object (क्लाइंट)  |
| PDO1 (tx)  | 0x180 + NodeID      | Process Data Object 1         |
| Heartbeat  | 0x700 + NodeID      | Heartbeat / Bootup            |

## हार्डवेयर आवश्यकताएँ

- SocketCAN समर्थन वाला **Linux** सिस्टम (`CONFIG_CAN` सक्षम)
- एक CAN इंटरफ़ेस (जैसे `can0`, वर्चुअल CAN के लिए `vcan0`)
- CAN-क्षम ट्रांसीवर हार्डवेयर (जैसे MCP2515, SJA1000, या USB-CAN एडाप्टर)

### वर्चुअल CAN इंटरफ़ेस सेट करना (परीक्षण के लिए)

```bash
sudo modprobe can
sudo modprobe can_raw
sudo modprobe vcan
sudo ip link add dev vcan0 type vcan
sudo ip link set up vcan0
```

## उपयोग

```go
package main

import (
    "fmt"
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/fieldbus/canopen"
)

func main() {
    p := canopen.New()
    codec, _ := p.NewCodec("can")

    // Read SDO object dictionary entry
    req := &kernel.Request{
        Function: "sdo_read",
        Metadata: map[string]any{
            "index": float64(0x1000), // Device type
            "sub":   float64(0),
        },
    }
    raw, _ := codec.Encode(req)
    fmt.Printf("SDO read frame: %X\n", raw)

    // NMT start remote node
    req2 := &kernel.Request{Function: "nmt_start"}
    raw2, _ := codec.Encode(req2)
    fmt.Printf("NMT start frame: %X\n", raw2)
}
```

## समर्थित फ़ंक्शन

| फ़ंक्शन     | विवरण                          |
|-------------|--------------------------------------|
| `sdo_read`  | ऑब्जेक्ट डिक्शनरी प्रविष्टि पढ़ें         |
| `sdo_write` | ऑब्जेक्ट डिक्शनरी प्रविष्टि लिखें        |
| `nmt_start` | रिमोट नोड प्रारंभ करें (NMT)              |
| `nmt_stop`  | रिमोट नोड रोकें (NMT)               |
| `nmt_reset` | रिमोट नोड रीसेट करें (NMT)              |
| `heartbeat` | heartbeat / bootup मैसेज भेजें      |

## परीक्षण

```bash
go test ./... -v
```

नोट: SocketCAN ड्राइवर परीक्षणों के लिए CAN हार्डवेयर या वर्चुअल CAN इंटरफ़ेस वाला Linux सिस्टम आवश्यक है।
