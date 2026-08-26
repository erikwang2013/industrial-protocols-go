# SDK CmdBridge EtherCAT

EtherCAT (Ethernet for Control Automation Technology) es un fieldbus de Ethernet industrial de alto rendimiento. Este paquete proporciona un codec EtherCAT mediante la utilidad de línea de comandos `ethercat`.

## Herramienta CLI

Utiliza la herramienta de línea de comandos del maestro IgH EtherCAT.

### Instalación

```bash
# Debian/Ubuntu
sudo apt-get install ethercat-master

# Desde el código fuente
git clone https://gitlab.com/etherlab.org/ethercat.git
cd ethercat
make && sudo make install
```

Verificar: `ethercat slaves`

## Uso

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/ethercat"
)

func main() {
    p := ethercat.New()
    codec, _ := p.NewCodec("cmd")

    // Subir SDO desde la dirección 0x1000
    req := &kernel.Request{
        Function: "upload",
        Address:  "0x1000",
        Count:    4,
    }
    raw, _ := codec.Encode(req)
    // raw = "upload 0x1000 4\n"
    _ = raw
}
```

## Driver

```go
b, c, err := ethercat.NewCmdDriver("/usr/bin/ethercat")
```

## Funciones soportadas

| Función   | Descripción            |
|------------|------------------------|
| `upload`   | Leer SDO de una dirección  |
| `download` | Escribir SDO en una dirección   |
| `slaves`   | Listar esclavos EtherCAT   |

## Pruebas

```bash
go test ./... -v
```
