# SERCOS I/II CmdBridge SDK

SERCOS I/II는 디지털 모션 제어용 SERCOS 인터페이스의 레거시 직렬 광섬유 버전입니다. 이 패키지는 `sercos_cli` 커맨드라인 유틸리티를 통한 SERCOS I/II 코덱을 제공합니다.

## CLI 도구

SERCOS 광섬유 인터페이스 CLI 유틸리티를 사용합니다.

### 설치

```bash
# SERCOS 인터페이스 카드 드라이버 및 도구
# 특정 SERCOS 마스터 카드에 대해서는 벤더 문서를 참조
```

확인: `sercos_cli --help`

## 사용법

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/sercos1"
)

func main() {
    p := sercos1.New()
    codec, _ := p.NewCodec("cmd")

    req := &kernel.Request{
        Function: "read",
        Address:  "0x0100",
        Count:    4,
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## 드라이버

```go
b, c, err := sercos1.NewCmdDriver("/usr/bin/sercos_cli")
```

## 지원 함수

| 함수 | 설명              |
|----------|--------------------------|
| `read`   | SERCOS IDN 읽기          |
| `write`  | SERCOS IDN 쓰기         |
| `status` | 드라이브 상태 조회       |

## 테스트

```bash
go test ./... -v
```
