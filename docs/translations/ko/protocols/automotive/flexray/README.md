# FlexRay CAN 프로토콜 SDK

FlexRay는 고속 결정적(deterministic) 자동차 통신 프로토콜입니다. 이 패키지는 CAN 버스 기반 FlexRay 코덱을 제공하며, CRC-16/XMODEM 무결성 검사가 포함된 사이클 기반 프레이밍을 지원합니다.

## 프로토콜 개요

FlexRay는 반복되는 통신 사이클을 가진 TDMA(시분할 다중 접속) 방식을 사용합니다. 각 사이클은 정적 세그먼트와 동적 세그먼트로 구성됩니다. 이 코덱은 FlexRay 프레임을 확장 29비트 CAN 프레임에 매핑합니다.

### 와이어 형식

FlexRay 사이클 페이로드:
- **Header**(2바이트): 사이클 번호(리틀 엔디언)
- **Status**(1바이트): bit 7=PPI(Payload Preamble Indicator), bit 6=NFI, bit 5=SYF, bit 4=SUF
- **Data**(N바이트): 페이로드(최대 254바이트)
- **CRC**(2바이트): header+status+data에 대한 CRC-16/XMODEM(리틀 엔디언)

CAN ID 인코딩(29비트 확장):
- Bits 28-24: 메시지 타입(0x01=프레임, 0x02=상태)
- Bits 23-10: 예약됨
- Bits 15-10: 슬롯 ID(6비트)
- Bits 9-0: 사이클 번호(10비트)

## 하드웨어 요구 사항

- SocketCAN을 지원하는 **Linux** 시스템
- Vector VN7600/VN7640 또는 Bosch FlexRay-CAN 어댑터
- 적절한 종단 처리(2.5V 바이어스)가 된 FlexRay 네트워크

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
    "github.com/erikwang2013/industrial-protocols-go/protocols/automotive/flexray"
)

func main() {
    p := flexray.New()
    codec, _ := p.NewCodec("can")

    // PPI가 설정된 사이클 5의 FlexRay 프레임 전송
    req := &kernel.Request{
        Function: "frame",
        Data:     []byte{0x42, 0x01},
        Metadata: map[string]any{
            "cycle": float64(5),
            "ppi":   true,
        },
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## 지원 함수

| 함수 | 설명                              |
|----------|------------------------------------------|
| `frame`  | FlexRay 프레임 페이로드 전송               |
| `status` | 슬롯 구성 조회(slot, cycle)   |

## 테스트

```bash
go test ./... -v
```
