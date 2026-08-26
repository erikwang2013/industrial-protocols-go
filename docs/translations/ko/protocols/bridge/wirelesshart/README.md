# WirelessHART CmdBridge SDK

WirelessHART(IEC 62591)는 HART 프로토콜을 기반으로 하는 무선 산업 네트워킹 표준입니다. 이 패키지는 `emerson_1410_cli`(Emerson 1410/1420 Wireless Gateway) 커맨드라인 유틸리티를 통한 WirelessHART 코덱을 제공합니다.

## CLI 도구

Emerson 1410/1420 Wireless Gateway CLI 유틸리티를 사용합니다.

### 설치

```bash
# Emerson Wireless Gateway 소프트웨어 및 도구 설치
# 1410/1420 Gateway 설정에 대해서는 Emerson 문서를 참조
```

확인: `emerson_1410_cli --help`

## 사용법

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/wirelesshart"
)

func main() {
    p := wirelesshart.New()
    codec, _ := p.NewCodec("cmd")

    req := &kernel.Request{
        Function: "read",
        Address:  "TT101",
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## 드라이버

```go
b, c, err := wirelesshart.NewCmdDriver("/usr/bin/emerson_1410_cli")
```

## 지원 함수

| 함수 | 설명                |
|----------|----------------------------|
| `read`   | 장치 파라미터 읽기      |
| `write`  | 장치 파라미터 쓰기     |
| `scan`   | 무선 장치 스캔  |

## 테스트

```bash
go test ./... -v
```
