# SDK del protocolo CANopen

CANopen es un protocolo de capa superior basado en CAN para sistemas de control embebidos. Este paquete proporciona un codec CANopen y un driver SocketCAN.

## Descripción del protocolo

CANopen utiliza identificadores CAN estándar de 11 bits con el siguiente conjunto de conexiones predefinido:

| Función    | ID CAN              | Descripción                   |
|------------|---------------------|-------------------------------|
| NMT        | 0x000               | Gestión de red            |
| SYNC       | 0x080               | Mensaje de sincronización       |
| SDO (tx)   | 0x580 + NodeID      | Objeto de datos de servicio (servidor)  |
| SDO (rx)   | 0x600 + NodeID      | Objeto de datos de servicio (cliente)  |
| PDO1 (tx)  | 0x180 + NodeID      | Objeto de datos de proceso 1         |
| Heartbeat  | 0x700 + NodeID      | Heartbeat / Bootup            |

## Requisitos de hardware

- Sistema **Linux** con soporte SocketCAN (`CONFIG_CAN` habilitado)
- Una interfaz CAN (p. ej. `can0`, `vcan0` para CAN virtual)
- Hardware transceptor compatible con CAN (p. ej. MCP2515, SJA1000 o adaptador USB-CAN)

### Configuración de una interfaz CAN virtual (para pruebas)

```bash
sudo modprobe can
sudo modprobe can_raw
sudo modprobe vcan
sudo ip link add dev vcan0 type vcan
sudo ip link set up vcan0
```

## Uso

```go
package main

import (
    "fmt"
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/fieldbus/canopen"
)

func main() {
    p := canopen.New()
    codec, _ := p.NewCodec("can")

    // Leer entrada del diccionario de objetos SDO
    req := &kernel.Request{
        Function: "sdo_read",
        Metadata: map[string]any{
            "index": float64(0x1000), // Tipo de dispositivo
            "sub":   float64(0),
        },
    }
    raw, _ := codec.Encode(req)
    fmt.Printf("SDO read frame: %X\n", raw)

    // NMT iniciar nodo remoto
    req2 := &kernel.Request{Function: "nmt_start"}
    raw2, _ := codec.Encode(req2)
    fmt.Printf("NMT start frame: %X\n", raw2)
}
```

## Funciones soportadas

| Función     | Descripción                          |
|-------------|--------------------------------------|
| `sdo_read`  | Leer entrada del diccionario de objetos         |
| `sdo_write` | Escribir entrada del diccionario de objetos        |
| `nmt_start` | Iniciar nodo remoto (NMT)              |
| `nmt_stop`  | Detener nodo remoto (NMT)               |
| `nmt_reset` | Reiniciar nodo remoto (NMT)              |
| `heartbeat` | Enviar mensaje de heartbeat / bootup      |

## Pruebas

```bash
go test ./... -v
```

Nota: las pruebas del driver SocketCAN requieren un sistema Linux con hardware CAN o una interfaz CAN virtual.
