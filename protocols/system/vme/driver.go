// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package vme

import (
	"fmt"
	"os"

	"github.com/erikwang2013/industrial-protocols-go/kernel/transport"
)

// VMEDriver opens a VME bus slot via procfs and wraps it as a
// transport.Transport for use with the VME protocol codec.
//
// The slot is the VME slot number (0-20). The driver opens
// /proc/vme/<slot> in read-write mode.
type VMEDriver struct {
	dev  *os.File
	slot string
}

// NewVMEDriver opens the procfs device node for the given VME slot.
// Requires root privileges and the vme_tsi148 kernel module loaded.
func NewVMEDriver(slot int) (*VMEDriver, error) {
	path := fmt.Sprintf("/proc/vme/%d", slot)
	f, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		return nil, fmt.Errorf("vme: open %s: %w", path, err)
	}
	return &VMEDriver{dev: f, slot: fmt.Sprint(slot)}, nil
}

// Transport returns a PipeTransport backed by the VME procfs file.
func (d *VMEDriver) Transport() transport.Transport {
	return transport.NewPipeTransport(d.dev, d.dev, d.slot)
}

// Close closes the VME procfs file.
func (d *VMEDriver) Close() error { return d.dev.Close() }
