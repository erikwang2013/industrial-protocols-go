# SDK du protocole PCI / PCIe

Implémente `kernel.Protocol` pour l'accès aux bus PCI et PCI Express via sysfs sous Linux.

## Protocole

- **Nom** : `pci`
- **Variantes** : `pci`
- **Port par défaut** : 0 (bus mappé en mémoire)
- **Transport** : sysfs `/sys/bus/pci/devices/<BDF>/config`

Le codec fait transiter les octets bruts entre l'application et l'espace
de configuration PCI. Les lectures et écritures vont directement vers les
registres de configuration du périphérique via le fichier de configuration sysfs.

## Exigences du noyau

Les modules noyau suivants doivent être chargés :

- `pcieport` -- pilote de port PCI Express
- `pci_sysfs` -- interface PCI sysfs (intégrée dans la plupart des noyaux)

Le système de fichiers sysfs doit être monté sur `/sys`. C'est le comportement
par défaut sur toutes les distributions Linux modernes.

Permissions requises :
- Accès root ou `CAP_SYS_ADMIN` pour l'accès à l'espace de configuration
- Le fichier de configuration appartient à root:root avec le mode 0600

## Pilote

```go
d, err := pci.NewPCIDriver("0000:00:1f.3")
if err != nil {
    log.Fatal(err)
}
defer d.Close()

tr := d.Transport()
// Use tr.Read/tr.Write for raw PCI config space access
```

## Format BDF

L'adresse de bus est une chaîne BDF (Bus:Device.Function) :
- `0000:00:1f.3` -- domaine 0000, bus 00, périphérique 1f, fonction 3
- `0000:01:00.0` -- domaine 0000, bus 01, périphérique 00, fonction 0
