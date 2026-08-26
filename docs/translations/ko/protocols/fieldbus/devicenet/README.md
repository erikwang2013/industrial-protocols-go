# DeviceNet 프로토콜 SDK

DeviceNet은 공장 자동화를 위한 CAN 기반 산업 네트워크 프로토콜입니다. 이 패키지는 SocketCAN 및 TCP 게이트웨이 드라이버를 모두 갖춘 DeviceNet 코덱을 제공합니다.

## 프로토콜 개요

DeviceNet은 CAN 상에서 Common Industrial Protocol(CIP)을 사용합니다. 사전 정의된 Master/Slave 연결 세트는 폴링 I/O 및 explicit messaging에 Group 2 메시지(CAN ID 0x400 + NodeID)를 사용합니다.

## 하드웨어 요구 사항

### CAN 모드(SocketCAN)
- SocketCAN(`CONFIG_CAN` 활성화)을 지원하는 **Linux** 시스템
- CAN 인터페이스(예: `can0`, `vcan0`)
- DeviceNet 지원 CAN 하드웨어(예: Anybus Communicator, HMS IXXAT)

### 게이트웨이 모드
- DeviceNet 게이트웨이에 대한 TCP/IP 연결
- 일반 텍스트 명령 프로토콜을 지원하는 게이트웨이(예: HMS Anybus, Hilscher netX)

### 가상 CAN 인터페이스 설정(테스트용)

```bash
sudo modprobe can
sudo modprobe can_raw
sudo modprobe vcan
sudo ip link add dev vcan0 type vcan
sudo ip link set up vcan0
```

## 사용법

### CAN 모드

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/fieldbus/devicenet"
)

func main() {
    p := devicenet.New()
    codec, _ := p.NewCodec("can")

    // Poll 요청
    req := &kernel.Request{
        Function: "poll",
        Data:     []byte{0x01, 0x00},
    }
    raw, _ := codec.Encode(req)

    // 연결 열기
    req2 := &kernel.Request{Function: "open"}
    raw2, _ := codec.Encode(req2)
    _ = raw
    _ = raw2
}
```

### 게이트웨이 모드

```go
p := devicenet.New()
codec, _ := p.NewCodec("gateway")

req := &kernel.Request{
    Function: "poll",
    Data:     []byte{0xAB, 0xCD},
}
raw, _ := codec.Encode(req)
// raw will be: "poll abcd\n"
```

## 지원 함수

| 함수 | 설명                    |
|----------|--------------------------------|
| `open`   | Explicit 연결 열기       |
| `poll`   | I/O 데이터 폴링(Group 2)        |

## 테스트

```bash
go test ./... -v
```

참고: SocketCAN 드라이버 테스트에는 Linux가 필요합니다. 게이트웨이 드라이버 테스트는 모든 OS에서 실행됩니다.
