# Modbus Plus प्रोटोकॉल SDK

Modbus Plus (MB+) Modicon (Schneider Electric) द्वारा विकसित एक उच्च-गति टोकन-पासिंग औद्योगिक नेटवर्क है।

## वेरिएंट

| वेरिएंट | ब्रिज प्रकार | विवरण |
|---------|-------------|-------------|
| `gateway` | GatewayBridge | SA85/BM85 एडाप्टर के लिए TCP कनेक्शन |
| `cmd` | CmdBridge | `sa85_cli` उपयोगिता के लिए CLI रैपर |

## गेटवे ड्राइवर

```go
driver, codec, err := modbusplus.NewGatewayDriver("192.168.0.160:2010")
handler, err := modbusplus.ReadyHandler("192.168.0.160:2010")
```

## Cmd ड्राइवर

```go
driver, codec, err := modbusplus.NewCmdDriver("/usr/bin/sa85_cli")
```

### CLI उपकरण इंस्टॉलेशन

```bash
# Install SA85 Modbus Plus driver and tools
# Refer to Schneider Electric documentation for SA85 adapter
```

सत्यापन: `sa85_cli --help`

## हार्डवेयर

| गेटवे | इंटरफ़ेस | डिफ़ॉल्ट IP |
|---------|-----------|-------------|
| Schneider SA85 | ISA Modbus Plus एडाप्टर | host-assigned |
| Schneider BM85 | Modbus Plus ब्रिज/मल्टीप्लेक्सर | 192.168.0.A0 |
| ProSoft MVI56-MBP | ControlLogix MB+ मॉड्यूल | host-assigned |

## वायरिंग

- MB+ ट्रंक के लिए BNC कनेक्टर वाली ट्विनैक्सियल केबल (RG-62)
- टर्मिनेटिंग रेसिस्टर (प्रत्येक सिरे पर 78 ओम)
- BM85 ब्रिज MB+ को ईथरनेट के माध्यम से TCP से जोड़ता है

## प्रोटोकॉल फ्रेम फॉर्मेट

- Magic (2 बाइट): `0x4D42`
- गंतव्य (1 बाइट): नोड पता
- कमांड (1 बाइट): 0x01=read, 0x02=write
- लंबाई (2 बाइट): पेलोड आकार (big-endian)
- पेलोड (N बाइट): डेटा
- CRC (2 बाइट): Modbus CRC-16 (little-endian)

## परीक्षण

```bash
go test ./... -v
```
