# ISA100.11a CmdBridge SDK

ISA100.11a प्रक्रिया स्वचालन के लिए एक वायरलेस औद्योगिक नेटवर्किंग मानक है। यह पैकेज `yfgw410_cli` (Yokogawa YFGW410 फील्ड वायरलेस गेटवे) कमांड-लाइन उपयोगिता के माध्यम से ISA100.11a कोडेक प्रदान करता है।

## CLI उपकरण

Yokogawa YFGW410 Field Wireless Gateway CLI का उपयोग करता है।

### इंस्टॉलेशन

```bash
# Install Yokogawa YFGW410 gateway software and tools
# Refer to Yokogawa documentation for Field Wireless Gateway setup
```

सत्यापन: `yfgw410_cli --help`

## उपयोग

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/isa100"
)

func main() {
    p := isa100.New()
    codec, _ := p.NewCodec("cmd")

    req := &kernel.Request{
        Function: "read",
        Address:  "DEV001",
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## ड्राइवर

```go
b, c, err := isa100.NewCmdDriver("/usr/bin/yfgw410_cli")
```

## समर्थित फ़ंक्शन

| फ़ंक्शन | विवरण                 |
|----------|-----------------------------|
| `read`   | डिवाइस विशेषता पढ़ें       |
| `write`  | डिवाइस विशेषता लिखें      |
| `list`   | प्रावधानित डिवाइस सूचीबद्ध करें    |

## परीक्षण

```bash
go test ./... -v
```
