# SDK GatewayBridge Lightbus

Protocole à fibre optique Beckhoff Lightbus via pont passerelle TCP.

## Matériel

| Passerelle | Interface | IP par défaut |
|---------|-----------|-------------|
| Beckhoff FC2001 | Carte PCI Lightbus | attribuée par l'hôte |
| Beckhoff BK2000 | Coupleur de bus Lightbus | 192.168.0.80 |
| Beckhoff FC9001 | Adaptateur Ethernet Lightbus | 192.168.0.81 |

## Câblage

- Topologie en anneau en fibre optique plastique (POF)
- FC2001/FC9001 relie l'anneau au pont TCP via Ethernet
- Chaque équipement possède des connecteurs fibre TX et RX

## Configuration IP de la passerelle

```go
driver, codec, err := lightbus.NewGatewayDriver("192.168.0.80:2008")
handler, err := lightbus.ReadyHandler("192.168.0.80:2008")
```
