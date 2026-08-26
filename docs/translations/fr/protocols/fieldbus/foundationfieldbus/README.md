# SDK GatewayBridge Foundation Fieldbus

Foundation Fieldbus H1/HSE via pont passerelle TCP.

## Matériel

| Passerelle | Interface | IP par défaut |
|---------|-----------|-------------|
| NI USB-8486 | Interface USB H1 | attribuée par l'hôte |
| Softing FFusb | Interface USB H1 | attribuée par l'hôte |
| P+F HD2-GTR-4PA | Passerelle H1 vers Ethernet | 192.168.0.20 |

## Câblage

- Tronc H1 (paire torsadée, blindée) avec terminaison aux deux extrémités
- Conditionneur d'alimentation de bus de terrain pour H1 (24 VDC, 350-500 mA par segment)
- La passerelle relie le segment H1 au pont TCP via Ethernet

## Configuration IP de la passerelle

```go
driver, codec, err := foundationfieldbus.NewGatewayDriver("192.168.0.20:2002")
handler, err := foundationfieldbus.ReadyHandler("192.168.0.20:2002")
```
