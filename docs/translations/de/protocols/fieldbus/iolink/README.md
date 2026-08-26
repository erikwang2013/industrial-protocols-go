# IO-Link-GatewayBridge-SDK

IO-Link-Protokoll über eine TCP-Gateway-Brücke.

## Hardware

| Gateway | Schnittstelle | Standard-IP |
|---------|-----------|-------------|
| ifm AL1332 | IO-Link-Master (EtherNet/IP) | 192.168.0.40 |
| Balluff BNI00AZ | IO-Link-Master (PROFINET) | 192.168.0.41 |
| SICK SIG200 | IO-Link-Master (Ethernet) | 192.168.0.42 |

## Verdrahtung

- M12-Steckverbinder (4-polig) für jeden IO-Link-Port: L+ (braun), L− (blau), C/Q (schwarz), ungenutzt (weiß)
- 24-VDC-Stromversorgung für Master und Geräte
- Ethernet am Master für die TCP-Brücke

## Gateway-IP-Konfiguration

```go
driver, codec, err := iolink.NewGatewayDriver("192.168.0.40:2004")
handler, err := iolink.ReadyHandler("192.168.0.40:2004")
```
