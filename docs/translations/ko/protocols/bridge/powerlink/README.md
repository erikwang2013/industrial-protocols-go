# POWERLINK CmdBridge SDK

POWERLINK(Ethernet POWERLINK)는 산업 자동화를 위한 실시간 이더넷 프로토콜입니다. 이 패키지는 `openPOWERLINK_demo` 커맨드라인 유틸리티를 통한 POWERLINK 코덱을 제공합니다.

## CLI 도구

openPOWERLINK 스택 데모 애플리케이션을 사용합니다.

### 설치

```bash
git clone https://github.com/OpenAutomationTechnologies/openPOWERLINK.git
cd openPOWERLINK
cmake -B build && cmake --build build
sudo cmake --install build
```

확인: `openPOWERLINK_demo --help`

## 사용법

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/powerlink"
)

func main() {
    p := powerlink.New()
    codec, _ := p.NewCodec("cmd")

    req := &kernel.Request{
        Function: "read",
        Address:  "0x2000",
        Count:    8,
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## 드라이버

```go
b, c, err := powerlink.NewCmdDriver("/usr/bin/openPOWERLINK_demo")
```

## 지원 함수

| 함수 | 설명                   |
|----------|-------------------------------|
| `read`   | 오브젝트 딕셔너리 항목 읽기  |
| `write`  | 오브젝트 딕셔너리 항목 쓰기 |
| `status` | 노드/NMT 상태 조회          |

## 테스트

```bash
go test ./... -v
```
