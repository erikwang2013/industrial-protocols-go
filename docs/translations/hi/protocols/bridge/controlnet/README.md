# ControlNet CmdBridge SDK

ControlNet एक रीयल-टाइम औद्योगिक नेटवर्क प्रोटोकॉल है जिसे Allen-Bradley (Rockwell Automation) ने उच्च-गति, समय-महत्वपूर्ण डेटा विनिमय के लिए विकसित किया है। यह पैकेज `1784-pcic-cli` कमांड-लाइन उपयोगिता के माध्यम से ControlNet कोडेक प्रदान करता है।

## CLI उपकरण

1784-PCIC ControlNet इंटरफ़ेस कार्ड CLI उपयोगिता का उपयोग करता है।

### इंस्टॉलेशन

```bash
# Install Rockwell 1784-PCIC driver and tools
# Refer to Rockwell Automation documentation for RSLinx Classic SDK
```

सत्यापन: `1784-pcic-cli --help`

## उपयोग

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/controlnet"
)

func main() {
    p := controlnet.New()
    codec, _ := p.NewCodec("cmd")

    req := &kernel.Request{
        Function: "read",
        Address:  "0x10",
        Count:    8,
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## ड्राइवर

```go
b, c, err := controlnet.NewCmdDriver("/usr/bin/1784-pcic-cli")
```

## समर्थित फ़ंक्शन

| फ़ंक्शन | विवरण               |
|----------|---------------------------|
| `read`   | ControlNet नोड से पढ़ें |
| `write`  | ControlNet नोड पर लिखें |
| `status` | PCIC कार्ड स्थिति पूछें |

## परीक्षण

```bash
go test ./... -v
```
