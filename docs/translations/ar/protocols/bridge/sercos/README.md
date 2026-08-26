# SDK جسر أوامر SERCOS III

SERCOS III (نظام الاتصال التسلسلي في الزمن الحقيقي) هو واجهة رقمية للتحكم في الحركة، تستخدم طوبولوجيا حلقية عبر الألياف الضوئية أو النحاس. توفر هذه الحزمة مُرمِّزًا (codec) لـ SERCOS III عبر أداة سطر الأوامر `netx_cli`.

## أداة سطر الأوامر

تستخدم أداة CLI الخاصة بـ Hilscher netX SERCOS III.

### التثبيت

```bash
# Install Hilscher netX driver and tools
# See https://www.hilscher.com/ for netX drivers
sudo apt-get install netx-driver
```

التحقق: `netx_cli --help`

## الاستخدام

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

## برنامج التشغيل

```go
b, c, err := sercos.NewCmdDriver("/usr/bin/netx_cli")
```

## الوظائف المدعومة

| الوظيفة | الوصف |
|----------|-----------------------------|
| `read`   | قراءة معامل SERCOS IDN/S |
| `write`  | كتابة معامل SERCOS IDN/S |
| `phase`  | ضبط مرحلة الاتصال |

## الاختبار

```bash
go test ./... -v
```
