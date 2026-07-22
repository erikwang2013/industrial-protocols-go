// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package bridge

import (
	"io"
	"os/exec"

	"github.com/erikwang2013/industrial-protocols-go/kernel/transport"
)

// Bridge 是最小化的硬件桥接接口。
// 每个硬件协议 SDK 实现具体的 Bridge（CmdBridge、GatewayBridge、CANBridge 等）。
type Bridge interface {
	Start() error
	Stop() error
	Transport() transport.Transport
}

// CmdBridge 包装一个外部命令行工具，通过 stdin/stdout 与硬件通信。
type CmdBridge struct {
	cmd     *exec.Cmd
	pipe    transport.Transport
	running bool
}

// NewCmdBridge 创建一个命令行桥接器。
func NewCmdBridge(name string, args ...string) *CmdBridge {
	return &CmdBridge{cmd: exec.Command(name, args...)}
}

// Start launches the underlying command and wires its stdin/stdout to a PipeTransport.
func (b *CmdBridge) Start() error {
	r, w := io.Pipe()
	r2, w2 := io.Pipe()
	b.cmd.Stdin = r
	b.cmd.Stdout = w2
	b.cmd.Stderr = nil
	if err := b.cmd.Start(); err != nil {
		return err
	}
	b.pipe = transport.NewPipeTransport(r2, w, b.cmd.Path)
	b.running = true
	return nil
}

// Stop closes the transport and kills the underlying process.
func (b *CmdBridge) Stop() error {
	b.running = false
	if b.pipe != nil {
		b.pipe.Close()
	}
	if b.cmd != nil && b.cmd.Process != nil {
		return b.cmd.Process.Kill()
	}
	return nil
}

// Transport returns the underlying PipeTransport. It returns nil before Start.
func (b *CmdBridge) Transport() transport.Transport { return b.pipe }
