# API 레퍼런스 — Industrial Protocols Go

> 프로젝트 개요는 [README](README.md) · 영문 버전: [API.en.md](../../API.en.md)

---

## 핵심 인터페이스

```go
// kernel/protocol.go — 모든 프로토콜 모듈이 구현해야 함
type Protocol interface {
    Name()        string              // "modbus"
    Variants()    []string            // ["tcp", "rtu", "ascii"]
    DefaultPort() int                 // 502
    NewCodec(variant string) (Codec, error)
}

// kernel/codec.go — 인코딩/디코딩
type Codec interface {
    Encode(req *Request) ([]byte, error)
    Decode(data []byte) (*Response, error)
}

// kernel/transport/transport.go — 전송 채널
type Transport interface {
    io.ReadWriter
    io.Closer
    Addr()  string
    Alive() bool
}

// kernel/pipeline/pipeline.go — 미들웨어
type Handler    func(ctx context.Context, req *Request) (*Response, error)
type Middleware func(next Handler) Handler
```

## 미들웨어 체인

```
Request → [Timeout] → [Retry] → [CircuitBreaker] → [Logger] → Codec.Encode → Transport.Write
Response ← Codec.Decode ← Transport.Read
```

내장 미들웨어:

| 미들웨어 | 생성자 | 설명 |
|--------|------|------|
| `Timeout` | `Timeout(d time.Duration)` | context 타임아웃 제어 |
| `Retry` | `Retry(max int, backoff BackoffFunc)` | linear / exponential / jitter 3가지 백오프 |
| `CircuitBreaker` | `NewCircuitBreaker(threshold int, cooldown time.Duration)` | closed → open → half-open 3상태 차단기 |
| `Logger` | `Logger(logger *log.Logger)` | 구조화된 요청 로그 |

## 사용 가이드

### 설치

Kernel 핵심 라이브러리는 필수 의존성이며, 프로토콜 모듈은 필요에 따라 import 합니다. 모든 모듈은 통일된 버전 `v1.1.2`를 사용합니다.

```bash
# 핵심 라이브러리(필수)
go get github.com/erikwang2013/industrial-protocols-go/kernel@v1.1.2
```

**산업용 이더넷(5):**

```bash
go get github.com/erikwang2013/industrial-protocols-go/protocols/ethernet/modbus@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/ethernet/bacnet@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/ethernet/ethernetip@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/ethernet/opcua@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/ethernet/profinet@v1.1.2
```

**필드버스(11):**

```bash
go get github.com/erikwang2013/industrial-protocols-go/protocols/fieldbus/hart@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/fieldbus/cclink@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/fieldbus/dnp3@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/fieldbus/iec61850@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/fieldbus/profibus@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/fieldbus/canopen@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/fieldbus/devicenet@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/fieldbus/foundationfieldbus@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/fieldbus/asinterface@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/fieldbus/iolink@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/fieldbus/cclinkie@v1.1.2
```

**IoT / 메시지(2):**

```bash
go get github.com/erikwang2013/industrial-protocols-go/protocols/iot/mqtt@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/iot/hartip@v1.1.2
```

**자동차 버스(5):**

```bash
go get github.com/erikwang2013/industrial-protocols-go/protocols/automotive/lin@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/automotive/kline@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/automotive/flexray@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/automotive/saej1850@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/automotive/most@v1.1.2
```

**빌딩 / 조명(2):**

```bash
go get github.com/erikwang2013/industrial-protocols-go/protocols/building/dali@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/building/lonworks@v1.1.2
```

**하드웨어 브리지(12):**

```bash
go get github.com/erikwang2013/industrial-protocols-go/protocols/bridge/ethercat@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/bridge/powerlink@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/bridge/sercos@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/bridge/sercos1@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/bridge/controlnet@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/bridge/interbus@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/bridge/worldfip@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/bridge/lightbus@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/bridge/modbusplus@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/bridge/isa100@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/bridge/wirelesshart@v1.1.2
```

**시스템 버스(3):**

```bash
go get github.com/erikwang2013/industrial-protocols-go/protocols/system/pci@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/system/vme@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/system/cpci@v1.1.2
```

> **pkg.go.dev 주소:** `https://pkg.go.dev/github.com/erikwang2013/industrial-protocols-go/` + 서브모듈 경로(예: `/kernel`, `/protocols/ethernet/modbus`)

### Modbus TCP 읽기/쓰기

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/kernel/pipeline"
    "github.com/erikwang2013/industrial-protocols-go/kernel/transport"
    "github.com/erikwang2013/industrial-protocols-go/protocols/ethernet/modbus"
)

func main() {
    tr, _ := transport.DialTCP("192.168.1.10:502")
    defer tr.Close()

    codec, _ := modbus.NewCodec("tcp")
    handler := func(ctx context.Context, req *kernel.Request) (*kernel.Response, error) {
        req.Metadata["unit_id"] = byte(1)
        raw, _ := codec.Encode(req)
        tr.Write(raw)
        buf := make([]byte, 256)
        n, _ := tr.Read(buf)
        return codec.Decode(buf[:n])
    }

    wrapped := pipeline.Chain(
        pipeline.Timeout(3*time.Second),
        pipeline.Retry(3, pipeline.ExponentialBackoff(100*time.Millisecond)),
    )(handler)

    resp, err := wrapped(context.Background(), &kernel.Request{
        Function: "read_holding_registers",
        Address:  "40001",
        Count:    2,
    })
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Registers: %x", resp.Data)
}
```

### ConnectionManager + 연결 풀 사용

```go
mgr := connection.NewManager(&connection.PooledStrategy{})

// 장치 등록(코드에서)
mgr.Register(&connection.DeviceConfig{
    Name: "plc-001", Network: "tcp", Addr: "192.168.1.10:502",
    Timeout: 3 * time.Second,
    Pool:    &connection.PoolConfig{MaxSize: 5},
})

// 또는 YAML 설정에서 로드
repo, _ := config.Load("devices.yaml")
for _, cfg := range repo.All() {
    mgr.Register(cfg)
}

// 연결 획득 → 사용 → 반환
tr, _ := mgr.Acquire("plc-001")
// ... 프로토콜 읽기/쓰기 ...
mgr.Release(tr)

// 헬스 체크
status := mgr.Health("plc-001")
```

### MQTT 발행/구독

```go
tr, _ := transport.DialTCP("broker.emqx.io:1883")
codec, _ := mqtt.New().NewCodec("tcp")

// CONNECT 전송
connect, _ := codec.Encode(&kernel.Request{
    Function: "connect",
    Metadata: map[string]any{"client_id": "goclient"},
})
tr.Write(connect)
buf := make([]byte, 256)
n, _ := tr.Read(buf)
codec.Decode(buf[:n]) // CONNACK

// 메시지 발행
pub, _ := codec.Encode(&kernel.Request{
    Function: "publish",
    Address:  "sensor/temperature",
    Data:     []byte("25.5"),
})
tr.Write(pub)

// 토픽 구독
sub, _ := codec.Encode(&kernel.Request{
    Function: "subscribe",
    Address:  "sensor/#",
})
tr.Write(sub)
```

### Vendor 프리셋 사용

```go
reg := vendor.NewRegistry()
vendor.RegisterSiemens(reg)
vendor.RegisterRockwell(reg)

profile, _ := reg.Find("Siemens", "S7-1200")
// profile.Defaults["protocol"] = "modbus"
// profile.Defaults["port"]     = 502
// profile.Defaults["endian"]   = "big"
```

### 프로토콜 간 게이트웨이 변환

```go
engine := gateway.New()
engine.Add(gateway.Rule{
    Src: modbus.NewProtocol(),
    Dst: mqtt.New(),
    Map: func(srcResp *kernel.Response) *kernel.Request {
        return &kernel.Request{
            Function: "publish",
            Address:  "plc/registers",
            Data:     srcResp.Data,
        }
    },
})

srcTr, _ := transport.DialTCP("192.168.1.10:502")
dstTr, _ := transport.DialTCP("broker.emqx.io:1883")
engine.Transform(ctx, srcTr, dstTr, &kernel.Request{
    Function: "read_holding_registers", Address: "40001", Count: 2,
})
```

### 커스텀 프로토콜

새 프로토콜을 구현하려면 두 파일만 있으면 됩니다:

```go
// myproto/myproto.go
type MyProtocol struct{}

func (p *MyProtocol) Name() string              { return "myproto" }
func (p *MyProtocol) Variants() []string        { return []string{"tcp"} }
func (p *MyProtocol) DefaultPort() int          { return 9000 }
func (p *MyProtocol) NewCodec(v string) (kernel.Codec, error) {
    return &myCodec{}, nil
}

type myCodec struct{}

func (c *myCodec) Encode(req *kernel.Request) ([]byte, error) {
    return []byte(req.Function + ":" + req.Address), nil
}

func (c *myCodec) Decode(data []byte) (*kernel.Response, error) {
    return &kernel.Response{Data: data}, nil
}
```

---

MIT — Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz
