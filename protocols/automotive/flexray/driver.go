package flexray

import (
	"context"
	"time"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
	"github.com/erikwang2013/industrial-protocols-go/kernel/bridge"
	"github.com/erikwang2013/industrial-protocols-go/kernel/pipeline"
)

// NewSocketCANDriver creates a FlexRay codec paired with a SocketCAN bridge
// for the given network interface (e.g. "can0").
func NewSocketCANDriver(iface string) (bridge.Bridge, kernel.Codec, error) {
	return bridge.NewCANBridge(iface), &flexrayCodec{slotID: 1, baseCycle: 0}, nil
}

// ReadyHandler returns a pipeline.Handler that starts the CAN bridge
// and provides a simple request/response loop with a 5-second timeout.
func ReadyHandler(iface string) (pipeline.Handler, error) {
	b := bridge.NewCANBridge(iface)
	c := &flexrayCodec{slotID: 1, baseCycle: 0}
	if err := b.Start(); err != nil {
		return nil, err
	}
	tr := b.Transport()
	return pipeline.Chain(
		pipeline.Timeout(5 * time.Second),
	)(func(ctx context.Context, req *kernel.Request) (*kernel.Response, error) {
		raw, _ := c.Encode(req)
		tr.Write(raw)
		buf := make([]byte, 64)
		n, _ := tr.Read(buf)
		return c.Decode(buf[:n])
	}), nil
}
