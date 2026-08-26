# SDK جسر أوامر SERCOS I/II

SERCOS I/II هو الإصدار التسلسلي القديم عبر الألياف الضوئية من واجهة SERCOS للتحكم الرقمي في الحركة. توفر هذه الحزمة مُرمِّزًا (codec) لـ SERCOS I/II عبر أداة سطر الأوامر `sercos_cli`.

## أداة سطر الأوامر

تستخدم أداة CLI لواجهة SERCOS عبر الألياف الضوئية.

### التثبيت

```bash
# SERCOS interface card driver and tools
# Refer to vendor documentation for the specific SERCOS master card
```

التحقق: `sercos_cli --help`

## الاستخدام

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

## برنامج التشغيل

```go
b, c, err := sercos1.NewCmdDriver("/usr/bin/sercos_cli")
```

## الوظائف المدعومة

| الوظيفة | الوصف |
|----------|--------------------------|
| `read`   | قراءة SERCOS IDN |
| `write`  | كتابة SERCOS IDN |
| `status` | الاستعلام عن حالة المشغّل |

## الاختبار

```bash
go test ./... -v
```
