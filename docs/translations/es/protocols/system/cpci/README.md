# SDK del protocolo CompactPCI

Implementa `kernel.Protocol` para el acceso al bus CompactPCI mediante sysfs de Linux.

## Protocolo

- **Name**: `cpci`
- **Variants**: `cpci`
- **Default Port**: 0 (bus mapeado en memoria)
- **Transport**: sysfs `/sys/bus/pci/devices/<BDF>/config`

CompactPCI (PICMG 2.0) utiliza la misma interfaz eléctrica y de software
que el PCI convencional. El codec pasa bytes en bruto entre la
aplicación y el espacio de configuración PCI mediante sysfs.

## Requisitos del kernel

Deben cargarse los siguientes módulos del kernel:

- `pcieport` -- driver de puerto PCI Express (para sistemas CPCIe híbridos)
- `pci_sysfs` -- interfaz PCI de sysfs (incorporada en la mayoría de kernels)
- `cpci_hotplug` -- controlador de conexión en caliente CompactPCI (opcional, para intercambio en caliente)

El sistema de archivos sysfs debe estar montado en `/sys`. Es el valor
predeterminado en todas las distribuciones modernas de Linux.

Permisos requeridos:
- Acceso root o `CAP_SYS_ADMIN` para el acceso al espacio de configuración
- El archivo de configuración pertenece a root:root con modo 0600

## Driver

```go
d, err := cpci.NewCPCIDriver("0000:02:00.0")
if err != nil {
    log.Fatal(err)
}
defer d.Close()

tr := d.Transport()
// Usar tr.Read/tr.Write para el acceso en bruto al espacio de configuración CPCI
```

## CompactPCI vs PCI

CompactPCI utiliza la enumeración de bus y el espacio de configuración PCI estándar.
Diferencias clave respecto al PCI de escritorio:

- **Formato 3U/6U**: mecánica Eurocard con conectores de pines y casquillos
- **Numeración de buses**: cada segmento de chasis CPCI recibe su propio número de bus PCI
- **Intercambio en caliente**: el intercambio en caliente PICMG 2.1 usa el modelo estándar de hotplug PCI
- **Slot de sistema**: el bus 0, dispositivo 0 es el controlador del slot de sistema

## Formato BDF

La dirección de bus es una cadena BDF (Bus:Device.Function):
- `0000:02:00.0` -- dominio 0000, bus 02, dispositivo 00, función 0
- `0000:02:08.0` -- dominio 0000, bus 02, dispositivo 08, función 0
