# SDK جسر أوامر ControlNet

ControlNet هو بروتوكول شبكة صناعية في الزمن الحقيقي طوّرته Allen-Bradley (Rockwell Automation) لتبادل البيانات عالي السرعة والحساس للوقت. توفر هذه الحزمة مُرمِّزًا (codec) لـ ControlNet عبر أداة سطر الأوامر `1784-pcic-cli`.

## أداة سطر الأوامر

تستخدم أداة CLI لبطاقة واجهة ControlNet 1784-PCIC.

### التثبيت

```bash
# Install Rockwell 1784-PCIC driver and tools
# Refer to Rockwell Automation documentation for RSLinx Classic SDK
```

التحقق: `1784-pcic-cli --help`

## الاستخدام

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

## برنامج التشغيل

```go
b, c, err := controlnet.NewCmdDriver("/usr/bin/1784-pcic-cli")
```

## الوظائف المدعومة

| الوظيفة | الوصف |
|----------|---------------------------|
| `read`   | قراءة من عقدة ControlNet |
| `write`  | كتابة إلى عقدة ControlNet |
| `status` | الاستعلام عن حالة بطاقة PCIC |

## الاختبار

```bash
go test ./... -v
```
