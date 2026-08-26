# SDK GatewayBridge WorldFIP

Protocole WorldFIP via pont passerelle TCP.

## Matériel

| Passerelle | Interface | IP par défaut |
|---------|-----------|-------------|
| Agent FIPIO | Agent de bus de terrain WorldFIP | 192.168.0.70 |
| Passerelle FIP (Alstom) | WorldFIP vers Ethernet | 192.168.0.71 |
| NI FIP-USB | Interface USB WorldFIP | attribuée par l'hôte |

## Câblage

- D-SUB 9 broches sur la passerelle vers le tronc WorldFIP (FIP1 = Data+, FIP2 = Data-)
- Terminaison de ligne (120 ohms) aux deux extrémités du bus
- Ethernet pour le pont TCP

## Configuration IP de la passerelle

```go
driver, codec, err := worldfip.NewGatewayDriver("192.168.0.70:2007")
handler, err := worldfip.ReadyHandler("192.168.0.70:2007")
```
