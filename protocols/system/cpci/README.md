# CompactPCI Protocol SDK

Implements `kernel.Protocol` for CompactPCI bus access via Linux sysfs.

## Protocol

- **Name**: `cpci`
- **Variants**: `cpci`
- **Default Port**: 0 (memory-mapped bus)
- **Transport**: sysfs `/sys/bus/pci/devices/<BDF>/config`

CompactPCI (PICMG 2.0) uses the same electrical and software interface
as conventional PCI. The codec passes through raw bytes between the
application and the PCI config space via sysfs.

## Kernel Requirements

The following kernel modules must be loaded:

- `pcieport` -- PCI Express port driver (for hybrid CPCIe systems)
- `pci_sysfs` -- sysfs PCI interface (built into most kernels)
- `cpci_hotplug` -- CompactPCI hotplug controller (optional, for hot-swap)

The sysfs filesystem must be mounted at `/sys`. This is the default on
all modern Linux distributions.

Required permissions:
- Root access or `CAP_SYS_ADMIN` for config space access
- The config file is owned by root:root with mode 0600

## Driver

```go
d, err := cpci.NewCPCIDriver("0000:02:00.0")
if err != nil {
    log.Fatal(err)
}
defer d.Close()

tr := d.Transport()
// Use tr.Read/tr.Write for raw CPCI config space access
```

## CompactPCI vs PCI

CompactPCI uses the standard PCI bus enumeration and config space.
Key differences from desktop PCI:

- **3U/6U form factor**: Eurocard mechanics with pin-and-socket connectors
- **Bus numbering**: Each CPCI chassis segment gets its own PCI bus number
- **Hot swap**: PICMG 2.1 hot swap uses the standard PCI hotplug model
- **System slot**: Bus 0, device 0 is the system slot controller

## BDF Format

The bus address is a BDF (Bus:Device.Function) string:
- `0000:02:00.0` -- domain 0000, bus 02, device 00, function 0
- `0000:02:08.0` -- domain 0000, bus 02, device 08, function 0
