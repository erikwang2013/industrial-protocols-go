# SDK Protokol CompactPCI

Mengimplementasikan `kernel.Protocol` untuk akses bus CompactPCI melalui sysfs Linux.

## Protokol

- **Nama**: `cpci`
- **Varian**: `cpci`
- **Port Default**: 0 (bus memory-mapped)
- **Transport**: sysfs `/sys/bus/pci/devices/<BDF>/config`

CompactPCI (PICMG 2.0) menggunakan antarmuka listrik dan perangkat lunak yang
sama dengan PCI konvensional. Codec meneruskan byte mentah antara aplikasi
dan ruang konfigurasi PCI melalui sysfs.

## Persyaratan Kernel

Modul kernel berikut harus dimuat:

- `pcieport` -- driver port PCI Express (untuk sistem CPCIe hibrida)
- `pci_sysfs` -- antarmuka sysfs PCI (bawaan sebagian besar kernel)
- `cpci_hotplug` -- pengontrol hotplug CompactPCI (opsional, untuk hot-swap)

Filesystem sysfs harus dipasang di `/sys`. Ini adalah default pada
semua distro Linux modern.

Izin yang diperlukan:
- Akses root atau `CAP_SYS_ADMIN` untuk akses ruang konfigurasi
- File konfigurasi dimiliki root:root dengan mode 0600

## Driver

```go
d, err := cpci.NewCPCIDriver("0000:02:00.0")
if err != nil {
    log.Fatal(err)
}
defer d.Close()

tr := d.Transport()
// Gunakan tr.Read/tr.Write untuk akses mentah ruang konfigurasi CPCI
```

## CompactPCI vs PCI

CompactPCI menggunakan enumerasi bus dan ruang konfigurasi PCI standar.
Perbedaan utama dari PCI desktop:

- **Faktor bentuk 3U/6U**: Mekanika Eurocard dengan konektor pin-and-socket
- **Penomoran bus**: Setiap segmen sasis CPCI mendapatkan nomor bus PCI sendiri
- **Hot swap**: Hot swap PICMG 2.1 menggunakan model hotplug PCI standar
- **Slot sistem**: Bus 0, device 0 adalah pengontrol slot sistem

## Format BDF

Alamat bus adalah string BDF (Bus:Device.Function):
- `0000:02:00.0` -- domain 0000, bus 02, device 00, fungsi 0
- `0000:02:08.0` -- domain 0000, bus 02, device 08, fungsi 0
