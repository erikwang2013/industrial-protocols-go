package isa100

import (
	"github.com/erikwang2013/industrial-protocols-go/kernel"
	"github.com/erikwang2013/industrial-protocols-go/kernel/bridge"
)

// NewCmdDriver creates an ISA100 codec paired with a CmdBridge
// wrapping the yfgw410_cli command-line utility.
func NewCmdDriver(toolPath string) (bridge.Bridge, kernel.Codec, error) {
	b := bridge.NewCmdBridge(toolPath, "--read", "--format", "hex")
	return b, &isa100Codec{}, nil
}
