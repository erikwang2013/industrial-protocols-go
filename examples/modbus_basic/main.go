// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package main

import (
	"context"
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

	handler := func(ctx context.Context, req *kernel.Request) (*kernel.Response, error) {
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

	resp, err := wrapped(context.Background(), &kernel.Request{
		Function: "read_holding_registers",
		Address:  "40001",
		Count:    2,
	})
	if err != nil {
		log.Fatalf("read: %v", err)
	}
	fmt.Printf("Holding registers: %x\n", resp.Data)
}
