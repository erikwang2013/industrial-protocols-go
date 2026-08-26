# SDK GatewayBridge PROFIBUS

Protocole PROFIBUS DP/PA via pont passerelle TCP.

## Matériel

| Passerelle | Interface | IP par défaut |
|---------|-----------|-------------|
| Anybus Communicator | Esclave PROFIBUS DP-V1 | 192.168.0.50 |
| Proxy Siemens CP 5611 | Maître PROFIBUS PCI/PCIe | attribuée par l'hôte |
| Passerelle de bus de terrain HMS | Anybus NP40 | 192.168.0.51 |

## Câblage

- DB9 femelle sur la passerelle vers le réseau PROFIBUS (ligne A verte, ligne B rouge)
- Résistance de terminaison ACTIVÉE aux deux extrémités du réseau
- Ethernet vers le port de gestion de la passerelle

## Configuration IP de la passerelle

```go
driver, codec, err := profibus.NewGatewayDriver("192.168.0.50:2000")
handler, err := profibus.ReadyHandler("192.168.0.50:2000")
```
