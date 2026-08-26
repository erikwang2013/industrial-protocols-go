# Rapport de test — industrial-protocols-go

Date : 2026-08-27
Portée : espace de travail Go multi-modules (go.work, 42 modules : kernel / protocols / examples)
Commandes : `go test ./...` et `go test -cover ./...` (exécutées de manière équivalente module par module, faute de go.mod à la racine de l'espace de travail ; résultats identiques)

## 1. Conclusion générale

- **Tout est vert** : les 42 modules et 52 packages contenant des tests passent tous (vérifié par plusieurs exécutions complètes, y compris pendant l'écriture de tests en parallèle).
- 106 fichiers de test, 659 fonctions de test.
- `go vet ./...` sans aucun avertissement sur tous les modules.
- Couverture moyenne des instructions **83,3 %** ; couverture de 95 % à 100 % pour les packages du noyau kernel.

## 2. Statistiques de test et couverture par module

### kernel (11 packages, couverture la plus élevée hors examples)

| Package | Couverture |
|---|---|
| kernel | 100% |
| kernel/bridge | 80.0% |
| kernel/config | 100% |
| kernel/connection | 95.3% |
| kernel/event | 100% |
| kernel/gateway | 100% |
| kernel/metrics | 100% |
| kernel/pipeline | 100% |
| kernel/security | 100% |
| kernel/transport | 100% |
| kernel/vendor | 100% |

### Modules protocols principaux

| Module | Couverture | Module | Couverture |
|---|---|---|---|
| automotive/flexray | 79.4% | fieldbus/canopen | 73.9% |
| automotive/kline | 93.3% | fieldbus/cclink | 93.8% |
| automotive/lin | 92.5% | fieldbus/cclinkie | 73.5% |
| automotive/most | 100% | fieldbus/devicenet | 73.2% |
| automotive/saej1850 | 84.2% | fieldbus/dnp3 | 88.9% |
| bridge/controlnet | 100% | fieldbus/hart | 96.6% |
| bridge/ethercat | 100% | fieldbus/iec61850 | 91.8% |
| bridge/interbus | 70.0% | iot/hartip | 88.6% |
| bridge/isa100 | 95.2% | iot/mqtt | 89.1% |
| bridge/powerlink | 97.8% | ethernet/bacnet | 95.2% |
| bridge/sercos | 97.9% | ethernet/ethernetip | 96.9% |
| bridge/sercos1 | 95.6% | ethernet/modbus | 67.7% |
| bridge/safej1850 | 100% | ethernet/opcua | 95.8% |
| bridge/wirelesshart | 95.2% | ethernet/profinet | 95.4% |
| building/dali | 92.9% | system/cpci | 78.6% |
| building/lonworks | **37.9%** | system/pci | 78.6% |
| fieldbus/asinterface | **37.9%** | system/vme | 78.6% |
| fieldbus/foundationfieldbus | **37.9%** | fieldbus/profibus | **35.7%** |
| fieldbus/iolink | **37.9%** | bridge/lightbus / modbusplus / worldfip | 72~73.5% |

> examples/modbus_basic n'a pas de fichier de test (programme d'exemple, comportement attendu).

## 3. Liste des corrections

### 3.1 Corrections de bugs dans le code source (découverts et vérifiés par tester-kernel / tester-protocols)

| Fichier:ligne | Problème | Correctif |
|---|---|---|
| kernel/gateway/engine.go:27 | `Transform` provoquait un panic par déréférencement nil pour les règles dont `Src`/`Dst`/`Map` étaient nil, aux appels `rule.Src.Name()` / `rule.Dst.Name()` / `rule.Map(...)` | Les règles incomplètement configurées sont désormais ignorées au début de la boucle (conforme à la sémantique existante « ignorer les règles invalides ») |
| protocols/ethernet/modbus/modbus.go:215 | `decodePDU` provoquait un panic hors limites en lisant `pdu[1]` pour un PDU d'exception d'un octet (p. ex. `0x81`) | Ajout d'une vérification `len(pdu) < 2` et retour d'une erreur d'analyse |
| protocols/ethernet/modbus/modbus.go:230 | `decodePDU` provoquait un panic par dépassement de la tranche lorsque le compteur d'octets (`pdu[1]`) dépassait les données réelles (trames malveillantes/corrompues) | Ajout d'une vérification des bornes `2+n > len(pdu)` et retour d'une erreur d'analyse |
| protocols/ethernet/opcua/opcua.go:63,74 | `encodeHello` réécrivait la longueur du message au décalage 8 (emplacement de la version du protocole), alors que la spécification exige le décalage 4 | Enregistrement préalable de `sizePos` (décalage 4) et réécriture au bon emplacement |
| protocols/ethernet/opcua/opcua.go:82,89 | Le champ de longueur de `encodeOpenSecureChannel` n'était jamais réécrit et restait à 0 | Réécriture en fin de trame selon la longueur réelle de la trame |
| protocols/ethernet/ethernetip/ethernetip.go:111,117 | `encodeReadTag` allouait 16+len(cipReq) octets mais écrivait la requête CIP à partir de [28:] ; les 12 octets manquants faisaient silencieusement perdre toute la requête CIP dans `copy` ; le champ de longueur était lui aussi en retrait de 12 | Allocation portée à 28+len(cipReq) et champ length défini sur 4+len(cipReq) |
| protocols/ethernet/profinet/profinet.go:72 | `Decode` provoquait un panic de dépassement de tranche pour un `blockLen` hors limites | Ajout de la vérification `len(data) < 10+blockLen` |
| protocols/fieldbus/canopen/canopen.go:102 | La vérification de longueur minimale de `Decode` était de 4, mais `unmarshalCAN` requiert `data[4:8]`, d'où un panic pour les trames de 4 à 7 octets | Longueur minimale portée à 8 (format de ligne fixé à 8 octets) |
| protocols/fieldbus/dnp3/dnp3.go:121 | `decodeTPDU` provoquait un panic de dépassement de tranche pour `length<5` (`apduEnd < apduStart`) ou une `length` excessive (`apduEnd > len(data)`) | Validation unifiée `length < 5 || apduEnd > len(data)` puis retour d'une erreur |
| protocols/fieldbus/dnp3/dnp3.go:157 | La variable de travail de `crc16DNP` n'appliquait pas le masque `& 0xFF` prévu par la spécification ; le vecteur connu `crc16DNP("123456789")` donnait 0x69FF (devrait être 0xEA82) | `temp := (crc ^ b) & 0xFF` |
| protocols/automotive/kline/kline.go:74 | La vérification `len<5` de `Decode` précédait l'évaluation de la trame de synchronisation fast-init, rejetant à tort la trame de synchronisation d'un octet (0x55) | Traiter d'abord les trames vides et de synchronisation, puis vérifier la longueur |

### 3.2 Corrections de tests (erreurs de compilation / assertions erronées / doublons / assertions obsolètes)

| Fichier:ligne | Problème | Correctif |
|---|---|---|
| protocols/automotive/kline/kline_extra_test.go:43 | L'affectation de la constante `0x68+0x10+0xF1+0x01` (362) à un byte ne compilait pas | Attente remplacée par les 8 bits de poids faible de la somme de contrôle, `0x6A` |
| protocols/automotive/saej1850/saej1850_extra_test.go | Appel de méthodes/champs non exportés sur l'interface `kernel.Codec`, échec de compilation | Ajout d'un helper `newCodec` affirmant le type `*j1850Codec` ; l'assertion de `TestEncodeMode01PID` porte désormais sur raw[4..6] (ID de 4 octets + 4 premiers octets de données) |
| protocols/bridge/sercos1/sercos1_extra_test.go:45 | `TestEncodeReadDefaults` dupliquait un test existant, échec de compilation | Renommé en `TestEncodeReadDefaultsFiber` (la couverture de la variante fiber est conservée) |
| protocols/building/dali/dali_extra_test.go:68 | `TestEncodeDirectArc` dupliquait un test existant, échec de compilation | Renommé en `TestEncodeDirectArcCmd` |
| protocols/bridge/powerlink/powerlink_extra_test.go:51 | L'attente était `write 0x1000 A`, la valeur réelle étant `0A` (`%X` fixe deux chiffres par octet, conformément à la convention `DEADBEEF` des tests existants) | Attente remplacée par `0A` |
| protocols/bridge/sercos/sercos_extra_test.go:62 | Comme ci-dessus, l'attente `1` devait être `01` | Attente remplacée par `01` |
| protocols/ethernet/opcua/opcua_extra_test.go:35 | Assertion de l'ancien comportement bugué (longueur au décalage 8), devenue obsolète après le correctif du code source | L'assertion vérifie désormais que le champ de version du protocole vaut 0 (la longueur est couverte par TestExtraHELMessageSizeAtOffset4) |
| protocols/fieldbus/canopen/canopen_extra_test.go:135 | Assertion inversée dans le test « prove-it » (échec quand `err != nil`, alors qu'une erreur d'analyse est précisément le résultat attendu) | Échec désormais quand `err == nil` |
| protocols/ethernet/profinet/profinet_extra_test.go:23 | Comme ci-dessus, assertion inversée | Échec désormais quand `err == nil` |
| protocols/fieldbus/dnp3/dnp3_extra_test.go:35,53 | Comme ci-dessus ×2 | Échec désormais quand `err == nil` |

> Note : le problème d'implémentation de l'interface `stubTransport` dans `kernel/security/tls_test.go` et les valeurs CRC attendues de `cclink/cclinkie` ont été corrigés par l'ingénieur de test dans le cadre de travaux parallèles ; ils n'ont pas été modifiés ici.

## 4. Risques résiduels

1. **Modules à faible couverture** : profibus (35,7 %), lonworks / asinterface / foundationfieldbus / iolink (37,9 %) — les tests ne couvrent que quelques chemins ; il est recommandé d'ajouter ultérieurement des branches Decode/Encode et des chemins d'erreur.
2. **CRC Modbus RTU non vérifié** : `decodeRTU` ne vérifie pas le CRC (le test `TestExtraRTUDecodesCorruptCRC` documente explicitement cette lacune et passe en l'état). Il est recommandé de l'ajouter pour l'interopérabilité avec de vrais équipements.
3. **Limitation du format de ligne CANopen** : `marshalCAN` ne transporte que les 4 premiers octets de données (documenté par `TestExtraSDOWritePayloadLostOnWire`) ; les charges utiles SDO d'écriture multi-octets sont perdues.
4. **Sauts liés au matériel** : cpci / pci / vme utilisent `t.Skip` en l'absence de matériel (comportement raisonnable).
5. **Fichiers temporaires résiduels dans le dépôt** : les fichiers non suivis à la racine — `main.go` (programme de sonde référençant un module `probe/` inexistant), `scripts/`, `docs/*.png` — ne font pas partie de cette livraison ; leur nettoyage est recommandé.
6. **Les tests sont encore en cours d'ajout en parallèle** : ce rapport s'appuie sur un instantané de la dernière exécution complète verte ; si l'ingénieur de test continue de soumettre des `*_extra_test.go`, une nouvelle régression complète sera nécessaire.
