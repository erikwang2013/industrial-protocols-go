# SDK PROFIBUS GatewayBridge

Protokol PROFIBUS DP/PA melalui jembatan gateway TCP.

## Perangkat Keras

| Gateway | Antarmuka | IP Default |
|---------|-----------|-------------|
| Anybus Communicator | Slave PROFIBUS DP-V1 | 192.168.0.50 |
| Proxy Siemens CP 5611 | Master PROFIBUS PCI/PCIe | ditetapkan host |
| HMS Fieldbus Gateway | Anybus NP40 | 192.168.0.51 |

## Kabel

- DB9 female pada gateway ke jaringan PROFIBUS (jalur A hijau, jalur B merah)
- Resistor terminasi ON di kedua ujung jaringan
- Ethernet ke port manajemen gateway

## Konfigurasi IP Gateway

```go
driver, codec, err := profibus.NewGatewayDriver("192.168.0.50:2000")
handler, err := profibus.ReadyHandler("192.168.0.50:2000")
```
