# WirelessHART CmdBridge SDK

WirelessHART (IEC 62591) HART प्रोटोकॉल पर आधारित एक वायरलेस औद्योगिक नेटवर्किंग मानक है। यह पैकेज `emerson_1410_cli` (Emerson 1410/1420 Wireless Gateway) कमांड-लाइन उपयोगिता के माध्यम से WirelessHART कोडेक प्रदान करता है।

## CLI उपकरण

Emerson 1410/1420 Wireless Gateway CLI उपयोगिता का उपयोग करता है।

### इंस्टॉलेशन

```bash
# Install Emerson Wireless Gateway software and tools
# Refer to Emerson documentation for 1410/1420 Gateway setup
```

सत्यापन: `emerson_1410_cli --help`

## उपयोग

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/wirelesshart"
)

func main() {
    p := wirelesshart.New()
    codec, _ := p.NewCodec("cmd")

    req := &kernel.Request{
        Function: "read",
        Address:  "TT101",
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## ड्राइवर

```go
b, c, err := wirelesshart.NewCmdDriver("/usr/bin/emerson_1410_cli")
```

## समर्थित फ़ंक्शन

| फ़ंक्शन | विवरण                |
|----------|----------------------------|
| `read`   | डिवाइस पैरामीटर पढ़ें      |
| `write`  | डिवाइस पैरामीटर लिखें     |
| `scan`   | वायरलेस डिवाइस खोजें  |

## परीक्षण

```bash
go test ./... -v
```
