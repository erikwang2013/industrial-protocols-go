# SDK جسر بوابة Interbus

بروتوكول Interbus عبر جسر بوابة TCP.

## الأجهزة

| البوابة | الواجهة | IP الافتراضي |
|---------|-----------|-------------|
| Phoenix Contact IBS S5 DSC/I-T | وحدة تحكم Interbus رئيسية | 192.168.0.60 |
| Phoenix Contact IBS PCI SC/I-T | Interbus رئيسي عبر PCI | يحدده المضيف |
| HMS Anybus Interbus | بوابة مدمجة | 192.168.0.61 |

## التوصيلات

- موصل D-SUB ذو 9 أسنان على البوابة إلى ناقل Interbus البعيد (داخل/خارج)
- الدرع موصول بـ FE في كلا الطرفين
- إيثرنت إلى البوابة لجسر TCP

## إعداد IP البوابة

```go
driver, codec, err := interbus.NewGatewayDriver("192.168.0.60:2006")
handler, err := interbus.ReadyHandler("192.168.0.60:2006")
```
