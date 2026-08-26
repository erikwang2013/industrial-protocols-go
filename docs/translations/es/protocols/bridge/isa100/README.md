# SDK CmdBridge ISA100.11a

ISA100.11a es un estándar de redes industriales inalámbricas para automatización de procesos. Este paquete proporciona un codec ISA100.11a mediante la utilidad de línea de comandos `yfgw410_cli` (gateway inalámbrico de campo Yokogawa YFGW410).

## Herramienta CLI

Utiliza la CLI del gateway inalámbrico de campo Yokogawa YFGW410.

### Instalación

```bash
# Instalar el software y las herramientas del gateway Yokogawa YFGW410
# Consulte la documentación de Yokogawa sobre la configuración del Field Wireless Gateway
```

Verificar: `yfgw410_cli --help`

## Uso

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/isa100"
)

func main() {
    p := isa100.New()
    codec, _ := p.NewCodec("cmd")

    req := &kernel.Request{
        Function: "read",
        Address:  "DEV001",
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## Driver

```go
b, c, err := isa100.NewCmdDriver("/usr/bin/yfgw410_cli")
```

## Funciones soportadas

| Función | Descripción                 |
|----------|-----------------------------|
| `read`   | Leer atributo de dispositivo       |
| `write`  | Escribir atributo de dispositivo      |
| `list`   | Listar dispositivos aprovisionados    |

## Pruebas

```bash
go test ./... -v
```
