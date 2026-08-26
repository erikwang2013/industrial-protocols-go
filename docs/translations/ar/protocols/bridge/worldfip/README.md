# SDK جسر بوابة WorldFIP

بروتوكول WorldFIP عبر جسر بوابة TCP.

## الأجهزة

| البوابة | الواجهة | IP الافتراضي |
|---------|-----------|-------------|
| FIPIO Agent | وكيل ناقل WorldFIP الميداني | 192.168.0.70 |
| FIP Gateway (Alstom) | WorldFIP إلى إيثرنت | 192.168.0.71 |
| NI FIP-USB | واجهة WorldFIP عبر USB | يحدده المضيف |

## التوصيلات

- موصل D-SUB ذو 9 أسنان على البوابة إلى جذع WorldFIP (FIP1 = Data+، FIP2 = Data-)
- إنهاء خط (120 أوم) في طرفي الناقل
- إيثرنت لجسر TCP

## إعداد IP البوابة

```go
driver, codec, err := worldfip.NewGatewayDriver("192.168.0.70:2007")
handler, err := worldfip.ReadyHandler("192.168.0.70:2007")
```
