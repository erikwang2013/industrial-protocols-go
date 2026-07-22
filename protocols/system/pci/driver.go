// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package pci

import (
	"fmt"
	"os"

	"github.com/erikwang2013/industrial-protocols-go/kernel/transport"
)

// PCIDriver opens a PCI device config space via sysfs and wraps it
// as a transport.Transport for use with the PCI protocol codec.
//
// The busAddr is the BDF (bus:device.function) identifier, e.g. "0000:00:1f.3".
// The driver opens /sys/bus/pci/devices/<busAddr>/config in read-write mode.
type PCIDriver struct {
	dev     *os.File
	busAddr string
}

// NewPCIDriver opens the sysfs config space for the given PCI BDF address.
// Requires root privileges or appropriate capabilities (CAP_SYS_ADMIN).
func NewPCIDriver(busAddr string) (*PCIDriver, error) {
	path := fmt.Sprintf("/sys/bus/pci/devices/%s/config", busAddr)
	f, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		return nil, fmt.Errorf("pci: open %s: %w", path, err)
	}
	return &PCIDriver{dev: f, busAddr: busAddr}, nil
}

// Transport returns a PipeTransport backed by the PCI config space file.
func (d *PCIDriver) Transport() transport.Transport {
	return transport.NewPipeTransport(d.dev, d.dev, d.busAddr)
}

// Close closes the PCI config space file.
func (d *PCIDriver) Close() error { return d.dev.Close() }
