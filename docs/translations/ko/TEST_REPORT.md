# 테스트 보고서 — industrial-protocols-go

날짜: 2026-08-27
범위: Go 멀티 모듈 워크스페이스(go.work, 42개 모듈: kernel / protocols / examples)
명령: `go test ./...` 및 `go test -cover ./...`(워크스페이스 루트에 go.mod가 없어 모듈별로 동일하게 실행했으며 결과는 동일함)

## 1. 종합 결론

- **전부 통과**: 42개 모듈, 테스트가 포함된 52개 패키지 모두 통과(병렬 테스트 작성 중에도 전체 재실행으로 여러 차례 검증).
- 테스트 파일 106개, 테스트 함수 659개.
- `go vet ./...` 전 모듈 경고 없음.
- 평균 문장 커버리지 **83.3%**; kernel 핵심 패키지 커버리지 95%~100%.

## 2. 모듈별 테스트 통계 및 커버리지

### kernel(11개 패키지, examples 제외하고 커버리지 최고)

| 패키지 | 커버리지 |
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

### protocols 주요 모듈

| 모듈 | 커버리지 | 모듈 | 커버리지 |
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

> examples/modbus_basic에는 테스트 파일이 없음(예제 프로그램으로 정상적인 상태).

## 3. 수정 내역

### 3.1 소스 코드 버그 수정(tester-kernel / tester-protocols가 발견 및 검증)

| 파일:라인 | 문제 | 수정 |
|---|---|---|
| kernel/gateway/engine.go:27 | `Transform`이 `Src`/`Dst`/`Map`이 nil인 규칙에 대해 `rule.Src.Name()` / `rule.Dst.Name()` / `rule.Map(...)`에서 nil 역참조 panic | 루프 시작 시 미완성 규칙 건너뛰기(기존 "잘못된 규칙 건너뛰기" 의미와 일치) |
| protocols/ethernet/modbus/modbus.go:215 | `decodePDU`가 1바이트 예외 PDU(예: `0x81`)에 대해 `pdu[1]` 읽기 범위 초과 panic | `len(pdu) < 2` 검사 추가, 파싱 오류 반환 |
| protocols/ethernet/modbus/modbus.go:230 | `decodePDU` 바이트 카운트(`pdu[1]`)가 실제 데이터를 초과할 때 슬라이스 범위 초과 panic(악의적/손상된 프레임) | `2+n > len(pdu)` 경계 검사 추가, 파싱 오류 반환 |
| protocols/ethernet/opcua/opcua.go:63,74 | `encodeHello`가 메시지 길이를 오프셋 8(프로토콜 버전 슬롯)에 역입력, 규격상 오프셋 4여야 함 | `sizePos`(오프셋 4)를 미리 기록해 올바른 위치에 역입력 |
| protocols/ethernet/opcua/opcua.go:82,89 | `encodeOpenSecureChannel`의 길이 필드가 역입력되지 않아 항상 0 | 끝에 실제 프레임 길이로 역입력 |
| protocols/ethernet/ethernetip/ethernetip.go:111,117 | `encodeReadTag`가 16+len(cipReq) 프레임을 할당하지만 CIP 요청을 [28:]에 기록, 초과된 12바이트로 `copy`가 CIP 요청 전체를 조용히 버림; 길이 필드도 12만큼 부족 | 할당을 28+len(cipReq)로, length 필드를 4+len(cipReq)로 수정 |
| protocols/ethernet/profinet/profinet.go:72 | `Decode`가 범위를 벗어난 `blockLen`으로 슬라이스 범위 초과 panic | `len(data) < 10+blockLen` 검사 추가 |
| protocols/fieldbus/canopen/canopen.go:102 | `Decode` 최소 길이 검사가 4지만 `unmarshalCAN`은 `data[4:8]` 필요, 4~7바이트 프레임에서 panic | 최소 길이를 8로 수정(라인 형식은 고정 8바이트) |
| protocols/fieldbus/dnp3/dnp3.go:121 | `decodeTPDU`가 `length<5`(`apduEnd < apduStart`) 또는 과도한 `length`(`apduEnd > len(data)`)에서 슬라이스 범위 초과 panic | `length < 5 || apduEnd > len(data)` 통합 검증 후 오류 반환 |
| protocols/fieldbus/dnp3/dnp3.go:157 | `crc16DNP` 작업 변수가 규격대로 `& 0xFF` 마스킹되지 않아 기지 벡터 `crc16DNP("123456789")`가 0x69FF로 산출(0xEA82여야 함) | `temp := (crc ^ b) & 0xFF` |
| protocols/automotive/kline/kline.go:74 | `Decode`의 `len<5` 검사가 fast-init 동기화 프레임 판단보다 먼저 실행되어 1바이트 동기화 프레임(0x55)이 오거부됨 | 빈 프레임과 동기화 프레임을 먼저 판단한 후 길이 검사 |

### 3.2 테스트 수정(컴파일 오류 / 어서션 오류 / 중복 이름 / 낡은 어서션 수정)

| 파일:라인 | 문제 | 수정 |
|---|---|---|
| protocols/automotive/kline/kline_extra_test.go:43 | 상수 `0x68+0x10+0xF1+0x01`(362)을 byte에 대입해 컴파일 실패 | 기대값을 체크섬 하위 8비트 `0x6A`로 수정 |
| protocols/automotive/saej1850/saej1850_extra_test.go | `kernel.Codec` 인터페이스에서 미공개 메서드/필드 호출로 컴파일 실패 | `newCodec` helper를 추가해 `*j1850Codec`로 단언; `TestEncodeMode01PID` 오프셋 단언을 raw[4..6]으로 수정(ID 4바이트 + 데이터 앞 4바이트) |
| protocols/bridge/sercos1/sercos1_extra_test.go:45 | `TestEncodeReadDefaults`가 기존 테스트와 이름 중복으로 컴파일 실패 | `TestEncodeReadDefaultsFiber`로 이름 변경(fiber 변형 커버리지 유지) |
| protocols/building/dali/dali_extra_test.go:68 | `TestEncodeDirectArc`가 기존 테스트와 이름 중복으로 컴파일 실패 | `TestEncodeDirectArcCmd`로 이름 변경 |
| protocols/bridge/powerlink/powerlink_extra_test.go:51 | 기대값 `write 0x1000 A`, 실제는 `0A`(`%X`가 바이트당 두 자리 고정, 기존 `DEADBEEF` 테스트 규약과 일치) | 기대값을 `0A`로 수정 |
| protocols/bridge/sercos/sercos_extra_test.go:62 | 위와 동일, 기대값 `1`은 `01`이어야 함 | 기대값을 `01`로 수정 |
| protocols/ethernet/opcua/opcua_extra_test.go:35 | 이전 버그 동작(오프셋 8에 길이)을 단언, 소스 수정 후 무효화됨 | 프로토콜 버전 필드가 0임을 단언하도록 수정(길이는 TestExtraHELMessageSizeAtOffset4가 커버) |
| protocols/fieldbus/canopen/canopen_extra_test.go:135 | "prove-it" 테스트 어서션 반전(`err != nil`이면 실패, 하지만 파싱 오류가 바로 기대 결과) | `err == nil`이면 실패로 수정 |
| protocols/ethernet/profinet/profinet_extra_test.go:23 | 위와 동일, 어서션 반전 | `err == nil`이면 실패로 수정 |
| protocols/fieldbus/dnp3/dnp3_extra_test.go:35,53 | 위와 동일 ×2 | `err == nil`이면 실패로 수정 |

> 참고: `kernel/security/tls_test.go`의 `stubTransport` 인터페이스 구현 문제와 `cclink/cclinkie`의 CRC 기대값은 테스트 엔지니어가 병렬 작업 중 자체적으로 수정했으며, 이번에는 변경하지 않음.

## 4. 잔여 리스크

1. **낮은 커버리지 모듈**: profibus(35.7%), lonworks / asinterface / foundationfieldbus / iolink(37.9%) — 테스트가 일부 경로만 커버, 추후 Decode/Encode 분기와 예외 경로 보완 권장.
2. **modbus RTU CRC 미검증**: `decodeRTU`가 CRC를 검증하지 않음(테스트 `TestExtraRTUDecodesCorruptCRC`가 이 GAP을 명시적으로 문서화하고 현 상태로 통과). 실제 장비와의 상호 운용 시 보완 권장.
3. **canopen 라인 형식 제한**: `marshalCAN`이 앞 4바이트 데이터만 전달(`TestExtraSDOWritePayloadLostOnWire`로 문서화), SDO 다중 바이트 페이로드 쓰기 시 손실됨.
4. **하드웨어 의존 스킵**: cpci / pci / vme는 하드웨어가 없는 환경에서 `t.Skip`(합리적).
5. **저장소 잔여 임시 파일**: 루트의 추적되지 않는 `main.go`(탐지 프로그램, 존재하지 않는 `probe/` 모듈 참조), `scripts/`, `docs/*.png`는 이번 인도물에 포함되지 않으므로 정리 권장.
6. **테스트가 아직 병렬로 진행 중**: 본 보고서는 마지막 전체 통과 테스트 스냅샷 기준. 테스트 엔지니어가 이후 계속 `*_extra_test.go`를 커밋한다면 전체 회귀 테스트를 다시 실행해야 함.
