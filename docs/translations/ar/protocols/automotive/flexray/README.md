# SDK بروتوكول FlexRay عبر CAN

FlexRay هو بروتوكول اتصال سيارات حتمي عالي السرعة. توفر هذه الحزمة مُرمِّزًا (codec) لـ FlexRay عبر ناقل CAN، مع إطارات مبنية على الدورات (cycles) وفحوصات سلامة CRC-16/XMODEM.

## نظرة عامة على البروتوكول

يستخدم FlexRay مخطط الوصول المتعدد بتقسيم الزمن (TDMA) مع دورات اتصال متكررة. تتكون كل دورة من مقاطع ثابتة ومتحركة. يرمّز هذا المرمّز إطارات FlexRay على إطارات CAN ممتدة بعنوان 29-bit.

### تنسيق السلك (Wire Format)

حمولة دورة FlexRay:
- **الترويسة** (2 بايت): رقم الدورة (little-endian)
- **الحالة** (1 بايت): bit 7=PPI (مؤشر ديباجة الحمولة)، bit 6=NFI، bit 5=SYF، bit 4=SUF
- **البيانات** (N بايت): الحمولة (بحد أقصى 254 بايت)
- **CRC** (2 بايت): CRC-16/XMODEM على الترويسة+الحالة+البيانات (little-endian)

ترميز معرّف CAN (ممتد 29-bit):
- البتات 28-24: نوع الرسالة (0x01=إطار، 0x02=حالة)
- البتات 23-10: محجوزة
- البتات 15-10: معرف الفتحة (6 بتات)
- البتات 9-0: رقم الدورة (10 بتات)

## متطلبات الأجهزة

- نظام **Linux** مع دعم SocketCAN
- محول Vector VN7600/VN7640 أو محول FlexRay-CAN من Bosch
- شبكة FlexRay مع إنهاء (termination) صحيح (انحياز 2.5V)

### إعداد واجهة CAN

```bash
sudo modprobe can
sudo modprobe can_raw
sudo ip link set can0 type can bitrate 500000
sudo ip link set up can0
```

## الاستخدام

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/automotive/flexray"
)

func main() {
    p := flexray.New()
    codec, _ := p.NewCodec("can")

    // Send a FlexRay frame in cycle 5 with PPI set
    req := &kernel.Request{
        Function: "frame",
        Data:     []byte{0x42, 0x01},
        Metadata: map[string]any{
            "cycle": float64(5),
            "ppi":   true,
        },
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## الوظائف المدعومة

| الوظيفة | الوصف |
|----------|------------------------------------------|
| `frame`  | إرسال حمولة إطار FlexRay |
| `status` | الاستعلام عن إعداد الفتحة (الفتحة، الدورة) |

## الاختبار

```bash
go test ./... -v
```
