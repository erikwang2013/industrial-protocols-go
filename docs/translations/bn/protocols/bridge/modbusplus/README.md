# Modbus Plus প্রোটোকল SDK

Modbus Plus (MB+) হল Modicon (Schneider Electric) কর্তৃক উন্নত একটি উচ্চ-গতির টোকেন-পাসিং শিল্প নেটওয়ার্ক।

## ভ্যারিয়েন্ট

| ভ্যারিয়েন্ট | ব্রিজ টাইপ | বিবরণ |
|---------|-------------|-------------|
| `gateway` | GatewayBridge | SA85/BM85 অ্যাডাপ্টারে TCP কানেকশন |
| `cmd` | CmdBridge | `sa85_cli` ইউটিলিটির CLI র‍্যাপার |

## গেটওয়ে ড্রাইভার

```go
driver, codec, err := modbusplus.NewGatewayDriver("192.168.0.160:2010")
handler, err := modbusplus.ReadyHandler("192.168.0.160:2010")
```

## Cmd ড্রাইভার

```go
driver, codec, err := modbusplus.NewCmdDriver("/usr/bin/sa85_cli")
```

### CLI টুল ইনস্টলেশন

```bash
# Install SA85 Modbus Plus driver and tools
# Refer to Schneider Electric documentation for SA85 adapter
```

যাচাই: `sa85_cli --help`

## হার্ডওয়্যার

| গেটওয়ে | ইন্টারফেস | ডিফল্ট IP |
|---------|-----------|-------------|
| Schneider SA85 | ISA Modbus Plus অ্যাডাপ্টার | হোস্ট-নির্ধারিত |
| Schneider BM85 | Modbus Plus ব্রিজ/মাল্টিপ্লেক্সার | 192.168.0.A0 |
| ProSoft MVI56-MBP | ControlLogix MB+ মডিউল | হোস্ট-নির্ধারিত |

## ওয়্যারিং

- MB+ ট্রাঙ্কের জন্য BNC কানেক্টরসহ টুইনঅ্যাক্সিয়াল ক্যাবল (RG-62)
- টার্মিনেটিং রেজিস্টর (প্রতিটি প্রান্তে ৭৮ ওহম)
- BM85 ব্রিজ ইথারনেটের মাধ্যমে MB+ কে TCP-তে সংযুক্ত করে

## প্রোটোকল ফ্রেম ফরম্যাট

- Magic (২ বাইট): `0x4D42`
- Destination (১ বাইট): নোড অ্যাড্রেস
- Command (১ বাইট): 0x01=read, 0x02=write
- Length (২ বাইট): পেলোড সাইজ (বিগ-এন্ডিয়ান)
- Payload (N বাইট): ডেটা
- CRC (২ বাইট): Modbus CRC-16 (লিটল-এন্ডিয়ান)

## টেস্টিং

```bash
go test ./... -v
```
