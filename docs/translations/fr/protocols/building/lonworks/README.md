# SDK GatewayBridge LonWorks

Protocole LonWorks (ANSI/CEA-709.1) via pont passerelle TCP.

## Matériel

| Passerelle | Interface | IP par défaut |
|---------|-----------|-------------|
| Echelon U60 | Interface réseau USB FT-10 | attribuée par l'hôte |
| Echelon U70 | Interface USB TP/XF-1250 | attribuée par l'hôte |
| Loytec L-IP | Routeur LonWorks/IP | 192.168.0.90 |

## Câblage

- FT-10 (Free Topology) : paire torsadée insensible à la polarité, jusqu'à 500 m en topologie libre
- TP/XF-1250 : topologie en bus avec terminaison de 105 ohms
- Ethernet sur le routeur L-IP pour le pont TCP

## Configuration IP de la passerelle

```go
driver, codec, err := lonworks.NewGatewayDriver("192.168.0.90:2009")
handler, err := lonworks.ReadyHandler("192.168.0.90:2009")
```
