package controlnet

import (
	"github.com/erikwang2013/industrial-protocols-go/kernel"
	"github.com/erikwang2013/industrial-protocols-go/kernel/bridge"
)

// NewCmdDriver creates a ControlNet codec paired with a CmdBridge
// wrapping the 1784-pcic-cli command-line utility.
func NewCmdDriver(toolPath string) (bridge.Bridge, kernel.Codec, error) {
	b := bridge.NewCmdBridge(toolPath, "--slot", "0", "--read")
	return b, &controlnetCodec{}, nil
}
