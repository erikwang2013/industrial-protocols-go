# PCI-/PCIe-Protokoll-SDK

Implementiert `kernel.Protocol` für den PCI- und PCI-Express-Buszugriff über Linux-sysfs.

## Protokoll

- **Name**: `pci`
- **Varianten**: `pci`
- **Standard-Port**: 0 (speichergemappter Bus)
- **Transport**: sysfs `/sys/bus/pci/devices/<BDF>/config`

Der Codec leitet Rohbytes zwischen der Anwendung und dem PCI-
Konfigurationsraum durch. Lese- und Schreibzugriffe gehen direkt über die
sysfs-Config-Datei zu den Konfigurationsregistern des Geräts.

## Kernel-Anforderungen

Die folgenden Kernelmodule müssen geladen sein:

- `pcieport` — PCI-Express-Porttreiber
- `pci_sysfs` — sysfs-PCI-Schnittstelle (in den meisten Kerneln enthalten)

Das sysfs-Dateisystem muss unter `/sys` eingehängt sein. Dies ist auf
allen modernen Linux-Distributionen der Standard.

Erforderliche Berechtigungen:
- Root-Zugriff oder `CAP_SYS_ADMIN` für den Zugriff auf den Konfigurationsraum
- Die Config-Datei gehört root:root mit Modus 0600

## Treiber

```go
d, err := pci.NewPCIDriver("0000:00:1f.3")
if err != nil {
    log.Fatal(err)
}
defer d.Close()

tr := d.Transport()
// Use tr.Read/tr.Write for raw PCI config space access
```

## BDF-Format

Die Busadresse ist eine BDF-Zeichenkette (Bus:Device.Function):
- `0000:00:1f.3` — Domain 0000, Bus 00, Gerät 1f, Funktion 3
- `0000:01:00.0` — Domain 0000, Bus 01, Gerät 00, Funktion 0
