# SDK GatewayBridge IO-Link

Protocole IO-Link via pont passerelle TCP.

## Matériel

| Passerelle | Interface | IP par défaut |
|---------|-----------|-------------|
| ifm AL1332 | Maître IO-Link (EtherNet/IP) | 192.168.0.40 |
| Balluff BNI00AZ | Maître IO-Link (PROFINET) | 192.168.0.41 |
| SICK SIG200 | Maître IO-Link (Ethernet) | 192.168.0.42 |

## Câblage

- Connecteur M12 (4 broches) pour chaque port IO-Link : L+ (marron), L- (bleu), C/Q (noir), inutilisé (blanc)
- Alimentation 24 VDC pour le maître et les équipements
- Ethernet sur le maître pour le pont TCP

## Configuration IP de la passerelle

```go
driver, codec, err := iolink.NewGatewayDriver("192.168.0.40:2004")
handler, err := iolink.ReadyHandler("192.168.0.40:2004")
```
