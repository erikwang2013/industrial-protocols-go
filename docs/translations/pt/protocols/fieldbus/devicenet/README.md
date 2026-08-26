# SDK do Protocolo DeviceNet

O DeviceNet é um protocolo de rede industrial baseado em CAN para automação de fábricas. Este pacote fornece um codec DeviceNet com drivers SocketCAN e gateway TCP.

## Visão Geral do Protocolo

O DeviceNet usa o Common Industrial Protocol (CIP) sobre CAN. O Master/Slave Connection Set predefinido usa mensagens do Grupo 2 (CAN ID 0x400 + NodeID) para I/O por polling e mensagens explícitas.

## Requisitos de Hardware

### Modo CAN (SocketCAN)
- Sistema **Linux** com suporte a SocketCAN (`CONFIG_CAN` habilitado)
- Uma interface CAN (ex.: `can0`, `vcan0`)
- Hardware CAN com capacidade DeviceNet (ex.: Anybus Communicator, HMS IXXAT)

### Modo Gateway
- Conectividade TCP/IP com o gateway DeviceNet
- Gateway com suporte a protocolo de comando em texto simples (ex.: HMS Anybus, Hilscher netX)

### Configurando uma interface CAN virtual (para testes)

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

    // Poll request
    req := &kernel.Request{
        Function: "poll",
        Data:     []byte{0x01, 0x00},
    }
    raw, _ := codec.Encode(req)

    // Open connection
    req2 := &kernel.Request{Function: "open"}
    raw2, _ := codec.Encode(req2)
    _ = raw
    _ = raw2
}
```

### Modo Gateway

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

## Funções Suportadas

| Função | Descrição                    |
|----------|--------------------------------|
| `open`   | Abrir conexão explícita       |
| `poll`   | Polling de dados I/O (Grupo 2) |

## Testes

```bash
go test ./... -v
```

Nota: os testes do driver SocketCAN exigem Linux. Os testes do driver de gateway rodam em qualquer sistema operacional.
