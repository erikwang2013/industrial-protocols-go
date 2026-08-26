# Testbericht — industrial-protocols-go

Datum: 2026-08-27
Umfang: Go-Multi-Modul-Workspace (go.work, 42 Module: kernel / protocols / examples)
Befehl: `go test ./...` und `go test -cover ./...` (da im Workspace-Root keine go.mod existiert, äquivalent pro Modul ausgeführt; Ergebnisse identisch)

## 1. Gesamtergebnis

- **Komplett grün**: Alle 42 Module und 52 Pakete mit Tests bestanden (mehrfach durch vollständige Neu-Läufe verifiziert, auch während parallel laufender Testschreibvorgänge).
- 106 Testdateien, 659 Testfunktionen.
- `go vet ./...` ohne Warnungen in allen Modulen.
- Durchschnittliche Statement-Abdeckung **83,3 %**; Kernel-Kernpakete 95 %–100 %.

## 2. Teststatistik und Abdeckung nach Modul

### kernel (11 Pakete, höchste Abdeckung außer examples)

| Paket | Abdeckung |
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

### protocols — Schlüsselmodule

| Modul | Abdeckung | Modul | Abdeckung |
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

> examples/modbus_basic hat keine Testdateien (Beispielprogramm, wie erwartet).

## 3. Liste der Fixes

### 3.1 Fehlerbehebungen im Quellcode (von tester-kernel / tester-protocols gefunden und verifiziert)

| Datei:Zeile | Problem | Fix |
|---|---|---|
| kernel/gateway/engine.go:27 | `Transform` löst bei Regeln mit nil-`Src`/`Dst`/`Map` einen Nil-Dereferenz-Panic an `rule.Src.Name()` / `rule.Dst.Name()` / `rule.Map(...)` aus | Am Schleifenanfang werden unvollständig konfigurierte Regeln übersprungen (konsistent mit der bestehenden „schlechte Regeln überspringen“-Semantik) |
| protocols/ethernet/modbus/modbus.go:215 | `decodePDU` liest bei 1-Byte-Exception-PDUs (z. B. `0x81`) außerhalb der Grenzen (`pdu[1]`) → Panic | Prüfung `len(pdu) < 2` ergänzt, gibt Parsing-Fehler zurück |
| protocols/ethernet/modbus/modbus.go:230 | `decodePDU`: Der Byte-Zähler (`pdu[1]`) übersteigt die tatsächlichen Daten → Slicing außerhalb der Grenzen (bösartige/beschädigte Frames) | Grenzprüfung `2+n > len(pdu)` ergänzt, gibt Parsing-Fehler zurück |
| protocols/ethernet/opcua/opcua.go:63,74 | `encodeHello` schreibt die Nachrichtenlänge an Offset 8 zurück (Slot der Protokollversion); die Spezifikation verlangt Offset 4 | `sizePos` (Offset 4) wird vorab notiert und an die richtige Stelle zurückgeschrieben |
| protocols/ethernet/opcua/opcua.go:82,89 | Das Längenfeld von `encodeOpenSecureChannel` wurde nie zurückgeschrieben und blieb stets 0 | Am Ende gemäß der tatsächlichen Frame-Länge zurückschreiben |
| protocols/ethernet/ethernetip/ethernetip.go:111,117 | `encodeReadTag` allokiert 16+len(cipReq), schreibt den CIP-Request aber nach [28:]; die überzähligen 12 Bytes führen dazu, dass `copy` den gesamten CIP-Request stillschweigend verwirft; das Längenfeld ist ebenfalls um 12 zu klein | Allokation auf 28+len(cipReq) geändert, Längenfeld auf 4+len(cipReq) |
| protocols/ethernet/profinet/profinet.go:72 | `Decode` löst bei übermäßigem `blockLen` einen Slicing-Panic aus | Prüfung `len(data) < 10+blockLen` ergänzt |
| protocols/fieldbus/canopen/canopen.go:102 | Die minimale Längenprüfung in `Decode` ist 4, aber `unmarshalCAN` benötigt `data[4:8]`; Frames mit 4–7 Bytes → Panic | Minimale Länge auf 8 geändert (Wire-Format fix 8 Bytes) |
| protocols/fieldbus/dnp3/dnp3.go:121 | `decodeTPDU` löst bei `length<5` (`apduEnd < apduStart`) oder übermäßigem `length` (`apduEnd > len(data)`) einen Slicing-Panic aus | Einheitliche Prüfung `length < 5 || apduEnd > len(data)` und Fehlerrückgabe |
| protocols/fieldbus/dnp3/dnp3.go:157 | Die Arbeitsvariable von `crc16DNP` wurde nicht wie spezifiziert mit `& 0xFF` maskiert; bekannter Vektor `crc16DNP("123456789")` ergab 0x69FF (sollte 0xEA82 sein) | `temp := (crc ^ b) & 0xFF` |
| protocols/automotive/kline/kline.go:74 | Die `len<5`-Prüfung in `Decode` läuft vor der Erkennung des Fast-Init-Sync-Frames; 1-Byte-Sync-Frames (0x55) wurden fälschlich abgelehnt | Erst leere Frames und Sync-Frames prüfen, dann die Länge |

### 3.2 Testkorrekturen (behobene Kompilierfehler / falsche Assertions / Namenskollisionen / veraltete Assertions)

| Datei:Zeile | Problem | Fix |
|---|---|---|
| protocols/automotive/kline/kline_extra_test.go:43 | Konstante `0x68+0x10+0xF1+0x01` (362) lässt sich nicht einer byte-Variable zuweisen → Kompilierfehler | Erwartung auf die niedrigen 8 Bit der Prüfsumme `0x6A` geändert |
| protocols/automotive/saej1850/saej1850_extra_test.go | Aufruf unexportierter Methoden/Felder über die `kernel.Codec`-Schnittstelle → Kompilierfehler | Neuer `newCodec`-Helper, der auf `*j1850Codec` prüft; die Offsets der Assertions in `TestEncodeMode01PID` auf raw[4..6] korrigiert (ID 4 Bytes + erste 4 Datenbytes) |
| protocols/bridge/sercos1/sercos1_extra_test.go:45 | `TestEncodeReadDefaults` kollidiert mit einem bestehenden Test → Kompilierfehler | Umbenannt in `TestEncodeReadDefaultsFiber` (Fiber-Variante bleibt abgedeckt) |
| protocols/building/dali/dali_extra_test.go:68 | `TestEncodeDirectArc` kollidiert mit einem bestehenden Test → Kompilierfehler | Umbenannt in `TestEncodeDirectArcCmd` |
| protocols/bridge/powerlink/powerlink_extra_test.go:51 | Erwartung `write 0x1000 A`, tatsächlich `0A` (`%X` formatiert Bytes fix zweistellig, konsistent mit der `DEADBEEF`-Konvention bestehender Tests) | Erwartung auf `0A` geändert |
| protocols/bridge/sercos/sercos_extra_test.go:62 | wie oben, Erwartung `1` sollte `01` sein | Erwartung auf `01` geändert |
| protocols/ethernet/opcua/opcua_extra_test.go:35 | Assertion für das alte Bug-Verhalten (Länge an Offset 8) wurde nach der Quellcode-Korrektur wertlos | Geändert: Assertion prüft das Protokollversionsfeld auf 0 (Länge wird von TestExtraHELMessageSizeAtOffset4 abgedeckt) |
| protocols/fieldbus/canopen/canopen_extra_test.go:135 | „Prove-it“-Test mit invertierter Assertion (`err != nil` → Fehlschlag, obwohl der Parsing-Fehler genau das erwartete Ergebnis ist) | Geändert: Fehlschlag bei `err == nil` |
| protocols/ethernet/profinet/profinet_extra_test.go:23 | wie oben, invertierte Assertion | Geändert: Fehlschlag bei `err == nil` |
| protocols/fieldbus/dnp3/dnp3_extra_test.go:35,53 | wie oben ×2 | Geändert: Fehlschlag bei `err == nil` |

> Hinweis: Die `stubTransport`-Schnittstellenimplementierung in `kernel/security/tls_test.go` und die CRC-Erwartungswerte von `cclink/cclinkie` wurden vom Testingenieur während der parallelen Arbeit selbst korrigiert; in diesem Durchgang nicht angefasst.

## 4. Verbleibende Risiken

1. **Module mit geringer Abdeckung**: profibus (35,7 %), lonworks / asinterface / foundationfieldbus / iolink (37,9 %) — die Tests decken nur wenige Pfade ab; es wird empfohlen, später Decode/Encode-Zweige und Fehlerpfade zu ergänzen.
2. **Modbus-RTU-CRC ungeprüft**: `decodeRTU` prüft die CRC nicht (der Test `TestExtraRTUDecodesCorruptCRC` dokumentiert diese Lücke ausdrücklich und besteht unverändert). Bei der Kommunikation mit echten Geräten sollte die Prüfung ergänzt werden.
3. **Einschränkung des CANopen-Wire-Formats**: `marshalCAN` überträgt nur die ersten 4 Datenbytes (`TestExtraSDOWritePayloadLostOnWire` dokumentiert dies); SDO-Schreibvorgänge mit mehreren Bytes Payload gehen verloren.
4. **Hardwareabhängige Tests übersprungen**: cpci / pci / vme verwenden ohne Hardware `t.Skip` (sinnvoll).
5. **Temporäre Dateien im Repository**: Nicht versionierte `main.go` im Root (Sondenprogramm, referenziert das nicht existierende `probe/`-Modul), `scripts/`, `docs/*.png` gehören nicht zu dieser Lieferung; Bereinigung empfohlen.
6. **Tests werden noch parallel eingereicht**: Dieser Bericht basiert auf einem Snapshot des letzten vollständig grünen Testlaufs; falls der Testingenieur weitere `*_extra_test.go`-Dateien einreicht, muss eine vollständige Regression erneut laufen.
