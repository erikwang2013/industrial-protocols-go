# SDK Protokol Modbus Plus

Modbus Plus (MB+) adalah jaringan industri token-passing berkecepatan tinggi yang dikembangkan oleh Modicon (Schneider Electric).

## Varian

| Varian | Tipe Bridge | Deskripsi |
|---------|-------------|-------------|
| `gateway` | GatewayBridge | Koneksi TCP ke adaptor SA85/BM85 |
| `cmd` | CmdBridge | Pembungkus CLI untuk utilitas `sa85_cli` |

## Driver Gateway

```go
driver, codec, err := modbusplus.NewGatewayDriver("192.168.0.160:2010")
handler, err := modbusplus.ReadyHandler("192.168.0.160:2010")
```

## Driver Cmd

```go
driver, codec, err := modbusplus.NewCmdDriver("/usr/bin/sa85_cli")
```

### Instalasi Alat CLI

```bash
# Install driver dan alat SA85 Modbus Plus
# Lihat dokumentasi Schneider Electric untuk adaptor SA85
```

Verifikasi: `sa85_cli --help`

## Perangkat Keras

| Gateway | Antarmuka | IP Default |
|---------|-----------|-------------|
| Schneider SA85 | Adaptor ISA Modbus Plus | ditetapkan host |
| Schneider BM85 | Bridge/multiplekser Modbus Plus | 192.168.0.A0 |
| ProSoft MVI56-MBP | Modul ControlLogix MB+ | ditetapkan host |

## Kabel

- Kabel twinaxial (RG-62) dengan konektor BNC untuk trunk MB+
- Resistor terminasi (78 ohm di setiap ujung)
- Bridge BM85 menghubungkan MB+ ke TCP via Ethernet

## Format Frame Protokol

- Magic (2 byte): `0x4D42`
- Tujuan (1 byte): Alamat node
- Perintah (1 byte): 0x01=baca, 0x02=tulis
- Panjang (2 byte): Ukuran payload (big-endian)
- Payload (N byte): Data
- CRC (2 byte): Modbus CRC-16 (little-endian)

## Pengujian

```bash
go test ./... -v
```
