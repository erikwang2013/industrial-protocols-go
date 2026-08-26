# SDK جسر بوابة Lightbus

بروتوكول Beckhoff Lightbus عبر الألياف الضوئية وجسر بوابة TCP.

## الأجهزة

| البوابة | الواجهة | IP الافتراضي |
|---------|-----------|-------------|
| Beckhoff FC2001 | بطاقة Lightbus عبر PCI | يحدده المضيف |
| Beckhoff BK2000 | مقرنة ناقل Lightbus | 192.168.0.80 |
| Beckhoff FC9001 | محول Lightbus عبر الإيثرنت | 192.168.0.81 |

## التوصيلات

- طوبولوجيا حلقة من الألياف الضوئية البلاستيكية (POF)
- FC2001/FC9001 يربط الحلقة بجسر TCP عبر الإيثرنت
- لكل جهاز موصلا ألياف TX وRX

## إعداد IP البوابة

```go
driver, codec, err := lightbus.NewGatewayDriver("192.168.0.80:2008")
handler, err := lightbus.ReadyHandler("192.168.0.80:2008")
```
