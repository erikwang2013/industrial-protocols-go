package powerlink

import (
	"github.com/erikwang2013/industrial-protocols-go/kernel"
	"github.com/erikwang2013/industrial-protocols-go/kernel/bridge"
)

// NewCmdDriver creates a POWERLINK codec paired with a CmdBridge
// wrapping the openPOWERLINK_demo command-line utility.
func NewCmdDriver(toolPath string) (bridge.Bridge, kernel.Codec, error) {
	b := bridge.NewCmdBridge(toolPath, "--iface", "eth0", "--read")
	return b, &powerlinkCodec{}, nil
}
