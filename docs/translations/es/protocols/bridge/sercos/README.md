# SDK CmdBridge SERCOS III

SERCOS III (SErial Real-time COmmunication System) es una interfaz digital para control de movimiento que utiliza una topología de anillo sobre fibra óptica o cobre. Este paquete proporciona un codec SERCOS III mediante la utilidad de línea de comandos `netx_cli`.

## Herramienta CLI

Utiliza la utilidad CLI SERCOS III de Hilscher netX.

### Instalación

```bash
# Instalar el driver y las herramientas de Hilscher netX
# Ver https://www.hilscher.com/ para los drivers netX
sudo apt-get install netx-driver
```

Verificar: `netx_cli --help`

## Uso

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/sercos"
)

func main() {
    p := sercos.New()
    codec, _ := p.NewCodec("cmd")

    req := &kernel.Request{
        Function: "read",
        Address:  "S-0-51",
        Count:    2,
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## Driver

```go
b, c, err := sercos.NewCmdDriver("/usr/bin/netx_cli")
```

## Funciones soportadas

| Función | Descripción                 |
|----------|-----------------------------|
| `read`   | Leer parámetro IDN/S de SERCOS |
| `write`  | Escribir parámetro IDN/S de SERCOS|
| `phase`  | Establecer la fase de comunicación     |

## Pruebas

```bash
go test ./... -v
```
