# SDK بروتوكول SAE J1850 عبر CAN

SAE J1850 هو معيار اتصال للمركبات يُستخدم للتشخيص على متن المركبة (OBD-II). توفر هذه الحزمة مُرمِّزًا (codec) لـ J1850 عبر ناقل CAN (ISO 15765-4 / CAN TP).

## نظرة عامة على البروتوكول

يستخدم SAE J1850 OBD-II عبر CAN معرّفات CAN ممتدة 29-bit بالبنية التالية:

### تنسيق معرّف CAN (29-bit)

| البتات | الحقل | الوصف |
|-----------|----------|------------------------------------|
| 28-26     | Priority | أولوية الرسالة (0-7، الافتراضي 6) |
| 25        | Ext ID   | دائمًا 1 للإطارات الممتدة |
| 24-16     | PF       | تنسيق المعامل (الترويسة) |
| 15-8      | PS       | المعامل المحدد (الهدف/المصدر) |
| 7-0       | SA       | عنوان المصدر |

### معرّفات OBD-II القياسية عبر CAN

| النوع | معرّف CAN (سداسي عشري) | الوصف |
|--------------------|-----------------|--------------------------------|
| طلب فيزيائي | 0x18DAxxF1 | طلب إلى وحدة تحكم إلكترونية محددة (xx=العنوان) |
| استجابة فيزيائية | 0x18DAF1xx | استجابة من وحدة التحكم الإلكترونية (xx=العنوان) |
| طلب وظيفي | 0x18DB33F1 | بث إلى جميع وحدات التحكم الإلكترونية |

### تنسيق إطار ISO 15765-2

الإطار الفردي: الربع العلوي من البايت 0 = طول البيانات (0-7)، والربع السفلي + البايتات المتبقية = بيانات التشخيص.

## متطلبات الأجهزة

- نظام **Linux** مع دعم SocketCAN
- محول J1850-CAN OBD-II (مثل محول USB-إلى-CAN متوافق مع ELM327، أو OBDLink SX، أو Macchina M2)
- مركبة بموصل OBD-II (معظم المركبات منذ 1996+)

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
    "github.com/erikwang2013/industrial-protocols-go/protocols/automotive/saej1850"
)

func main() {
    p := saej1850.New()
    codec, _ := p.NewCodec("can")

    // Mode $01 PID $0C: Engine RPM
    req := &kernel.Request{
        Function: "mode01",
        Metadata: map[string]any{"pid": float64(0x0C)},
    }
    raw, _ := codec.Encode(req)
    _ = raw

    // Mode $03: Request emission-related DTCs
    req2 := &kernel.Request{Function: "mode03"}
    raw2, _ := codec.Encode(req2)
    _ = raw2
}
```

## الوظائف المدعومة

| الوظيفة | وضع OBD-II | الوصف |
|----------------|-------------|-------------------------------------|
| `mode01`       | $01         | طلب بيانات مجموعة نقل الحركة الحالية |
| `mode03`       | $03         | طلب رموز الأعطال المتعلقة بالانبعاثات |
| `mode0A`       | $0A         | طلب رموز الأعطال الدائمة |
| `diag_request` | مخصص | طلب تشخيص عام |
| `diag_response`| -           | إطار استجابة التشخيص |
| `broadcast`    | -           | بث وظيفي إلى جميع وحدات التحكم الإلكترونية |

## الاختبار

```bash
go test ./... -v
```
