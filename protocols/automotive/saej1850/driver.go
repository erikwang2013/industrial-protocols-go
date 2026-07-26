// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package saej1850

import (
	"context"
	"time"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
	"github.com/erikwang2013/industrial-protocols-go/kernel/bridge"
	"github.com/erikwang2013/industrial-protocols-go/kernel/pipeline"
)

// NewSocketCANDriver creates a SAE J1850 codec paired with a SocketCAN bridge
// for the given network interface (e.g. "can0").
func NewSocketCANDriver(iface string) (bridge.Bridge, kernel.Codec, error) {
	return bridge.NewCANBridge(iface), &j1850Codec{
		targetAddr: 0x01,
		sourceAddr: 0xF1,
	}, nil
}

// ReadyHandler returns a pipeline.Handler that starts the CAN bridge
// and provides a simple request/response loop with a 5-second timeout.
func ReadyHandler(iface string) (pipeline.Handler, error) {
	b := bridge.NewCANBridge(iface)
	c := &j1850Codec{targetAddr: 0x01, sourceAddr: 0xF1}
	if err := b.Start(); err != nil {
		return nil, err
	}
	tr := b.Transport()
	return pipeline.Chain(
		pipeline.Timeout(5 * time.Second),
	)(func(ctx context.Context, req *kernel.Request) (*kernel.Response, error) {
		raw, _ := c.Encode(req)
		if _, err := tr.Write(raw); err != nil {
		return nil, err
	}
		buf := make([]byte, 64)
		n, err := tr.Read(buf)
	if err != nil {
		return nil, err
	}
		return c.Decode(buf[:n])
	}), nil
}
