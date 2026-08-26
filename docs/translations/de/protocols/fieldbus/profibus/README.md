# PROFIBUS-GatewayBridge-SDK

PROFIBUS-DP/PA-Protokoll über eine TCP-Gateway-Brücke.

## Hardware

| Gateway | Schnittstelle | Standard-IP |
|---------|-----------|-------------|
| Anybus Communicator | PROFIBUS-DP-V1-Slave | 192.168.0.50 |
| Siemens-CP-5611-Proxy | PCI/PCIe-PROFIBUS-Master | vom Host vergeben |
| HMS Fieldbus Gateway | Anybus NP40 | 192.168.0.51 |

## Verdrahtung

- DB9-Buchse am Gateway zum PROFIBUS-Netzwerk (A-Leitung grün, B-Leitung rot)
- Abschlusswiderstand an beiden Netzwerkenden EIN
- Ethernet zum Verwaltungsport des Gateways

## Gateway-IP-Konfiguration

```go
driver, codec, err := profibus.NewGatewayDriver("192.168.0.50:2000")
handler, err := profibus.ReadyHandler("192.168.0.50:2000")
```
