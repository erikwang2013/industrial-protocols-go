// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package cpci

import (
	"fmt"
	"os"

	"github.com/erikwang2013/industrial-protocols-go/kernel/transport"
)

// CPCIDriver opens a CompactPCI device config space via sysfs and wraps
// it as a transport.Transport for use with the CPCI protocol codec.
//
// CompactPCI uses the same PCI sysfs interface as conventional PCI.
// The busAddr is the BDF (bus:device.function) identifier, e.g. "0000:02:00.0".
// For CompactPCI systems, the bus number identifies the CPCI segment.
type CPCIDriver struct {
	dev     *os.File
	busAddr string
}

// NewCPCIDriver opens the sysfs config space for the given CPCI BDF address.
// Requires root privileges or appropriate capabilities (CAP_SYS_ADMIN).
func NewCPCIDriver(busAddr string) (*CPCIDriver, error) {
	path := fmt.Sprintf("/sys/bus/pci/devices/%s/config", busAddr)
	f, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		return nil, fmt.Errorf("cpci: open %s: %w", path, err)
	}
	return &CPCIDriver{dev: f, busAddr: busAddr}, nil
}

// Transport returns a PipeTransport backed by the CPCI config space file.
func (d *CPCIDriver) Transport() transport.Transport {
	return transport.NewPipeTransport(d.dev, d.dev, d.busAddr)
}

// Close closes the CPCI config space file.
func (d *CPCIDriver) Close() error { return d.dev.Close() }
