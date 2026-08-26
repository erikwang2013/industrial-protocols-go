# Laporan Pengujian — industrial-protocols-go

Tanggal: 2026-08-27
Ruang lingkup: Workspace multi-modul Go (go.work, 42 modul: kernel / protocols / examples)
Perintah: `go test ./...` dan `go test -cover ./...` (karena tidak ada go.mod di akar workspace, dijalankan setara per modul, hasil identik)

## 1. Kesimpulan Umum

- **Semua hijau**: 42 modul, 52 paket berisi pengujian semuanya lolos (terverifikasi dengan beberapa kali pengulangan penuh, termasuk selama periode penulisan pengujian paralel).
- File pengujian 106, fungsi pengujian 659.
- `go vet ./...` di seluruh modul tanpa peringatan.
- Rata-rata cakupan pernyataan **83,3%**; cakupan paket inti kernel 95%~100%.

## 2. Statistik dan Cakupan Pengujian per Modul

### kernel (11 paket, cakupan tertinggi selain examples)

| Paket | Cakupan |
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

### Modul penting protocols

| Modul | Cakupan | Modul | Cakupan |
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

> examples/modbus_basic tidak memiliki file pengujian (program contoh, sesuai harapan).

## 3. Daftar Perbaikan

### 3.1 Perbaikan bug sumber (ditemukan dan diverifikasi oleh tester-kernel / tester-protocols)

| File:Baris | Masalah | Perbaikan |
|---|---|---|
| kernel/gateway/engine.go:27 | `Transform` mengalami nil dereference panic di `rule.Src.Name()` / `rule.Dst.Name()` / `rule.Map(...)` untuk aturan dengan `Src`/`Dst`/`Map` bernilai nil | Lewati aturan yang belum dikonfigurasi lengkap di awal loop (konsisten dengan semantik "lewati aturan buruk" yang sudah ada) |
| protocols/ethernet/modbus/modbus.go:215 | `decodePDU` membaca `pdu[1]` di luar batas untuk PDU pengecualian 1 byte (mis. `0x81`) sehingga panic | Tambahkan pemeriksaan `len(pdu) < 2`, kembalikan error parsing |
| protocols/ethernet/modbus/modbus.go:230 | `decodePDU` slice out-of-bounds panic ketika penghitung byte (`pdu[1]`) melebihi data aktual (frame jahat/rusak) | Tambahkan pemeriksaan batas `2+n > len(pdu)`, kembalikan error parsing |
| protocols/ethernet/opcua/opcua.go:63,74 | `encodeHello` menulis ulang panjang pesan ke offset 8 (slot versi protokol), sedangkan spesifikasi mensyaratkan offset 4 | Catat `sizePos` (offset 4) lebih awal, tulis kembali ke posisi yang benar |
| protocols/ethernet/opcua/opcua.go:82,89 | Field panjang `encodeOpenSecureChannel` tidak pernah ditulis ulang, selalu 0 | Tulis ulang di akhir berdasarkan panjang frame aktual |
| protocols/ethernet/ethernetip/ethernetip.go:111,117 | `encodeReadTag` mengalokasikan frame 16+len(cipReq) tetapi permintaan CIP ditulis di [28:], kelebihan 12 byte menyebabkan `copy` diam-diam membuang seluruh permintaan CIP; field panjang juga kurang 12 | Alokasi diubah menjadi 28+len(cipReq), field length diubah menjadi 4+len(cipReq) |
| protocols/ethernet/profinet/profinet.go:72 | `Decode` slice out-of-bounds panic untuk `blockLen` di luar batas | Tambahkan pemeriksaan `len(data) < 10+blockLen` |
| protocols/fieldbus/canopen/canopen.go:102 | Pemeriksaan panjang minimum `Decode` adalah 4, tetapi `unmarshalCAN` membutuhkan `data[4:8]`, frame 4~7 byte menyebabkan panic | Panjang minimum diubah menjadi 8 (format wire tetap 8 byte) |
| protocols/fieldbus/dnp3/dnp3.go:121 | `decodeTPDU` slice out-of-bounds panic untuk `length<5` (`apduEnd < apduStart`) atau `length` terlalu panjang (`apduEnd > len(data)`) | Validasi seragam `length < 5 || apduEnd > len(data)` lalu kembalikan error |
| protocols/fieldbus/dnp3/dnp3.go:157 | Variabel kerja `crc16DNP` tidak dimask `& 0xFF` sesuai spesifikasi, vektor yang diketahui `crc16DNP("123456789")` menghasilkan 0x69FF (seharusnya 0xEA82) | `temp := (crc ^ b) & 0xFF` |
| protocols/automotive/kline/kline.go:74 | Pemeriksaan `len<5` pada `Decode` mendahului penentuan frame sinkronisasi fast-init, frame sinkronisasi 1 byte (0x55) ditolak secara keliru | Periksa frame kosong dan frame sinkronisasi terlebih dahulu, baru kemudian panjang |

### 3.2 Perbaikan pengujian (memperbaiki error kompilasi / error asersi / nama duplikat / asersi usang)

| File:Baris | Masalah | Perbaikan |
|---|---|---|
| protocols/automotive/kline/kline_extra_test.go:43 | Konstanta `0x68+0x10+0xF1+0x01`(362) diberikan ke byte gagal kompilasi | Harapan diubah menjadi 8 bit rendah checksum `0x6A` |
| protocols/automotive/saej1850/saej1850_extra_test.go | Memanggil metode/field tidak diekspor pada antarmuka `kernel.Codec`, gagal kompilasi | Tambahkan helper `newCodec` yang mengasertasi sebagai `*j1850Codec`; asersi offset `TestEncodeMode01PID` dikoreksi menjadi raw[4..6] (ID 4 byte + 4 byte awal data) |
| protocols/bridge/sercos1/sercos1_extra_test.go:45 | `TestEncodeReadDefaults` duplikat nama dengan pengujian yang sudah ada, gagal kompilasi | Ganti nama menjadi `TestEncodeReadDefaultsFiber` (mempertahankan cakupan varian fiber) |
| protocols/building/dali/dali_extra_test.go:68 | `TestEncodeDirectArc` duplikat nama dengan pengujian yang sudah ada, gagal kompilasi | Ganti nama menjadi `TestEncodeDirectArcCmd` |
| protocols/bridge/powerlink/powerlink_extra_test.go:51 | Harapan `write 0x1000 A`, aktual `0A` (`%X` memformat dua digit per byte, konsisten dengan konvensi pengujian `DEADBEEF` yang sudah ada) | Harapan diubah menjadi `0A` |
| protocols/bridge/sercos/sercos_extra_test.go:62 | Sama seperti di atas, harapan `1` harusnya `01` | Harapan diubah menjadi `01` |
| protocols/ethernet/opcua/opcua_extra_test.go:35 | Mengasertasi perilaku bug lama (panjang di offset 8), tidak berlaku lagi setelah perbaikan sumber | Diubah untuk mengasertasi field versi protokol 0 (panjang dicakup oleh TestExtraHELMessageSizeAtOffset4) |
| protocols/fieldbus/canopen/canopen_extra_test.go:135 | Pengujian "prove-it" asersi terbalik (gagal ketika `err != nil`, padahal error parsing justru hasil yang diharapkan) | Diubah menjadi gagal saat `err == nil` |
| protocols/ethernet/profinet/profinet_extra_test.go:23 | Sama seperti di atas, asersi terbalik | Diubah menjadi gagal saat `err == nil` |
| protocols/fieldbus/dnp3/dnp3_extra_test.go:35,53 | Sama seperti di atas ×2 | Diubah menjadi gagal saat `err == nil` |

> Catatan: masalah implementasi antarmuka `stubTransport` di `kernel/security/tls_test.go` dan nilai harapan CRC `cclink/cclinkie` diperbaiki sendiri oleh insinyur penguji dalam pekerjaan paralel, tidak diubah pada kesempatan ini.

## 4. Risiko Tersisa

1. **Modul cakupan rendah**: profibus (35.7%), lonworks / asinterface / foundationfieldbus / iolink (37.9%) — pengujian hanya mencakup sedikit jalur, disarankan menambah cabang Decode/Encode dan jalur pengecualian di kemudian hari.
2. **CRC modbus RTU tidak diverifikasi**: `decodeRTU` tidak memverifikasi CRC (pengujian `TestExtraRTUDecodesCorruptCRC` mendokumentasikan GAP ini secara eksplisit dan lolos sesuai kondisi saat ini). Disarankan melengkapinya saat interoperasi dengan perangkat nyata.
3. **Keterbatasan format wire canopen**: `marshalCAN` hanya membawa 4 byte data pertama (`TestExtraSDOWritePayloadLostOnWire` mendokumentasikannya), penulisan SDO multi-byte akan kehilangan payload.
4. **Lewati karena ketergantungan perangkat keras**: cpci / pci / vme melakukan `t.Skip` di lingkungan tanpa perangkat keras (wajar).
5. **Sisa file sementara di repositori**: `main.go` tak terlacak di akar (program percobaan, mereferensikan modul `probe/` yang tidak ada), `scripts/`, `docs/*.png`, bukan bagian dari pengiriman ini, disarankan dibersihkan.
6. **Pengujian masih ditambahkan secara paralel**: laporan ini berdasarkan snapshot pengujian hijau penuh terakhir; jika insinyur penguji terus mengirimkan `*_extra_test.go` berikutnya, perlu menjalankan ulang regresi penuh.
