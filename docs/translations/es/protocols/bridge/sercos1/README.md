# SDK CmdBridge SERCOS I/II

SERCOS I/II es la versión heredada de interfaz serial de fibra óptica del sistema SERCOS para control de movimiento digital. Este paquete proporciona un codec SERCOS I/II mediante la utilidad de línea de comandos `sercos_cli`.

## Herramienta CLI

Utiliza una utilidad CLI de interfaz de fibra óptica SERCOS.

### Instalación

```bash
# Driver y herramientas de la tarjeta de interfaz SERCOS
# Consulte la documentación del fabricante para la tarjeta maestra SERCOS específica
```

Verificar: `sercos_cli --help`

## Uso

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/sercos1"
)

func main() {
    p := sercos1.New()
    codec, _ := p.NewCodec("cmd")

    req := &kernel.Request{
        Function: "read",
        Address:  "0x0100",
        Count:    4,
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## Driver

```go
b, c, err := sercos1.NewCmdDriver("/usr/bin/sercos_cli")
```

## Funciones soportadas

| Función | Descripción              |
|----------|--------------------------|
| `read`   | Leer IDN SERCOS          |
| `write`  | Escribir IDN SERCOS         |
| `status` | Consultar el estado del accionamiento       |

## Pruebas

```bash
go test ./... -v
```
