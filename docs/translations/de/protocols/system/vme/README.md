# VME-/VPX-Protokoll-SDK

Implementiert `kernel.Protocol` für den VMEbus- und VPX-Buszugriff über Linux-procfs.

## Protokoll

- **Name**: `vme`
- **Varianten**: `vme`
- **Standard-Port**: 0 (speichergemappter Bus)
- **Transport**: procfs `/proc/vme/<slot>`

Der Codec leitet Rohbytes zwischen der Anwendung und dem VME-
Adressraum durch. Lese- und Schreibzugriffe gehen direkt über die procfs-
Schnittstelle des VME-Kerneltreibers auf den Bus.

## Kernel-Anforderungen

Die folgenden Kernelmodule müssen geladen sein:

- `vme_tsi148` — Tundra-TSI148-VME-Bridge-Treiber (am häufigsten)
  - Ebenfalls unterstützt: `vme_ca91cx42` (Universe II), `vme_user`

Das procfs-Dateisystem muss unter `/proc` eingehängt sein.

Erforderliche Berechtigungen:
- Root-Zugriff ist erforderlich, um `/proc/vme/<slot>` zum Lesen/Schreiben zu öffnen
- Die Geräteknoten gehören root:root

## Treiber

```go
d, err := vme.NewVMEDriver(0) // slot 0
if err != nil {
    log.Fatal(err)
}
defer d.Close()

tr := d.Transport()
// Use tr.Read/tr.Write for raw VME bus access
```

## VME-Adressierungsmodi

Der Codec leitet alle Address-Modifier- und Adressbytes
transparent durch. Anwendungen sollten die Adressinformationen
dem Daten-Payload voranstellen:

- **A16**: 16-Bit-Short-I/O-Adressraum
- **A24**: 24-Bit-Standard-Adressraum
- **A32**: 32-Bit-erweiterter Adressraum

## VPX-Kompatibilität

VPX-Systeme (VITA 46), die eine VME-kompatible procfs-Schnittstelle
bereitstellen, können diesen Treiber verwenden. Die Slot-Nummerierung
ist identisch.
