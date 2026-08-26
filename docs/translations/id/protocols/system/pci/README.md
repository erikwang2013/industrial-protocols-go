# SDK Protokol PCI / PCIe

Mengimplementasikan `kernel.Protocol` untuk akses bus PCI dan PCI Express melalui sysfs Linux.

## Protokol

- **Nama**: `pci`
- **Varian**: `pci`
- **Port Default**: 0 (bus memory-mapped)
- **Transport**: sysfs `/sys/bus/pci/devices/<BDF>/config`

Codec meneruskan byte mentah antara aplikasi dan ruang konfigurasi PCI.
Pembacaan dan penulisan langsung ke register konfigurasi perangkat
melalui file config sysfs.

## Persyaratan Kernel

Modul kernel berikut harus dimuat:

- `pcieport` -- driver port PCI Express
- `pci_sysfs` -- antarmuka sysfs PCI (bawaan sebagian besar kernel)

Filesystem sysfs harus dipasang di `/sys`. Ini adalah default pada
semua distro Linux modern.

Izin yang diperlukan:
- Akses root atau `CAP_SYS_ADMIN` untuk akses ruang konfigurasi
- File konfigurasi dimiliki root:root dengan mode 0600

## Driver

```go
d, err := pci.NewPCIDriver("0000:00:1f.3")
if err != nil {
    log.Fatal(err)
}
defer d.Close()

tr := d.Transport()
// Gunakan tr.Read/tr.Write untuk akses mentah ruang konfigurasi PCI
```

## Format BDF

Alamat bus adalah string BDF (Bus:Device.Function):
- `0000:00:1f.3` -- domain 0000, bus 00, device 1f, fungsi 3
- `0000:01:00.0` -- domain 0000, bus 01, device 00, fungsi 0
