# SDK CmdBridge WirelessHART

WirelessHART (IEC 62591) es un estándar de redes industriales inalámbricas basado en el protocolo HART. Este paquete proporciona un codec WirelessHART mediante la utilidad de línea de comandos `emerson_1410_cli` (gateway inalámbrico Emerson 1410/1420).

## Herramienta CLI

Utiliza la utilidad CLI del gateway inalámbrico Emerson 1410/1420.

### Instalación

```bash
# Instalar el software y las herramientas de Emerson Wireless Gateway
# Consulte la documentación de Emerson sobre la configuración del gateway 1410/1420
```

Verificar: `emerson_1410_cli --help`

## Uso

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/wirelesshart"
)

func main() {
    p := wirelesshart.New()
    codec, _ := p.NewCodec("cmd")

    req := &kernel.Request{
        Function: "read",
        Address:  "TT101",
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## Driver

```go
b, c, err := wirelesshart.NewCmdDriver("/usr/bin/emerson_1410_cli")
```

## Funciones soportadas

| Función | Descripción                |
|----------|----------------------------|
| `read`   | Leer parámetro de dispositivo      |
| `write`  | Escribir parámetro de dispositivo     |
| `scan`   | Escanear dispositivos inalámbricos  |

## Pruebas

```bash
go test ./... -v
```
