// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package modbus

import (
	"context"
	"fmt"
	"time"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
	"github.com/erikwang2013/industrial-protocols-go/kernel/pipeline"
	"github.com/erikwang2013/industrial-protocols-go/kernel/transport"
)

// TCPClient is a high-level Modbus TCP client with optional middleware.
type TCPClient struct {
	handle pipeline.Handler
	tr     transport.Transport
}

// NewTCPClient dials a Modbus TCP server and returns a client with the given unit ID.
func NewTCPClient(address string, unitID byte) (*TCPClient, error) {
	tr, err := transport.DialTCP(address)
	if err != nil {
		return nil, err
	}
	codec, _ := NewCodec("tcp")
	handle := func(ctx context.Context, req *kernel.Request) (*kernel.Response, error) {
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

// WithTimeout wraps the handler with a timeout middleware.
func (c *TCPClient) WithTimeout(d time.Duration) *TCPClient {
	c.handle = pipeline.Chain(pipeline.Timeout(d))(c.handle)
	return c
}

// WithRetry wraps the handler with a retry middleware.
func (c *TCPClient) WithRetry(max int, backoff pipeline.BackoffFunc) *TCPClient {
	c.handle = pipeline.Chain(pipeline.Retry(max, backoff))(c.handle)
	return c
}

// Close shuts down the underlying transport.
func (c *TCPClient) Close() error {
	return c.tr.Close()
}

// ReadCoils reads coil status from the device.
func (c *TCPClient) ReadCoils(address uint16, count uint16) ([]bool, error) {
	req := &kernel.Request{Function: "read_coils", Address: fmt.Sprintf("%d", address), Count: int(count)}
	resp, err := c.handle(context.Background(), req)
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

// ReadHoldingRegisters reads holding registers from the device.
func (c *TCPClient) ReadHoldingRegisters(address uint16, count uint16) ([]uint16, error) {
	req := &kernel.Request{Function: "read_holding_registers", Address: fmt.Sprintf("%d", address), Count: int(count)}
	resp, err := c.handle(context.Background(), req)
	if err != nil {
		return nil, err
	}
	regs := make([]uint16, count)
	for i := 0; i < len(resp.Data)/2 && i < int(count); i++ {
		regs[i] = uint16(resp.Data[i*2])<<8 | uint16(resp.Data[i*2+1])
	}
	return regs, nil
}
