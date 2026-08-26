# SDK del protocolo DeviceNet

DeviceNet es un protocolo de red industrial basado en CAN para la automatización de fábricas. Este paquete proporciona un codec DeviceNet con drivers SocketCAN y de gateway TCP.

## Descripción del protocolo

DeviceNet utiliza el protocolo industrial común (CIP) sobre CAN. El conjunto de conexiones maestro/esclavo predefinido utiliza mensajes de Grupo 2 (ID CAN 0x400 + NodeID) para I/O por polling y mensajería explícita.

## Requisitos de hardware

### Modo CAN (SocketCAN)
- Sistema **Linux** con soporte SocketCAN (`CONFIG_CAN` habilitado)
- Una interfaz CAN (p. ej. `can0`, `vcan0`)
- Hardware CAN compatible con DeviceNet (p. ej. Anybus Communicator, HMS IXXAT)

### Modo gateway
- Conectividad TCP/IP al gateway DeviceNet
- Gateway con protocolo de comandos en texto plano (p. ej. HMS Anybus, Hilscher netX)

### Configuración de una interfaz CAN virtual (para pruebas)

```bash
sudo modprobe can
sudo modprobe can_raw
sudo modprobe vcan
sudo ip link add dev vcan0 type vcan
sudo ip link set up vcan0
```

## Uso

### Modo CAN

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/fieldbus/devicenet"
)

func main() {
    p := devicenet.New()
    codec, _ := p.NewCodec("can")

    // Solicitud de polling
    req := &kernel.Request{
        Function: "poll",
        Data:     []byte{0x01, 0x00},
    }
    raw, _ := codec.Encode(req)

    // Abrir conexión
    req2 := &kernel.Request{Function: "open"}
    raw2, _ := codec.Encode(req2)
    _ = raw
    _ = raw2
}
```

### Modo gateway

```go
p := devicenet.New()
codec, _ := p.NewCodec("gateway")

req := &kernel.Request{
    Function: "poll",
    Data:     []byte{0xAB, 0xCD},
}
raw, _ := codec.Encode(req)
// raw will be: "poll abcd\n"
```

## Funciones soportadas

| Función | Descripción                    |
|----------|--------------------------------|
| `open`   | Abrir conexión explícita       |
| `poll`   | Hacer polling de datos I/O (Grupo 2)        |

## Pruebas

```bash
go test ./... -v
```

Nota: las pruebas del driver SocketCAN requieren Linux. Las pruebas del driver de gateway se ejecutan en cualquier sistema operativo.
