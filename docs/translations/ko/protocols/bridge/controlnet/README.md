# ControlNet CmdBridge SDK

ControlNet은 Allen-Bradley(Rockwell Automation)가 개발한 고속·시간 결정적 데이터 교환을 위한 실시간 산업 네트워크 프로토콜입니다. 이 패키지는 `1784-pcic-cli` 커맨드라인 유틸리티를 통한 ControlNet 코덱을 제공합니다.

## CLI 도구

1784-PCIC ControlNet 인터페이스 카드 CLI 유틸리티를 사용합니다.

### 설치

```bash
# Rockwell 1784-PCIC 드라이버 및 도구 설치
# RSLinx Classic SDK에 대해서는 Rockwell Automation 문서를 참조
```

확인: `1784-pcic-cli --help`

## 사용법

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/controlnet"
)

func main() {
    p := controlnet.New()
    codec, _ := p.NewCodec("cmd")

    req := &kernel.Request{
        Function: "read",
        Address:  "0x10",
        Count:    8,
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## 드라이버

```go
b, c, err := controlnet.NewCmdDriver("/usr/bin/1784-pcic-cli")
```

## 지원 함수

| 함수 | 설명               |
|----------|---------------------------|
| `read`   | ControlNet 노드에서 읽기 |
| `write`  | ControlNet 노드에 쓰기  |
| `status` | PCIC 카드 상태 조회    |

## 테스트

```bash
go test ./... -v
```
