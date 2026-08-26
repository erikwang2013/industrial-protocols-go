# SDK جسر أوامر EtherCAT

EtherCAT (الإيثرنت لتقنية أتمتة التحكم) هو ناقل ميدان إيثرنت صناعي عالي الأداء. توفر هذه الحزمة مُرمِّزًا (codec) لـ EtherCAT عبر أداة سطر الأوامر `ethercat`.

## أداة سطر الأوامر

تستخدم أداة سطر الأوامر الخاصة بـ IgH EtherCAT Master.

### التثبيت

```bash
# Debian/Ubuntu
sudo apt-get install ethercat-master

# From source
git clone https://gitlab.com/etherlab.org/ethercat.git
cd ethercat
make && sudo make install
```

التحقق: `ethercat slaves`

## الاستخدام

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/ethercat"
)

func main() {
    p := ethercat.New()
    codec, _ := p.NewCodec("cmd")

    // Upload SDO from address 0x1000
    req := &kernel.Request{
        Function: "upload",
        Address:  "0x1000",
        Count:    4,
    }
    raw, _ := codec.Encode(req)
    // raw = "upload 0x1000 4\n"
    _ = raw
}
```

## برنامج التشغيل

```go
b, c, err := ethercat.NewCmdDriver("/usr/bin/ethercat")
```

## الوظائف المدعومة

| الوظيفة | الوصف |
|------------|------------------------|
| `upload`   | قراءة SDO من العنوان |
| `download` | كتابة SDO إلى العنوان |
| `slaves`   | سرد عُبيد EtherCAT |

## الاختبار

```bash
go test ./... -v
```
