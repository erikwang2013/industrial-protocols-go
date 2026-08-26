# SDK IO-Link GatewayBridge

Protokol IO-Link melalui jembatan gateway TCP.

## Perangkat Keras

| Gateway | Antarmuka | IP Default |
|---------|-----------|-------------|
| ifm AL1332 | IO-Link Master (EtherNet/IP) | 192.168.0.40 |
| Balluff BNI00AZ | IO-Link Master (PROFINET) | 192.168.0.41 |
| SICK SIG200 | IO-Link Master (Ethernet) | 192.168.0.42 |

## Kabel

- Konektor M12 (4 pin) untuk setiap port IO-Link: L+ (coklat), L- (biru), C/Q (hitam), tidak terpakai (putih)
- Catu daya 24 VDC untuk master dan perangkat
- Ethernet pada master untuk jembatan TCP

## Konfigurasi IP Gateway

```go
driver, codec, err := iolink.NewGatewayDriver("192.168.0.40:2004")
handler, err := iolink.ReadyHandler("192.168.0.40:2004")
```
