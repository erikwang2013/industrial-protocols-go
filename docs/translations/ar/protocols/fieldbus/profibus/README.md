# SDK جسر بوابة PROFIBUS

بروتوكول PROFIBUS DP/PA عبر جسر بوابة TCP.

## الأجهزة

| البوابة | الواجهة | IP الافتراضي |
|---------|-----------|-------------|
| Anybus Communicator | تابع PROFIBUS DP-V1 | 192.168.0.50 |
| وكيل Siemens CP 5611 | رئيسي PROFIBUS عبر PCI/PCIe | يحدده المضيف |
| HMS Fieldbus Gateway | Anybus NP40 | 192.168.0.51 |

## التوصيلات

- DB9 أنثى على البوابة إلى شبكة PROFIBUS (الخط A أخضر، الخط B أحمر)
- مقاوم إنهاء مفعّل في طرفي الشبكة
- إيثرنت إلى منفذ إدارة البوابة

## إعداد IP البوابة

```go
driver, codec, err := profibus.NewGatewayDriver("192.168.0.50:2000")
handler, err := profibus.ReadyHandler("192.168.0.50:2000")
```
