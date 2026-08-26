# Foundation-Fieldbus-GatewayBridge-SDK

Foundation Fieldbus H1/HSE über eine TCP-Gateway-Brücke.

## Hardware

| Gateway | Schnittstelle | Standard-IP |
|---------|-----------|-------------|
| NI USB-8486 | USB-H1-Schnittstelle | vom Host vergeben |
| Softing FFusb | USB-H1-Schnittstelle | vom Host vergeben |
| P+F HD2-GTR-4PA | H1-zu-Ethernet-Gateway | 192.168.0.20 |

## Verdrahtung

- H1-Stammbus (Twisted-Pair, geschirmt) mit Abschluss an beiden Enden
- Feldbus-Spannungsversorgung für H1 (24 VDC, 350-500 mA pro Segment)
- Das Gateway verbindet das H1-Segment über Ethernet mit der TCP-Brücke

## Gateway-IP-Konfiguration

```go
driver, codec, err := foundationfieldbus.NewGatewayDriver("192.168.0.20:2002")
handler, err := foundationfieldbus.ReadyHandler("192.168.0.20:2002")
```
