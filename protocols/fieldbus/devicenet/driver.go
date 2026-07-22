package devicenet

import (
	"context"
	"time"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
	"github.com/erikwang2013/industrial-protocols-go/kernel/bridge"
	"github.com/erikwang2013/industrial-protocols-go/kernel/pipeline"
)

// NewSocketCANDriver creates a DeviceNet codec paired with a SocketCAN bridge
// for the given network interface (e.g. "can0", "vcan0").
func NewSocketCANDriver(iface string, macID byte) (bridge.Bridge, kernel.Codec, error) {
	return bridge.NewCANBridge(iface), &dnCodec{macID: macID}, nil
}

// NewGatewayDriver creates a DeviceNet gateway codec paired with a TCP GatewayBridge.
func NewGatewayDriver(address string) (bridge.Bridge, kernel.Codec, error) {
	return bridge.NewGatewayBridge(address), &dnGatewayCodec{}, nil
}

// ReadyHandler returns a pipeline.Handler that starts the given bridge driver
// and provides a simple request/response loop with a 5-second timeout.
func ReadyHandler(driver func() (bridge.Bridge, kernel.Codec, error)) (pipeline.Handler, error) {
	b, codec, err := driver()
	if err != nil {
		return nil, err
	}
	if err := b.Start(); err != nil {
		return nil, err
	}
	tr := b.Transport()
	return pipeline.Chain(
		pipeline.Timeout(5 * time.Second),
	)(func(ctx context.Context, req *kernel.Request) (*kernel.Response, error) {
		raw, _ := codec.Encode(req)
		tr.Write(raw)
		buf := make([]byte, 64)
		n, _ := tr.Read(buf)
		return codec.Decode(buf[:n])
	}), nil
}
