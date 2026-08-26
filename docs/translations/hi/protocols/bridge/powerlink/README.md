# POWERLINK CmdBridge SDK

POWERLINK (Ethernet POWERLINK) औद्योगिक स्वचालन के लिए एक रीयल-टाइम ईथरनेट प्रोटोकॉल है। यह पैकेज `openPOWERLINK_demo` कमांड-लाइन उपयोगिता के माध्यम से POWERLINK कोडेक प्रदान करता है।

## CLI उपकरण

openPOWERLINK स्टैक डेमो एप्लिकेशन का उपयोग करता है।

### इंस्टॉलेशन

```bash
git clone https://github.com/OpenAutomationTechnologies/openPOWERLINK.git
cd openPOWERLINK
cmake -B build && cmake --build build
sudo cmake --install build
```

सत्यापन: `openPOWERLINK_demo --help`

## उपयोग

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/powerlink"
)

func main() {
    p := powerlink.New()
    codec, _ := p.NewCodec("cmd")

    req := &kernel.Request{
        Function: "read",
        Address:  "0x2000",
        Count:    8,
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## ड्राइवर

```go
b, c, err := powerlink.NewCmdDriver("/usr/bin/openPOWERLINK_demo")
```

## समर्थित फ़ंक्शन

| फ़ंक्शन | विवरण                   |
|----------|-------------------------------|
| `read`   | ऑब्जेक्ट डिक्शनरी प्रविष्टि पढ़ें  |
| `write`  | ऑब्जेक्ट डिक्शनरी प्रविष्टि लिखें |
| `status` | नोड/NMT स्थिति पूछें          |

## परीक्षण

```bash
go test ./... -v
```
