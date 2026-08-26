# SDK AS-Interface GatewayBridge

Protokol AS-Interface (ASi) melalui jembatan gateway TCP.

## Perangkat Keras

| Gateway | Antarmuka | IP Default |
|---------|-----------|-------------|
| Bihl+Wiedemann BWU2044 | Master ASi (Ethernet) | 192.168.0.30 |
| P+F VBG-ENX-K20-DMD-EV | Gateway ASi | 192.168.0.31 |
| ifm AC1375 | Kontroler ASi E | 192.168.0.32 |

## Kabel

- Kabel ASi kuning (daya + data) dari gateway ke slave
- Kabel daya bantu hitam (24 VDC untuk aktuator) opsional
- Ethernet pada gateway untuk jembatan TCP

## Konfigurasi IP Gateway

```go
driver, codec, err := asinterface.NewGatewayDriver("192.168.0.30:2003")
handler, err := asinterface.ReadyHandler("192.168.0.30:2003")
```
