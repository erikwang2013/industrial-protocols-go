# SERCOS I/II CmdBridge SDK

SERCOS I/II डिजिटल मोशन कंट्रोल के लिए SERCOS इंटरफ़ेस का लीगेसी सीरियल फाइबर ऑप्टिक संस्करण है। यह पैकेज `sercos_cli` कमांड-लाइन उपयोगिता के माध्यम से SERCOS I/II कोडेक प्रदान करता है।

## CLI उपकरण

SERCOS फाइबर ऑप्टिक इंटरफ़ेस CLI उपयोगिता का उपयोग करता है।

### इंस्टॉलेशन

```bash
# SERCOS interface card driver and tools
# Refer to vendor documentation for the specific SERCOS master card
```

सत्यापन: `sercos_cli --help`

## उपयोग

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/sercos1"
)

func main() {
    p := sercos1.New()
    codec, _ := p.NewCodec("cmd")

    req := &kernel.Request{
        Function: "read",
        Address:  "0x0100",
        Count:    4,
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## ड्राइवर

```go
b, c, err := sercos1.NewCmdDriver("/usr/bin/sercos_cli")
```

## समर्थित फ़ंक्शन

| फ़ंक्शन | विवरण              |
|----------|--------------------------|
| `read`   | SERCOS IDN पढ़ें          |
| `write`  | SERCOS IDN लिखें         |
| `status` | ड्राइव स्थिति पूछें       |

## परीक्षण

```bash
go test ./... -v
```
