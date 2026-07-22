package sercos1

import (
	"github.com/erikwang2013/industrial-protocols-go/kernel"
	"github.com/erikwang2013/industrial-protocols-go/kernel/bridge"
)

// NewCmdDriver creates a SERCOS I/II codec paired with a CmdBridge
// wrapping the sercos_cli command-line utility.
func NewCmdDriver(toolPath string) (bridge.Bridge, kernel.Codec, error) {
	b := bridge.NewCmdBridge(toolPath, "--fiber", "--read")
	return b, &sercos1Codec{}, nil
}
