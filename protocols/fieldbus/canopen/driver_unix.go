//go:build !linux

package canopen

import (
	"fmt"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
	"github.com/erikwang2013/industrial-protocols-go/kernel/bridge"
	"github.com/erikwang2013/industrial-protocols-go/kernel/pipeline"
)

// NewSocketCANDriver returns an error on non-Linux platforms since SocketCAN is Linux-only.
func NewSocketCANDriver(iface string, nodeID byte) (bridge.Bridge, kernel.Codec, error) {
	return nil, nil, fmt.Errorf("canopen: SocketCAN requires Linux")
}

// ReadyHandler returns an error on non-Linux platforms.
func ReadyHandler(iface string, nodeID byte) (pipeline.Handler, error) {
	return nil, fmt.Errorf("canopen: SocketCAN requires Linux")
}
