# SDK بروتوكول MOST عبر المنفذ التسلسلي

MOST (نظام النقل الموجه للوسائط) هو تقنية شبكات وسائط متعددة عالية السرعة تُستخدم أساسًا في أنظمة الترفيه بالسيارات. توفر هذه الحزمة مُرمِّزًا (codec) لـ MOST عبر محول تسلسلي يستخدم واجهة أوامر AT.

## نظرة عامة على البروتوكول

يستخدم MOST اتصالًا تسلسليًا متزامنًا فوق طبقة فيزيائية من الألياف الضوئية. يتصل هذا التنفيذ عبر محول تسلسلي يعرض واجهة أوامر AT بمعدل بود 115200.

### أوامر AT

| الأمر | الوصف |
|--------------------|------------------------------|
| `AT+READ=<addr>,<len>`  | قراءة بايتات من العنوان |
| `AT+WRITE=<addr>,<hex>` | كتابة بايتات سداسية عشرية إلى العنوان |
| `AT+STATUS`             | الاستعلام عن حالة الحلقة/الشبكة |

### تنسيق الاستجابة

- `+OK:<hex_data>` — استجابة ناجحة
- `+ERR:<code>` — استجابة خطأ

## متطلبات الأجهزة

- شبكة ألياف ضوئية MOST مع إنهاء صحيح
- محول MOST-إلى-تسلسلي (مثل محول MOST150 USB)
- منفذ تسلسلي بمعدل 115200 بود، 8N1

## الاستخدام

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/automotive/most"
)

func main() {
    p := most.New()
    codec, _ := p.NewCodec("serial")

    // Read 4 bytes from address 0x0100
    req := &kernel.Request{
        Function: "read",
        Address:  "0x0100",
        Count:    4,
    }
    raw, _ := codec.Encode(req)
    // raw = "AT+READ=0x0100,4\r\n"
    _ = raw

    // Write data to address 0x0200
    req2 := &kernel.Request{
        Function: "write",
        Address:  "0x0200",
        Data:     []byte{0xDE, 0xAD, 0xBE, 0xEF},
    }
    raw2, _ := codec.Encode(req2)
    // raw2 = "AT+WRITE=0x0200,DEADBEEF\r\n"
    _ = raw2
}
```

## برنامج التشغيل

```go
b, c, err := most.NewSerialDriver("/dev/ttyUSB0")
```

## الوظائف المدعومة

| الوظيفة | الوصف |
|----------|-----------------------------|
| `read`   | قراءة من عنوان MOST |
| `write`  | كتابة بيانات إلى عنوان MOST |
| `status` | الاستعلام عن حالة الحلقة/الشبكة |

## الاختبار

```bash
go test ./... -v
```
