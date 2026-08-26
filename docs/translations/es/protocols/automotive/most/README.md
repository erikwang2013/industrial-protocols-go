# SDK del protocolo serial MOST

MOST (Media Oriented Systems Transport) es una tecnología de red multimedia de alta velocidad utilizada principalmente en sistemas de infoentretenimiento automotriz. Este paquete proporciona un codec MOST sobre un adaptador serial con interfaz de comandos AT.

## Descripción del protocolo

MOST utiliza comunicación serial síncrona sobre capa física de fibra óptica. Esta implementación se conecta mediante un adaptador serial que expone una interfaz de comandos AT a 115200 baudios.

### Comandos AT

| Comando            | Descripción                  |
|--------------------|------------------------------|
| `AT+READ=<addr>,<len>`  | Leer bytes de una dirección |
| `AT+WRITE=<addr>,<hex>` | Escribir bytes hex en una dirección |
| `AT+STATUS`             | Consultar estado del anillo/red |

### Formato de respuesta

- `+OK:<hex_data>` -- respuesta correcta
- `+ERR:<code>` -- respuesta de error

## Requisitos de hardware

- Red de fibra óptica MOST con terminación adecuada
- Adaptador MOST a serial (p. ej. adaptador USB MOST150)
- Puerto serial a 115200 baudios, 8N1

## Uso

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/automotive/most"
)

func main() {
    p := most.New()
    codec, _ := p.NewCodec("serial")

    // Leer 4 bytes de la dirección 0x0100
    req := &kernel.Request{
        Function: "read",
        Address:  "0x0100",
        Count:    4,
    }
    raw, _ := codec.Encode(req)
    // raw = "AT+READ=0x0100,4\r\n"
    _ = raw

    // Escribir datos en la dirección 0x0200
    req2 := &kernel.Request{
        Function: "write",
        Address:  "0x0200",
        Data:     []byte{0xDE, 0xAD, 0xBE, 0xEF},
    }
    raw2, _ := codec.Encode(req2)
    // raw2 = "AT+WRITE=0x0200,DEADBEEF\r\n"
    _ = raw2
}
```

## Driver

```go
b, c, err := most.NewSerialDriver("/dev/ttyUSB0")
```

## Funciones soportadas

| Función | Descripción                 |
|----------|-----------------------------|
| `read`   | Leer de una dirección MOST    |
| `write`  | Escribir datos en una dirección MOST |
| `status` | Consultar estado del anillo/red   |

## Pruebas

```bash
go test ./... -v
```
