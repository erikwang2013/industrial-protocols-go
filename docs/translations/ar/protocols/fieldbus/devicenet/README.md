# SDK بروتوكول DeviceNet

DeviceNet هو بروتوكول شبكة صناعية قائم على CAN لأتمتة المصانع. توفر هذه الحزمة مُرمِّزًا (codec) لـ DeviceNet مع برنامجي تشغيل: SocketCAN وبوابة TCP.

## نظرة عامة على البروتوكول

يستخدم DeviceNet البروتوكول الصناعي المشترك (CIP) فوق CAN. تستخدم مجموعة الاتصالات المعرفة مسبقًا الرئيسية/التابعة رسائل المجموعة 2 (معرّف CAN 0x400 + NodeID) لـ I/O المقرون والرسائل الصريحة.

## متطلبات الأجهزة

### وضع CAN (SocketCAN)
- نظام **Linux** مع دعم SocketCAN (بتفعيل `CONFIG_CAN`)
- واجهة CAN (مثل `can0`، `vcan0`)
- أجهزة CAN تدعم DeviceNet (مثل Anybus Communicator أو HMS IXXAT)

### وضع البوابة
- اتصال TCP/IP ببوابة DeviceNet
- بوابة تدعم بروتوكول أوامر نصية عادي (مثل HMS Anybus أو Hilscher netX)

### إعداد واجهة CAN افتراضية (للاختبار)

```bash
sudo modprobe can
sudo modprobe can_raw
sudo modprobe vcan
sudo ip link add dev vcan0 type vcan
sudo ip link set up vcan0
```

## الاستخدام

### وضع CAN

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/fieldbus/devicenet"
)

func main() {
    p := devicenet.New()
    codec, _ := p.NewCodec("can")

    // Poll request
    req := &kernel.Request{
        Function: "poll",
        Data:     []byte{0x01, 0x00},
    }
    raw, _ := codec.Encode(req)

    // Open connection
    req2 := &kernel.Request{Function: "open"}
    raw2, _ := codec.Encode(req2)
    _ = raw
    _ = raw2
}
```

### وضع البوابة

```go
p := devicenet.New()
codec, _ := p.NewCodec("gateway")

req := &kernel.Request{
    Function: "poll",
    Data:     []byte{0xAB, 0xCD},
}
raw, _ := codec.Encode(req)
// raw will be: "poll abcd\n"
```

## الوظائف المدعومة

| الوظيفة | الوصف |
|----------|--------------------------------|
| `open`   | فتح اتصال صريح |
| `poll`   | استقصاء بيانات I/O (المجموعة 2) |

## الاختبار

```bash
go test ./... -v
```

ملاحظة: تتطلب اختبارات برنامج تشغيل SocketCAN نظام Linux. أما اختبارات برنامج تشغيل البوابة فتعمل على أي نظام تشغيل.
