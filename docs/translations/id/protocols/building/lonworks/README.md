# SDK LonWorks GatewayBridge

Protokol LonWorks (ANSI/CEA-709.1) melalui jembatan gateway TCP.

## Perangkat Keras

| Gateway | Antarmuka | IP Default |
|---------|-----------|-------------|
| Echelon U60 | Antarmuka jaringan USB FT-10 | ditetapkan host |
| Echelon U70 | Antarmuka USB TP/XF-1250 | ditetapkan host |
| Loytec L-IP | Router LonWorks/IP | 192.168.0.90 |

## Kabel

- FT-10 (Free Topology): twisted pair tidak sensitif polaritas, hingga 500 m free topology
- TP/XF-1250: topologi bus dengan terminator 105 ohm
- Ethernet pada router L-IP untuk jembatan TCP

## Konfigurasi IP Gateway

```go
driver, codec, err := lonworks.NewGatewayDriver("192.168.0.90:2009")
handler, err := lonworks.ReadyHandler("192.168.0.90:2009")
```
