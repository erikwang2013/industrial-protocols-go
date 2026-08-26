# SDK del protocolo PCI / PCIe

Implementa `kernel.Protocol` para el acceso al bus PCI y PCI Express mediante sysfs de Linux.

## Protocolo

- **Name**: `pci`
- **Variants**: `pci`
- **Default Port**: 0 (bus mapeado en memoria)
- **Transport**: sysfs `/sys/bus/pci/devices/<BDF>/config`

El codec pasa bytes en bruto entre la aplicación y el espacio de
configuración PCI. Las lecturas y escrituras van directamente a los
registros de configuración del dispositivo mediante el archivo de configuración de sysfs.

## Requisitos del kernel

Deben cargarse los siguientes módulos del kernel:

- `pcieport` -- driver de puerto PCI Express
- `pci_sysfs` -- interfaz PCI de sysfs (incorporada en la mayoría de kernels)

El sistema de archivos sysfs debe estar montado en `/sys`. Es el valor
predeterminado en todas las distribuciones modernas de Linux.

Permisos requeridos:
- Acceso root o `CAP_SYS_ADMIN` para el acceso al espacio de configuración
- El archivo de configuración pertenece a root:root con modo 0600

## Driver

```go
d, err := pci.NewPCIDriver("0000:00:1f.3")
if err != nil {
    log.Fatal(err)
}
defer d.Close()

tr := d.Transport()
// Usar tr.Read/tr.Write para el acceso en bruto al espacio de configuración PCI
```

## Formato BDF

La dirección de bus es una cadena BDF (Bus:Device.Function):
- `0000:00:1f.3` -- dominio 0000, bus 00, dispositivo 1f, función 3
- `0000:01:00.0` -- dominio 0000, bus 01, dispositivo 00, función 0
