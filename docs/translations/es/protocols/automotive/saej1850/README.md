# SDK del protocolo CAN SAE J1850

SAE J1850 es un estándar de comunicación vehicular utilizado para el diagnóstico a bordo (OBD-II). Este paquete proporciona un codec J1850 sobre bus CAN (ISO 15765-4 / CAN TP).

## Descripción del protocolo

SAE J1850 OBD-II sobre CAN utiliza identificadores CAN extendidos de 29 bits con la siguiente estructura:

### Formato de ID CAN (29 bits)

| Bits       | Campo    | Descripción                        |
|-----------|----------|------------------------------------|
| 28-26     | Prioridad | Prioridad del mensaje (0-7, por defecto 6)  |
| 25        | Ext ID   | Siempre 1 para tramas extendidas       |
| 24-16     | PF       | Formato de parámetro (cabecera)          |
| 15-8      | PS       | Parámetro específico (destino/origen) |
| 7-0       | SA       | Dirección de origen                     |

### IDs CAN OBD-II estándar

| Tipo                | ID CAN (hex)    | Descripción                    |
|--------------------|-----------------|--------------------------------|
| Solicitud física   | 0x18DAxxF1      | Solicitud a una ECU concreta (xx=dir) |
| Respuesta física   | 0x18DAF1xx      | Respuesta de una ECU (xx=dir)   |
| Solicitud funcional | 0x18DB33F1      | Difusión a todas las ECUs         |

### Formato de trama ISO 15765-2

Trama única: nibble alto del byte 0 = longitud de datos (0-7), nibble bajo + bytes restantes = datos de diagnóstico.

## Requisitos de hardware

- Sistema **Linux** con soporte SocketCAN
- Adaptador J1850-CAN OBD-II (p. ej. USB-to-CAN compatible con ELM327, OBDLink SX, Macchina M2)
- Vehículo con conector OBD-II (la mayoría de vehículos de 1996 en adelante)

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
    "github.com/erikwang2013/industrial-protocols-go/protocols/automotive/saej1850"
)

func main() {
    p := saej1850.New()
    codec, _ := p.NewCodec("can")

    // Modo $01 PID $0C: RPM del motor
    req := &kernel.Request{
        Function: "mode01",
        Metadata: map[string]any{"pid": float64(0x0C)},
    }
    raw, _ := codec.Encode(req)
    _ = raw

    // Modo $03: solicitar DTCs relacionados con emisiones
    req2 := &kernel.Request{Function: "mode03"}
    raw2, _ := codec.Encode(req2)
    _ = raw2
}
```

## Funciones soportadas

| Función        | Modo OBD-II | Descripción                         |
|----------------|-------------|-------------------------------------|
| `mode01`       | $01         | Solicitar datos actuales del tren motriz     |
| `mode03`       | $03         | Solicitar DTCs relacionados con emisiones       |
| `mode0A`       | $0A         | Solicitar DTCs permanentes              |
| `diag_request` | Personalizado | Solicitud de diagnóstico genérica          |
| `diag_response`| -           | Trama de respuesta de diagnóstico           |
| `broadcast`    | -           | Difusión funcional a todas las ECUs    |

## Pruebas

```bash
go test ./... -v
```
