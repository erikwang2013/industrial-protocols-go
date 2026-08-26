# SDK do Protocolo FlexRay sobre CAN

O FlexRay é um protocolo de comunicação automotiva determinístico de alta velocidade. Este pacote fornece um codec FlexRay sobre barramento CAN, com enquadramento baseado em ciclos e verificação de integridade CRC-16/XMODEM.

## Visão Geral do Protocolo

O FlexRay usa um esquema de acesso múltiplo por divisão de tempo (TDMA) com ciclos de comunicação repetidos. Cada ciclo é composto por segmentos estáticos e dinâmicos. Este codec mapeia quadros FlexRay em quadros CAN estendidos de 29 bits.

### Formato de Linha

Payload de ciclo FlexRay:
- **Header** (2 bytes): número do ciclo (little-endian)
- **Status** (1 byte): bit 7=PPI (Payload Preamble Indicator), bit 6=NFI, bit 5=SYF, bit 4=SUF
- **Data** (N bytes): payload (máx. 254 bytes)
- **CRC** (2 bytes): CRC-16/XMODEM sobre header+status+data (little-endian)

Codificação do ID CAN (estendido de 29 bits):
- Bits 28-24: tipo de mensagem (0x01=quadro, 0x02=status)
- Bits 23-10: reservados
- Bits 15-10: ID do slot (6 bits)
- Bits 9-0: número do ciclo (10 bits)

## Requisitos de Hardware

- Sistema **Linux** com suporte a SocketCAN
- Vector VN7600/VN7640 ou adaptador Bosch FlexRay-CAN
- Rede FlexRay com terminação adequada (bias de 2,5V)

### Configurando a interface CAN

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

    // Send a FlexRay frame in cycle 5 with PPI set
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

## Funções Suportadas

| Função | Descrição                              |
|----------|------------------------------------------|
| `frame`  | Enviar payload de quadro FlexRay         |
| `status` | Consultar configuração de slot (slot, cycle) |

## Testes

```bash
go test ./... -v
```
