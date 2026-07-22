// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package most

import (
	"github.com/erikwang2013/industrial-protocols-go/kernel"
	"github.com/erikwang2013/industrial-protocols-go/kernel/bridge"
)

// NewSerialDriver creates a MOST codec paired with a SerialBridge
// for the given serial device (e.g. "/dev/ttyUSB0") at 115200 baud.
func NewSerialDriver(dev string) (bridge.Bridge, kernel.Codec, error) {
	b := bridge.NewSerialBridge(dev, 115200)
	return b, &mostCodec{}, nil
}
