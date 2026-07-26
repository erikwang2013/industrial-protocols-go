// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package profibus

import (
	"context"
	"time"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
	"github.com/erikwang2013/industrial-protocols-go/kernel/bridge"
	"github.com/erikwang2013/industrial-protocols-go/kernel/pipeline"
)

// NewGatewayDriver creates a PROFIBUS GatewayBridge + codec for the given
// TCP address (e.g. "192.168.1.100:2000").
func NewGatewayDriver(address string) (bridge.Bridge, kernel.Codec, error) {
	return bridge.NewGatewayBridge(address), &profibusCodec{}, nil
}

// ReadyHandler returns a pipeline.Handler that starts the bridge and provides
// a request/response loop with a 5-second timeout and 3 retries.
func ReadyHandler(address string) (pipeline.Handler, error) {
	b := bridge.NewGatewayBridge(address)
	codec := &profibusCodec{}
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
		if _, err := tr.Write(raw); err != nil {
		return nil, err
	}
		buf := make([]byte, 4096)
		n, err := tr.Read(buf)
	if err != nil {
		return nil, err
	}
		return codec.Decode(buf[:n])
	}), nil
}
