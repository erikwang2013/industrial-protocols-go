# ISA100.11a CmdBridge SDK

ISA100.11a는 프로세스 자동화를 위한 무선 산업 네트워킹 표준입니다. 이 패키지는 `yfgw410_cli`(Yokogawa YFGW410 필드 무선 게이트웨이) 커맨드라인 유틸리티를 통한 ISA100.11a 코덱을 제공합니다.

## CLI 도구

Yokogawa YFGW410 Field Wireless Gateway CLI를 사용합니다.

### 설치

```bash
# Yokogawa YFGW410 게이트웨이 소프트웨어 및 도구 설치
# Field Wireless Gateway 설정에 대해서는 Yokogawa 문서를 참조
```

확인: `yfgw410_cli --help`

## 사용법

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/isa100"
)

func main() {
    p := isa100.New()
    codec, _ := p.NewCodec("cmd")

    req := &kernel.Request{
        Function: "read",
        Address:  "DEV001",
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## 드라이버

```go
b, c, err := isa100.NewCmdDriver("/usr/bin/yfgw410_cli")
```

## 지원 함수

| 함수 | 설명                 |
|----------|-----------------------------|
| `read`   | 장치 속성 읽기       |
| `write`  | 장치 속성 쓰기      |
| `list`   | 프로비저닝된 장치 목록    |

## 테스트

```bash
go test ./... -v
```
