# VME / VPX Protocol SDK

Implements `kernel.Protocol` for VMEbus and VPX bus access via Linux procfs.

## Protocol

- **Name**: `vme`
- **Variants**: `vme`
- **Default Port**: 0 (memory-mapped bus)
- **Transport**: procfs `/proc/vme/<slot>`

The codec passes through raw bytes between the application and the VME
address space. Reads and writes go directly to the bus via the procfs
interface provided by the VME kernel driver.

## Kernel Requirements

The following kernel module must be loaded:

- `vme_tsi148` -- Tundra TSI148 VME bridge driver (most common)
  - Also supported: `vme_ca91cx42` (Universe II), `vme_user`

The procfs filesystem must be mounted at `/proc`.

Required permissions:
- Root access is required to open `/proc/vme/<slot>` for read-write
- The device nodes are owned by root:root

## Driver

```go
d, err := vme.NewVMEDriver(0) // slot 0
if err != nil {
    log.Fatal(err)
}
defer d.Close()

tr := d.Transport()
// Use tr.Read/tr.Write for raw VME bus access
```

## VME Addressing Modes

The codec passes through all address modifier and address bytes
transparently. Applications should prepend addressing information
to the data payload:

- **A16**: 16-bit short I/O address space
- **A24**: 24-bit standard address space
- **A32**: 32-bit extended address space

## VPX Compatibility

VPX (VITA 46) systems that expose a VME-compatible procfs interface
can use this driver. The slot numbering is the same.
