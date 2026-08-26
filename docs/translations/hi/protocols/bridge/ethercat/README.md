# EtherCAT CmdBridge SDK

EtherCAT (Ethernet for Control Automation Technology) एक उच्च-प्रदर्शन औद्योगिक ईथरनेट फील्डबस है। यह पैकेज `ethercat` कमांड-लाइन उपयोगिता के माध्यम से EtherCAT कोडेक प्रदान करता है।

## CLI उपकरण

IgH EtherCAT Master कमांड-लाइन उपकरण का उपयोग करता है।

### इंस्टॉलेशन

```bash
# Debian/Ubuntu
sudo apt-get install ethercat-master

# From source
git clone https://gitlab.com/etherlab.org/ethercat.git
cd ethercat
make && sudo make install
```

सत्यापन: `ethercat slaves`

## उपयोग

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/ethercat"
)

func main() {
    p := ethercat.New()
    codec, _ := p.NewCodec("cmd")

    // Upload SDO from address 0x1000
    req := &kernel.Request{
        Function: "upload",
        Address:  "0x1000",
        Count:    4,
    }
    raw, _ := codec.Encode(req)
    // raw = "upload 0x1000 4\n"
    _ = raw
}
```

## ड्राइवर

```go
b, c, err := ethercat.NewCmdDriver("/usr/bin/ethercat")
```

## समर्थित फ़ंक्शन

| फ़ंक्शन   | विवरण            |
|------------|------------------------|
| `upload`   | पते से SDO पढ़ें  |
| `download` | पते पर SDO लिखें   |
| `slaves`   | EtherCAT स्लेव सूचीबद्ध करें |

## परीक्षण

```bash
go test ./... -v
```
