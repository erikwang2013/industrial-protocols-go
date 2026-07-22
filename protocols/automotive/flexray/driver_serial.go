package flexray

import (
	"github.com/erikwang2013/industrial-protocols-go/kernel"
	"github.com/erikwang2013/industrial-protocols-go/kernel/bridge"
)

// NewSerialDriver creates a FlexRay codec paired with a SerialBridge
// for the given serial device at 115200 baud.
func NewSerialDriver(dev string) (bridge.Bridge, kernel.Codec, error) {
	b := bridge.NewSerialBridge(dev, 115200)
	return b, &flexrayCodec{slotID: 1, baseCycle: 0}, nil
}
