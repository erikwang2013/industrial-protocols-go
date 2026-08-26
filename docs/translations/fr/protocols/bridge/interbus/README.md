# SDK GatewayBridge Interbus

Protocole Interbus via pont passerelle TCP.

## Matériel

| Passerelle | Interface | IP par défaut |
|---------|-----------|-------------|
| Phoenix Contact IBS S5 DSC/I-T | Contrôleur maître Interbus | 192.168.0.60 |
| Phoenix Contact IBS PCI SC/I-T | Maître Interbus PCI | attribuée par l'hôte |
| HMS Anybus Interbus | Passerelle embarquée | 192.168.0.61 |

## Câblage

- D-SUB 9 broches sur la passerelle vers le bus distant Interbus (entrant/sortant)
- Blindage connecté à la FE (terre fonctionnelle) aux deux extrémités
- Ethernet vers la passerelle pour le pont TCP

## Configuration IP de la passerelle

```go
driver, codec, err := interbus.NewGatewayDriver("192.168.0.60:2006")
handler, err := interbus.ReadyHandler("192.168.0.60:2006")
```
