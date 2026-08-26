# SDK جسر أوامر ISA100.11a

ISA100.11a هو معيار شبكات صناعية لاسلكية لأتمتة العمليات. توفر هذه الحزمة مُرمِّزًا (codec) لـ ISA100.11a عبر أداة سطر الأوامر `yfgw410_cli` (بوابة يوكوغاوا الميدانية اللاسلكية YFGW410).

## أداة سطر الأوامر

تستخدم CLI الخاصة ببوابة Yokogawa YFGW410 Field Wireless Gateway.

### التثبيت

```bash
# Install Yokogawa YFGW410 gateway software and tools
# Refer to Yokogawa documentation for Field Wireless Gateway setup
```

التحقق: `yfgw410_cli --help`

## الاستخدام

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

## برنامج التشغيل

```go
b, c, err := isa100.NewCmdDriver("/usr/bin/yfgw410_cli")
```

## الوظائف المدعومة

| الوظيفة | الوصف |
|----------|-----------------------------|
| `read`   | قراءة خاصية جهاز |
| `write`  | كتابة خاصية جهاز |
| `list`   | سرد الأجهزة الموفّرة |

## الاختبار

```bash
go test ./... -v
```
