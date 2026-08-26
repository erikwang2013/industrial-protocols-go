# SDK Lightbus GatewayBridge

Protokol serat optik Beckhoff Lightbus melalui jembatan gateway TCP.

## Perangkat Keras

| Gateway | Antarmuka | IP Default |
|---------|-----------|-------------|
| Beckhoff FC2001 | Kartu PCI Lightbus | ditetapkan host |
| Beckhoff BK2000 | Bus coupler Lightbus | 192.168.0.80 |
| Beckhoff FC9001 | Adaptor Ethernet Lightbus | 192.168.0.81 |

## Kabel

- Topologi cincin serat optik plastik (POF)
- FC2001/FC9001 menghubungkan cincin ke jembatan TCP via Ethernet
- Setiap perangkat memiliki konektor serat TX dan RX

## Konfigurasi IP Gateway

```go
driver, codec, err := lightbus.NewGatewayDriver("192.168.0.80:2008")
handler, err := lightbus.ReadyHandler("192.168.0.80:2008")
```
