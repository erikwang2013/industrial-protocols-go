# SDK du protocole CompactPCI

Implémente `kernel.Protocol` pour l'accès au bus CompactPCI via sysfs sous Linux.

## Protocole

- **Nom** : `cpci`
- **Variantes** : `cpci`
- **Port par défaut** : 0 (bus mappé en mémoire)
- **Transport** : sysfs `/sys/bus/pci/devices/<BDF>/config`

CompactPCI (PICMG 2.0) utilise la même interface électrique et logicielle
que le PCI conventionnel. Le codec fait transiter les octets bruts entre
l'application et l'espace de configuration PCI via sysfs.

## Exigences du noyau

Les modules noyau suivants doivent être chargés :

- `pcieport` -- pilote de port PCI Express (pour les systèmes CPCIe hybrides)
- `pci_sysfs` -- interface PCI sysfs (intégrée dans la plupart des noyaux)
- `cpci_hotplug` -- contrôleur d'insertion à chaud CompactPCI (optionnel, pour le hot-swap)

Le système de fichiers sysfs doit être monté sur `/sys`. C'est le comportement
par défaut sur toutes les distributions Linux modernes.

Permissions requises :
- Accès root ou `CAP_SYS_ADMIN` pour l'accès à l'espace de configuration
- Le fichier de configuration appartient à root:root avec le mode 0600

## Pilote

```go
d, err := cpci.NewCPCIDriver("0000:02:00.0")
if err != nil {
    log.Fatal(err)
}
defer d.Close()

tr := d.Transport()
// Use tr.Read/tr.Write for raw CPCI config space access
```

## CompactPCI vs PCI

CompactPCI utilise l'énumération de bus et l'espace de configuration PCI
standard. Principales différences avec le PCI de bureau :

- **Format 3U/6U** : mécanique Eurocard avec connecteurs à broches et douilles
- **Numérotation des bus** : chaque segment de châssis CPCI reçoit son propre numéro de bus PCI
- **Insertion à chaud** : le hot swap PICMG 2.1 utilise le modèle d'hotplug PCI standard
- **Emplacement système** : bus 0, périphérique 0 est le contrôleur de l'emplacement système

## Format BDF

L'adresse de bus est une chaîne BDF (Bus:Device.Function) :
- `0000:02:00.0` -- domaine 0000, bus 02, périphérique 00, fonction 0
- `0000:02:08.0` -- domaine 0000, bus 02, périphérique 08, fonction 0
