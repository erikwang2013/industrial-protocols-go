# SDK جسر أوامر POWERLINK

POWERLINK (Ethernet POWERLINK) هو بروتوكول إيثرنت في الزمن الحقيقي لأتمتة الصناعة. توفر هذه الحزمة مُرمِّزًا (codec) لـ POWERLINK عبر أداة سطر الأوامر `openPOWERLINK_demo`.

## أداة سطر الأوامر

تستخدم تطبيق العرض التوضيحي لمكدس openPOWERLINK.

### التثبيت

```bash
git clone https://github.com/OpenAutomationTechnologies/openPOWERLINK.git
cd openPOWERLINK
cmake -B build && cmake --build build
sudo cmake --install build
```

التحقق: `openPOWERLINK_demo --help`

## الاستخدام

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

## برنامج التشغيل

```go
b, c, err := powerlink.NewCmdDriver("/usr/bin/openPOWERLINK_demo")
```

## الوظائف المدعومة

| الوظيفة | الوصف |
|----------|-------------------------------|
| `read`   | قراءة إدخال في قاموس الكائنات |
| `write`  | كتابة إدخال في قاموس الكائنات |
| `status` | الاستعلام عن حالة العقدة/NMT |

## الاختبار

```bash
go test ./... -v
```
