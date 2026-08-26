# SDK del protocolo CAN FlexRay

FlexRay es un protocolo de comunicación automotriz determinista de alta velocidad. Este paquete proporciona un codec FlexRay sobre bus CAN, con tramado basado en ciclos y comprobación de integridad CRC-16/XMODEM.

## Descripción del protocolo

FlexRay utiliza un esquema de acceso múltiple por división de tiempo (TDMA) con ciclos de comunicación repetidos. Cada ciclo consta de segmentos estáticos y dinámicos. Este codec mapea las tramas FlexRay sobre tramas CAN extendidas de 29 bits.

### Formato de línea

Carga útil de un ciclo FlexRay:
- **Cabecera** (2 bytes): número de ciclo (little-endian)
- **Estado** (1 byte): bit 7=PPI (Payload Preamble Indicator), bit 6=NFI, bit 5=SYF, bit 4=SUF
- **Datos** (N bytes): carga útil (máx. 254 bytes)
- **CRC** (2 bytes): CRC-16/XMODEM sobre cabecera+estado+datos (little-endian)

Codificación del ID CAN (extendido de 29 bits):
- Bits 28-24: Tipo de mensaje (0x01=trama, 0x02=estado)
- Bits 23-10: Reservados
- Bits 15-10: Slot ID (6 bits)
- Bits 9-0: Número de ciclo (10 bits)

## Requisitos de hardware

- Sistema **Linux** con soporte SocketCAN
- Vector VN7600/VN7640 o adaptador FlexRay-CAN de Bosch
- Red FlexRay con terminación adecuada (bias de 2,5 V)

### Configuración de la interfaz CAN

```bash
sudo modprobe can
sudo modprobe can_raw
sudo ip link set can0 type can bitrate 500000
sudo ip link set up can0
```

## Uso

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/automotive/flexray"
)

func main() {
    p := flexray.New()
    codec, _ := p.NewCodec("can")

    // Enviar una trama FlexRay en el ciclo 5 con PPI activado
    req := &kernel.Request{
        Function: "frame",
        Data:     []byte{0x42, 0x01},
        Metadata: map[string]any{
            "cycle": float64(5),
            "ppi":   true,
        },
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## Funciones soportadas

| Función | Descripción                              |
|----------|------------------------------------------|
| `frame`  | Enviar carga útil de trama FlexRay        |
| `status` | Consultar configuración de slot (slot, ciclo) |

## Pruebas

```bash
go test ./... -v
```
