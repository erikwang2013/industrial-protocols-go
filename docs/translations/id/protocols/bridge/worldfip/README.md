# SDK WorldFIP GatewayBridge

Protokol WorldFIP melalui jembatan gateway TCP.

## Perangkat Keras

| Gateway | Antarmuka | IP Default |
|---------|-----------|-------------|
| FIPIO Agent | Agen fieldbus WorldFIP | 192.168.0.70 |
| FIP Gateway (Alstom) | WorldFIP ke Ethernet | 192.168.0.71 |
| NI FIP-USB | Antarmuka WorldFIP USB | ditetapkan host |

## Kabel

- D-SUB 9-pin pada gateway ke trunk WorldFIP (FIP1 = Data+, FIP2 = Data-)
- Terminator saluran (120 ohm) di kedua ujung bus
- Ethernet untuk jembatan TCP

## Konfigurasi IP Gateway

```go
driver, codec, err := worldfip.NewGatewayDriver("192.168.0.70:2007")
handler, err := worldfip.ReadyHandler("192.168.0.70:2007")
```
