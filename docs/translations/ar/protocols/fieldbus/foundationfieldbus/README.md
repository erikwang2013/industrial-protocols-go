# SDK جسر بوابة Foundation Fieldbus

Foundation Fieldbus H1/HSE عبر جسر بوابة TCP.

## الأجهزة

| البوابة | الواجهة | IP الافتراضي |
|---------|-----------|-------------|
| NI USB-8486 | واجهة H1 عبر USB | يحدده المضيف |
| Softing FFusb | واجهة H1 عبر USB | يحدده المضيف |
| P+F HD2-GTR-4PA | بوابة H1 إلى إيثرنت | 192.168.0.20 |

## التوصيلات

- جذع H1 (زوج مجدول، محمي) مع إنهاء في كلا الطرفين
- مكيف طاقة ناقل لـ H1 (24 فولت تيار مستمر، 350-500 مللي أمبير لكل مقطع)
- البوابة تربط مقطع H1 بجسر TCP عبر الإيثرنت

## إعداد IP البوابة

```go
driver, codec, err := foundationfieldbus.NewGatewayDriver("192.168.0.20:2002")
handler, err := foundationfieldbus.ReadyHandler("192.168.0.20:2002")
```
