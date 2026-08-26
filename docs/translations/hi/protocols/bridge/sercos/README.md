# SERCOS III CmdBridge SDK

SERCOS III (SErial Real-time COmmunication System) मोशन कंट्रोल के लिए एक डिजिटल इंटरफ़ेस है, जो फाइबर ऑप्टिक या कॉपर पर रिंग टोपोलॉजी का उपयोग करता है। यह पैकेज `netx_cli` कमांड-लाइन उपयोगिता के माध्यम से SERCOS III कोडेक प्रदान करता है।

## CLI उपकरण

Hilscher netX SERCOS III CLI उपयोगिता का उपयोग करता है।

### इंस्टॉलेशन

```bash
# Install Hilscher netX driver and tools
# See https://www.hilscher.com/ for netX drivers
sudo apt-get install netx-driver
```

सत्यापन: `netx_cli --help`

## उपयोग

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/sercos"
)

func main() {
    p := sercos.New()
    codec, _ := p.NewCodec("cmd")

    req := &kernel.Request{
        Function: "read",
        Address:  "S-0-51",
        Count:    2,
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## ड्राइवर

```go
b, c, err := sercos.NewCmdDriver("/usr/bin/netx_cli")
```

## समर्थित फ़ंक्शन

| फ़ंक्शन | विवरण                 |
|----------|-----------------------------|
| `read`   | SERCOS IDN/S पैरामीटर पढ़ें |
| `write`  | SERCOS IDN/S पैरामीटर लिखें |
| `phase`  | संचार फेज़ सेट करें     |

## परीक्षण

```bash
go test ./... -v
```
