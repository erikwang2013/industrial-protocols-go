# SDK do Protocolo CANopen

O CANopen é um protocolo de camada superior baseado em CAN para sistemas de controle embarcados. Este pacote fornece um codec CANopen e um driver SocketCAN.

## Visão Geral do Protocolo

O CANopen usa identificadores CAN padrão de 11 bits com o seguinte conjunto de conexões predefinido:

| Função    | ID CAN              | Descrição                   |
|------------|---------------------|-------------------------------|
| NMT        | 0x000               | Gerenciamento de rede        |
| SYNC       | 0x080               | Mensagem de sincronização    |
| SDO (tx)   | 0x580 + NodeID      | Service Data Object (servidor) |
| SDO (rx)   | 0x600 + NodeID      | Service Data Object (cliente) |
| PDO1 (tx)  | 0x180 + NodeID      | Process Data Object 1        |
| Heartbeat  | 0x700 + NodeID      | Heartbeat / Bootup           |

## Requisitos de Hardware

- Sistema **Linux** com suporte a SocketCAN (`CONFIG_CAN` habilitado)
- Uma interface CAN (ex.: `can0`, `vcan0` para CAN virtual)
- Hardware transceptor com capacidade CAN (ex.: MCP2515, SJA1000 ou adaptador USB-CAN)

### Configurando uma interface CAN virtual (para testes)

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

    // Read SDO object dictionary entry
    req := &kernel.Request{
        Function: "sdo_read",
        Metadata: map[string]any{
            "index": float64(0x1000), // Device type
            "sub":   float64(0),
        },
    }
    raw, _ := codec.Encode(req)
    fmt.Printf("SDO read frame: %X\n", raw)

    // NMT start remote node
    req2 := &kernel.Request{Function: "nmt_start"}
    raw2, _ := codec.Encode(req2)
    fmt.Printf("NMT start frame: %X\n", raw2)
}
```

## Funções Suportadas

| Função     | Descrição                          |
|-------------|--------------------------------------|
| `sdo_read`  | Ler entrada do dicionário de objetos |
| `sdo_write` | Gravar entrada do dicionário de objetos |
| `nmt_start` | Iniciar nó remoto (NMT)            |
| `nmt_stop`  | Parar nó remoto (NMT)              |
| `nmt_reset` | Reiniciar nó remoto (NMT)          |
| `heartbeat` | Enviar mensagem de heartbeat/bootup |

## Testes

```bash
go test ./... -v
```

Nota: os testes do driver SocketCAN exigem um sistema Linux com hardware CAN ou uma interface CAN virtual.
