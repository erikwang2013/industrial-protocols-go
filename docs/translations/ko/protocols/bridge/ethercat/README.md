# EtherCAT CmdBridge SDK

EtherCAT(Ethernet for Control Automation Technology)은 고성능 산업용 이더넷 필드버스입니다. 이 패키지는 `ethercat` 커맨드라인 유틸리티를 통한 EtherCAT 코덱을 제공합니다.

## CLI 도구

IgH EtherCAT Master 커맨드라인 도구를 사용합니다.

### 설치

```bash
# Debian/Ubuntu
sudo apt-get install ethercat-master

# 소스에서 빌드
git clone https://gitlab.com/etherlab.org/ethercat.git
cd ethercat
make && sudo make install
```

확인: `ethercat slaves`

## 사용법

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/ethercat"
)

func main() {
    p := ethercat.New()
    codec, _ := p.NewCodec("cmd")

    // 주소 0x1000에서 SDO 업로드
    req := &kernel.Request{
        Function: "upload",
        Address:  "0x1000",
        Count:    4,
    }
    raw, _ := codec.Encode(req)
    // raw = "upload 0x1000 4\n"
    _ = raw
}
```

## 드라이버

```go
b, c, err := ethercat.NewCmdDriver("/usr/bin/ethercat")
```

## 지원 함수

| 함수   | 설명            |
|------------|------------------------|
| `upload`   | 주소에서 SDO 읽기  |
| `download` | 주소에 SDO 쓰기   |
| `slaves`   | EtherCAT 슬레이브 목록   |

## 테스트

```bash
go test ./... -v
```
