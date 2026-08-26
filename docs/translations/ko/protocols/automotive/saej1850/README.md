# SAE J1850 CAN 프로토콜 SDK

SAE J1850은 OBD-II(온보드 진단)에 사용되는 차량 통신 표준입니다. 이 패키지는 CAN 버스(ISO 15765-4 / CAN TP) 기반 J1850 코덱을 제공합니다.

## 프로토콜 개요

SAE J1850 OBD-II over CAN은 다음과 같은 구조의 29비트 확장 CAN 식별자를 사용합니다:

### CAN ID 형식(29비트)

| Bits       | 필드    | 설명                        |
|-----------|----------|------------------------------------|
| 28-26     | Priority | 메시지 우선순위(0-7, 기본값 6)  |
| 25        | Ext ID   | 확장 프레임의 경우 항상 1       |
| 24-16     | PF       | Parameter Format(헤더)          |
| 15-8      | PS       | Parameter Specific(대상/출발지) |
| 7-0       | SA       | 출발지 주소                     |

### 표준 OBD-II CAN ID

| 유형                | CAN ID(16진수)    | 설명                    |
|--------------------|-----------------|--------------------------------|
| Physical Request   | 0x18DAxxF1      | 특정 ECU에 요청(xx=주소) |
| Physical Response  | 0x18DAF1xx      | ECU의 응답(xx=주소)   |
| Functional Request | 0x18DB33F1      | 모든 ECU에 브로드캐스트         |

### ISO 15765-2 프레임 형식

단일 프레임: 바이트 0의 상위 니블 = 데이터 길이(0-7), 하위 니블 + 나머지 바이트 = 진단 데이터.

## 하드웨어 요구 사항

- SocketCAN을 지원하는 **Linux** 시스템
- J1850-CAN OBD-II 어댑터(예: ELM327 호환 USB-to-CAN, OBDLink SX, Macchina M2)
- OBD-II 커넥터가 있는 차량(1996년 이후 대부분의 차량)

### CAN 인터페이스 설정

```bash
sudo modprobe can
sudo modprobe can_raw
sudo ip link set can0 type can bitrate 500000
sudo ip link set up can0
```

## 사용법

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/automotive/saej1850"
)

func main() {
    p := saej1850.New()
    codec, _ := p.NewCodec("can")

    // Mode $01 PID $0C: 엔진 RPM
    req := &kernel.Request{
        Function: "mode01",
        Metadata: map[string]any{"pid": float64(0x0C)},
    }
    raw, _ := codec.Encode(req)
    _ = raw

    // Mode $03: 배출가스 관련 DTC 요청
    req2 := &kernel.Request{Function: "mode03"}
    raw2, _ := codec.Encode(req2)
    _ = raw2
}
```

## 지원 함수

| 함수        | OBD-II 모드 | 설명                         |
|----------------|-------------|-------------------------------------|
| `mode01`       | $01         | 현재 파워트레인 데이터 요청     |
| `mode03`       | $03         | 배출가스 관련 DTC 요청       |
| `mode0A`       | $0A         | 영구 DTC 요청              |
| `diag_request` | 커스텀      | 범용 진단 요청          |
| `diag_response`| -           | 진단 응답 프레임           |
| `broadcast`    | -           | 모든 ECU에 기능적 브로드캐스트    |

## 테스트

```bash
go test ./... -v
```
