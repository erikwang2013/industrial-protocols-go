# LonWorks-GatewayBridge-SDK

LonWorks-Protokoll (ANSI/CEA-709.1) über eine TCP-Gateway-Brücke.

## Hardware

| Gateway | Schnittstelle | Standard-IP |
|---------|-----------|-------------|
| Echelon U60 | USB-FT-10-Netzwerkschnittstelle | vom Host vergeben |
| Echelon U70 | USB-TP/XF-1250-Schnittstelle | vom Host vergeben |
| Loytec L-IP | LonWorks/IP-Router | 192.168.0.90 |

## Verdrahtung

- FT-10 (Free Topology): polaritätsunabhängige Twisted-Pair-Verkabelung, bis zu 500 m freie Topologie
- TP/XF-1250: Bustopologie mit 105-Ohm-Abschluss
- Ethernet am L-IP-Router für die TCP-Brücke

## Gateway-IP-Konfiguration

```go
driver, codec, err := lonworks.NewGatewayDriver("192.168.0.90:2009")
handler, err := lonworks.ReadyHandler("192.168.0.90:2009")
```
