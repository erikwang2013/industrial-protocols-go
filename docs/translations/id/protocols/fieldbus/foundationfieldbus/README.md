# SDK Foundation Fieldbus GatewayBridge

Foundation Fieldbus H1/HSE melalui jembatan gateway TCP.

## Perangkat Keras

| Gateway | Antarmuka | IP Default |
|---------|-----------|-------------|
| NI USB-8486 | Antarmuka USB H1 | ditetapkan host |
| Softing FFusb | Antarmuka USB H1 | ditetapkan host |
| P+F HD2-GTR-4PA | Gateway H1 ke Ethernet | 192.168.0.20 |

## Kabel

- Trunk H1 (twisted-pair, berpelindung) dengan terminator di kedua ujung
- Power conditioner fieldbus untuk H1 (24 VDC, 350-500 mA per segmen)
- Gateway menghubungkan segmen H1 ke jembatan TCP via Ethernet

## Konfigurasi IP Gateway

```go
driver, codec, err := foundationfieldbus.NewGatewayDriver("192.168.0.20:2002")
handler, err := foundationfieldbus.ReadyHandler("192.168.0.20:2002")
```
