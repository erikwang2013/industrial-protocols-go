# SDK CmdBridge POWERLINK

POWERLINK (Ethernet POWERLINK) es un protocolo Ethernet en tiempo real para automatización industrial. Este paquete proporciona un codec POWERLINK mediante la utilidad de línea de comandos `openPOWERLINK_demo`.

## Herramienta CLI

Utiliza la aplicación de demostración de la pila openPOWERLINK.

### Instalación

```bash
git clone https://github.com/OpenAutomationTechnologies/openPOWERLINK.git
cd openPOWERLINK
cmake -B build && cmake --build build
sudo cmake --install build
```

Verificar: `openPOWERLINK_demo --help`

## Uso

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/powerlink"
)

func main() {
    p := powerlink.New()
    codec, _ := p.NewCodec("cmd")

    req := &kernel.Request{
        Function: "read",
        Address:  "0x2000",
        Count:    8,
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## Driver

```go
b, c, err := powerlink.NewCmdDriver("/usr/bin/openPOWERLINK_demo")
```

## Funciones soportadas

| Función | Descripción                   |
|----------|-------------------------------|
| `read`   | Leer entrada del diccionario de objetos  |
| `write`  | Escribir entrada del diccionario de objetos |
| `status` | Consultar estado del nodo/NMT          |

## Pruebas

```bash
go test ./... -v
```
