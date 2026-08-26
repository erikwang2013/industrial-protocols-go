[English](../../README.en.md) | [中文](../../README.md) | [한국어](../ko/README.md) | [Русский](../ru/README.md) | [Deutsch](../de/README.md) | Français | [Español](../es/README.md) | [Português](../pt/README.md) | [हिन्दी](../hi/README.md) | [العربية](../ar/README.md) | [বাংলা](../bn/README.md) | [Bahasa Indonesia](../id/README.md) | [日本語](../ja/README.md)

# Industrial Protocols Go

Ensemble de protocoles de communication réseau industrielle en Go — architecture en couches + middleware, couvrant 40 protocoles industriels : 14 implémentations purement logicielles + 26 SDK matériels (tous avec pilotes).

> Référence à l'implémentation PHP : [github.com/erikwang2013/industrial-protocols](https://github.com/erikwang2013/industrial-protocols)

---

## Conception du projet

### Architecture en couches

```
┌──────────────────────────────────────────────────┐
│                User Application                   │
├──────────────────────────────────────────────────┤
│  Pipeline Middleware                               │
│  Retry · Timeout · CircuitBreaker · Logger        │
├──────────────────────────────────────────────────┤
│  Kernel                                           │
│  ConnectionManager · ConfigRepository              │
│  GatewayEngine · Bridge · Vendor · Event          │
│  Metrics · Security (TLS)                         │
├──────────────────────────────────────────────────┤
│  Codec (Encode/Decode)           Transport        │
│  Module indépendant par protocole  (TCP/UDP/Serial)│
│                                                │
│  14 Pure-Soft   +   26 Hardware SDKs (with drivers)│
└──────────────────────────────────────────────────┘
```

### Philosophie de conception

**Micro-noyau + SDK de protocole.** Le noyau ne définit que les interfaces et les préoccupations transversales (pool de connexions, reprise, disjoncteur, événements, métriques), sans inclure aucune implémentation de protocole concrète. Chaque protocole est un module Go indépendant, importé selon les besoins.

**Découplage par couches.** Quatre niveaux d'abstraction :

| Couche | Responsabilité | Types principaux |
|----|------|---------|
| **Transport** | Canal de communication de bas niveau | Interface `Transport` — `TCPTransport`, `UDPTransport`, `PipeTransport` |
| **Codec** | Encodage/décodage du protocole | Interface `Codec` — `Encode(*Request) ([]byte, error)` / `Decode([]byte) (*Response, error)` |
| **Pipeline** | Middleware transversal | `Middleware` — `Timeout`, `Retry`, `CircuitBreaker`, `Logger` |
| **Kernel** | Cycle de vie des équipements | `ConnectionManager`, `ConfigRepository` |

**Les auteurs de protocoles n'ont qu'à implémenter les deux interfaces `Protocol` + `Codec`.** La couche de transport, le pool de connexions, la reprise, le délai d'attente et le disjoncteur sont tous réutilisés depuis le kernel.

### Modules du noyau

| Module | Chemin | Description |
|------|------|------|
| ConnectionManager | `kernel/connection/` | Enregistrement des équipements, pool de connexions, contrôle de santé ; prend en charge les trois stratégies Lazy/Eager/Pooled |
| ConfigRepository | `kernel/config/` | Chargement de la configuration des équipements YAML/JSON |
| GatewayEngine | `kernel/gateway/` | Moteur de règles de conversion inter-protocoles (Modbus→MQTT, etc.) |
| Bridge | `kernel/bridge/` | Pontage de processus externes (communication stdin/stdout) |
| Vendor | `kernel/vendor/` | Préréglages des paramètres fabricants (Siemens S7, Rockwell AB, etc.) |
| Event | `kernel/event/` | Bus d'événements basé sur les canaux |
| Metrics | `kernel/metrics/` | Interface de collecte de métriques (Prometheus + noop) |
| Security | `kernel/security/` | Enveloppe de sécurité TLS au niveau transport |

---

## Protocoles pris en charge

### Ethernet industriel (5/5 terminés)

| Protocole | Transport | Port | Fonctionnalités |
|------|------|------|------|
| **Modbus** | TCP, RTU | 502 | FC 01/02/03/04/05/06/16, CRC16, MBAP |
| **BACnet** | UDP | 47808 | Découverte d'équipements Who-Is/I-Am, ReadProperty |
| **EtherNet/IP** | TCP | 44818 | ENIP RegisterSession + CIP Read Tag |
| **OPC UA** | TCP | 4840 | Binary HEL + OpenSecureChannel |
| **PROFINET** | UDP | 34964 | DCP Identify/Set + Record Data Read/Write |

### Bus de terrain (11/11 — 4 logiciels purs + 7 SDK matériels)

| Protocole | Transport | Port | Fonctionnalités |
|------|------|------|------|
| **HART** | Serial FSK | — | Trames courtes/longues, Command 0/3, contrôle XOR |
| **CC-Link** | RS-485 | — | Interrogation maître/esclave (polling), CRC-16/XMODEM |
| **DNP3** | TCP, Serial | 20000 | Segmentation/réassemblage de la couche transport, interrogation Class 0, CRC-16/DNP |
| **IEC 61850** | TCP | 102 | MMS Initiate/Conclude, Read/Write, BER-TLV |
| **PROFIBUS** | Serial | — | Implémentation purement logicielle — nécessite le matériel CP 5611 |
| **CANopen** | CAN, Gateway | — | Lecture/écriture SDO, démarrage/arrêt/réinitialisation NMT, Heartbeat |
| **DeviceNet** | CAN, Gateway | — | Explicit Messaging, Poll, connexions E/S |
| **Foundation Fieldbus** | Serial | — | Implémentation purement logicielle — nécessite une carte d'interface FF H1 |
| **AS-Interface** | Serial | — | Implémentation purement logicielle — nécessite une passerelle ASi |
| **IO-Link** | Serial | — | Implémentation purement logicielle — nécessite un maître IO-Link |
| **CC-Link IE** | Ethernet | — | Implémentation purement logicielle — nécessite une passerelle CC-Link IE |

### IoT / Messagerie (2/2 terminés)

| Protocole | Transport | Port | Fonctionnalités |
|------|------|------|------|
| **MQTT** | TCP, WS | 1883 | 3.1.1 CONNECT/PUBLISH/SUBSCRIBE/PING |
| **HART-IP** | TCP, UDP | 5094 | HART over TCP, réutilise le codec HART |

### Bus automobiles (5/5 — 2 logiciels purs + 3 SDK matériels)

| Protocole | Transport | Débit en bauds | Fonctionnalités |
|------|------|--------|------|
| **LIN** | UART | — | Trames maître/esclave, contrôle PID, somme de contrôle Classic/Enhanced |
| **K-Line** | Serial | 10400 | ISO 9141/14230, 5-baud Fast Init, OBD-II SID 01/03/09 |
| **FlexRay** | CAN, Serial | — | Encodage/décodage Slot/Frame, CRC-16/XMODEM, SocketCAN + SerialBridge |
| **SAE J1850** | CAN | — | Encodage/décodage PWM/VPW, Mode 01/03/0A, CRC-8 |
| **MOST** | Serial | — | Encodage/décodage de trames hexadécimales, Read/Write/Status, SerialBridge |

### Bâtiment / Éclairage (2/2 — 1 logiciel pur + 1 SDK matériel)

| Protocole | Transport | Fonctionnalités |
|------|------|------|
| **DALI** | Serial | Trame aller 16 bits, trame retour 8 bits, commandes standard (Off/Max/Dim) |
| **LonWorks** | Serial | Codec purement logiciel — nécessite une puce Neuron ou une passerelle |

### Pontage matériel (12/12 — tous avec pilotes CmdBridge/SerialBridge)

| Protocole | Type de SDK | Fonctionnalités |
|------|----------|------|
| **EtherCAT** | CmdBridge + hex-codec | Sous-commandes Upload/Download/Slaves, encodage/décodage hexadécimal CoE |
| **POWERLINK** | CmdBridge + hex-codec | Sous-commandes Read/Write/Status, trames SoC/Preq/Pres |
| **SERCOS III** | CmdBridge + hex-codec | Sous-commandes Read/Write/Phase, transitions de phase (NRT→CP4) |
| **SERCOS I/II** | CmdBridge + hex-codec | Read/Write/Status, topologie en anneau à fibre optique |
| **ControlNet** | CmdBridge + hex-codec | Read/Write/Status, ordonnancement CTDMA, producteur/consommateur |
| **Interbus** | CmdBridge | Read/Decode, topologie en anneau, trames IBS CMD |
| **WorldFIP** | CmdBridge | Read/Write/Decode, producteur/consommateur, arbitre de bus |
| **Lightbus** | CmdBridge | Read/Decode, interconnexion par fibre optique, 32 nœuds |
| **Modbus Plus** | CmdBridge + hex-codec | Read/Write/Decode, passage de jeton (token passing), peer-to-peer |
| **ISA100.11a** | CmdBridge + hex-codec | Read/Write/List, sans fil 6LoWPAN, maillé (mesh) |
| **WirelessHART** | CmdBridge + hex-codec | Read/Write/Scan, MAC TSMP, auto-organisé |
| **SAE J1850** | CmdBridge | Diagnostic automobile converti vers CmdBridge, nécessite une interface J1850 |

### Bus système (3/3 — tous avec pilotes sysfs/procfs)

| Protocole | Type de SDK | Fonctionnalités |
|------|----------|------|
| **PCI/PCIe** | sysfs — `/sys/bus/pci/devices/<BDF>/config` | Lecture/écriture de l'espace de configuration, PipeTransport, nécessite CAP_SYS_ADMIN |
| **VME/VPX** | procfs — `/proc/vme/<slot>` | Adressage A16/A24/A32, PipeTransport, nécessite le module vme_tsi148 |
| **CompactPCI** | sysfs — `/sys/bus/pci/devices/<BDF>/config` | PICMG 2.0, insertion à chaud (hot-swap), 3U/6U, nécessite cpci_hotplug |

---

## Guide d'utilisation

La référence API et les exemples d'utilisation (installation, lecture/écriture Modbus, pool de connexions, MQTT, préréglages Vendor, conversion par passerelle, protocoles personnalisés) ont été déplacés vers un document dédié :

> **[API.md](API.md) — référence API complète et guide d'utilisation**
>
> Version anglaise : [API.en.md](../../API.en.md)

---

## Structure du projet

```
industrial-protocols-go/
├── go.work                       # l'espace de travail agrège tous les modules
├── kernel/                       # modules du noyau
│   ├── connection/               # ConnectionManager + pool de connexions
│   ├── config/                   # ConfigRepository (YAML)
│   ├── gateway/                  # GatewayEngine conversion de protocoles
│   ├── bridge/                   # Bridge pontage de processus externes
│   ├── vendor/                   # Vendor préréglages fabricants
│   ├── event/                    # bus d'événements
│   ├── metrics/                  # interface de métriques
│   ├── security/                 # sécurité TLS
│   ├── pipeline/                 # chaîne de middleware
│   └── transport/                # transports TCP/UDP/Pipe
│
├── protocols/                    # 40 modules de protocoles
│   ├── ethernet/                 # 5 Ethernet industriel (tous terminés)
│   ├── fieldbus/                 # 11 bus de terrain (4 logiciels purs + 7 SDK matériels)
│   ├── iot/                      # 2 IoT/messagerie (tous terminés)
│   ├── automotive/               # 5 bus automobiles (2 logiciels purs + 3 SDK matériels)
│   ├── building/                 # 2 bâtiment/éclairage (1 logiciel pur + 1 SDK matériel)
│   ├── bridge/                   # 12 pontages matériels (CmdBridge/SerialBridge)
│   └── system/                   # 3 bus système (pilotes sysfs/procfs)
│
├── examples/modbus_basic/        # exemple Modbus TCP
├── _tools/                       # Makefile + scripts auxiliaires
└── docs/superpowers/             # documents de conception
```

## Tests

```bash
make test         # tests complets
make test-unit    # tests unitaires uniquement
make vet && make fmt
```

---

## Nous soutenir

Si ce projet vous a été utile, n'hésitez pas à faire un don pour nous aider à poursuivre sa maintenance.

### Alipay / WeChat

| Alipay | WeChat |
|--------|------|
| ![Alipay](../../alipay.png) | ![WeChat](../../weixinpay.png) |

### Virement international (transfert bancaire international)

Transfert bancaire pour les utilisateurs situés hors de la Chine continentale :

**Informations du bénéficiaire :**

| Élément | Valeur |
|------|------|
| Nom du bénéficiaire | WANG KEXUN |
| Numéro de compte du bénéficiaire | 881015918251 |

**Banque du bénéficiaire :**

| Élément | Valeur |
|------|------|
| Nom de la banque | ZA Bank Limited |
| SWIFT Code | AABLHKHHXXX |
| Numéro de banque | 387 |
| Adresse de la banque | Core F, Cyberport 3, 100 Cyberport Road, Hong Kong |

**Banque correspondante pour les virements internationaux (banque intermédiaire, si nécessaire) :**

> Veuillez noter qu'il s'agit des informations de la banque correspondante (banque intermédiaire) pour les virements internationaux, et non de celles de la banque du bénéficiaire. Renseignez-vous auprès de votre banque pour savoir si les informations de la banque correspondante sont requises.

- **Pour les virements en dollars de Hong Kong (HKD), renminbi (RMB) et dollars américains (USD)** — la banque correspondante est Citibank :
  - Nom de la banque : Citibank N.A. Hong Kong
  - SWIFT Code : CITIHKHXXXX
  - Numéro de banque : 006
  - Nom de la succursale : Hong Kong Branch
  - Numéro de succursale : 391
  - Adresse de la banque : Citibank Tower, Citibank Plaza, 3 Garden Road, Central, Hong Kong
- **Pour les virements dans d'autres devises** — la banque correspondante est BNY Mellon :
  - Nom de la banque : THE BANK OF NEW YORK MELLON
  - SWIFT Code : IRVTUS3NXXX
  - Adresse de la banque : THE BANK OF NEW YORK MELLON, 240 GREENWICH STREET, NEW YORK, United States

---

## License

MIT — Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz
