// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package modbusplus

import (
	"context"
	"time"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
	"github.com/erikwang2013/industrial-protocols-go/kernel/bridge"
	"github.com/erikwang2013/industrial-protocols-go/kernel/pipeline"
)

// NewGatewayDriver creates a Modbus Plus GatewayBridge + codec for the given
// TCP address (e.g. "192.168.1.100:2010").
func NewGatewayDriver(address string) (bridge.Bridge, kernel.Codec, error) {
	return bridge.NewGatewayBridge(address), &modbusplusCodec{}, nil
}

// NewCmdDriver creates a Modbus Plus codec paired with a CmdBridge
// wrapping the sa85_cli command-line utility.
func NewCmdDriver(toolPath string) (bridge.Bridge, kernel.Codec, error) {
	b := bridge.NewCmdBridge(toolPath, "--node", "1", "--read")
	return b, &modbusplusCodec{}, nil
}

// ReadyHandler returns a pipeline.Handler that starts the bridge and provides
// a request/response loop with a 5-second timeout and 3 retries.
func ReadyHandler(address string) (pipeline.Handler, error) {
	b := bridge.NewGatewayBridge(address)
	codec := &modbusplusCodec{}
	if err := b.Start(); err != nil {
		return nil, err
	}
	tr := b.Transport()
	return pipeline.Chain(
		pipeline.Timeout(5*time.Second),
		pipeline.Retry(3, pipeline.LinearBackoff(100*time.Millisecond)),
	)(func(ctx context.Context, req *kernel.Request) (*kernel.Response, error) {
		raw, err := codec.Encode(req)
		if err != nil {
			return nil, err
		}
		tr.Write(raw)
		buf := make([]byte, 4096)
		n, _ := tr.Read(buf)
		return codec.Decode(buf[:n])
	}), nil
}
