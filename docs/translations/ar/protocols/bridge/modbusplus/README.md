# SDK بروتوكول Modbus Plus

Modbus Plus (MB+) هو شبكة صناعية عالية السرعة بتمرير الرمز (token) طوّرتها Modicon (Schneider Electric).

## الأنواع (Variants)

| النوع | نوع الجسر | الوصف |
|---------|-------------|-------------|
| `gateway` | GatewayBridge | اتصال TCP بمحول SA85/BM85 |
| `cmd` | CmdBridge | غلاف CLI لأداة `sa85_cli` |

## برنامج تشغيل البوابة

```go
driver, codec, err := modbusplus.NewGatewayDriver("192.168.0.160:2010")
handler, err := modbusplus.ReadyHandler("192.168.0.160:2010")
```

## برنامج تشغيل الأوامر

```go
driver, codec, err := modbusplus.NewCmdDriver("/usr/bin/sa85_cli")
```

### تثبيت أداة سطر الأوامر

```bash
# Install SA85 Modbus Plus driver and tools
# Refer to Schneider Electric documentation for SA85 adapter
```

التحقق: `sa85_cli --help`

## الأجهزة

| البوابة | الواجهة | IP الافتراضي |
|---------|-----------|-------------|
| Schneider SA85 | محول Modbus Plus عبر ISA | يحدده المضيف |
| Schneider BM85 | جسر/مضاعِف Modbus Plus | 192.168.0.A0 |
| ProSoft MVI56-MBP | وحدة MB+ لـ ControlLogix | يحدده المضيف |

## التوصيلات

- كبل ثنائي المحور (RG-62) بموصلات BNC للجذع MB+
- مقاوم إنهاء (78 أوم في كل طرف)
- جسر BM85 يربط MB+ بـ TCP عبر الإيثرنت

## تنسيق إطار البروتوكول

- Magic (2 بايت): `0x4D42`
- الوجهة (1 بايت): عنوان العقدة
- الأمر (1 بايت): 0x01=قراءة، 0x02=كتابة
- الطول (2 بايت): حجم الحمولة (big-endian)
- الحمولة (N بايت): البيانات
- CRC (2 بايت): Modbus CRC-16 (little-endian)

## الاختبار

```bash
go test ./... -v
```
