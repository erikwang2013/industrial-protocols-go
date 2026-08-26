# WorldFIP-GatewayBridge-SDK

WorldFIP-Protokoll über eine TCP-Gateway-Brücke.

## Hardware

| Gateway | Schnittstelle | Standard-IP |
|---------|-----------|-------------|
| FIPIO Agent | WorldFIP-Feldbus-Agent | 192.168.0.70 |
| FIP Gateway (Alstom) | WorldFIP zu Ethernet | 192.168.0.71 |
| NI FIP-USB | USB-WorldFIP-Schnittstelle | vom Host vergeben |

## Verdrahtung

- 9-poliger D-SUB am Gateway zum WorldFIP-Stammbus (FIP1 = Data+, FIP2 = Data−)
- Leitungsabschluss (120 Ohm) an beiden Busenden
- Ethernet für die TCP-Brücke

## Gateway-IP-Konfiguration

```go
driver, codec, err := worldfip.NewGatewayDriver("192.168.0.70:2007")
handler, err := worldfip.ReadyHandler("192.168.0.70:2007")
```
