# SDK جسر بوابة IO-Link

بروتوكول IO-Link عبر جسر بوابة TCP.

## الأجهزة

| البوابة | الواجهة | IP الافتراضي |
|---------|-----------|-------------|
| ifm AL1332 | IO-Link Master (EtherNet/IP) | 192.168.0.40 |
| Balluff BNI00AZ | IO-Link Master (PROFINET) | 192.168.0.41 |
| SICK SIG200 | IO-Link Master (Ethernet) | 192.168.0.42 |

## التوصيلات

- موصل M12 (4 أسنان) لكل منفذ IO-Link: L+ (بني)، L- (أزرق)، C/Q (أسود)، غير مستخدم (أبيض)
- مصدر طاقة 24 فولت تيار مستمر للرئيسي والأجهزة
- إيثرنت على الرئيسي لجسر TCP

## إعداد IP البوابة

```go
driver, codec, err := iolink.NewGatewayDriver("192.168.0.40:2004")
handler, err := iolink.ReadyHandler("192.168.0.40:2004")
```
