# SDK جسر بوابة AS-Interface

بروتوكول AS-Interface (ASi) عبر جسر بوابة TCP.

## الأجهزة

| البوابة | الواجهة | IP الافتراضي |
|---------|-----------|-------------|
| Bihl+Wiedemann BWU2044 | ASi Master (Ethernet) | 192.168.0.30 |
| P+F VBG-ENX-K20-DMD-EV | بوابة ASi | 192.168.0.31 |
| ifm AC1375 | وحدة تحكم ASi | 192.168.0.32 |

## التوصيلات

- كبل ASi أصفر (طاقة + بيانات) من البوابة إلى العبيد
- كبل طاقة مساعد أسود (24 فولت تيار مستمر للمشغلات) اختياري
- إيثرنت على البوابة لجسر TCP

## إعداد IP البوابة

```go
driver, codec, err := asinterface.NewGatewayDriver("192.168.0.30:2003")
handler, err := asinterface.ReadyHandler("192.168.0.30:2003")
```
