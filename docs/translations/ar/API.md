# مرجع API — Industrial Protocols Go

> نظرة عامة على المشروع في [README](README.md) · النسخة الإنجليزية: [API.en.md](../../API.en.md)

---

## الواجهات الأساسية

```go
// kernel/protocol.go — يجب أن تنفذها كل وحدة بروتوكول
type Protocol interface {
    Name()        string              // "modbus"
    Variants()    []string            // ["tcp", "rtu", "ascii"]
    DefaultPort() int                 // 502
    NewCodec(variant string) (Codec, error)
}

// kernel/codec.go — الترميز وفك الترميز
type Codec interface {
    Encode(req *Request) ([]byte, error)
    Decode(data []byte) (*Response, error)
}

// kernel/transport/transport.go — قناة النقل
type Transport interface {
    io.ReadWriter
    io.Closer
    Addr()  string
    Alive() bool
}

// kernel/pipeline/pipeline.go — الوسيطات
type Handler    func(ctx context.Context, req *Request) (*Response, error)
type Middleware func(next Handler) Handler
```

## سلسلة الوسيطات

```
Request → [Timeout] → [Retry] → [CircuitBreaker] → [Logger] → Codec.Encode → Transport.Write
Response ← Codec.Decode ← Transport.Read
```

الوسيطات المدمجة:

| الوسيطة | الإنشاء | الوصف |
|--------|------|------|
| `Timeout` | `Timeout(d time.Duration)` | التحكم في المهلة عبر context |
| `Retry` | `Retry(max int, backoff BackoffFunc)` | ثلاث أنواع من التباطؤ: linear / exponential / jitter |
| `CircuitBreaker` | `NewCircuitBreaker(threshold int, cooldown time.Duration)` | قاطع دائرة بثلاث حالات closed → open → half-open |
| `Logger` | `Logger(logger *log.Logger)` | تسجيل طلبات منظم |

## دليل الاستخدام

### التثبيت

مكتبة Kernel الأساسية هي تبعية إلزامية، ووحدات البروتوكول تُستورد عند الحاجة. جميع الوحدات بنسخة موحدة `v1.1.2`.

```bash
# المكتبة الأساسية (إلزامية)
go get github.com/erikwang2013/industrial-protocols-go/kernel@v1.1.2
```

**الإيثرنت الصناعي (5):**

```bash
go get github.com/erikwang2013/industrial-protocols-go/protocols/ethernet/modbus@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/ethernet/bacnet@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/ethernet/ethernetip@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/ethernet/opcua@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/ethernet/profinet@v1.1.2
```

**ناقل الميدان (11):**

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

**IoT / الرسائل (2):**

```bash
go get github.com/erikwang2013/industrial-protocols-go/protocols/iot/mqtt@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/iot/hartip@v1.1.2
```

**ناقل السيارات (5):**

```bash
go get github.com/erikwang2013/industrial-protocols-go/protocols/automotive/lin@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/automotive/kline@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/automotive/flexray@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/automotive/saej1850@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/automotive/most@v1.1.2
```

**المباني / الإضاءة (2):**

```bash
go get github.com/erikwang2013/industrial-protocols-go/protocols/building/dali@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/building/lonworks@v1.1.2
```

**الجسور الإلكترونية (12):**

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

**ناقل النظام (3):**

```bash
go get github.com/erikwang2013/industrial-protocols-go/protocols/system/pci@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/system/vme@v1.1.2
go get github.com/erikwang2013/industrial-protocols-go/protocols/system/cpci@v1.1.2
```

> **عنوان pkg.go.dev:** `https://pkg.go.dev/github.com/erikwang2013/industrial-protocols-go/` + مسار الوحدة الفرعية (مثل `/kernel` أو `/protocols/ethernet/modbus`)

### قراءة/كتابة Modbus TCP

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

### استخدام ConnectionManager + تجمع الاتصالات

```go
mgr := connection.NewManager(&connection.PooledStrategy{})

// تسجيل جهاز (من الكود)
mgr.Register(&connection.DeviceConfig{
    Name: "plc-001", Network: "tcp", Addr: "192.168.1.10:502",
    Timeout: 3 * time.Second,
    Pool:    &connection.PoolConfig{MaxSize: 5},
})

// أو التحميل من إعداد YAML
repo, _ := config.Load("devices.yaml")
for _, cfg := range repo.All() {
    mgr.Register(cfg)
}

// الحصول على اتصال → استخدام → إرجاع
tr, _ := mgr.Acquire("plc-001")
// ... قراءة/كتابة البروتوكول ...
mgr.Release(tr)

// فحص الصحة
status := mgr.Health("plc-001")
```

### نشر واشتراك MQTT

```go
tr, _ := transport.DialTCP("broker.emqx.io:1883")
codec, _ := mqtt.New().NewCodec("tcp")

// إرسال CONNECT
connect, _ := codec.Encode(&kernel.Request{
    Function: "connect",
    Metadata: map[string]any{"client_id": "goclient"},
})
tr.Write(connect)
buf := make([]byte, 256)
n, _ := tr.Read(buf)
codec.Decode(buf[:n]) // CONNACK

// نشر رسالة
pub, _ := codec.Encode(&kernel.Request{
    Function: "publish",
    Address:  "sensor/temperature",
    Data:     []byte("25.5"),
})
tr.Write(pub)

// الاشتراك في موضوع
sub, _ := codec.Encode(&kernel.Request{
    Function: "subscribe",
    Address:  "sensor/#",
})
tr.Write(sub)
```

### استخدام إعدادات Vendor المسبقة

```go
reg := vendor.NewRegistry()
vendor.RegisterSiemens(reg)
vendor.RegisterRockwell(reg)

profile, _ := reg.Find("Siemens", "S7-1200")
// profile.Defaults["protocol"] = "modbus"
// profile.Defaults["port"]     = 502
// profile.Defaults["endian"]   = "big"
```

### تحويل البوابة بين البروتوكولات

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

### بروتوكول مخصص

يتطلب تنفيذ بروتوكول جديد ملفين فقط:

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

MIT — حقوق النشر (c) 2026 erik <erik@erik.xyz> — https://erik.xyz
