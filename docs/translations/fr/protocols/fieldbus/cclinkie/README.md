# SDK GatewayBridge CC-Link IE Field

Protocole CC-Link IE Field via pont passerelle TCP.

## Matériel

| Passerelle | Interface | IP par défaut |
|---------|-----------|-------------|
| Mitsubishi MELSEC IQ-R | Maître CC-Link IE Field | 192.168.3.40 |
| Mitsubishi RJ71GF11-T2 | Module CC-Link IE Field | attribuée par l'hôte |
| HMS Anybus CC-Link IE | Passerelle embarquée | 192.168.0.52 |

## Câblage

- RJ45 Ethernet pour CC-Link IE Field (topologie en anneau 1 Gbps ou en étoile)
- Port de gestion sur un réseau séparé
- La passerelle relie le réseau de terrain au pont TCP

## Configuration IP de la passerelle

```go
driver, codec, err := cclinkie.NewGatewayDriver("192.168.3.40:2005")
handler, err := cclinkie.ReadyHandler("192.168.3.40:2005")
```
