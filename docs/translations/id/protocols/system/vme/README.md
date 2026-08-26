# SDK Protokol VME / VPX

Mengimplementasikan `kernel.Protocol` untuk akses bus VMEbus dan VPX melalui procfs Linux.

## Protokol

- **Nama**: `vme`
- **Varian**: `vme`
- **Port Default**: 0 (bus memory-mapped)
- **Transport**: procfs `/proc/vme/<slot>`

Codec meneruskan byte mentah antara aplikasi dan ruang alamat VME.
Pembacaan dan penulisan langsung ke bus melalui antarmuka procfs
yang disediakan oleh driver kernel VME.

## Persyaratan Kernel

Modul kernel berikut harus dimuat:

- `vme_tsi148` -- driver bridge VME Tundra TSI148 (paling umum)
  - Juga didukung: `vme_ca91cx42` (Universe II), `vme_user`

Filesystem procfs harus dipasang di `/proc`.

Izin yang diperlukan:
- Akses root diperlukan untuk membuka `/proc/vme/<slot>` untuk baca-tulis
- Node perangkat dimiliki root:root

## Driver

```go
d, err := vme.NewVMEDriver(0) // slot 0
if err != nil {
    log.Fatal(err)
}
defer d.Close()

tr := d.Transport()
// Gunakan tr.Read/tr.Write untuk akses mentah bus VME
```

## Mode Pengalamatan VME

Codec meneruskan semua byte address modifier dan address secara
transparan. Aplikasi harus menambahkan informasi pengalamatan
di awal payload data:

- **A16**: Ruang alamat I/O pendek 16-bit
- **A24**: Ruang alamat standar 24-bit
- **A32**: Ruang alamat extended 32-bit

## Kompatibilitas VPX

Sistem VPX (VITA 46) yang mengekspos antarmuka procfs kompatibel VME
dapat menggunakan driver ini. Penomoran slot sama.
