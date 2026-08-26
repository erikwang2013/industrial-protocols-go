# Lightbus-GatewayBridge-SDK

Beckhoff-Lightbus-Glasfaserprotokoll über eine TCP-Gateway-Brücke.

## Hardware

| Gateway | Schnittstelle | Standard-IP |
|---------|-----------|-------------|
| Beckhoff FC2001 | Lightbus-PCI-Karte | vom Host vergeben |
| Beckhoff BK2000 | Lightbus-Buskuppler | 192.168.0.80 |
| Beckhoff FC9001 | Lightbus-Ethernet-Adapter | 192.168.0.81 |

## Verdrahtung

- Ringtopologie mit Kunststoff-Lichtwellenleiter (POF)
- FC2001/FC9001 verbindet den Ring über Ethernet mit der TCP-Brücke
- Jedes Gerät besitzt TX- und RX-Glasfaseranschlüsse

## Gateway-IP-Konfiguration

```go
driver, codec, err := lightbus.NewGatewayDriver("192.168.0.80:2008")
handler, err := lightbus.ReadyHandler("192.168.0.80:2008")
```
