# تقرير الاختبار — industrial-protocols-go

التاريخ: 2026-08-27
النطاق: مساحة عمل Go متعددة الوحدات (go.work، 42 وحدة: kernel / protocols / examples)
الأمر: `go test ./...` و `go test -cover ./...` (نظرًا لعدم وجود go.mod في جذر مساحة العمل، تم التنفيذ بالتوازي وحدةً بوحدة، والنتائج متطابقة)

## 1. الاستنتاج العام

- **أخضر بالكامل**: جميع الوحدات الـ 42 والحزم الـ 52 التي تحتوي على اختبارات نجحت (تم التحقق عبر عدة عمليات تشغيل كاملة، بما في ذلك أثناء كتابة اختبارات متوازية).
- عدد ملفات الاختبار 106، وعدد دوال الاختبار 659.
- `go vet ./...` دون أي تحذيرات في جميع الوحدات.
- متوسط تغطية العبارات **83.3%**؛ تغطية حزم kernel الأساسية بين 95% و100%.

## 2. إحصائيات الاختبار والتغطية لكل وحدة

### kernel (11 حزمة، الأعلى تغطية باستثناء examples)

| الحزمة | التغطية |
|---|---|
| kernel | 100% |
| kernel/bridge | 80.0% |
| kernel/config | 100% |
| kernel/connection | 95.3% |
| kernel/event | 100% |
| kernel/gateway | 100% |
| kernel/metrics | 100% |
| kernel/pipeline | 100% |
| kernel/security | 100% |
| kernel/transport | 100% |
| kernel/vendor | 100% |

### الوحدات الرئيسية في protocols

| الوحدة | التغطية | الوحدة | التغطية |
|---|---|---|---|
| automotive/flexray | 79.4% | fieldbus/canopen | 73.9% |
| automotive/kline | 93.3% | fieldbus/cclink | 93.8% |
| automotive/lin | 92.5% | fieldbus/cclinkie | 73.5% |
| automotive/most | 100% | fieldbus/devicenet | 73.2% |
| automotive/saej1850 | 84.2% | fieldbus/dnp3 | 88.9% |
| bridge/controlnet | 100% | fieldbus/hart | 96.6% |
| bridge/ethercat | 100% | fieldbus/iec61850 | 91.8% |
| bridge/interbus | 70.0% | iot/hartip | 88.6% |
| bridge/isa100 | 95.2% | iot/mqtt | 89.1% |
| bridge/powerlink | 97.8% | ethernet/bacnet | 95.2% |
| bridge/sercos | 97.9% | ethernet/ethernetip | 96.9% |
| bridge/sercos1 | 95.6% | ethernet/modbus | 67.7% |
| bridge/safej1850 | 100% | ethernet/opcua | 95.8% |
| bridge/wirelesshart | 95.2% | ethernet/profinet | 95.4% |
| building/dali | 92.9% | system/cpci | 78.6% |
| building/lonworks | **37.9%** | system/pci | 78.6% |
| fieldbus/asinterface | **37.9%** | system/vme | 78.6% |
| fieldbus/foundationfieldbus | **37.9%** | fieldbus/profibus | **35.7%** |
| fieldbus/iolink | **37.9%** | bridge/lightbus / modbusplus / worldfip | 72~73.5% |

> examples/modbus_basic لا يحتوي على ملفات اختبار (برنامج مثال، وهذا متوقع).

## 3. قائمة الإصلاحات

### 3.1 إصلاح أخطاء المصدر (اكتُشفت وتم التحقق منها بواسطة tester-kernel / tester-protocols)

| الملف:السطر | المشكلة | الإصلاح |
|---|---|---|
| kernel/gateway/engine.go:27 | `Transform` يسبب انهيار nil dereference عند `rule.Src.Name()` / `rule.Dst.Name()` / `rule.Map(...)` للقواعد التي تكون فيها `Src`/`Dst`/`Map` قيمة nil | تخطي القواعد غير المكتملة الإعداد في بداية الحلقة (بما يتوافق مع دلالة "تخطي القواعد السيئة" الحالية) |
| protocols/ethernet/modbus/modbus.go:215 | `decodePDU` يسبب انهيار تجاوز حدود المصفوفة عند قراءة `pdu[1]` من PDU غير طبيعي من 1 بايت (مثل `0x81`) | إضافة فحص `len(pdu) < 2` وإرجاع خطأ تحليل |
| protocols/ethernet/modbus/modbus.go:230 | تجاوز حدود المصفوفة عند فك الترميز عندما يتجاوز عدد البايتات (`pdu[1]`) البيانات الفعلية (إطار خبيث/تالف) | إضافة فحص حد `2+n > len(pdu)` وإرجاع خطأ تحليل |
| protocols/ethernet/opcua/opcua.go:63,74 | `encodeHello` يعيد كتابة طول الرسالة عند الإزاحة 8 (موضع إصدار البروتوكول)، بينما تتطلب المواصفات الإزاحة 4 | تسجيل `sizePos` مسبقًا (الإزاحة 4) وإعادة الكتابة في الموضع الصحيح |
| protocols/ethernet/opcua/opcua.go:82,89 | حقل الطول في `encodeOpenSecureChannel` لم يُعاد كتابته أبدًا، وظل صفرًا دائمًا | إعادة الكتابة في النهاية وفقًا لطول الإطار الفعلي |
| protocols/ethernet/ethernetip/ethernetip.go:111,117 | `encodeReadTag` يخصص إطارًا بحجم 16+len(cipReq) لكن طلب CIP يُكتب عند [28:]، والـ 12 بايت الزائدة تسبب إسقاط `copy` لطلب CIP كاملًا بصمت؛ وحقل الطول ينقصه 12 بالمقابل | التخصيص أصبح 28+len(cipReq)، وحقل length أصبح 4+len(cipReq) |
| protocols/ethernet/profinet/profinet.go:72 | `Decode` يسبب انهيار تجاوز حدود المصفوفة عند `blockLen` خارج النطاق | إضافة فحص `len(data) < 10+blockLen` |
| protocols/fieldbus/canopen/canopen.go:102 | الحد الأدنى للطول في `Decode` هو 4، لكن `unmarshalCAN` يتطلب `data[4:8]`، فالإطارات من 4 إلى 7 بايت تسبب انهيارًا | تغيير الحد الأدنى إلى 8 (تنسيق الخط ثابت 8 بايت) |
| protocols/fieldbus/dnp3/dnp3.go:121 | `decodeTPDU` يسبب انهيار تجاوز حدود المصفوفة عند `length<5` (عندما يكون `apduEnd < apduStart`) أو عند `length` كبير جدًا (عندما يكون `apduEnd > len(data)`) | توحيد التحقق `length < 5 || apduEnd > len(data)` ثم إرجاع خطأ |
| protocols/fieldbus/dnp3/dnp3.go:157 | متغير العمل في `crc16DNP` لم يُقنّع بـ `& 0xFF` وفق المواصفات؛ المتجه المعروف `crc16DNP("123456789")` يعطي 0x69FF (يجب أن يكون 0xEA82) | `temp := (crc ^ b) & 0xFF` |
| protocols/automotive/kline/kline.go:74 | فحص `len<5` في `Decode` يسبق الحكم على إطار التزامن سريع التهيئة، فإطار التزامن من 1 بايت (0x55) يُرفض خطأً | فحص الإطار الفارغ وإطار التزامن أولًا، ثم فحص الطول |

### 3.2 إصلاحات الاختبارات (تصحيح أخطاء الترجمة / أخطاء التأكيد / تكرار الأسماء / تأكيدات قديمة)

| الملف:السطر | المشكلة | الإصلاح |
|---|---|---|
| protocols/automotive/kline/kline_extra_test.go:43 | الثابت `0x68+0x10+0xF1+0x01`(362) يُسند إلى byte فتفشل الترجمة | تغيير التوقع إلى 8 بتات السفلى من مجموع الاختبار `0x6A` |
| protocols/automotive/saej1850/saej1850_extra_test.go | استدعاء دوال/حقول غير مُصدَّرة على واجهة `kernel.Codec`، فشل الترجمة | إضافة helper `newCodec` للتأكيد على أنه `*j1850Codec`؛ تصحيح إزاحة التأكيد في `TestEncodeMode01PID` إلى raw[4..6] (معرف 4 بايتات + أول 4 بايتات من البيانات) |
| protocols/bridge/sercos1/sercos1_extra_test.go:45 | `TestEncodeReadDefaults` يتعارض اسمه مع اختبار قائم، فشل الترجمة | إعادة التسمية إلى `TestEncodeReadDefaultsFiber` (مع الإبقاء على تغطية متغير fiber) |
| protocols/building/dali/dali_extra_test.go:68 | `TestEncodeDirectArc` يتعارض اسمه مع اختبار قائم، فشل الترجمة | إعادة التسمية إلى `TestEncodeDirectArcCmd` |
| protocols/bridge/powerlink/powerlink_extra_test.go:51 | التوقع `write 0x1000 A` بينما الفعلي `0A` (`%X` يثبّت الرقمين لكل بايت، بما يتوافق مع اتفاقية الاختبار القائم `DEADBEEF`) | تغيير التوقع إلى `0A` |
| protocols/bridge/sercos/sercos_extra_test.go:62 | كما سبق، التوقع `1` يجب أن يكون `01` | تغيير التوقع إلى `01` |
| protocols/ethernet/opcua/opcua_extra_test.go:35 | التأكيد على سلوك الخطأ القديم (الطول عند الإزاحة 8) لم يعد صالحًا بعد إصلاح المصدر | تغيير التأكيد إلى أن حقل إصدار البروتوكول هو 0 (الطول تغطيه TestExtraHELMessageSizeAtOffset4) |
| protocols/fieldbus/canopen/canopen_extra_test.go:135 | تأكيد "prove-it" معكوس (`err != nil` يسبب الفشل، بينما خطأ التحليل هو النتيجة المتوقعة) | تغييره إلى الفشل عند `err == nil` |
| protocols/ethernet/profinet/profinet_extra_test.go:23 | كما سبق، تأكيد معكوس | تغييره إلى الفشل عند `err == nil` |
| protocols/fieldbus/dnp3/dnp3_extra_test.go:35,53 | كما سبق ×2 | تغييره إلى الفشل عند `err == nil` |

> ملاحظة: مشكلة تنفيذ واجهة `stubTransport` في `kernel/security/tls_test.go` وقيم CRC المتوقعة في `cclink/cclinkie` صُححت ذاتيًا من قِبل مهندس الاختبار أثناء العمل المتوازي، ولم تُعدل في هذه الدفعة.

## 4. المخاطر المتبقية

1. **الوحدات منخفضة التغطية**: profibus (35.7%)، lonworks / asinterface / foundationfieldbus / iolink (37.9%) — الاختبارات تغطي مسارات قليلة فقط، ويُوصى لاحقًا بإضافة فروع Decode/Encode والمسارات غير الطبيعية.
2. **CRC الخاص بـ modbus RTU غير مُتحقق منه**: `decodeRTU` لا يتحقق من CRC (اختبار `TestExtraRTUDecodesCorruptCRC` يوثّق هذه الفجوة صراحةً ويمر كما هو). يُوصى بإضافته عند التكامل مع أجهزة حقيقية.
3. **قيد تنسيق خط canopen**: `marshalCAN` يحمل أول 4 بايتات من البيانات فقط (`TestExtraSDOWritePayloadLostOnWire` يوثّق ذلك)، فتُفقد حمولات SDO متعددة البايتات عند الكتابة.
4. **تخطي الاختبارات المعتمدة على الأجهزة**: cpci / pci / vme تستخدم `t.Skip` في بيئة بلا أجهزة (منطقي).
5. **ملفات مؤقتة متبقية في المستودع**: ملف `main.go` غير المتتبع في الجذر (برنامج استكشاف يشير إلى وحدة `probe/` غير موجودة)، و`scripts/`، و`docs/*.png` — لا تنتمي لهذا التسليم، ويُوصى بتنظيفها.
6. **الاختبارات لا تزال تُكتب بالتوازي**: يعتمد هذا التقرير على آخر لقطة خضراء كاملة؛ إذا استمر مهندس الاختبار بإضافة ملفات `*_extra_test.go` لاحقًا، فيجب إعادة تشغيل الانحدار الشامل.
