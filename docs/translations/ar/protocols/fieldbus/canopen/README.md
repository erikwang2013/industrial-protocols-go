# SDK بروتوكول CANopen

CANopen هو بروتوكول طبقة أعلى قائم على CAN لأنظمة التحكم المدمجة. توفر هذه الحزمة مُرمِّزًا (codec) لـ CANopen وبرنامج تشغيل SocketCAN.

## نظرة عامة على البروتوكول

يستخدم CANopen معرّفات CAN قياسية 11-bit مع مجموعة الاتصالات المعرفة مسبقًا التالية:

| الوظيفة | معرّف CAN | الوصف |
|------------|---------------------|-------------------------------|
| NMT        | 0x000               | إدارة الشبكة |
| SYNC       | 0x080               | رسالة المزامنة |
| SDO (tx)   | 0x580 + NodeID      | كائن خدمة البيانات (خادم) |
| SDO (rx)   | 0x600 + NodeID      | كائن خدمة البيانات (عميل) |
| PDO1 (tx)  | 0x180 + NodeID      | كائن بيانات العملية 1 |
| Heartbeat  | 0x700 + NodeID      | نبضة القلب / بدء التشغيل |

## متطلبات الأجهزة

- نظام **Linux** مع دعم SocketCAN (بتفعيل `CONFIG_CAN`)
- واجهة CAN (مثل `can0`، أو `vcan0` للـ CAN الافتراضي)
- أجهزة إرسال/استقبال تدعم CAN (مثل MCP2515 أو SJA1000 أو محول USB-CAN)

### إعداد واجهة CAN افتراضية (للاختبار)

```bash
sudo modprobe can
sudo modprobe can_raw
sudo modprobe vcan
sudo ip link add dev vcan0 type vcan
sudo ip link set up vcan0
```

## الاستخدام

```go
package main

import (
    "fmt"
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/fieldbus/canopen"
)

func main() {
    p := canopen.New()
    codec, _ := p.NewCodec("can")

    // Read SDO object dictionary entry
    req := &kernel.Request{
        Function: "sdo_read",
        Metadata: map[string]any{
            "index": float64(0x1000), // Device type
            "sub":   float64(0),
        },
    }
    raw, _ := codec.Encode(req)
    fmt.Printf("SDO read frame: %X\n", raw)

    // NMT start remote node
    req2 := &kernel.Request{Function: "nmt_start"}
    raw2, _ := codec.Encode(req2)
    fmt.Printf("NMT start frame: %X\n", raw2)
}
```

## الوظائف المدعومة

| الوظيفة | الوصف |
|-------------|--------------------------------------|
| `sdo_read`  | قراءة إدخال في قاموس الكائنات |
| `sdo_write` | كتابة إدخال في قاموس الكائنات |
| `nmt_start` | تشغيل عقدة بعيدة (NMT) |
| `nmt_stop`  | إيقاف عقدة بعيدة (NMT) |
| `nmt_reset` | إعادة ضبط عقدة بعيدة (NMT) |
| `heartbeat` | إرسال رسالة نبضة القلب / بدء التشغيل |

## الاختبار

```bash
go test ./... -v
```

ملاحظة: تتطلب اختبارات برنامج تشغيل SocketCAN نظام Linux مع أجهزة CAN أو واجهة CAN افتراضية.
