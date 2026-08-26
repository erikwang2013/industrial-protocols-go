# PCI / PCIe Protocol SDK

Implements `kernel.Protocol` for PCI and PCI Express bus access via Linux sysfs.

## Protocol

- **Name**: `pci`
- **Variants**: `pci`
- **Default Port**: 0 (memory-mapped bus)
- **Transport**: sysfs `/sys/bus/pci/devices/<BDF>/config`

The codec passes through raw bytes between the application and the PCI
config space. Reads and writes go directly to the device's configuration
registers via the sysfs config file.

## Kernel Requirements

The following kernel modules must be loaded:

- `pcieport` -- PCI Express port driver
- `pci_sysfs` -- sysfs PCI interface (built into most kernels)

The sysfs filesystem must be mounted at `/sys`. This is the default on
all modern Linux distributions.

Required permissions:
- Root access or `CAP_SYS_ADMIN` for config space access
- The config file is owned by root:root with mode 0600

## Driver

```go
d, err := pci.NewPCIDriver("0000:00:1f.3")
if err != nil {
    log.Fatal(err)
}
defer d.Close()

tr := d.Transport()
// Use tr.Read/tr.Write for raw PCI config space access
```

## BDF Format

The bus address is a BDF (Bus:Device.Function) string:
- `0000:00:1f.3` -- domain 0000, bus 00, device 1f, function 3
- `0000:01:00.0` -- domain 0000, bus 01, device 00, function 0
