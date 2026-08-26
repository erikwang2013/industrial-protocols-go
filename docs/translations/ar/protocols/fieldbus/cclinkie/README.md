# SDK جسر بوابة CC-Link IE Field

بروتوكول CC-Link IE Field عبر جسر بوابة TCP.

## الأجهزة

| البوابة | الواجهة | IP الافتراضي |
|---------|-----------|-------------|
| Mitsubishi MELSEC IQ-R | CC-Link IE Field رئيسي | 192.168.3.40 |
| Mitsubishi RJ71GF11-T2 | وحدة CC-Link IE Field | يحدده المضيف |
| HMS Anybus CC-Link IE | بوابة مدمجة | 192.168.0.52 |

## التوصيلات

- إيثرنت RJ45 لـ CC-Link IE Field (طوبولوجيا حلقة أو نجمة بسرعة 1 جيجابت/ثانية)
- منفذ الإدارة على شبكة منفصلة
- البوابة تربط الشبكة الميدانية بجسر TCP

## إعداد IP البوابة

```go
driver, codec, err := cclinkie.NewGatewayDriver("192.168.3.40:2005")
handler, err := cclinkie.ReadyHandler("192.168.3.40:2005")
```
