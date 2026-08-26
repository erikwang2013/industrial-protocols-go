# MOST 직렬 프로토콜 SDK

MOST(Media Oriented Systems Transport)는 주로 자동차 인포테인먼트 시스템에서 사용되는 고속 멀티미디어 네트워크 기술입니다. 이 패키지는 AT-command 인터페이스를 사용하는 직렬 어댑터를 통해 MOST 코덱을 제공합니다.

## 프로토콜 개요

MOST는 광섬유 물리 계층을 통한 동기식 직렬 통신을 사용합니다. 이 구현은 115200 보드레이트에서 AT-command 인터페이스를 노출하는 직렬 어댑터를 통해 연결됩니다.

### AT 명령

| 명령            | 설명                  |
|--------------------|------------------------------|
| `AT+READ=<addr>,<len>`  | 주소에서 바이트 읽기 |
| `AT+WRITE=<addr>,<hex>` | 주소에 16진수 바이트 쓰기 |
| `AT+STATUS`             | 링/네트워크 상태 조회 |

### 응답 형식

- `+OK:<hex_data>` -- 성공 응답
- `+ERR:<code>` -- 오류 응답

## 하드웨어 요구 사항

- 적절한 종단 처리가 된 MOST 광섬유 네트워크
- MOST-to-직렬 어댑터(예: MOST150 USB 어댑터)
- 115200 보드레이트, 8N1 직렬 포트

## 사용법

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/automotive/most"
)

func main() {
    p := most.New()
    codec, _ := p.NewCodec("serial")

    // 주소 0x0100에서 4바이트 읽기
    req := &kernel.Request{
        Function: "read",
        Address:  "0x0100",
        Count:    4,
    }
    raw, _ := codec.Encode(req)
    // raw = "AT+READ=0x0100,4\r\n"
    _ = raw

    // 주소 0x0200에 데이터 쓰기
    req2 := &kernel.Request{
        Function: "write",
        Address:  "0x0200",
        Data:     []byte{0xDE, 0xAD, 0xBE, 0xEF},
    }
    raw2, _ := codec.Encode(req2)
    // raw2 = "AT+WRITE=0x0200,DEADBEEF\r\n"
    _ = raw2
}
```

## 드라이버

```go
b, c, err := most.NewSerialDriver("/dev/ttyUSB0")
```

## 지원 함수

| 함수 | 설명                 |
|----------|-----------------------------|
| `read`   | MOST 주소에서 읽기    |
| `write`  | MOST 주소에 데이터 쓰기 |
| `status` | 링/네트워크 상태 조회   |

## 테스트

```bash
go test ./... -v
```
