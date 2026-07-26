// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

//go:build linux

package canopen

import (
	"context"
	"time"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
	"github.com/erikwang2013/industrial-protocols-go/kernel/bridge"
	"github.com/erikwang2013/industrial-protocols-go/kernel/pipeline"
)

// NewSocketCANDriver creates a CANopen codec paired with a SocketCAN bridge
// for the given network interface (e.g. "can0", "vcan0").
func NewSocketCANDriver(iface string, nodeID byte) (bridge.Bridge, kernel.Codec, error) {
	return bridge.NewCANBridge(iface), &canopenCodec{nodeID: nodeID}, nil
}

// ReadyHandler returns a pipeline.Handler that starts the CAN bridge
// and provides a simple request/response loop with a 5-second timeout.
func ReadyHandler(iface string, nodeID byte) (pipeline.Handler, error) {
	b := bridge.NewCANBridge(iface)
	c := &canopenCodec{nodeID: nodeID}
	if err := b.Start(); err != nil {
		return nil, err
	}
	tr := b.Transport()
	return pipeline.Chain(
		pipeline.Timeout(5 * time.Second),
	)(func(ctx context.Context, req *kernel.Request) (*kernel.Response, error) {
			raw, err := c.Encode(req)
			if err != nil {
				return nil, err
			}
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
