# SDK CC-Link IE Field GatewayBridge

Protokol CC-Link IE Field melalui jembatan gateway TCP.

## Perangkat Keras

| Gateway | Antarmuka | IP Default |
|---------|-----------|-------------|
| Mitsubishi MELSEC IQ-R | Master CC-Link IE Field | 192.168.3.40 |
| Mitsubishi RJ71GF11-T2 | Modul CC-Link IE Field | ditetapkan host |
| HMS Anybus CC-Link IE | Gateway tertanam | 192.168.0.52 |

## Kabel

- RJ45 Ethernet untuk CC-Link IE Field (topologi cincin atau bintang 1 Gbps)
- Port manajemen di jaringan terpisah
- Gateway menghubungkan jaringan lapangan ke jembatan TCP

## Konfigurasi IP Gateway

```go
driver, codec, err := cclinkie.NewGatewayDriver("192.168.3.40:2005")
handler, err := cclinkie.ReadyHandler("192.168.3.40:2005")
```
