# SDK Interbus GatewayBridge

Protokol Interbus melalui jembatan gateway TCP.

## Perangkat Keras

| Gateway | Antarmuka | IP Default |
|---------|-----------|-------------|
| Phoenix Contact IBS S5 DSC/I-T | Kontroler master Interbus | 192.168.0.60 |
| Phoenix Contact IBS PCI SC/I-T | Master Interbus PCI | ditetapkan host |
| HMS Anybus Interbus | Gateway tertanam | 192.168.0.61 |

## Kabel

- D-SUB 9-pin pada gateway ke remote bus Interbus (masuk/keluar)
- Pelindung terhubung ke FE di kedua ujung
- Ethernet ke gateway untuk jembatan TCP

## Konfigurasi IP Gateway

```go
driver, codec, err := interbus.NewGatewayDriver("192.168.0.60:2006")
handler, err := interbus.ReadyHandler("192.168.0.60:2006")
```
