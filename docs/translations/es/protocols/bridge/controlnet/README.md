# SDK CmdBridge ControlNet

ControlNet es un protocolo de red industrial en tiempo real desarrollado por Allen-Bradley (Rockwell Automation) para el intercambio de datos de alta velocidad y crítico en tiempo. Este paquete proporciona un codec ControlNet mediante la utilidad de línea de comandos `1784-pcic-cli`.

## Herramienta CLI

Utiliza la utilidad CLI de la tarjeta de interfaz ControlNet 1784-PCIC.

### Instalación

```bash
# Instalar el driver y las herramientas de Rockwell 1784-PCIC
# Consulte la documentación de Rockwell Automation sobre el SDK de RSLinx Classic
```

Verificar: `1784-pcic-cli --help`

## Uso

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/controlnet"
)

func main() {
    p := controlnet.New()
    codec, _ := p.NewCodec("cmd")

    req := &kernel.Request{
        Function: "read",
        Address:  "0x10",
        Count:    8,
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## Driver

```go
b, c, err := controlnet.NewCmdDriver("/usr/bin/1784-pcic-cli")
```

## Funciones soportadas

| Función | Descripción               |
|----------|---------------------------|
| `read`   | Leer de un nodo ControlNet |
| `write`  | Escribir en un nodo ControlNet  |
| `status` | Consultar el estado de la tarjeta PCIC    |

## Pruebas

```bash
go test ./... -v
```
