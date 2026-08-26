# SERCOS III CmdBridge SDK

SERCOS III(SErial Real-time COmmunication System)은 광섬유 또는 동축 케이블을 통한 링 토폴로지를 사용하는 모션 제어용 디지털 인터페이스입니다. 이 패키지는 `netx_cli` 커맨드라인 유틸리티를 통한 SERCOS III 코덱을 제공합니다.

## CLI 도구

Hilscher netX SERCOS III CLI 유틸리티를 사용합니다.

### 설치

```bash
# Hilscher netX 드라이버 및 도구 설치
# netX 드라이버에 대해서는 https://www.hilscher.com/ 참조
sudo apt-get install netx-driver
```

확인: `netx_cli --help`

## 사용법

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/sercos"
)

func main() {
    p := sercos.New()
    codec, _ := p.NewCodec("cmd")

    req := &kernel.Request{
        Function: "read",
        Address:  "S-0-51",
        Count:    2,
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## 드라이버

```go
b, c, err := sercos.NewCmdDriver("/usr/bin/netx_cli")
```

## 지원 함수

| 함수 | 설명                 |
|----------|-----------------------------|
| `read`   | SERCOS IDN/S 파라미터 읽기 |
| `write`  | SERCOS IDN/S 파라미터 쓰기|
| `phase`  | 통신 단계 설정     |

## 테스트

```bash
go test ./... -v
```
