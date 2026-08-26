# SDK GatewayBridge AS-Interface

Protocole AS-Interface (ASi) via pont passerelle TCP.

## Matériel

| Passerelle | Interface | IP par défaut |
|---------|-----------|-------------|
| Bihl+Wiedemann BWU2044 | Maître ASi (Ethernet) | 192.168.0.30 |
| P+F VBG-ENX-K20-DMD-EV | Passerelle ASi | 192.168.0.31 |
| ifm AC1375 | Contrôleur ASi E | 192.168.0.32 |

## Câblage

- Câble ASi jaune (alimentation + données) de la passerelle vers les esclaves
- Câble d'alimentation auxiliaire noir (24 VDC pour les actionneurs) en option
- Ethernet sur la passerelle pour le pont TCP

## Configuration IP de la passerelle

```go
driver, codec, err := asinterface.NewGatewayDriver("192.168.0.30:2003")
handler, err := asinterface.ReadyHandler("192.168.0.30:2003")
```
