# Interbus-GatewayBridge-SDK

Interbus-Protokoll über eine TCP-Gateway-Brücke.

## Hardware

| Gateway | Schnittstelle | Standard-IP |
|---------|-----------|-------------|
| Phoenix Contact IBS S5 DSC/I-T | Interbus-Master-Controller | 192.168.0.60 |
| Phoenix Contact IBS PCI SC/I-T | PCI-Interbus-Master | vom Host vergeben |
| HMS Anybus Interbus | Eingebettetes Gateway | 192.168.0.61 |

## Verdrahtung

- 9-poliger D-SUB am Gateway zum Interbus-Remote-Bus (eingehend/ausgehend)
- Schirm an beiden Enden mit FE verbunden
- Ethernet zum Gateway für die TCP-Brücke

## Gateway-IP-Konfiguration

```go
driver, codec, err := interbus.NewGatewayDriver("192.168.0.60:2006")
handler, err := interbus.ReadyHandler("192.168.0.60:2006")
```
