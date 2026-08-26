# Modbus Plus 프로토콜 SDK

Modbus Plus(MB+)는 Modicon(Schneider Electric)이 개발한 고속 토큰 전달 방식 산업 네트워크입니다.

## 변형

| 변형 | 브리지 유형 | 설명 |
|---------|-------------|-------------|
| `gateway` | GatewayBridge | SA85/BM85 어댑터에 대한 TCP 연결 |
| `cmd` | CmdBridge | `sa85_cli` 유틸리티용 CLI 래퍼 |

## 게이트웨이 드라이버

```go
driver, codec, err := modbusplus.NewGatewayDriver("192.168.0.160:2010")
handler, err := modbusplus.ReadyHandler("192.168.0.160:2010")
```

## Cmd 드라이버

```go
driver, codec, err := modbusplus.NewCmdDriver("/usr/bin/sa85_cli")
```

### CLI 도구 설치

```bash
# SA85 Modbus Plus 드라이버 및 도구 설치
# SA85 어댑터에 대해서는 Schneider Electric 문서를 참조
```

확인: `sa85_cli --help`

## 하드웨어

| 게이트웨이 | 인터페이스 | 기본 IP |
|---------|-----------|-------------|
| Schneider SA85 | ISA Modbus Plus 어댑터 | 호스트 할당 |
| Schneider BM85 | Modbus Plus 브리지/멀티플렉서 | 192.168.0.A0 |
| ProSoft MVI56-MBP | ControlLogix MB+ 모듈 | 호스트 할당 |

## 배선

- MB+ 트렁크용 BNC 커넥터가 있는 트윈액시얼 케이블(RG-62)
- 종단 저항(각 끝에 78옴)
- BM85 브리지가 이더넷을 통해 MB+를 TCP에 연결

## 프로토콜 프레임 형식

- Magic(2바이트): `0x4D42`
- Destination(1바이트): 노드 주소
- Command(1바이트): 0x01=read, 0x02=write
- Length(2바이트): 페이로드 크기(빅 엔디언)
- Payload(N바이트): 데이터
- CRC(2바이트): Modbus CRC-16(리틀 엔디언)

## 테스트

```bash
go test ./... -v
```
