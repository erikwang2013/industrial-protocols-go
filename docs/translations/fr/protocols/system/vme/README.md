# SDK du protocole VME / VPX

Implémente `kernel.Protocol` pour l'accès aux bus VME et VPX via procfs sous Linux.

## Protocole

- **Nom** : `vme`
- **Variantes** : `vme`
- **Port par défaut** : 0 (bus mappé en mémoire)
- **Transport** : procfs `/proc/vme/<slot>`

Le codec fait transiter les octets bruts entre l'application et l'espace
d'adressage VME. Les lectures et écritures vont directement au bus via
l'interface procfs fournie par le pilote noyau VME.

## Exigences du noyau

Le module noyau suivant doit être chargé :

- `vme_tsi148` -- pilote de pont VME Tundra TSI148 (le plus courant)
  - Également pris en charge : `vme_ca91cx42` (Universe II), `vme_user`

Le système de fichiers procfs doit être monté sur `/proc`.

Permissions requises :
- L'accès root est requis pour ouvrir `/proc/vme/<slot>` en lecture-écriture
- Les nœuds de périphérique appartiennent à root:root

## Pilote

```go
d, err := vme.NewVMEDriver(0) // slot 0
if err != nil {
    log.Fatal(err)
}
defer d.Close()

tr := d.Transport()
// Use tr.Read/tr.Write for raw VME bus access
```

## Modes d'adressage VME

Le codec fait transiter tous les octets de modificateur d'adresse et d'adresse
de manière transparente. Les applications doivent préfixer les informations
d'adressage à la charge utile des données :

- **A16** : espace d'adressage E/S court 16 bits
- **A24** : espace d'adressage standard 24 bits
- **A32** : espace d'adressage étendu 32 bits

## Compatibilité VPX

Les systèmes VPX (VITA 46) qui exposent une interface procfs compatible VME
peuvent utiliser ce pilote. La numérotation des slots est identique.
