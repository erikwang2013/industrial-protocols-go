# CompactPCI-Protokoll-SDK

Implementiert `kernel.Protocol` für den CompactPCI-Buszugriff über Linux-sysfs.

## Protokoll

- **Name**: `cpci`
- **Varianten**: `cpci`
- **Standard-Port**: 0 (speichergemappter Bus)
- **Transport**: sysfs `/sys/bus/pci/devices/<BDF>/config`

CompactPCI (PICMG 2.0) verwendet dieselbe elektrische und softwaretechnische Schnittstelle
wie herkömmliches PCI. Der Codec leitet Rohbytes zwischen der
Anwendung und dem PCI-Konfigurationsraum über sysfs durch.

## Kernel-Anforderungen

Die folgenden Kernelmodule müssen geladen sein:

- `pcieport` — PCI-Express-Porttreiber (für Hybrid-CPCIe-Systeme)
- `pci_sysfs` — sysfs-PCI-Schnittstelle (in den meisten Kerneln enthalten)
- `cpci_hotplug` — CompactPCI-Hotplug-Controller (optional, für Hot-Swap)

Das sysfs-Dateisystem muss unter `/sys` eingehängt sein. Dies ist auf
allen modernen Linux-Distributionen der Standard.

Erforderliche Berechtigungen:
- Root-Zugriff oder `CAP_SYS_ADMIN` für den Zugriff auf den Konfigurationsraum
- Die Config-Datei gehört root:root mit Modus 0600

## Treiber

```go
d, err := cpci.NewCPCIDriver("0000:02:00.0")
if err != nil {
    log.Fatal(err)
}
defer d.Close()

tr := d.Transport()
// Use tr.Read/tr.Write for raw CPCI config space access
```

## CompactPCI vs. PCI

CompactPCI verwendet die standardmäßige PCI-Bus-Enumeration und den standardmäßigen Konfigurationsraum.
Hauptunterschiede zum Desktop-PCI:

- **3U/6U-Formfaktor**: Eurocard-Mechanik mit Stift-/Buchsensteckern
- **Busnummerierung**: Jedes CPCI-Chassis-Segment erhält eine eigene PCI-Busnummer
- **Hot-Swap**: PICMG-2.1-Hot-Swap nutzt das standardmäßige PCI-Hotplug-Modell
- **System-Slot**: Bus 0, Gerät 0 ist der System-Slot-Controller

## BDF-Format

Die Busadresse ist eine BDF-Zeichenkette (Bus:Device.Function):
- `0000:02:00.0` — Domain 0000, Bus 02, Gerät 00, Funktion 0
- `0000:02:08.0` — Domain 0000, Bus 02, Gerät 08, Funktion 0
