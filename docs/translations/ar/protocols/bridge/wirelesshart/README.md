# SDK جسر أوامر WirelessHART

WirelessHART (IEC 62591) هو معيار شبكات صناعية لاسلكية مبني على بروتوكول HART. توفر هذه الحزمة مُرمِّزًا (codec) لـ WirelessHART عبر أداة سطر الأوامر `emerson_1410_cli` (بوابة Emerson 1410/1420 اللاسلكية).

## أداة سطر الأوامر

تستخدم أداة CLI الخاصة بالبوابة اللاسلكية Emerson 1410/1420.

### التثبيت

```bash
# Install Emerson Wireless Gateway software and tools
# Refer to Emerson documentation for 1410/1420 Gateway setup
```

التحقق: `emerson_1410_cli --help`

## الاستخدام

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

## برنامج التشغيل

```go
b, c, err := wirelesshart.NewCmdDriver("/usr/bin/emerson_1410_cli")
```

## الوظائف المدعومة

| الوظيفة | الوصف |
|----------|----------------------------|
| `read`   | قراءة معامل جهاز |
| `write`  | كتابة معامل جهاز |
| `scan`   | المسح بحثًا عن الأجهزة اللاسلكية |

## الاختبار

```bash
go test ./... -v
```
