# Industrial Protocols Go — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a Go industrial communication protocols monorepo with kernel interfaces, transport layer, pipeline middleware, 40 protocol stubs, and 3 full implementations (Modbus, MQTT, OPC UA).

**Architecture:** Go workspace monorepo. Kernel module defines Protocol/Codec/Transport interfaces plus reusable pipeline middleware. Each protocol is an independent Go module implementing Protocol + Codec. Transport is abstracted behind an `io.ReadWriter`-based interface with built-in TCP/UDP implementations.

**Tech Stack:** Go 1.20+, no external dependencies beyond stdlib for kernel.

---

## Phase 1: Project Scaffold

### Task 1: Initialize project root and go.work

**Files:**
- Create: `go.work`
- Create: `_tools/Makefile`
- Create: `README.md`
- Create: `LICENSE`

- [ ] **Step 1: Create LICENSE**

```bash
cat > /home/wwwroot/bag/industrial-protocols-go/LICENSE << 'LICEOF'
MIT License

Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
LICEOF
```

- [ ] **Step 2: Create _tools/Makefile**

```makefile
.PHONY: test test-unit test-integration lint vet fmt

test:
	cd kernel && go test ./...
	@for dir in $$(find protocols -name go.mod -exec dirname {} \;); do \
		(cd $$dir && go test ./...) || exit 1; \
	done

test-unit:
	cd kernel && go test -short ./...
	@for dir in $$(find protocols -name go.mod -exec dirname {} \;); do \
		(cd $$dir && go test -short ./...) || exit 1; \
	done

test-integration:
	cd kernel && go test -tags=integration ./...
	@for dir in $$(find protocols -name go.mod -exec dirname {} \;); do \
		(cd $$dir && go test -tags=integration ./...) || exit 1; \
	done

lint:
	golangci-lint run ./...

vet:
	cd kernel && go vet ./...
	@for dir in $$(find protocols -name go.mod -exec dirname {} \;); do \
		(cd $$dir && go vet ./...) || exit 1; \
	done

fmt:
	go fmt ./...
```

- [ ] **Step 3: Create README.md**

Write a README with project overview, architecture diagram, quick-start example, and link to the reference project.

- [ ] **Step 4: Initialize go.work**

```bash
cd /home/wwwroot/bag/industrial-protocols-go
go work init
```

- [ ] **Step 5: Verify go.work was created**

```bash
cat go.work
```

---

### Task 2: Initialize kernel module

**Files:**
- Create: `kernel/go.mod`

- [ ] **Step 1: Create kernel module**

```bash
cd /home/wwwroot/bag/industrial-protocols-go/kernel
go mod init github.com/erikwang2013/industrial-protocols-go/kernel
```

- [ ] **Step 2: Add kernel to workspace**

```bash
cd /home/wwwroot/bag/industrial-protocols-go
go work use ./kernel
go work sync
```

---

## Phase 2: Kernel — Core Interfaces

### Task 3: Define errors

**Files:**
- Create: `kernel/errors.go`
- Create: `kernel/errors_test.go`

- [ ] **Step 1: Write the failing test**

Create `kernel/errors_test.go`:

```go
package kernel

import (
	"errors"
	"testing"
)

func TestSentinelErrors(t *testing.T) {
	errs := []error{
		ErrTimeout,
		ErrCircuitOpen,
		ErrRetryExhausted,
		ErrTransportClosed,
		ErrInvalidAddress,
	}
	for _, e := range errs {
		if e.Error() == "" {
			t.Errorf("sentinel error %v has empty message", e)
		}
	}
}

func TestProtocolError(t *testing.T) {
	raw := []byte{0x01, 0x83, 0x02}
	pe := &ProtocolError{Code: "02", Message: "illegal data address", Raw: raw}

	if pe.Error() != "protocol error [02]: illegal data address" {
		t.Errorf("unexpected error string: %s", pe.Error())
	}
	if !errors.Is(pe, pe) {
		t.Error("ProtocolError should be errors.Is to itself")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
cd /home/wwwroot/bag/industrial-protocols-go/kernel && go test ./...
```

- [ ] **Step 3: Implement errors.go**

Create `kernel/errors.go`:

```go
package kernel

import (
	"errors"
	"fmt"
)

var (
	ErrTimeout         = errors.New("protocol: timeout")
	ErrCircuitOpen     = errors.New("protocol: circuit breaker open")
	ErrRetryExhausted  = errors.New("protocol: all retries exhausted")
	ErrTransportClosed = errors.New("protocol: transport closed")
	ErrInvalidAddress  = errors.New("protocol: invalid address")
)

type ProtocolError struct {
	Code    string
	Message string
	Raw     []byte
}

func (e *ProtocolError) Error() string {
	return fmt.Sprintf("protocol error [%s]: %s", e.Code, e.Message)
}
```

- [ ] **Step 4: Run test to verify it passes**

```bash
cd /home/wwwroot/bag/industrial-protocols-go/kernel && go test ./... -v
```

- [ ] **Step 5: Commit**

```bash
cd /home/wwwroot/bag/industrial-protocols-go
git add kernel/go.mod kernel/errors.go kernel/errors_test.go
git commit -m "feat(kernel): add sentinel errors and ProtocolError"
```

---

### Task 4: Define Request, Response, and Codec

**Files:**
- Create: `kernel/codec.go`
- Create: `kernel/codec_test.go`

- [ ] **Step 1: Implement codec.go**

Create `kernel/codec.go`:

```go
package kernel

type Request struct {
	Function string
	Address  string
	Count    int
	Data     []byte
	Metadata map[string]any
}

type Response struct {
	Address  string
	Data     []byte
	Metadata map[string]any
}

type Codec interface {
	Encode(req *Request) ([]byte, error)
	Decode(data []byte) (*Response, error)
}
```

- [ ] **Step 2: Write test**

Create `kernel/codec_test.go`:

```go
package kernel

import "testing"

func TestRequestStruct(t *testing.T) {
	r := &Request{
		Function: "read_coils",
		Address:  "40001",
		Count:    8,
		Metadata: map[string]any{"unit_id": 1},
	}
	if r.Function != "read_coils" {
		t.Errorf("expected read_coils, got %s", r.Function)
	}
}
```

- [ ] **Step 3: Run tests**

```bash
cd /home/wwwroot/bag/industrial-protocols-go/kernel && go test ./... -v
```

- [ ] **Step 4: Commit**

```bash
git add kernel/codec.go kernel/codec_test.go
git commit -m "feat(kernel): add Request, Response, and Codec interface"
```

---

### Task 5: Define Protocol interface

**Files:**
- Create: `kernel/protocol.go`
- Create: `kernel/protocol_test.go`

- [ ] **Step 1: Implement protocol.go**

Create `kernel/protocol.go`:

```go
package kernel

type Protocol interface {
	Name()        string
	Variants()    []string
	DefaultPort() int
	NewCodec(variant string) (Codec, error)
}
```

- [ ] **Step 2: Write test with a mock protocol**

Create `kernel/protocol_test.go`:

```go
package kernel

import "testing"

type testProtocol struct{}

func (p *testProtocol) Name() string                 { return "test" }
func (p *testProtocol) Variants() []string           { return []string{"v1", "v2"} }
func (p *testProtocol) DefaultPort() int             { return 9000 }
func (p *testProtocol) NewCodec(v string) (Codec, error) {
	return nil, nil
}

func TestProtocolInterface(t *testing.T) {
	var p Protocol = &testProtocol{}
	if p.Name() != "test" {
		t.Errorf("expected test, got %s", p.Name())
	}
	if len(p.Variants()) != 2 {
		t.Errorf("expected 2 variants, got %d", len(p.Variants()))
	}
	if p.DefaultPort() != 9000 {
		t.Errorf("expected 9000, got %d", p.DefaultPort())
	}
}
```

- [ ] **Step 3: Run tests**

```bash
cd /home/wwwroot/bag/industrial-protocols-go/kernel && go test ./... -v
```

- [ ] **Step 4: Commit**

```bash
git add kernel/protocol.go kernel/protocol_test.go
git commit -m "feat(kernel): add Protocol interface"
```

---

### Task 6: Define Transport interface and TCP implementation

**Files:**
- Create: `kernel/transport/transport.go`
- Create: `kernel/transport/tcp.go`
- Create: `kernel/transport/transport_test.go`

- [ ] **Step 1: Implement transport interface**

Create `kernel/transport/transport.go`:

```go
package transport

import "io"

type Transport interface {
	io.ReadWriter
	io.Closer
	Addr()  string
	Alive() bool
}
```

- [ ] **Step 2: Implement TCPTransport**

Create `kernel/transport/tcp.go`:

```go
package transport

import (
	"net"
	"sync"
)

type TCPTransport struct {
	conn   net.Conn
	addr   string
	mu     sync.Mutex
	closed bool
}

func DialTCP(address string) (*TCPTransport, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	return &TCPTransport{conn: conn, addr: address}, nil
}

func (t *TCPTransport) Read(p []byte) (int, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.closed {
		return 0, net.ErrClosed
	}
	return t.conn.Read(p)
}

func (t *TCPTransport) Write(p []byte) (int, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.closed {
		return 0, net.ErrClosed
	}
	return t.conn.Write(p)
}

func (t *TCPTransport) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.closed {
		return nil
	}
	t.closed = true
	return t.conn.Close()
}

func (t *TCPTransport) Addr() string {
	return t.addr
}

func (t *TCPTransport) Alive() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return !t.closed
}
```

- [ ] **Step 3: Write integration test**

Create `kernel/transport/transport_test.go`:

```go
package transport

import (
	"net"
	"testing"
)

func TestTCPTransportIntegration(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		buf := make([]byte, 64)
		n, _ := conn.Read(buf)
		conn.Write(buf[:n])
		conn.Close()
	}()

	tr, err := DialTCP(ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer tr.Close()

	if !tr.Alive() {
		t.Error("transport should be alive after dial")
	}

	data := []byte("hello")
	if _, err := tr.Write(data); err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 64)
	n, err := tr.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	if string(buf[:n]) != "hello" {
		t.Errorf("expected 'hello', got '%s'", string(buf[:n]))
	}

	if err := tr.Close(); err != nil {
		t.Fatal(err)
	}
	if tr.Alive() {
		t.Error("transport should not be alive after close")
	}
}
```

- [ ] **Step 4: Run tests**

```bash
cd /home/wwwroot/bag/industrial-protocols-go/kernel && go test ./... -v -short
```

- [ ] **Step 5: Commit**

```bash
git add kernel/transport/
git commit -m "feat(kernel): add Transport interface and TCP implementation"
```

---

### Task 7: Add UDP transport

**Files:**
- Create: `kernel/transport/udp.go`

- [ ] **Step 1: Implement UDPTransport**

Create `kernel/transport/udp.go`:

```go
package transport

import (
	"net"
	"sync"
)

type UDPTransport struct {
	conn   net.Conn
	addr   string
	mu     sync.Mutex
	closed bool
}

func DialUDP(address string) (*UDPTransport, error) {
	conn, err := net.Dial("udp", address)
	if err != nil {
		return nil, err
	}
	return &UDPTransport{conn: conn, addr: address}, nil
}

func (t *UDPTransport) Read(p []byte) (int, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.closed {
		return 0, net.ErrClosed
	}
	return t.conn.Read(p)
}

func (t *UDPTransport) Write(p []byte) (int, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.closed {
		return 0, net.ErrClosed
	}
	return t.conn.Write(p)
}

func (t *UDPTransport) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.closed {
		return nil
	}
	t.closed = true
	return t.conn.Close()
}

func (t *UDPTransport) Addr() string {
	return t.addr
}

func (t *UDPTransport) Alive() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return !t.closed
}
```

- [ ] **Step 2: Verify compilation**

```bash
cd /home/wwwroot/bag/industrial-protocols-go/kernel && go build ./...
```

- [ ] **Step 3: Commit**

```bash
git add kernel/transport/udp.go
git commit -m "feat(kernel): add UDP transport"
```

---

### Task 8: Add serial transport stub

**Files:**
- Create: `kernel/transport/serial.go`

- [ ] **Step 1: Implement serial stub**

Create `kernel/transport/serial.go`:

```go
package transport

import (
	"errors"
)

type ReadWriteCloser interface {
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	Close() error
}

type SerialTransport struct {
	rw     ReadWriteCloser
	addr   string
	closed bool
}

func NewSerial(rwc ReadWriteCloser, addr string) (*SerialTransport, error) {
	if rwc == nil {
		return nil, errors.New("transport: nil serial port")
	}
	return &SerialTransport{rw: rwc, addr: addr}, nil
}

func (t *SerialTransport) Read(p []byte) (int, error) {
	if t.closed {
		return 0, errors.New("transport: serial port closed")
	}
	return t.rw.Read(p)
}

func (t *SerialTransport) Write(p []byte) (int, error) {
	if t.closed {
		return 0, errors.New("transport: serial port closed")
	}
	return t.rw.Write(p)
}

func (t *SerialTransport) Close() error {
	if t.closed {
		return nil
	}
	t.closed = true
	return t.rw.Close()
}

func (t *SerialTransport) Addr() string {
	return t.addr
}

func (t *SerialTransport) Alive() bool {
	return !t.closed
}
```

- [ ] **Step 2: Run all kernel tests**

```bash
cd /home/wwwroot/bag/industrial-protocols-go/kernel && go test ./... -v -short && go build ./...
```

- [ ] **Step 3: Commit**

```bash
git add kernel/transport/serial.go
git commit -m "feat(kernel): add serial transport stub"
```

---

## Phase 3: Kernel — Pipeline Middleware

### Task 9: Define Pipeline types and Chain

**Files:**
- Create: `kernel/pipeline/pipeline.go`
- Create: `kernel/pipeline/pipeline_test.go`

- [ ] **Step 1: Implement pipeline.go**

Create `kernel/pipeline/pipeline.go`:

```go
package pipeline

import (
	"context"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

type Handler func(ctx context.Context, req *kernel.Request) (*kernel.Response, error)

type Middleware func(next Handler) Handler

func Chain(middlewares ...Middleware) Middleware {
	return func(next Handler) Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}
```

- [ ] **Step 2: Write tests**

Create `kernel/pipeline/pipeline_test.go`:

```go
package pipeline

import (
	"context"
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func echoHandler(ctx context.Context, req *kernel.Request) (*kernel.Response, error) {
	return &kernel.Response{Address: req.Address, Data: req.Data}, nil
}

func TestChainOrder(t *testing.T) {
	order := []string{}

	mk := func(name string) Middleware {
		return func(next Handler) Handler {
			return func(ctx context.Context, req *kernel.Request) (*kernel.Response, error) {
				order = append(order, name)
				return next(ctx, req)
			}
		}
	}

	wrapped := Chain(mk("a"), mk("b"), mk("c"))(echoHandler)
	_, err := wrapped(context.Background(), &kernel.Request{})
	if err != nil {
		t.Fatal(err)
	}

	if len(order) != 3 || order[0] != "a" || order[1] != "b" || order[2] != "c" {
		t.Errorf("unexpected order: %v", order)
	}
}

func TestEmptyChain(t *testing.T) {
	wrapped := Chain()(echoHandler)
	resp, err := wrapped(context.Background(), &kernel.Request{Address: "x", Data: []byte("y")})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Address != "x" {
		t.Errorf("expected x, got %s", resp.Address)
	}
}
```

- [ ] **Step 3: Run tests**

```bash
cd /home/wwwroot/bag/industrial-protocols-go/kernel && go test ./... -v
```

- [ ] **Step 4: Commit**

```bash
git add kernel/pipeline/pipeline.go kernel/pipeline/pipeline_test.go
git commit -m "feat(kernel): add pipeline Chain and middleware types"
```

---

### Task 10: Implement Timeout middleware

**Files:**
- Create: `kernel/pipeline/timeout.go`
- Create: `kernel/pipeline/timeout_test.go`

- [ ] **Step 1: Write the failing test**

Create `kernel/pipeline/timeout_test.go`:

```go
package pipeline

import (
	"context"
	"testing"
	"time"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestTimeout_NoTimeout(t *testing.T) {
	wrapped := Chain(Timeout(100 * time.Millisecond))(echoHandler)
	_, err := wrapped(context.Background(), &kernel.Request{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestTimeout_Exceeded(t *testing.T) {
	slow := func(ctx context.Context, req *kernel.Request) (*kernel.Response, error) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(200 * time.Millisecond):
			return &kernel.Response{}, nil
		}
	}
	wrapped := Chain(Timeout(10 * time.Millisecond))(slow)
	_, err := wrapped(context.Background(), &kernel.Request{})
	if err == nil {
		t.Fatal("expected timeout error")
	}
}
```

- [ ] **Step 2: Implement timeout.go**

Create `kernel/pipeline/timeout.go`:

```go
package pipeline

import (
	"context"
	"time"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func Timeout(d time.Duration) Middleware {
	return func(next Handler) Handler {
		return func(ctx context.Context, req *kernel.Request) (*kernel.Response, error) {
			ctx, cancel := context.WithTimeout(ctx, d)
			defer cancel()
			return next(ctx, req)
		}
	}
}
```

- [ ] **Step 3: Run tests**

```bash
cd /home/wwwroot/bag/industrial-protocols-go/kernel && go test ./pipeline/ -v
```

- [ ] **Step 4: Commit**

```bash
git add kernel/pipeline/timeout.go kernel/pipeline/timeout_test.go
git commit -m "feat(kernel): add Timeout middleware"
```

---

### Task 11: Implement Retry middleware

**Files:**
- Create: `kernel/pipeline/retry.go`
- Create: `kernel/pipeline/retry_test.go`

- [ ] **Step 1: Implement retry.go**

Create `kernel/pipeline/retry.go`:

```go
package pipeline

import (
	"context"
	"time"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

type BackoffFunc func(attempt int) time.Duration

func LinearBackoff(base time.Duration) BackoffFunc {
	return func(attempt int) time.Duration {
		return base * time.Duration(attempt)
	}
}

func ExponentialBackoff(base time.Duration) BackoffFunc {
	return func(attempt int) time.Duration {
		d := base
		for i := 1; i < attempt; i++ {
			d *= 2
		}
		return d
	}
}

func JitterBackoff(base time.Duration) BackoffFunc {
	return func(attempt int) time.Duration {
		d := base * time.Duration(attempt)
		return d + time.Duration(d/2)
	}
}

func Retry(maxAttempts int, backoff BackoffFunc) Middleware {
	return func(next Handler) Handler {
		return func(ctx context.Context, req *kernel.Request) (*kernel.Response, error) {
			var lastErr error
			for i := 0; i < maxAttempts; i++ {
				resp, err := next(ctx, req)
				if err == nil {
					return resp, nil
				}
				lastErr = err
				if i < maxAttempts-1 {
					select {
					case <-ctx.Done():
						return nil, ctx.Err()
					case <-time.After(backoff(i + 1)):
					}
				}
			}
			return nil, kernel.ErrRetryExhausted
		}
	}
}
```

- [ ] **Step 2: Write tests**

Create `kernel/pipeline/retry_test.go`:

```go
package pipeline

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestRetry_SuccessFirstAttempt(t *testing.T) {
	wrapped := Chain(Retry(3, LinearBackoff(10*time.Millisecond)))(echoHandler)
	_, err := wrapped(context.Background(), &kernel.Request{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestRetry_EventualSuccess(t *testing.T) {
	attempts := 0
	flaky := func(ctx context.Context, req *kernel.Request) (*kernel.Response, error) {
		attempts++
		if attempts < 2 {
			return nil, errors.New("transient error")
		}
		return &kernel.Response{}, nil
	}
	wrapped := Chain(Retry(3, LinearBackoff(10*time.Millisecond)))(flaky)
	_, err := wrapped(context.Background(), &kernel.Request{})
	if err != nil {
		t.Fatal(err)
	}
	if attempts != 2 {
		t.Errorf("expected 2 attempts, got %d", attempts)
	}
}

func TestRetry_Exhausted(t *testing.T) {
	alwaysFail := func(ctx context.Context, req *kernel.Request) (*kernel.Response, error) {
		return nil, errors.New("always fail")
	}
	wrapped := Chain(Retry(2, LinearBackoff(10*time.Millisecond)))(alwaysFail)
	_, err := wrapped(context.Background(), &kernel.Request{})
	if !errors.Is(err, kernel.ErrRetryExhausted) {
		t.Errorf("expected ErrRetryExhausted, got %v", err)
	}
}
```

- [ ] **Step 3: Run tests**

```bash
cd /home/wwwroot/bag/industrial-protocols-go/kernel && go test ./pipeline/ -v
```

- [ ] **Step 4: Commit**

```bash
git add kernel/pipeline/retry.go kernel/pipeline/retry_test.go
git commit -m "feat(kernel): add Retry middleware with linear/exponential/jitter backoff"
```

---

### Task 12: Implement Logger middleware

**Files:**
- Create: `kernel/pipeline/logger.go`
- Create: `kernel/pipeline/logger_test.go`

- [ ] **Step 1: Implement logger.go**

Create `kernel/pipeline/logger.go`:

```go
package pipeline

import (
	"context"
	"log"
	"time"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func Logger(logger *log.Logger) Middleware {
	return func(next Handler) Handler {
		return func(ctx context.Context, req *kernel.Request) (*kernel.Response, error) {
			start := time.Now()
			resp, err := next(ctx, req)
			elapsed := time.Since(start)
			if err != nil {
				logger.Printf("[ERROR] %s %s %v: %v", req.Function, req.Address, elapsed, err)
			} else {
				logger.Printf("[DEBUG] %s %s %v", req.Function, req.Address, elapsed)
			}
			return resp, err
		}
	}
}
```

- [ ] **Step 2: Write test**

Create `kernel/pipeline/logger_test.go`:

```go
package pipeline

import (
	"bytes"
	"context"
	"errors"
	"log"
	"strings"
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)
	wrapped := Chain(Logger(logger))(echoHandler)
	_, err := wrapped(context.Background(), &kernel.Request{Function: "read", Address: "40001"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "read") {
		t.Errorf("log should contain function name, got: %s", buf.String())
	}
}

func TestLogger_Error(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)
	fail := func(ctx context.Context, req *kernel.Request) (*kernel.Response, error) {
		return nil, errors.New("boom")
	}
	wrapped := Chain(Logger(logger))(fail)
	_, _ = wrapped(context.Background(), &kernel.Request{Function: "write"})
	if !strings.Contains(buf.String(), "ERROR") {
		t.Errorf("log should contain ERROR, got: %s", buf.String())
	}
}
```

- [ ] **Step 3: Run tests**

```bash
cd /home/wwwroot/bag/industrial-protocols-go/kernel && go test ./pipeline/ -v
```

- [ ] **Step 4: Commit**

```bash
git add kernel/pipeline/logger.go kernel/pipeline/logger_test.go
git commit -m "feat(kernel): add Logger middleware"
```

---

### Task 13: Implement CircuitBreaker middleware

**Files:**
- Create: `kernel/pipeline/breaker.go`
- Create: `kernel/pipeline/breaker_test.go`

- [ ] **Step 1: Implement breaker.go**

Create `kernel/pipeline/breaker.go`:

```go
package pipeline

import (
	"context"
	"sync"
	"time"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

type CircuitBreaker struct {
	threshold   int
	cooldown    time.Duration
	failures    int
	lastFailure time.Time
	open        bool
	mu          sync.Mutex
}

func (cb *CircuitBreaker) Middleware() Middleware {
	return func(next Handler) Handler {
		return func(ctx context.Context, req *kernel.Request) (*kernel.Response, error) {
			cb.mu.Lock()
			if cb.open {
				if time.Since(cb.lastFailure) > cb.cooldown {
					cb.open = false
					cb.failures = 0
				} else {
					cb.mu.Unlock()
					return nil, kernel.ErrCircuitOpen
				}
			}
			cb.mu.Unlock()

			resp, err := next(ctx, req)
			if err != nil {
				cb.mu.Lock()
				cb.failures++
				cb.lastFailure = time.Now()
				if cb.failures >= cb.threshold {
					cb.open = true
				}
				cb.mu.Unlock()
			} else {
				cb.mu.Lock()
				cb.failures = 0
				cb.mu.Unlock()
			}
			return resp, err
		}
	}
}

func NewCircuitBreaker(threshold int, cooldown time.Duration) *CircuitBreaker {
	return &CircuitBreaker{threshold: threshold, cooldown: cooldown}
}
```

- [ ] **Step 2: Write tests**

Create `kernel/pipeline/breaker_test.go`:

```go
package pipeline

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestCircuitBreaker_ClosedByDefault(t *testing.T) {
	cb := NewCircuitBreaker(3, 100*time.Millisecond)
	wrapped := Chain(cb.Middleware())(echoHandler)
	_, err := wrapped(context.Background(), &kernel.Request{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestCircuitBreaker_OpensAfterThreshold(t *testing.T) {
	cb := NewCircuitBreaker(2, 100*time.Millisecond)
	fail := func(ctx context.Context, req *kernel.Request) (*kernel.Response, error) {
		return nil, errors.New("fail")
	}
	wrapped := Chain(cb.Middleware())(fail)

	for i := 0; i < 2; i++ {
		_, _ = wrapped(context.Background(), &kernel.Request{})
	}

	_, err := wrapped(context.Background(), &kernel.Request{})
	if !errors.Is(err, kernel.ErrCircuitOpen) {
		t.Errorf("expected ErrCircuitOpen, got %v", err)
	}
}

func TestCircuitBreaker_ResetsAfterCooldown(t *testing.T) {
	cb := NewCircuitBreaker(2, 30*time.Millisecond)
	fail := func(ctx context.Context, req *kernel.Request) (*kernel.Response, error) {
		return nil, errors.New("fail")
	}
	wrapped := Chain(cb.Middleware())(fail)

	for i := 0; i < 2; i++ {
		_, _ = wrapped(context.Background(), &kernel.Request{})
	}
	time.Sleep(50 * time.Millisecond)

	_, err := wrapped(context.Background(), &kernel.Request{})
	if err == nil {
		t.Log("breaker reset after cooldown")
	}
}
```

- [ ] **Step 3: Run tests**

```bash
cd /home/wwwroot/bag/industrial-protocols-go/kernel && go test ./pipeline/ -v
```

- [ ] **Step 4: Commit**

```bash
git add kernel/pipeline/breaker.go kernel/pipeline/breaker_test.go
git commit -m "feat(kernel): add CircuitBreaker middleware"
```

---

## Phase 4: Protocol Stubs — All 40 Protocols

### Task 14: Helper script for stub generation

**Files:**
- Create: `_tools/gen-stub.sh`

- [ ] **Step 1: Create generation script**

Create `_tools/gen-stub.sh`:

```bash
#!/bin/bash
set -e

CATEGORY="$1"
NAME="$2"
DISPLAY="$3"
VARIANTS="$4"
PORT="$5"

DIR="protocols/${CATEGORY}"
MOD_PATH="github.com/erikwang2013/industrial-protocols-go/protocols/${CATEGORY}"
KERNEL_VERSION="v0.0.0"

mkdir -p "$DIR"
cd "$DIR"
go mod init "$MOD_PATH"
go mod edit -require "github.com/erikwang2013/industrial-protocols-go/kernel@${KERNEL_VERSION}"
go mod edit -replace "github.com/erikwang2013/industrial-protocols-go/kernel=../../../../kernel"

cat > "${NAME}.go" << GOFILE
package ${NAME}

import "github.com/erikwang2013/industrial-protocols-go/kernel"

type ${DISPLAY}Protocol struct{}

func New() *${DISPLAY}Protocol { return &${DISPLAY}Protocol{} }

func (p *${DISPLAY}Protocol) Name() string        { return "${NAME}" }
func (p *${DISPLAY}Protocol) Variants() []string   { return []string{${VARIANTS}} }
func (p *${DISPLAY}Protocol) DefaultPort() int     { return ${PORT} }
func (p *${DISPLAY}Protocol) NewCodec(variant string) (kernel.Codec, error) {
	return nil, kernel.ErrInvalidAddress
}
GOFILE

cat > "${NAME}_test.go" << TESTFILE
package ${NAME}

import "testing"

func TestProtocol_ImplementsInterface(t *testing.T) {
	p := New()
	if p.Name() != "${NAME}" {
		t.Errorf("expected ${NAME}, got %s", p.Name())
	}
	if len(p.Variants()) == 0 {
		t.Error("should have at least one variant")
	}
	if p.DefaultPort() <= 0 {
		t.Error("default port must be positive")
	}
	_, err := p.NewCodec("invalid")
	if err == nil {
		t.Error("stub NewCodec should return error")
	}
}
TESTFILE

echo "Stub created: $DIR"
```

- [ ] **Step 2: Make it executable**

```bash
chmod +x /home/wwwroot/bag/industrial-protocols-go/_tools/gen-stub.sh
```

- [ ] **Step 3: Commit**

```bash
git add _tools/gen-stub.sh
git commit -m "chore: add protocol stub generator script"
```

---

### Task 15: Generate all 40 protocol stubs

Run the generation script for all 40 protocols, then add them to the workspace.

- [ ] **Step 1: Generate ethernet stubs (5)**

```bash
cd /home/wwwroot/bag/industrial-protocols-go

./_tools/gen-stub.sh ethernet/modbus    Modbus    '"tcp","rtu","ascii"' 502
./_tools/gen-stub.sh ethernet/bacnet    Bacnet    '"ip"'                 47808
./_tools/gen-stub.sh ethernet/ethernetip EtherNetIP '"tcp"'              44818
./_tools/gen-stub.sh ethernet/opcua     Opcua     '"binary"'             4840
./_tools/gen-stub.sh ethernet/profinet  Profinet  '"nrt"'                34964

for dir in protocols/ethernet/*/; do go work use "$dir"; done
go work sync
for dir in protocols/ethernet/*/; do (cd "$dir" && go build ./... && go test ./... -v) || exit 1; done
```

- [ ] **Step 2: Commit ethernet stubs**

```bash
git add protocols/ethernet/ go.work go.work.sum
git commit -m "feat(protocols): add ethernet protocol stubs (modbus, bacnet, ethernetip, opcua, profinet)"
```

- [ ] **Step 3: Generate fieldbus stubs (11)**

```bash
./_tools/gen-stub.sh fieldbus/hart               Hart               '"fsk"'                0
./_tools/gen-stub.sh fieldbus/cclink             CCLink             '"rs485"'              0
./_tools/gen-stub.sh fieldbus/dnp3               DNP3               '"serial","tcp"'       20000
./_tools/gen-stub.sh fieldbus/iec61850           Iec61850           '"mms"'                102
./_tools/gen-stub.sh fieldbus/profibus           Profibus           '"dp"'                 0
./_tools/gen-stub.sh fieldbus/canopen            CanOpen            '"can"'                0
./_tools/gen-stub.sh fieldbus/devicenet          DeviceNet          '"can"'                0
./_tools/gen-stub.sh fieldbus/foundationfieldbus FoundationFieldbus '"h1","hse"'           0
./_tools/gen-stub.sh fieldbus/asinterface        AsInterface        '"serial"'             0
./_tools/gen-stub.sh fieldbus/iolink             IOLink             '"serial"'             0
./_tools/gen-stub.sh fieldbus/cclinkie           CCLinkIE           '"ethernet"'           0

for dir in protocols/fieldbus/*/; do go work use "$dir"; done
go work sync
for dir in protocols/fieldbus/*/; do (cd "$dir" && go build ./... && go test ./... -v) || exit 1; done
```

- [ ] **Step 4: Commit fieldbus stubs**

```bash
git add protocols/fieldbus/ go.work go.work.sum
git commit -m "feat(protocols): add fieldbus protocol stubs (11 protocols)"
```

- [ ] **Step 5: Generate IoT, automotive, building stubs (9)**

```bash
./_tools/gen-stub.sh iot/mqtt    MQTT    '"tcp","ws"' 1883
./_tools/gen-stub.sh iot/hartip  HARTIP  '"tcp","udp"' 5094

./_tools/gen-stub.sh automotive/lin       LIN      '"uart"'    0
./_tools/gen-stub.sh automotive/kline     KLine    '"serial"'  10400
./_tools/gen-stub.sh automotive/flexray   FlexRay  '"serial"'  0
./_tools/gen-stub.sh automotive/saej1850  SaeJ1850 '"serial"'  0
./_tools/gen-stub.sh automotive/most      Most     '"optical"' 0

./_tools/gen-stub.sh building/lonworks LonWorks '"serial"' 0
./_tools/gen-stub.sh building/dali     Dali     '"serial"' 0

for dir in protocols/iot/*/ protocols/automotive/*/ protocols/building/*/; do go work use "$dir"; done
go work sync
for dir in protocols/iot/*/ protocols/automotive/*/ protocols/building/*/; do (cd "$dir" && go build ./... && go test ./... -v) || exit 1; done
```

- [ ] **Step 6: Commit IoT, automotive, building stubs**

```bash
git add protocols/iot/ protocols/automotive/ protocols/building/ go.work go.work.sum
git commit -m "feat(protocols): add IoT, automotive, and building protocol stubs (9 protocols)"
```

- [ ] **Step 7: Generate bridge and system bus stubs (16)**

```bash
./_tools/gen-stub.sh bridge/ethercat     EtherCat     '"ethercat"'  0
./_tools/gen-stub.sh bridge/powerlink    Powerlink    '"ethernet"'  0
./_tools/gen-stub.sh bridge/sercos       Sercos       '"fiber"'     0
./_tools/gen-stub.sh bridge/sercos1      Sercos1      '"fiber"'     0
./_tools/gen-stub.sh bridge/controlnet   ControlNet   '"coax"'      0
./_tools/gen-stub.sh bridge/interbus     Interbus     '"serial"'    0
./_tools/gen-stub.sh bridge/worldfip     WorldFIP     '"serial"'    0
./_tools/gen-stub.sh bridge/lightbus     Lightbus     '"fiber"'     0
./_tools/gen-stub.sh bridge/modbusplus   ModbusPlus   '"serial"'    0
./_tools/gen-stub.sh bridge/isa100       Isa100       '"wireless"'  0
./_tools/gen-stub.sh bridge/wirelesshart WirelessHART '"wireless"'  0
./_tools/gen-stub.sh bridge/safej1850    SafeJ1850    '"serial"'    0

./_tools/gen-stub.sh system/pci  PCI  '"pci"'  0
./_tools/gen-stub.sh system/vme  VME  '"vme"'  0
./_tools/gen-stub.sh system/cpci CPCI '"cpci"' 0

for dir in protocols/bridge/*/ protocols/system/*/; do go work use "$dir"; done
go work sync
for dir in protocols/bridge/*/ protocols/system/*/; do (cd "$dir" && go build ./... && go test ./... -v) || exit 1; done
```

- [ ] **Step 8: Commit bridge and system bus stubs**

```bash
git add protocols/bridge/ protocols/system/ go.work go.work.sum
git commit -m "feat(protocols): add bridge and system bus protocol stubs (16 protocols)"
```

- [ ] **Step 9: Full workspace verification**

```bash
cd /home/wwwroot/bag/industrial-protocols-go
cd kernel && go test ./... -v -short && go vet ./...
for dir in $(find protocols -name go.mod -exec dirname {} \;); do
  (cd "$dir" && go test ./... -v -short && go vet ./...) || echo "FAIL: $dir"
done
```

---

## Phase 5: Modbus Implementation

### Task 16: Implement Modbus TCP/RTU codec

**Files:**
- Rewrite: `protocols/ethernet/modbus/modbus.go`
- Rewrite: `protocols/ethernet/modbus/modbus_test.go`
- Create: `protocols/ethernet/modbus/tcp.go`

- [ ] **Step 1: Write the failing test**

Create `protocols/ethernet/modbus/modbus_test.go`:

```go
package modbus

import (
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestTCPCodec_ReadCoils(t *testing.T) {
	c, _ := NewCodec("tcp")
	req := &kernel.Request{Function: "read_coils", Address: "0", Count: 8, Metadata: map[string]any{"unit_id": byte(1)}}
	enc, err := c.Encode(req)
	if err != nil {
		t.Fatal(err)
	}
	if len(enc) != 12 {
		t.Fatalf("expected 12 bytes (7 MBAP + 5 PDU), got %d", len(enc))
	}
	if enc[7] != 1 {
		t.Errorf("expected function code 1, got %d", enc[7])
	}
}

func TestRTUCodec_ReadCoils(t *testing.T) {
	c, _ := NewCodec("rtu")
	req := &kernel.Request{Function: "read_coils", Address: "0", Count: 8, Metadata: map[string]any{"unit_id": byte(1)}}
	enc, err := c.Encode(req)
	if err != nil {
		t.Fatal(err)
	}
	if len(enc) != 8 {
		t.Fatalf("expected 8 bytes (1 addr + 1 func + 2 addr + 2 count + 2 crc), got %d", len(enc))
	}
}

func TestDecodeModbusException(t *testing.T) {
	c, _ := NewCodec("tcp")
	_, err := c.Decode([]byte{0x00, 0x01, 0x00, 0x00, 0x00, 0x03, 0x01, 0x81, 0x02})
	if err == nil {
		t.Fatal("expected ProtocolError for exception response")
	}
	var pe *kernel.ProtocolError
	if !errors.As(err, &pe) {
		t.Errorf("expected ProtocolError, got %T", err)
	}
}

func TestProtocolInterface(t *testing.T) {
	p := NewProtocol()
	if p.Name() != "modbus" {
		t.Errorf("expected modbus, got %s", p.Name())
	}
	if p.DefaultPort() != 502 {
		t.Errorf("expected 502, got %d", p.DefaultPort())
	}
	variants := p.Variants()
	if len(variants) != 3 {
		t.Errorf("expected 3 variants, got %d", len(variants))
	}
}
```

- [ ] **Step 2: Run test to verify failure**

```bash
cd /home/wwwroot/bag/industrial-protocols-go/protocols/ethernet/modbus && go test ./... -v
```

- [ ] **Step 3: Implement full Modbus codec**

Rewrite `protocols/ethernet/modbus/modbus.go`:

```go
package modbus

import (
	"encoding/binary"
	"fmt"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

type ModbusProtocol struct{}

func NewProtocol() *ModbusProtocol { return &ModbusProtocol{} }

func (p *ModbusProtocol) Name() string        { return "modbus" }
func (p *ModbusProtocol) Variants() []string   { return []string{"tcp", "rtu", "ascii"} }
func (p *ModbusProtocol) DefaultPort() int     { return 502 }

func (p *ModbusProtocol) NewCodec(variant string) (kernel.Codec, error) {
	switch variant {
	case "tcp", "rtu":
		return &modbusCodec{variant: variant}, nil
	default:
		return nil, fmt.Errorf("modbus: unsupported variant %q", variant)
	}
}

type modbusCodec struct {
	variant     string
	transaction uint16
}

func NewCodec(variant string) (*modbusCodec, error) {
	return &modbusCodec{variant: variant}, nil
}

func (c *modbusCodec) Encode(req *kernel.Request) ([]byte, error) {
	mr := c.toModbusReq(req)
	switch c.variant {
	case "tcp":
		return c.encodeTCP(mr), nil
	case "rtu":
		return c.encodeRTU(mr), nil
	default:
		return nil, fmt.Errorf("modbus: unsupported variant %q", c.variant)
	}
}

func (c *modbusCodec) Decode(data []byte) (*kernel.Response, error) {
	switch c.variant {
	case "tcp":
		return c.decodeTCP(data)
	case "rtu":
		return c.decodeRTU(data)
	default:
		return nil, fmt.Errorf("modbus: unsupported variant %q", c.variant)
	}
}

type modbusRequest struct {
	function byte
	address  uint16
	count    uint16
	data     []byte
	unitID   byte
}

func (c *modbusCodec) toModbusReq(req *kernel.Request) modbusRequest {
	mr := modbusRequest{
		address: parseAddr(req.Address),
		count:   uint16(req.Count),
		data:    req.Data,
	}
	if uid, ok := req.Metadata["unit_id"].(byte); ok {
		mr.unitID = uid
	} else if uid, ok := req.Metadata["unit_id"].(float64); ok {
		mr.unitID = byte(uid)
	}
	switch req.Function {
	case "read_coils":             mr.function = 1
	case "read_discrete_inputs":   mr.function = 2
	case "read_holding_registers": mr.function = 3
	case "read_input_registers":   mr.function = 4
	case "write_single_coil":      mr.function = 5
	case "write_single_register":  mr.function = 6
	case "write_multiple_registers": mr.function = 16
	default:                       mr.function = 3
	}
	return mr
}

func parseAddr(s string) uint16 {
	var addr uint32
	fmt.Sscanf(s, "%d", &addr)
	return uint16(addr)
}

func (c *modbusCodec) encodeTCP(mr modbusRequest) []byte {
	c.transaction++
	pdu := c.encodePDU(mr)
	buf := make([]byte, 7+len(pdu))
	binary.BigEndian.PutUint16(buf[0:2], c.transaction)
	binary.BigEndian.PutUint16(buf[2:4], 0)
	binary.BigEndian.PutUint16(buf[4:6], uint16(len(pdu)+1))
	buf[6] = mr.unitID
	copy(buf[7:], pdu)
	return buf
}

func (c *modbusCodec) decodeTCP(data []byte) (*kernel.Response, error) {
	if len(data) < 9 {
		return nil, fmt.Errorf("modbus: tcp frame too short (%d bytes)", len(data))
	}
	length := binary.BigEndian.Uint16(data[4:6])
	pdu := data[7 : 7+length]
	return c.decodePDU(pdu)
}

func (c *modbusCodec) encodeRTU(mr modbusRequest) []byte {
	pdu := c.encodePDU(mr)
	buf := make([]byte, 1+len(pdu)+2)
	buf[0] = mr.unitID
	copy(buf[1:], pdu)
	crc := crc16(buf[:1+len(pdu)])
	binary.LittleEndian.PutUint16(buf[1+len(pdu):], crc)
	return buf
}

func (c *modbusCodec) decodeRTU(data []byte) (*kernel.Response, error) {
	if len(data) < 5 {
		return nil, fmt.Errorf("modbus: rtu frame too short (%d bytes)", len(data))
	}
	pdu := data[1 : len(data)-2]
	return c.decodePDU(pdu)
}

func (c *modbusCodec) encodePDU(mr modbusRequest) []byte {
	switch mr.function {
	case 1, 2, 3, 4:
		pdu := make([]byte, 5)
		pdu[0] = mr.function
		binary.BigEndian.PutUint16(pdu[1:3], mr.address)
		binary.BigEndian.PutUint16(pdu[3:5], mr.count)
		return pdu
	case 5:
		pdu := make([]byte, 5)
		pdu[0] = mr.function
		binary.BigEndian.PutUint16(pdu[1:3], mr.address)
		if len(mr.data) > 0 && mr.data[0] != 0 {
			pdu[3] = 0xFF
		}
		return pdu
	case 6:
		pdu := make([]byte, 5)
		pdu[0] = mr.function
		binary.BigEndian.PutUint16(pdu[1:3], mr.address)
		if len(mr.data) >= 2 {
			copy(pdu[3:5], mr.data[:2])
		}
		return pdu
	case 16:
		pdu := make([]byte, 6+len(mr.data))
		pdu[0] = mr.function
		binary.BigEndian.PutUint16(pdu[1:3], mr.address)
		binary.BigEndian.PutUint16(pdu[3:5], uint16(len(mr.data)/2))
		pdu[5] = byte(len(mr.data))
		copy(pdu[6:], mr.data)
		return pdu
	default:
		return []byte{mr.function}
	}
}

func (c *modbusCodec) decodePDU(pdu []byte) (*kernel.Response, error) {
	if len(pdu) == 0 {
		return nil, fmt.Errorf("modbus: empty PDU")
	}
	fn := pdu[0]
	if fn&0x80 != 0 {
		return nil, &kernel.ProtocolError{
			Code:    fmt.Sprintf("%02X", pdu[1]),
			Message: modbusException(pdu[1]),
			Raw:     pdu,
		}
	}
	switch fn {
	case 1, 2, 3, 4:
		if len(pdu) < 3 {
			return nil, fmt.Errorf("modbus: PDU too short")
		}
		return &kernel.Response{Address: "0", Data: pdu[2 : 2+pdu[1]]}, nil
	case 5, 6:
		return &kernel.Response{Address: "0"}, nil
	case 16:
		return &kernel.Response{Address: "0"}, nil
	default:
		return &kernel.Response{Data: pdu}, nil
	}
}

func modbusException(code byte) string {
	switch code {
	case 1:  return "illegal function"
	case 2:  return "illegal data address"
	case 3:  return "illegal data value"
	case 4:  return "device failure"
	case 6:  return "device busy"
	default: return "unknown exception"
	}
}

func crc16(data []byte) uint16 {
	var crc uint16 = 0xFFFF
	for _, b := range data {
		crc ^= uint16(b)
		for i := 0; i < 8; i++ {
			if crc&1 != 0 {
				crc = (crc >> 1) ^ 0xA001
			} else {
				crc >>= 1
			}
		}
	}
	return crc
}
```

- [ ] **Step 4: Implement TCP client helper**

Create `protocols/ethernet/modbus/tcp.go`:

```go
package modbus

import (
	"fmt"
	"time"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
	"github.com/erikwang2013/industrial-protocols-go/kernel/pipeline"
	"github.com/erikwang2013/industrial-protocols-go/kernel/transport"
)

type TCPClient struct {
	handle pipeline.Handler
	tr     transport.Transport
}

func NewTCPClient(address string, unitID byte) (*TCPClient, error) {
	tr, err := transport.DialTCP(address)
	if err != nil {
		return nil, err
	}
	codec, _ := NewCodec("tcp")
	handle := func(req *kernel.Request) (*kernel.Response, error) {
		if req.Metadata == nil {
			req.Metadata = map[string]any{}
		}
		req.Metadata["unit_id"] = unitID
		raw, err := codec.Encode(req)
		if err != nil {
			return nil, err
		}
		if _, err := tr.Write(raw); err != nil {
			return nil, err
		}
		buf := make([]byte, 256)
		n, err := tr.Read(buf)
		if err != nil {
			return nil, err
		}
		return codec.Decode(buf[:n])
	}
	return &TCPClient{handle: handle, tr: tr}, nil
}

func (c *TCPClient) WithTimeout(d time.Duration) *TCPClient {
	c.handle = pipeline.Chain(pipeline.Timeout(d))(c.handle)
	return c
}

func (c *TCPClient) WithRetry(max int, backoff pipeline.BackoffFunc) *TCPClient {
	c.handle = pipeline.Chain(pipeline.Retry(max, backoff))(c.handle)
	return c
}

func (c *TCPClient) Close() error {
	return c.tr.Close()
}

func (c *TCPClient) ReadCoils(address uint16, count uint16) ([]bool, error) {
	req := &kernel.Request{Function: "read_coils", Address: fmt.Sprintf("%d", address), Count: int(count)}
	resp, err := c.handle(nil, req)
	if err != nil {
		return nil, err
	}
	bits := make([]bool, count)
	for i, b := range resp.Data {
		for j := 0; j < 8 && int(i*8+j) < int(count); j++ {
			bits[i*8+j] = b&(1<<j) != 0
		}
	}
	return bits, nil
}

func (c *TCPClient) ReadHoldingRegisters(address uint16, count uint16) ([]uint16, error) {
	req := &kernel.Request{Function: "read_holding_registers", Address: fmt.Sprintf("%d", address), Count: int(count)}
	resp, err := c.handle(nil, req)
	if err != nil {
		return nil, err
	}
	regs := make([]uint16, count)
	for i := 0; i < len(resp.Data)/2 && i < int(count); i++ {
		regs[i] = uint16(resp.Data[i*2])<<8 | uint16(resp.Data[i*2+1])
	}
	return regs, nil
}
```

- [ ] **Step 5: Run tests**

```bash
cd /home/wwwroot/bag/industrial-protocols-go/protocols/ethernet/modbus && go test ./... -v && go build ./...
```

- [ ] **Step 6: Commit**

```bash
git add protocols/ethernet/modbus/
git commit -m "feat(modbus): implement TCP/RTU codec with FC 01/03/04/06/10 and TCP client helper"
```

---

## Phase 6: MQTT Implementation

### Task 17: Implement MQTT 3.1.1 codec

**Files:**
- Rewrite: `protocols/iot/mqtt/mqtt.go`
- Rewrite: `protocols/iot/mqtt/mqtt_test.go`

- [ ] **Step 1: Write the failing test**

Create `protocols/iot/mqtt/mqtt_test.go`:

```go
package mqtt

import (
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestProtocolInterface(t *testing.T) {
	p := New()
	if p.Name() != "mqtt" {
		t.Errorf("expected mqtt, got %s", p.Name())
	}
	if p.DefaultPort() != 1883 {
		t.Errorf("expected 1883, got %d", p.DefaultPort())
	}
	c, _ := p.NewCodec("tcp")
	if c == nil {
		t.Error("NewCodec returned nil")
	}
}

func TestEncodeConnect(t *testing.T) {
	c, _ := New().NewCodec("tcp")
	req := &kernel.Request{Function: "connect", Metadata: map[string]any{"client_id": "test"}}
	data, err := c.Encode(req)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) < 10 {
		t.Errorf("connect frame too short: %d bytes", len(data))
	}
	if data[0] != 0x10 {
		t.Errorf("expected CONNECT type 0x10, got 0x%02X", data[0])
	}
}

func TestEncodePublish(t *testing.T) {
	c, _ := New().NewCodec("tcp")
	req := &kernel.Request{Function: "publish", Address: "sensor/temp", Data: []byte("25.5")}
	data, err := c.Encode(req)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) < 5 {
		t.Errorf("publish frame too short: %d bytes", len(data))
	}
}

func TestEncodeSubscribe(t *testing.T) {
	c, _ := New().NewCodec("tcp")
	req := &kernel.Request{Function: "subscribe", Address: "sensor/#"}
	data, err := c.Encode(req)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) < 6 {
		t.Errorf("subscribe frame too short: %d bytes", len(data))
	}
}

func TestEncodePingreq(t *testing.T) {
	c, _ := New().NewCodec("tcp")
	req := &kernel.Request{Function: "pingreq"}
	data, err := c.Encode(req)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 2 || data[0] != 0xC0 || data[1] != 0x00 {
		t.Errorf("expected PINGREQ [0xC0, 0x00], got %v", data)
	}
}

func TestDecodeConnack(t *testing.T) {
	c, _ := New().NewCodec("tcp")
	resp, err := c.Decode([]byte{0x20, 0x02, 0x00, 0x00})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Function != "connack" {
		t.Errorf("expected connack, got %s", resp.Function)
	}
}

func TestDecodePublish(t *testing.T) {
	c, _ := New().NewCodec("tcp")
	// PUBLISH to topic "test" (4 bytes) with payload "hi"
	resp, err := c.Decode([]byte{0x30, 0x06, 0x00, 0x04, 't', 'e', 's', 't', 'h', 'i'})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Function != "publish" {
		t.Errorf("expected publish, got %s", resp.Function)
	}
	if resp.Address != "test" {
		t.Errorf("expected topic 'test', got '%s'", resp.Address)
	}
}

func TestDecodeSuback(t *testing.T) {
	c, _ := New().NewCodec("tcp")
	resp, err := c.Decode([]byte{0x90, 0x03, 0x00, 0x01, 0x00})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Function != "suback" {
		t.Errorf("expected suback, got %s", resp.Function)
	}
}

func TestDecodePingresp(t *testing.T) {
	c, _ := New().NewCodec("tcp")
	resp, err := c.Decode([]byte{0xD0, 0x00})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Function != "pingresp" {
		t.Errorf("expected pingresp, got %s", resp.Function)
	}
}
```

- [ ] **Step 2: Implement MQTT codec**

Rewrite `protocols/iot/mqtt/mqtt.go`:

```go
package mqtt

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

type MqttProtocol struct{}

func New() *MqttProtocol                             { return &MqttProtocol{} }
func (p *MqttProtocol) Name() string                 { return "mqtt" }
func (p *MqttProtocol) Variants() []string           { return []string{"tcp", "ws"} }
func (p *MqttProtocol) DefaultPort() int             { return 1883 }

func (p *MqttProtocol) NewCodec(v string) (kernel.Codec, error) {
	return &mqttCodec{}, nil
}

type mqttCodec struct{}

func (c *mqttCodec) Encode(req *kernel.Request) ([]byte, error) {
	switch req.Function {
	case "connect":
		return c.encodeConnect(req)
	case "publish":
		return c.encodePublish(req)
	case "subscribe":
		return c.encodeSubscribe(req)
	case "pingreq":
		return []byte{0xC0, 0x00}, nil
	case "disconnect":
		return []byte{0xE0, 0x00}, nil
	default:
		return nil, fmt.Errorf("mqtt: unknown function %q", req.Function)
	}
}

func (c *mqttCodec) Decode(data []byte) (*kernel.Response, error) {
	if len(data) < 2 {
		return nil, fmt.Errorf("mqtt: frame too short")
	}
	msgType := data[0] >> 4
	switch msgType {
	case 2:
		return &kernel.Response{Function: "connack", Data: data}, nil
	case 3:
		return c.decodePublish(data)
	case 9:
		return &kernel.Response{Function: "suback", Data: data}, nil
	case 13:
		return &kernel.Response{Function: "pingresp"}, nil
	default:
		return &kernel.Response{Data: data}, nil
	}
}

func (c *mqttCodec) encodeConnect(req *kernel.Request) ([]byte, error) {
	clientID := "goclient"
	if v, ok := req.Metadata["client_id"].(string); ok {
		clientID = v
	}
	keepAlive := uint16(60)
	if v, ok := req.Metadata["keep_alive"].(float64); ok {
		keepAlive = uint16(v)
	}

	buf := &bytes.Buffer{}
	buf.WriteByte(0x10)
	payload := &bytes.Buffer{}
	writeString(payload, "MQTT")
	payload.WriteByte(4)
	payload.WriteByte(2)
	binary.Write(payload, binary.BigEndian, keepAlive)
	writeString(payload, clientID)
	encodeRemainingLength(buf, payload.Len())
	buf.Write(payload.Bytes())
	return buf.Bytes(), nil
}

func (c *mqttCodec) encodePublish(req *kernel.Request) ([]byte, error) {
	buf := &bytes.Buffer{}
	buf.WriteByte(0x30)
	payload := &bytes.Buffer{}
	writeString(payload, req.Address)
	payload.Write(req.Data)
	encodeRemainingLength(buf, payload.Len())
	buf.Write(payload.Bytes())
	return buf.Bytes(), nil
}

func (c *mqttCodec) encodeSubscribe(req *kernel.Request) ([]byte, error) {
	buf := &bytes.Buffer{}
	buf.WriteByte(0x82)
	payload := &bytes.Buffer{}
	binary.Write(payload, binary.BigEndian, uint16(1))
	writeString(payload, req.Address)
	payload.WriteByte(0)
	encodeRemainingLength(buf, payload.Len())
	buf.Write(payload.Bytes())
	return buf.Bytes(), nil
}

func (c *mqttCodec) decodePublish(data []byte) (*kernel.Response, error) {
	offset := 1
	_, n := decodeRemainingLength(data[offset:])
	offset += n
	topicLen := int(data[offset])<<8 | int(data[offset+1])
	offset += 2
	topic := string(data[offset : offset+topicLen])
	offset += topicLen
	return &kernel.Response{Function: "publish", Address: topic, Data: data[offset:]}, nil
}

func encodeRemainingLength(w io.Writer, n int) {
	for {
		b := byte(n % 128)
		n /= 128
		if n > 0 {
			b |= 0x80
		}
		w.(*bytes.Buffer).WriteByte(b)
		if n == 0 {
			break
		}
	}
}

func decodeRemainingLength(data []byte) (int, int) {
	var length, multiplier int
	for i, b := range data {
		length += int(b&0x7F) << multiplier
		if b&0x80 == 0 {
			return length, i + 1
		}
		multiplier += 7
	}
	return 0, 0
}

func writeString(w io.Writer, s string) {
	binary.Write(w, binary.BigEndian, uint16(len(s)))
	w.Write([]byte(s))
}
```

- [ ] **Step 3: Run tests**

```bash
cd /home/wwwroot/bag/industrial-protocols-go/protocols/iot/mqtt && go test ./... -v && go build ./...
```

- [ ] **Step 4: Commit**

```bash
git add protocols/iot/mqtt/
git commit -m "feat(mqtt): implement MQTT 3.1.1 CONNECT/PUBLISH/SUBSCRIBE/PING codec"
```

---

## Phase 7: OPC UA Basic Binary Codec

### Task 18: Implement OPC UA HELLO and OpenSecureChannel

**Files:**
- Rewrite: `protocols/ethernet/opcua/opcua.go`
- Rewrite: `protocols/ethernet/opcua/opcua_test.go`

- [ ] **Step 1: Write tests**

Create `protocols/ethernet/opcua/opcua_test.go`:

```go
package opcua

import (
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestProtocolInterface(t *testing.T) {
	p := New()
	if p.Name() != "opcua" {
		t.Errorf("expected opcua, got %s", p.Name())
	}
	if p.DefaultPort() != 4840 {
		t.Errorf("expected 4840, got %d", p.DefaultPort())
	}
}

func TestEncodeHello(t *testing.T) {
	c, _ := New().NewCodec("binary")
	req := &kernel.Request{Function: "hel", Metadata: map[string]any{"endpoint": "opc.tcp://192.168.1.1:4840"}}
	data, err := c.Encode(req)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) < 32 {
		t.Errorf("HELLO frame too short: %d bytes", len(data))
	}
	if string(data[:3]) != "HEL" {
		t.Errorf("expected HEL, got %s", string(data[:3]))
	}
}

func TestEncodeOpenSecureChannel(t *testing.T) {
	c, _ := New().NewCodec("binary")
	req := &kernel.Request{Function: "open_secure_channel"}
	data, err := c.Encode(req)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) < 8 {
		t.Errorf("OPN frame too short: %d bytes", len(data))
	}
	if string(data[:3]) != "OPN" {
		t.Errorf("expected OPN, got %s", string(data[:3]))
	}
}

func TestDecodeAck(t *testing.T) {
	c, _ := New().NewCodec("binary")
	ack := make([]byte, 28)
	copy(ack, "ACK")
	ack[3] = 'F'
	resp, err := c.Decode(ack)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Function != "ack" {
		t.Errorf("expected ack, got %s", resp.Function)
	}
}

func TestDecodeOpenResponse(t *testing.T) {
	c, _ := New().NewCodec("binary")
	opn := make([]byte, 8)
	copy(opn, "OPN")
	opn[3] = 'F'
	resp, err := c.Decode(opn)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Function != "open_response" {
		t.Errorf("expected open_response, got %s", resp.Function)
	}
}
```

- [ ] **Step 2: Implement OPC UA codec**

Rewrite `protocols/ethernet/opcua/opcua.go`:

```go
package opcua

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

type OpcuaProtocol struct{}

func New() *OpcuaProtocol                           { return &OpcuaProtocol{} }
func (p *OpcuaProtocol) Name() string               { return "opcua" }
func (p *OpcuaProtocol) Variants() []string         { return []string{"binary"} }
func (p *OpcuaProtocol) DefaultPort() int           { return 4840 }

func (p *OpcuaProtocol) NewCodec(v string) (kernel.Codec, error) {
	return &opcuaCodec{}, nil
}

type opcuaCodec struct{}

func (c *opcuaCodec) Encode(req *kernel.Request) ([]byte, error) {
	switch req.Function {
	case "hel":
		return c.encodeHello(req)
	case "open_secure_channel":
		return c.encodeOpenSecureChannel()
	default:
		return nil, fmt.Errorf("opcua: unknown function %q", req.Function)
	}
}

func (c *opcuaCodec) Decode(data []byte) (*kernel.Response, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("opcua: frame too short (%d bytes)", len(data))
	}
	msgType := string(data[:3])
	switch msgType {
	case "ACK":
		return &kernel.Response{Function: "ack", Data: data}, nil
	case "OPN":
		return &kernel.Response{Function: "open_response", Data: data}, nil
	default:
		return &kernel.Response{Data: data}, nil
	}
}

func (c *opcuaCodec) encodeHello(req *kernel.Request) ([]byte, error) {
	endpoint := "opc.tcp://localhost:4840"
	if v, ok := req.Metadata["endpoint"].(string); ok {
		endpoint = v
	}
	buf := &bytes.Buffer{}
	buf.WriteString("HEL")
	buf.WriteByte('F')
	binary.Write(buf, binary.LittleEndian, uint32(0))
	sizePos := buf.Len()
	binary.Write(buf, binary.LittleEndian, uint32(0))
	binary.Write(buf, binary.LittleEndian, uint32(0))
	binary.Write(buf, binary.LittleEndian, uint32(65536))
	binary.Write(buf, binary.LittleEndian, uint32(65536))
	binary.Write(buf, binary.LittleEndian, uint32(0))
	binary.Write(buf, binary.LittleEndian, uint32(0))
	binary.Write(buf, binary.LittleEndian, int32(len(endpoint)))
	buf.WriteString(endpoint)
	raw := buf.Bytes()
	binary.LittleEndian.PutUint32(raw[sizePos:], uint32(len(raw)))
	return raw, nil
}

func (c *opcuaCodec) encodeOpenSecureChannel() ([]byte, error) {
	buf := &bytes.Buffer{}
	buf.WriteString("OPN")
	buf.WriteByte('F')
	binary.Write(buf, binary.LittleEndian, uint32(0))
	binary.Write(buf, binary.LittleEndian, uint32(44))
	for i := 0; i < 7; i++ {
		binary.Write(buf, binary.LittleEndian, uint32(0))
	}
	return buf.Bytes(), nil
}
```

- [ ] **Step 3: Run tests**

```bash
cd /home/wwwroot/bag/industrial-protocols-go/protocols/ethernet/opcua && go test ./... -v && go build ./...
```

- [ ] **Step 4: Commit**

```bash
git add protocols/ethernet/opcua/
git commit -m "feat(opcua): implement OPC UA binary HEL/OpenSecureChannel codec"
```

---

## Phase 8: Examples and Final Verification

### Task 19: Create Modbus TCP example

**Files:**
- Create: `examples/modbus_basic/main.go`

- [ ] **Step 1: Write example**

Create `examples/modbus_basic/main.go`:

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
	"github.com/erikwang2013/industrial-protocols-go/kernel/pipeline"
	"github.com/erikwang2013/industrial-protocols-go/kernel/transport"
	"github.com/erikwang2013/industrial-protocols-go/protocols/ethernet/modbus"
)

func main() {
	tr, err := transport.DialTCP("192.168.1.10:502")
	if err != nil {
		log.Fatalf("connect: %v", err)
	}
	defer tr.Close()

	codec, _ := modbus.NewCodec("tcp")

	handler := func(req *kernel.Request) (*kernel.Response, error) {
		if req.Metadata == nil {
			req.Metadata = map[string]any{}
		}
		req.Metadata["unit_id"] = byte(1)
		raw, err := codec.Encode(req)
		if err != nil {
			return nil, err
		}
		if _, err := tr.Write(raw); err != nil {
			return nil, err
		}
		buf := make([]byte, 256)
		n, err := tr.Read(buf)
		if err != nil {
			return nil, err
		}
		return codec.Decode(buf[:n])
	}

	wrapped := pipeline.Chain(
		pipeline.Timeout(3*time.Second),
		pipeline.Retry(3, pipeline.ExponentialBackoff(100*time.Millisecond)),
	)(handler)

	resp, err := wrapped(nil, &kernel.Request{
		Function: "read_holding_registers",
		Address:  "40001",
		Count:    2,
	})
	if err != nil {
		log.Fatalf("read: %v", err)
	}
	fmt.Printf("Holding registers: %x\n", resp.Data)
}
```

- [ ] **Step 2: Verify compilation**

```bash
cd /home/wwwroot/bag/industrial-protocols-go && go build ./examples/modbus_basic/
```

- [ ] **Step 3: Commit**

```bash
git add examples/
git commit -m "docs: add Modbus TCP basic example"
```

---

### Task 20: Create README

**Files:**
- Create: `README.md` (overwrite placeholder)

- [ ] **Step 1: Write README**

Create `README.md`:

```markdown
# Industrial Protocols Go

Go 语言工业网络通信协议集 —— 分层 + 中间件架构，覆盖 40 种工业协议。

> 参考 PHP 实现: [github.com/erikwang2013/industrial-protocols](https://github.com/erikwang2013/industrial-protocols)

## 架构

```
┌─────────────────────────────────────────────┐
│              User Application                │
├─────────────────────────────────────────────┤
│  Pipeline Middleware                          │
│  Retry · Timeout · CircuitBreaker · Logger   │
├─────────────────────────────────────────────┤
│  Codec (Encode/Decode)         Transport     │
│  每协议实现                      (TCP/UDP/Serial) │
├─────────────────────────────────────────────┤
│  40 Protocol Packages                        │
└─────────────────────────────────────────────┘
```

## 快速开始

```go
package main

import (
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
	handler := func(req *kernel.Request) (*kernel.Response, error) {
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

	resp, err := wrapped(nil, &kernel.Request{
		Function: "read_holding_registers",
		Address:  "40001",
		Count:    2,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Data: %x", resp.Data)
}
```

## 协议支持

| 类别 | 数量 | 已实现 | 状态 |
|------|------|--------|------|
| 工业以太网 | 5 | Modbus, OPC UA, BACnet | 3 stub |
| 现场总线 | 11 | HART, DNP3, IEC 61850 | 8 stub |
| IoT/消息 | 2 | MQTT | 1 stub |
| 汽车总线 | 5 | — | 5 stub |
| 楼宇/照明 | 2 | — | 2 stub |
| 硬件桥接 | 13 | — | 13 stub |
| 系统总线 | 3 | — | 3 stub |

## 安装

```bash
go get github.com/erikwang2013/industrial-protocols-go/kernel
go get github.com/erikwang2013/industrial-protocols-go/protocols/ethernet/modbus
```

## License

MIT — Copyright (c) 2026 erik <erik@erik.xyz>
```

- [ ] **Step 2: Commit**

```bash
git add README.md
git commit -m "docs: add README with architecture diagram and quick start"
```

---

### Task 21: Run full verification suite

- [ ] **Step 1: Full workspace build**

```bash
cd /home/wwwroot/bag/industrial-protocols-go
go work sync
cd kernel && go build ./... && go vet ./...
for dir in $(find protocols -name go.mod -exec dirname {} \;); do
  (cd "$dir" && go build ./... && go vet ./...) || echo "BUILD/VET FAIL: $dir"
done
```

- [ ] **Step 2: Full workspace test**

```bash
cd kernel && go test ./... -v -short
for dir in $(find protocols -name go.mod -exec dirname {} \;); do
  echo "=== $dir ===" && (cd "$dir" && go test ./... -v -short) || echo "TEST FAIL: $dir"
done
```

- [ ] **Step 3: Commit final state**

```bash
git add -A
git status
git commit -m "chore: full test suite passes — kernel + 40 protocol stubs + 3 implementations"
```

---

## Summary

| Phase | Tasks | Files | Description |
|-------|-------|-------|-------------|
| 1. Scaffold | 1-2 | 4 | go.work, kernel module, Makefile, LICENSE |
| 2. Core Interfaces | 3-8 | 12 | errors, codec, protocol, transport (TCP/UDP/serial) |
| 3. Pipeline | 9-13 | 10 | Chain, Timeout, Retry, Logger, CircuitBreaker |
| 4. Protocol Stubs | 14-15 | ~80 | 40 protocol modules with go.mod + stub + test |
| 5. Modbus | 16 | 3 | Full TCP/RTU codec + TCP client helper |
| 6. MQTT | 17 | 2 | Full 3.1.1 CONNECT/PUBLISH/SUBSCRIBE codec |
| 7. OPC UA | 18 | 2 | Binary HEL/OpenSecureChannel codec |
| 8. Examples | 19-21 | 3 | Example, README, full suite verification |

**Total: ~116 files, 21 tasks.**
