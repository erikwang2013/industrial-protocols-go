// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package ethercat

import (
	"github.com/erikwang2013/industrial-protocols-go/kernel"
	"github.com/erikwang2013/industrial-protocols-go/kernel/bridge"
)

// NewCmdDriver creates an EtherCAT codec paired with a CmdBridge
// wrapping the ethercat command-line utility.
func NewCmdDriver(toolPath string) (bridge.Bridge, kernel.Codec, error) {
	b := bridge.NewCmdBridge(toolPath, "slaves", "--format", "hex")
	return b, &ethercatCodec{}, nil
}
