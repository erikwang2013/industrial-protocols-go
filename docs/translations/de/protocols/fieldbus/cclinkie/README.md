# CC-Link-IE-Field-GatewayBridge-SDK

CC-Link-IE-Field-Protokoll über eine TCP-Gateway-Brücke.

## Hardware

| Gateway | Schnittstelle | Standard-IP |
|---------|-----------|-------------|
| Mitsubishi MELSEC IQ-R | CC-Link-IE-Field-Master | 192.168.3.40 |
| Mitsubishi RJ71GF11-T2 | CC-Link-IE-Field-Modul | vom Host vergeben |
| HMS Anybus CC-Link IE | Eingebettetes Gateway | 192.168.0.52 |

## Verdrahtung

- RJ45-Ethernet für CC-Link IE Field (1-Gbps-Ring- oder Sterntopologie)
- Verwaltungsport in einem separaten Netzwerk
- Das Gateway verbindet das Feldnetzwerk mit der TCP-Brücke

## Gateway-IP-Konfiguration

```go
driver, codec, err := cclinkie.NewGatewayDriver("192.168.3.40:2005")
handler, err := cclinkie.ReadyHandler("192.168.3.40:2005")
```
