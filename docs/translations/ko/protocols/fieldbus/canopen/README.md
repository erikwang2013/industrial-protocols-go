# CANopen 프로토콜 SDK

CANopen은 임베디드 제어 시스템을 위한 CAN 기반 상위 계층 프로토콜입니다. 이 패키지는 CANopen 코덱과 SocketCAN 드라이버를 제공합니다.

## 프로토콜 개요

CANopen은 다음과 같은 사전 정의된 연결 세트와 함께 표준 11비트 CAN 식별자를 사용합니다:

| 기능    | CAN ID              | 설명                   |
|------------|---------------------|-------------------------------|
| NMT        | 0x000               | 네트워크 관리            |
| SYNC       | 0x080               | 동기화 메시지       |
| SDO (tx)   | 0x580 + NodeID      | Service Data Object(서버)  |
| SDO (rx)   | 0x600 + NodeID      | Service Data Object(클라이언트)  |
| PDO1 (tx)  | 0x180 + NodeID      | Process Data Object 1         |
| Heartbeat  | 0x700 + NodeID      | Heartbeat / Bootup            |

## 하드웨어 요구 사항

- SocketCAN(`CONFIG_CAN` 활성화)을 지원하는 **Linux** 시스템
- CAN 인터페이스(예: `can0`, 가상 CAN용 `vcan0`)
- CAN 지원 트랜시버 하드웨어(예: MCP2515, SJA1000 또는 USB-CAN 어댑터)

### 가상 CAN 인터페이스 설정(테스트용)

```bash
sudo modprobe can
sudo modprobe can_raw
sudo modprobe vcan
sudo ip link add dev vcan0 type vcan
sudo ip link set up vcan0
```

## 사용법

```go
package main

import (
    "fmt"
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/fieldbus/canopen"
)

func main() {
    p := canopen.New()
    codec, _ := p.NewCodec("can")

    // SDO 오브젝트 딕셔너리 항목 읽기
    req := &kernel.Request{
        Function: "sdo_read",
        Metadata: map[string]any{
            "index": float64(0x1000), // Device type
            "sub":   float64(0),
        },
    }
    raw, _ := codec.Encode(req)
    fmt.Printf("SDO read frame: %X\n", raw)

    // NMT 원격 노드 시작
    req2 := &kernel.Request{Function: "nmt_start"}
    raw2, _ := codec.Encode(req2)
    fmt.Printf("NMT start frame: %X\n", raw2)
}
```

## 지원 함수

| 함수     | 설명                          |
|-------------|--------------------------------------|
| `sdo_read`  | 오브젝트 딕셔너리 항목 읽기         |
| `sdo_write` | 오브젝트 딕셔너리 항목 쓰기        |
| `nmt_start` | 원격 노드 시작(NMT)              |
| `nmt_stop`  | 원격 노드 정지(NMT)               |
| `nmt_reset` | 원격 노드 리셋(NMT)              |
| `heartbeat` | Heartbeat / Bootup 메시지 전송      |

## 테스트

```bash
go test ./... -v
```

참고: SocketCAN 드라이버 테스트에는 CAN 하드웨어 또는 가상 CAN 인터페이스가 있는 Linux 시스템이 필요합니다.
