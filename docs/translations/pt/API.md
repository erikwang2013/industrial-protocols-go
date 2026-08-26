# Referência da API — Industrial Protocols Go

> Visão geral do projeto em [README](README.md) · Versão em inglês: [API.en.md](../../API.en.md)

---

## Interfaces principais

```go
// kernel/protocol.go — 每个协议模块必须实现
type Protocol interface {
    Name()        string              // "modbus"
    Variants()    []string            // ["tcp", "rtu", "ascii"]
    DefaultPort() int                 // 502
    NewCodec(variant string) (Codec, error)
}

// kernel/codec.go — 编解码
type Codec interface {
    Encode(req *Request) ([]byte, error)
    Decode(data []byte) (*Response, error)
}

// kernel/transport/transport.go — 传输信道
type Transport interface {
    io.ReadWriter
    io.Closer
    Addr()  string
    Alive() bool
}

// kernel/pipeline/pipeline.go — 中间件
type Handler    func(ctx context.Context, req *Request) (*Response, error)
type Middleware func(next Handler) Handler
```

## Cadeia de Middlewares

```
Request → [Timeout] → [Retry] → [CircuitBreaker] → [Logger] → Codec.Encode → Transport.Write
Response ← Codec.Decode ← Transport.Read
```

Middlewares integrados:

| Middleware | Construtor | Descrição |
|--------|------|------|
| `Timeout` | `Timeout(d time.Duration)` | controle de timeout via context |
| `Retry` | `Retry(max int, backoff BackoffFunc)` | três tipos de backoff: linear / exponential / jitter |
| `CircuitBreaker` | `NewCircuitBreaker(threshold int, cooldown time.Duration)` | circuit breaker de três estados: closed → open → half-open |
| `Logger` | `Logger(logger *log.Logger)` | registro estruturado de requisições |

## Guia de Uso

### Instalação

A biblioteca principal do kernel é uma dependência obrigatória; os módulos de protocolo são importados conforme a necessidade. Todos os módulos usam a versão unificada `v1.1.2`.

```bash
# 核心库（必须）
go get github.com/erikwang2013/industrial-protocols-go/kernel@v1.1.2
```

**Ethernet industrial (5):**

```bash
go get github.com/erikwang2013/industrial-protocols-go/protocols/ethernet/modbus@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/ethernet/bacnet@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/ethernet/ethernetip@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/ethernet/opcua@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/ethernet/profinet@v1.1.2
```

**Fieldbus (11):**

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

**IoT / Mensagens (2):**

```bash
go get github.com/erikwang2013/industrial-protocols-go/protocols/iot/mqtt@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/iot/hartip@v1.1.2
```

**Barramento automotivo (5):**

```bash
go get github.com/erikwang2013/industrial-protocols-go/protocols/automotive/lin@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/automotive/kline@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/automotive/flexray@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/automotive/saej1850@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/automotive/most@v1.1.2
```

**Edifícios / Iluminação (2):**

```bash
go get github.com/erikwang2013/industrial-protocols-go/protocols/building/dali@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/building/lonworks@v1.1.2
```

**Pontes de hardware (12):**

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

**Barramento de sistema (3):**

```bash
go get github.com/erikwang2013/industrial-protocols-go/protocols/system/pci@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/system/vme@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/system/cpci@v1.1.2
```

> **Endereço pkg.go.dev:** `https://pkg.go.dev/github.com/erikwang2013/industrial-protocols-go/` + caminho do submódulo (ex.: `/kernel`, `/protocols/ethernet/modbus`)

### Leitura/escrita Modbus TCP

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

### Usando ConnectionManager + pool de conexões

```go
mgr := connection.NewManager(&connection.PooledStrategy{})

// 注册设备（从代码）
mgr.Register(&connection.DeviceConfig{
    Name: "plc-001", Network: "tcp", Addr: "192.168.1.10:502",
    Timeout: 3 * time.Second,
    Pool:    &connection.PoolConfig{MaxSize: 5},
})

// 或从 YAML 配置加载
repo, _ := config.Load("devices.yaml")
for _, cfg := range repo.All() {
    mgr.Register(cfg)
}

// 获取连接 → 使用 → 归还
tr, _ := mgr.Acquire("plc-001")
// ... 协议读写 ...
mgr.Release(tr)

// 健康检查
status := mgr.Health("plc-001")
```

### Publicação e assinatura MQTT

```go
tr, _ := transport.DialTCP("broker.emqx.io:1883")
codec, _ := mqtt.New().NewCodec("tcp")

// 发送 CONNECT
connect, _ := codec.Encode(&kernel.Request{
    Function: "connect",
    Metadata: map[string]any{"client_id": "goclient"},
})
tr.Write(connect)
buf := make([]byte, 256)
n, _ := tr.Read(buf)
codec.Decode(buf[:n]) // CONNACK

// 发布消息
pub, _ := codec.Encode(&kernel.Request{
    Function: "publish",
    Address:  "sensor/temperature",
    Data:     []byte("25.5"),
})
tr.Write(pub)

// 订阅主题
sub, _ := codec.Encode(&kernel.Request{
    Function: "subscribe",
    Address:  "sensor/#",
})
tr.Write(sub)
```

### Usando presets Vendor

```go
reg := vendor.NewRegistry()
vendor.RegisterSiemens(reg)
vendor.RegisterRockwell(reg)

profile, _ := reg.Find("Siemens", "S7-1200")
// profile.Defaults["protocol"] = "modbus"
// profile.Defaults["port"]     = 502
// profile.Defaults["endian"]   = "big"
```

### Conversão entre protocolos por gateway

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

### Protocolo personalizado

Implementar um novo protocolo requer apenas dois arquivos:

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
