# SDK du protocole Modbus Plus

Modbus Plus (MB+) est un réseau industriel à passage de jeton à haut débit développé par Modicon (Schneider Electric).

## Variantes

| Variante | Type de pont | Description |
|---------|-------------|-------------|
| `gateway` | GatewayBridge | Connexion TCP vers un adaptateur SA85/BM85 |
| `cmd` | CmdBridge | Enveloppe CLI pour l'utilitaire `sa85_cli` |

## Pilote de passerelle

```go
driver, codec, err := modbusplus.NewGatewayDriver("192.168.0.160:2010")
handler, err := modbusplus.ReadyHandler("192.168.0.160:2010")
```

## Pilote Cmd

```go
driver, codec, err := modbusplus.NewCmdDriver("/usr/bin/sa85_cli")
```

### Installation de l'outil CLI

```bash
# Install SA85 Modbus Plus driver and tools
# Refer to Schneider Electric documentation for SA85 adapter
```

Vérification : `sa85_cli --help`

## Matériel

| Passerelle | Interface | IP par défaut |
|---------|-----------|-------------|
| Schneider SA85 | Adaptateur ISA Modbus Plus | attribuée par l'hôte |
| Schneider BM85 | Pont/multiplexeur Modbus Plus | 192.168.0.A0 |
| ProSoft MVI56-MBP | Module MB+ ControlLogix | attribuée par l'hôte |

## Câblage

- Câble twinaxial (RG-62) avec connecteurs BNC pour le tronc MB+
- Résistance de terminaison (78 ohms à chaque extrémité)
- Le pont BM85 relie MB+ à TCP via Ethernet

## Format de trame du protocole

- Magic (2 octets) : `0x4D42`
- Destination (1 octet) : adresse du nœud
- Commande (1 octet) : 0x01=lecture, 0x02=écriture
- Longueur (2 octets) : taille de la charge utile (big-endian)
- Charge utile (N octets) : données
- CRC (2 octets) : Modbus CRC-16 (little-endian)

## Tests

```bash
go test ./... -v
```
