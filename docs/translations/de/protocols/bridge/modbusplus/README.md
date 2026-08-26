# Modbus-Plus-Protokoll-SDK

Modbus Plus (MB+) ist ein schnelles industrielles Netzwerk mit Token-Passing, entwickelt von Modicon (Schneider Electric).

## Varianten

| Variante | Brückentyp | Beschreibung |
|---------|-------------|-------------|
| `gateway` | GatewayBridge | TCP-Verbindung zum SA85/BM85-Adapter |
| `cmd` | CmdBridge | CLI-Wrapper für das Programm `sa85_cli` |

## Gateway-Treiber

```go
driver, codec, err := modbusplus.NewGatewayDriver("192.168.0.160:2010")
handler, err := modbusplus.ReadyHandler("192.168.0.160:2010")
```

## Cmd-Treiber

```go
driver, codec, err := modbusplus.NewCmdDriver("/usr/bin/sa85_cli")
```

### Installation des CLI-Tools

```bash
# Install SA85 Modbus Plus driver and tools
# Refer to Schneider Electric documentation for SA85 adapter
```

Verifizieren: `sa85_cli --help`

## Hardware

| Gateway | Schnittstelle | Standard-IP |
|---------|-----------|-------------|
| Schneider SA85 | ISA-Modbus-Plus-Adapter | vom Host vergeben |
| Schneider BM85 | Modbus-Plus-Brücke/Multiplexer | 192.168.0.A0 |
| ProSoft MVI56-MBP | ControlLogix-MB+-Modul | vom Host vergeben |

## Verdrahtung

- Twinax-Kabel (RG-62) mit BNC-Steckern für den MB+-Stammbus
- Abschlusswiderstand (78 Ohm an jedem Ende)
- Die BM85-Brücke verbindet MB+ über Ethernet mit TCP

## Frame-Format des Protokolls

- Magic (2 Bytes): `0x4D42`
- Ziel (1 Byte): Knotenadresse
- Befehl (1 Byte): 0x01=Lesen, 0x02=Schreiben
- Länge (2 Bytes): Payload-Größe (Big-Endian)
- Payload (N Bytes): Daten
- CRC (2 Bytes): Modbus-CRC-16 (Little-Endian)

## Testen

```bash
go test ./... -v
```
