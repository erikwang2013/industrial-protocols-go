# AS-Interface-GatewayBridge-SDK

AS-Interface-Protokoll (ASi) über eine TCP-Gateway-Brücke.

## Hardware

| Gateway | Schnittstelle | Standard-IP |
|---------|-----------|-------------|
| Bihl+Wiedemann BWU2044 | ASi-Master (Ethernet) | 192.168.0.30 |
| P+F VBG-ENX-K20-DMD-EV | ASi-Gateway | 192.168.0.31 |
| ifm AC1375 | ASi-ControllerE | 192.168.0.32 |

## Verdrahtung

- Gelbes ASi-Kabel (Strom + Daten) vom Gateway zu den Slaves
- Optionales schwarzes Hilfsspannungskabel (24 VDC für Aktuatoren)
- Ethernet am Gateway für die TCP-Brücke

## Gateway-IP-Konfiguration

```go
driver, codec, err := asinterface.NewGatewayDriver("192.168.0.30:2003")
handler, err := asinterface.ReadyHandler("192.168.0.30:2003")
```
