# SDK del protocolo VME / VPX

Implementa `kernel.Protocol` para el acceso al bus VME y VPX mediante procfs de Linux.

## Protocolo

- **Name**: `vme`
- **Variants**: `vme`
- **Default Port**: 0 (bus mapeado en memoria)
- **Transport**: procfs `/proc/vme/<slot>`

El codec pasa bytes en bruto entre la aplicación y el espacio de
direcciones VME. Las lecturas y escrituras van directamente al bus mediante la
interfaz procfs proporcionada por el driver VME del kernel.

## Requisitos del kernel

Debe cargarse el siguiente módulo del kernel:

- `vme_tsi148` -- driver de puente VME Tundra TSI148 (el más común)
  - También compatible: `vme_ca91cx42` (Universe II), `vme_user`

El sistema de archivos procfs debe estar montado en `/proc`.

Permisos requeridos:
- Se requiere acceso root para abrir `/proc/vme/<slot>` en modo lectura-escritura
- Los nodos de dispositivo pertenecen a root:root

## Driver

```go
d, err := vme.NewVMEDriver(0) // slot 0
if err != nil {
    log.Fatal(err)
}
defer d.Close()

tr := d.Transport()
// Usar tr.Read/tr.Write para el acceso en bruto al bus VME
```

## Modos de direccionamiento VME

El codec pasa de forma transparente todos los bytes de modificador de
dirección y de dirección. Las aplicaciones deben anteponer la información
de direccionamiento a la carga útil de datos:

- **A16**: espacio de direcciones I/O corto de 16 bits
- **A24**: espacio de direcciones estándar de 24 bits
- **A32**: espacio de direcciones extendido de 32 bits

## Compatibilidad VPX

Los sistemas VPX (VITA 46) que exponen una interfaz procfs compatible
con VME pueden utilizar este driver. La numeración de slots es la misma.
