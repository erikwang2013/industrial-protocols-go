# Test Report — industrial-protocols-go

Date: 2026-08-27
Scope: Go multi-module workspace (go.work, 42 modules: kernel / protocols / examples)
Commands: `go test ./...` and `go test -cover ./...` (executed per module with equivalent results, since the workspace root has no go.mod)

## 1. Overall Conclusion

- **All green**: all 42 modules and 52 packages containing tests passed (verified with multiple full reruns, including during parallel test writes).
- 106 test files, 659 test functions.
- `go vet ./...` passes for all modules with no warnings.
- Average statement coverage **83.3%**; kernel core packages at 95%~100%.

## 2. Per-Module Test Statistics and Coverage

### kernel (11 packages, highest coverage outside examples)

| Package | Coverage |
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

### protocols key modules

| Module | Coverage | Module | Coverage |
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

> examples/modbus_basic has no test files (example program, as expected).

## 3. Fix List

### 3.1 Source code bug fixes (found and verified by tester-kernel / tester-protocols)

| File:line | Issue | Fix |
|---|---|---|
| kernel/gateway/engine.go:27 | `Transform` nil-pointer panics at `rule.Src.Name()` / `rule.Dst.Name()` / `rule.Map(...)` for rules with nil `Src`/`Dst`/`Map` | Skip incompletely configured rules at loop start (consistent with existing "skip bad rules" semantics) |
| protocols/ethernet/modbus/modbus.go:215 | `decodePDU` reads `pdu[1]` out of bounds for 1-byte exception PDUs (e.g. `0x81`) | Added `len(pdu) < 2` check, returns parse error |
| protocols/ethernet/modbus/modbus.go:230 | `decodePDU` byte count (`pdu[1]`) exceeding actual data causes slice out-of-bounds panic (malicious/corrupt frames) | Added `2+n > len(pdu)` boundary check, returns parse error |
| protocols/ethernet/opcua/opcua.go:63,74 | `encodeHello` writes the message length back at offset 8 (protocol version slot); the spec requires offset 4 | Record `sizePos` (offset 4) early, write back at the correct position |
| protocols/ethernet/opcua/opcua.go:82,89 | `encodeOpenSecureChannel` length field never written back, always 0 | Write back the actual frame length at the end |
| protocols/ethernet/ethernetip/ethernetip.go:111,117 | `encodeReadTag` allocates 16+len(cipReq) but the CIP request is written at [28:]; the extra 12 bytes make `copy` silently drop the entire CIP request; length field also 12 short | Allocation changed to 28+len(cipReq), length field changed to 4+len(cipReq) |
| protocols/ethernet/profinet/profinet.go:72 | `Decode` slices out of bounds for out-of-range `blockLen` | Added `len(data) < 10+blockLen` check |
| protocols/fieldbus/canopen/canopen.go:102 | `Decode` minimum length check is 4, but `unmarshalCAN` needs `data[4:8]`; 4~7-byte frames panic | Minimum length changed to 8 (wire format fixed at 8 bytes) |
| protocols/fieldbus/dnp3/dnp3.go:121 | `decodeTPDU` slices out of bounds for `length<5` (`apduEnd < apduStart`) or oversized `length` (`apduEnd > len(data)`) | Unified validation `length < 5 || apduEnd > len(data)` before returning an error |
| protocols/fieldbus/dnp3/dnp3.go:157 | `crc16DNP` working variable not masked `& 0xFF` per spec; known vector `crc16DNP("123456789")` gives 0x69FF (should be 0xEA82) | `temp := (crc ^ b) & 0xFF` |
| protocols/automotive/kline/kline.go:74 | `Decode`'s `len<5` check runs before the fast-init sync frame check; 1-byte sync frames (0x55) are wrongly rejected | Check empty frame and sync frame first, then length |

### 3.2 Test fixes (compilation errors / wrong assertions / duplicate names / stale assertions)

| File:line | Issue | Fix |
|---|---|---|
| protocols/automotive/kline/kline_extra_test.go:43 | Constant `0x68+0x10+0xF1+0x01`(362) assigned to a byte fails to compile | Expected value changed to low 8 bits of checksum `0x6A` |
| protocols/automotive/saej1850/saej1850_extra_test.go | Calls unexported methods/fields on the `kernel.Codec` interface, fails to compile | Added `newCodec` helper asserting `*j1850Codec`; `TestEncodeMode01PID` offset corrected to raw[4..6] (4-byte ID + first 4 data bytes) |
| protocols/bridge/sercos1/sercos1_extra_test.go:45 | `TestEncodeReadDefaults` duplicates an existing test, fails to compile | Renamed `TestEncodeReadDefaultsFiber` (keeps fiber variant coverage) |
| protocols/building/dali/dali_extra_test.go:68 | `TestEncodeDirectArc` duplicates an existing test, fails to compile | Renamed `TestEncodeDirectArcCmd` |
| protocols/bridge/powerlink/powerlink_extra_test.go:51 | Expected `write 0x1000 A`, actual is `0A` (`%X` formats two fixed digits per byte, consistent with existing `DEADBEEF` convention) | Expected changed to `0A` |
| protocols/bridge/sercos/sercos_extra_test.go:62 | Same as above, expected `1` should be `01` | Expected changed to `01` |
| protocols/ethernet/opcua/opcua_extra_test.go:35 | Asserts old buggy behavior (length at offset 8), stale after source fix | Changed to assert protocol version field is 0 (length covered by TestExtraHELMessageSizeAtOffset4) |
| protocols/fieldbus/canopen/canopen_extra_test.go:135 | "prove-it" test assertion inverted (fails when `err != nil`, but a parse error is the expected result) | Changed to fail on `err == nil` |
| protocols/ethernet/profinet/profinet_extra_test.go:23 | Same as above, assertion inverted | Changed to fail on `err == nil` |
| protocols/fieldbus/dnp3/dnp3_extra_test.go:35,53 | Same as above ×2 | Changed to fail on `err == nil` |

> Note: `kernel/security/tls_test.go` `stubTransport` interface implementation issues and the `cclink/cclinkie` CRC expected values were fixed independently by the test engineers during parallel work; not modified here.

## 4. Remaining Risks

1. **Low-coverage modules**: profibus (35.7%), lonworks / asinterface / foundationfieldbus / iolink (37.9%) — tests cover only a few paths; recommend adding Decode/Encode branches and error paths later.
2. **Modbus RTU CRC not validated**: `decodeRTU` does not check CRC (test `TestExtraRTUDecodesCorruptCRC` explicitly documents this GAP and passes as-is). Recommend adding it when interop with real devices matters.
3. **CANopen wire format limitation**: `marshalCAN` carries only the first 4 data bytes (`TestExtraSDOWritePayloadLostOnWire` documents this); multi-byte SDO write payloads are lost.
4. **Hardware-dependent skips**: cpci / pci / vme use `t.Skip` in environments without hardware (reasonable).
5. **Leftover temp files in repo**: untracked root `main.go` (probe program referencing nonexistent `probe/` module), `scripts/`, `docs/*.png` — not part of this delivery; recommend cleanup.
6. **Tests still being landed in parallel**: this report is based on the last full green snapshot; if test engineers continue committing `*_extra_test.go`, a full regression rerun is needed.
