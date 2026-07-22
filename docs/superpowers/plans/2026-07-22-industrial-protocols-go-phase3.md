# Industrial Protocols Go — Phase 3 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add 3 new bridge types (Gateway, CAN, Serial) to kernel/bridge, then upgrade 26 stub protocols to hardware SDKs with codec + driver.

**Architecture:** Extend `kernel/bridge/` with `GatewayBridge`, `CANBridge`, `SerialBridge`. Each protocol gets a `driver.go` factory + protocol codec. All follow the unified `ReadyHandler` pattern.

**Tech Stack:** Go 1.20+, stdlib; Linux SocketCAN (build tag `linux` for CANBridge).

---

## Phase 3A: Bridge Module Extensions

### Task 1: Implement GatewayBridge

**Files:** Create `kernel/bridge/gateway.go`, `kernel/bridge/gateway_test.go`

- [ ] **Step 1: Implement + test + commit**

```go
// gateway.go
package bridge

import (
	"net"

	"github.com/erikwang2013/industrial-protocols-go/kernel/transport"
)

type GatewayBridge struct {
	conn    net.Conn
	addr    string
	running bool
}

func NewGatewayBridge(address string) *GatewayBridge {
	return &GatewayBridge{addr: address}
}

func (b *GatewayBridge) Start() error {
	conn, err := net.Dial("tcp", b.addr)
	if err != nil {
		return err
	}
	b.conn = conn
	b.running = true
	return nil
}

func (b *GatewayBridge) Stop() error {
	b.running = false
	if b.conn != nil {
		return b.conn.Close()
	}
	return nil
}

func (b *GatewayBridge) Transport() transport.Transport {
	if b.conn == nil {
		return nil
	}
	conn := b.conn
	return &gatewayTransport{conn: conn, addr: b.addr}
}

type gatewayTransport struct {
	conn net.Conn
	addr string
}

func (t *gatewayTransport) Read(p []byte) (int, error)  { return t.conn.Read(p) }
func (t *gatewayTransport) Write(p []byte) (int, error) { return t.conn.Write(p) }
func (t *gatewayTransport) Close() error                { return t.conn.Close() }
func (t *gatewayTransport) Addr() string                { return t.addr }
func (t *gatewayTransport) Alive() bool                 { return true }
```

```go
// gateway_test.go
package bridge

import (
	"net"
	"testing"
)

func TestGatewayBridgeEcho(t *testing.T) {
	ln, _ := net.Listen("tcp", "127.0.0.1:0")
	defer ln.Close()
	go func() {
		conn, _ := ln.Accept()
		buf := make([]byte, 64)
		n, _ := conn.Read(buf)
		conn.Write(buf[:n])
	}()

	b := NewGatewayBridge(ln.Addr().String())
	if err := b.Start(); err != nil {
		t.Fatal(err)
	}
	defer b.Stop()

	tr := b.Transport()
	tr.Write([]byte("test"))
	buf := make([]byte, 64)
	n, _ := tr.Read(buf)
	if string(buf[:n]) != "test" {
		t.Errorf("expected 'test', got '%s'", string(buf[:n]))
	}
}
```

```bash
cd kernel && go test ./bridge/ -v -short
git add kernel/bridge/gateway.go kernel/bridge/gateway_test.go
git commit -m "feat(bridge): add GatewayBridge for TCP gateway hardware"
```

---

### Task 2: Implement CANBridge (Linux only)

**Files:** Create `kernel/bridge/can.go`, `kernel/bridge/can_test.go`

- [ ] **Step 1: Implement + test + commit**

```go
// can.go (build tag: linux)
//go:build linux

package bridge

import (
	"encoding/binary"
	"fmt"
	"syscall"
	"unsafe"

	"github.com/erikwang2013/industrial-protocols-go/kernel/transport"
)

type CANFrame struct {
	ID   uint32
	Data []byte
	Ext  bool
}

type CANBridge struct {
	iface   string
	fd      int
	running bool
}

func NewCANBridge(iface string) *CANBridge {
	return &CANBridge{iface: iface}
}

func (b *CANBridge) Start() error {
	fd, err := syscall.Socket(syscall.AF_CAN, syscall.SOCK_RAW, syscall.CAN_RAW)
	if err != nil {
		return fmt.Errorf("can: socket: %w", err)
	}
	// Bind to interface
	var addr [16]byte
	copy(addr[:], b.iface)
	if _, _, err := syscall.Syscall(syscall.SYS_BIND, uintptr(fd), uintptr(unsafe.Pointer(&addr[0])), 0); err != 0 {
		return fmt.Errorf("can: bind %s: %w", b.iface, err)
	}
	b.fd = fd
	b.running = true
	return nil
}

func (b *CANBridge) Stop() error {
	b.running = false
	if b.fd > 0 {
		return syscall.Close(b.fd)
	}
	return nil
}

func (b *CANBridge) Transport() transport.Transport {
	return &canTransport{fd: b.fd, iface: b.iface}
}

type canTransport struct {
	fd    int
	iface string
}

func (t *canTransport) Read(p []byte) (int, error) {
	n, err := syscall.Read(t.fd, p)
	return n, err
}

func (t *canTransport) Write(p []byte) (int, error) {
	n, err := syscall.Write(t.fd, p)
	return n, err
}

func (t *canTransport) Close() error                { return syscall.Close(t.fd) }
func (t *canTransport) Addr() string                { return t.iface }
func (t *canTransport) Alive() bool                 { return t.fd > 0 }
```

```go
// can_test.go
//go:build linux

package bridge

import "testing"

func TestNewCANBridge(t *testing.T) {
	b := NewCANBridge("vcan0")
	if b.iface != "vcan0" {
		t.Errorf("expected vcan0, got %s", b.iface)
	}
	if b.running {
		t.Error("not running before Start()")
	}
}

func TestCANFrame(t *testing.T) {
	f := CANFrame{ID: 0x123, Data: []byte{1, 2, 3}, Ext: true}
	if f.ID != 0x123 {
		t.Errorf("expected 0x123, got 0x%X", f.ID)
	}
}
```

```bash
cd kernel && go test ./bridge/ -v -short
git add kernel/bridge/can.go kernel/bridge/can_test.go
git commit -m "feat(bridge): add CANBridge for Linux SocketCAN (build tag linux)"
```

---

### Task 3: Implement SerialBridge

**Files:** Create `kernel/bridge/serial.go`, `kernel/bridge/serial_test.go`

- [ ] **Step 1: Implement + test + commit**

```go
// serial.go
package bridge

import (
	"io"

	"github.com/erikwang2013/industrial-protocols-go/kernel/transport"
)

type SerialBridge struct {
	rw      io.ReadWriteCloser
	dev     string
	baud    int
	running bool
}

func NewSerialBridge(dev string, baud int) *SerialBridge {
	return &SerialBridge{dev: dev, baud: baud}
}

func (b *SerialBridge) Start() error {
	b.running = true
	return nil
}

func (b *SerialBridge) Stop() error {
	b.running = false
	if b.rw != nil {
		return b.rw.Close()
	}
	return nil
}

func (b *SerialBridge) Transport() transport.Transport {
	if b.rw == nil {
		return nil
	}
	return transport.NewPipeTransport(b.rw, b.rw, b.dev)
}

func (b *SerialBridge) SetPort(rw io.ReadWriteCloser) {
	b.rw = rw
}
```

```go
// serial_test.go
package bridge

import (
	"bytes"
	"testing"
)

type mockSerial struct{ bytes.Buffer }

func (m *mockSerial) Close() error { return nil }

func TestSerialBridgeMock(t *testing.T) {
	b := NewSerialBridge("/dev/ttyUSB0", 115200)
	if b.running { t.Error("not running before Start()") }

	rw := &mockSerial{}
	b.SetPort(rw)
	if err := b.Start(); err != nil { t.Fatal(err) }

	tr := b.Transport()
	if tr == nil { t.Fatal("transport is nil") }
	tr.Write([]byte("AT\r"))
	b.Stop()
}
```

```bash
cd kernel && go test ./bridge/ -v -short
git add kernel/bridge/serial.go kernel/bridge/serial_test.go
git commit -m "feat(bridge): add SerialBridge for serial-port hardware"
```

---

## Phase 3B: CAN Protocol SDKs (3)

### Task 4: CANopen SDK

**Files:** Rewrite `protocols/fieldbus/canopen/`

- [ ] **Step 1: Implement codec (CANopen over CAN: NMT/SDO/PDO framing with 11-bit IDs, node addressing) + SocketCAN driver + test + README + commit**

```bash
cd protocols/fieldbus/canopen && go test ./... -v
git add protocols/fieldbus/canopen/
git commit -m "feat(protocols): implement CANopen SDK with SocketCAN driver"
```

---

### Task 5: DeviceNet SDK

**Files:** Rewrite `protocols/fieldbus/devicenet/`

- [ ] **Step 1: Implement DeviceNet over CAN codec + GatewayBridge driver (Anybus Scanner) + test + README + commit**

```bash
cd protocols/fieldbus/devicenet && go test ./... -v
git add protocols/fieldbus/devicenet/
git commit -m "feat(protocols): implement DeviceNet SDK with Anybus gateway driver"
```

---

### Task 6: FlexRay CAN + SAE J1850 CAN SDKs

**Files:** Rewrite `protocols/automotive/flexray/` and `protocols/automotive/saej1850/`

- [ ] **Step 1: Implement both CAN-based codecs + SocketCAN drivers + tests + READMEs + commit**

```bash
cd protocols/automotive/flexray && go test ./... -v
cd protocols/automotive/saej1850 && go test ./... -v
git add protocols/automotive/flexray/ protocols/automotive/saej1850/
git commit -m "feat(protocols): implement FlexRay CAN and SAE J1850 CAN SDKs"
```

---

## Phase 3C: GatewayBridge Protocol SDKs (10)

### Task 7: Ethernet/Fieldbus Gateway SDKs (5)

**Files:** Rewrite `protocols/fieldbus/profibus/`, `protocols/fieldbus/cclinkie/`, `protocols/bridge/interbus/`, `protocols/bridge/worldfip/`, `protocols/bridge/lightbus/`

- [ ] **Step 1: Implement each with GatewayBridge driver + codec + test + README + commit**

```bash
for dir in protocols/fieldbus/profibus protocols/fieldbus/cclinkie protocols/bridge/interbus protocols/bridge/worldfip protocols/bridge/lightbus; do
  cd "$dir" && go test ./... -v
done
git add protocols/fieldbus/profibus/ protocols/fieldbus/cclinkie/ protocols/bridge/interbus/ protocols/bridge/worldfip/ protocols/bridge/lightbus/
git commit -m "feat(protocols): add GatewayBridge SDKs (PROFIBUS, CC-LinkIE, Interbus, WorldFIP, Lightbus)"
```

---

### Task 8: Fieldbus Gateway SDKs (3)

**Files:** Rewrite `protocols/fieldbus/foundationfieldbus/`, `protocols/fieldbus/asinterface/`, `protocols/fieldbus/iolink/`

- [ ] **Step 1: Implement each with GatewayBridge driver + codec + test + README + commit**

```bash
for dir in protocols/fieldbus/foundationfieldbus protocols/fieldbus/asinterface protocols/fieldbus/iolink; do
  cd "$dir" && go test ./... -v
done
git add protocols/fieldbus/foundationfieldbus/ protocols/fieldbus/asinterface/ protocols/fieldbus/iolink/
git commit -m "feat(protocols): add GatewayBridge SDKs (Foundation Fieldbus, AS-Interface, IO-Link)"
```

---

### Task 9: LonWorks + ModbusPlus Gateway SDKs (2)

**Files:** Rewrite `protocols/building/lonworks/`, `protocols/bridge/modbusplus/`

- [ ] **Step 1: Implement both with GatewayBridge driver + codec + test + README + commit**

```bash
cd protocols/building/lonworks && go test ./... -v
cd protocols/bridge/modbusplus && go test ./... -v
git add protocols/building/lonworks/ protocols/bridge/modbusplus/
git commit -m "feat(protocols): add GatewayBridge SDKs (LonWorks, ModbusPlus)"
```

---

## Phase 3D: SerialBridge + CmdBridge Protocol SDKs (10)

### Task 10: SerialBridge SDKs (2)

**Files:** Rewrite `protocols/automotive/flexray/` (serial variant as additional driver), `protocols/automotive/most/`

- [ ] **Step 1: Implement FlexRay serial driver + MOST serial codec + tests + READMEs + commit**

```bash
cd protocols/automotive/flexray && go test ./... -v
cd protocols/automotive/most && go test ./... -v
git add protocols/automotive/flexray/ protocols/automotive/most/
git commit -m "feat(protocols): add SerialBridge SDKs (FlexRay serial, MOST)"
```

---

### Task 11: CmdBridge SDKs Batch 1 (4)

**Files:** Rewrite `protocols/bridge/ethercat/`, `protocols/bridge/powerlink/`, `protocols/bridge/sercos/`, `protocols/bridge/sercos1/`

- [ ] **Step 1: Implement each with CmdBridge driver + codec + test + README + commit**

```bash
for dir in protocols/bridge/ethercat protocols/bridge/powerlink protocols/bridge/sercos protocols/bridge/sercos1; do
  cd "$dir" && go test ./... -v
done
git add protocols/bridge/ethercat/ protocols/bridge/powerlink/ protocols/bridge/sercos/ protocols/bridge/sercos1/
git commit -m "feat(protocols): add CmdBridge SDKs (EtherCAT, POWERLINK, SERCOS III, SERCOS I/II)"
```

---

### Task 12: CmdBridge SDKs Batch 2 (4)

**Files:** Rewrite `protocols/bridge/controlnet/`, `protocols/bridge/isa100/`, `protocols/bridge/wirelesshart/`

- [ ] **Step 1: Implement each with CmdBridge driver + codec + test + README + commit**

```bash
for dir in protocols/bridge/controlnet protocols/bridge/isa100 protocols/bridge/wirelesshart; do
  cd "$dir" && go test ./... -v
done
git add protocols/bridge/controlnet/ protocols/bridge/isa100/ protocols/bridge/wirelesshart/
git commit -m "feat(protocols): add CmdBridge SDKs (ControlNet, ISA100, WirelessHART)"
```

---

## Phase 3E: SystemBus SDKs (3)

### Task 13: PCI, VME, CompactPCI

**Files:** Rewrite `protocols/system/pci/`, `protocols/system/vme/`, `protocols/system/cpci/`

- [ ] **Step 1: Implement each with sysfs/procfs driver + codec + test + README + commit**

```bash
for dir in protocols/system/pci protocols/system/vme protocols/system/cpci; do
  cd "$dir" && go test ./... -v
done
git add protocols/system/
git commit -m "feat(protocols): add SystemBus SDKs (PCI/PCIe, VME/VPX, CompactPCI)"
```

---

## Phase 3F: Final Verification

### Task 14: Full workspace verification + README update

- [ ] **Step 1: Build, vet, test all**

```bash
cd /home/wwwroot/bag/industrial-protocols-go
go work sync
cd kernel && go build ./... && go vet ./...
for dir in $(find protocols -name go.mod -exec dirname {} \;); do
  (cd "$dir" && go build ./... && go vet ./... && go test ./... -v -short) || echo "FAIL: $dir"
done
```

- [ ] **Step 2: Update README protocol tables (all 40 now implemented), commit**

```bash
git add README.md
git commit -m "chore: Phase 3 complete — 26 hardware SDKs, all 40 protocols with drivers"
```

---

## Summary

| Phase | Tasks | Files (est.) |
|-------|-------|-------------|
| 3A: Bridges | 3 | 6 |
| 3B: CAN SDKs | 3 | ~12 |
| 3C: Gateway SDKs | 3 | ~50 |
| 3D: Serial+Cmd SDKs | 3 | ~20 |
| 3E: SystemBus | 1 | ~6 |
| 3F: Verify | 1 | — |
| **Total** | **14 tasks** | **~95 files** |
