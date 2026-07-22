// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package gateway

import (
	"context"
	"fmt"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
	"github.com/erikwang2013/industrial-protocols-go/kernel/transport"
)

type GatewayEngine struct {
	rules []Rule
}

func New() *GatewayEngine { return &GatewayEngine{} }

func (e *GatewayEngine) Add(rule Rule) { e.rules = append(e.rules, rule) }

func (e *GatewayEngine) Transform(ctx context.Context, srcTransport, dstTransport transport.Transport, req *kernel.Request) (*kernel.Response, error) {
	for _, rule := range e.rules {
		srcName := rule.Src.Name()
		srcCodec, err := rule.Src.NewCodec(srcName)
		if err != nil {
			continue
		}

		srcRaw, err := srcCodec.Encode(req)
		if err != nil {
			return nil, fmt.Errorf("gateway: src encode: %w", err)
		}
		if _, err := srcTransport.Write(srcRaw); err != nil {
			return nil, fmt.Errorf("gateway: src write: %w", err)
		}
		buf := make([]byte, 4096)
		n, err := srcTransport.Read(buf)
		if err != nil {
			return nil, fmt.Errorf("gateway: src read: %w", err)
		}
		srcResp, err := srcCodec.Decode(buf[:n])
		if err != nil {
			return nil, fmt.Errorf("gateway: src decode: %w", err)
		}

		dstName := rule.Dst.Name()
		dstCodec, err := rule.Dst.NewCodec(dstName)
		if err != nil {
			return nil, fmt.Errorf("gateway: dst codec: %w", err)
		}
		dstReq := rule.Map(srcResp)
		dstRaw, err := dstCodec.Encode(dstReq)
		if err != nil {
			return nil, fmt.Errorf("gateway: dst encode: %w", err)
		}
		if _, err := dstTransport.Write(dstRaw); err != nil {
			return nil, fmt.Errorf("gateway: dst write: %w", err)
		}
		n2, err := dstTransport.Read(buf)
		if err != nil {
			return nil, fmt.Errorf("gateway: dst read: %w", err)
		}
		return dstCodec.Decode(buf[:n2])
	}
	return nil, fmt.Errorf("gateway: no matching rule found")
}
