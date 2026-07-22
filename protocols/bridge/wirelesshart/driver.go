package wirelesshart

import (
	"github.com/erikwang2013/industrial-protocols-go/kernel"
	"github.com/erikwang2013/industrial-protocols-go/kernel/bridge"
)

// NewCmdDriver creates a WirelessHART codec paired with a CmdBridge
// wrapping the emerson_1410_cli command-line utility.
func NewCmdDriver(toolPath string) (bridge.Bridge, kernel.Codec, error) {
	b := bridge.NewCmdBridge(toolPath, "--read", "--format", "hex")
	return b, &wirelesshartCodec{}, nil
}
