# SDK do Protocolo SAE J1850 sobre CAN

O SAE J1850 é um padrão de comunicação veicular usado para diagnóstico de bordo (OBD-II). Este pacote fornece um codec J1850 sobre barramento CAN (ISO 15765-4 / CAN TP).

## Visão Geral do Protocolo

O OBD-II SAE J1850 sobre CAN usa identificadores CAN estendidos de 29 bits com a seguinte estrutura:

### Formato do ID CAN (29 bits)

| Bits       | Campo    | Descrição                        |
|-----------|----------|------------------------------------|
| 28-26     | Priority | Prioridade da mensagem (0-7, padrão 6) |
| 25        | Ext ID   | Sempre 1 para quadros estendidos  |
| 24-16     | PF       | Formato de Parâmetro (Header)     |
| 15-8      | PS       | Parâmetro Específico (destino/origem) |
| 7-0       | SA       | Endereço de origem                |

### IDs CAN OBD-II Padrão

| Tipo                | ID CAN (hex)    | Descrição                    |
|--------------------|-----------------|--------------------------------|
| Solicitação física | 0x18DAxxF1      | Solicitação a ECU específica (xx=addr) |
| Resposta física    | 0x18DAF1xx      | Resposta da ECU (xx=addr)    |
| Solicitação funcional | 0x18DB33F1      | Transmissão para todas as ECUs |

### Formato de Quadro ISO 15765-2

Quadro único: o nibble alto do byte 0 = tamanho dos dados (0-7); nibble baixo + bytes restantes = dados de diagnóstico.

## Requisitos de Hardware

- Sistema **Linux** com suporte a SocketCAN
- Adaptador J1850-CAN OBD-II (ex.: USB-para-CAN compatível com ELM327, OBDLink SX, Macchina M2)
- Veículo com conector OBD-II (maioria dos veículos de 1996 em diante)

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
    "github.com/erikwang2013/industrial-protocols-go/protocols/automotive/saej1850"
)

func main() {
    p := saej1850.New()
    codec, _ := p.NewCodec("can")

    // Mode $01 PID $0C: Engine RPM
    req := &kernel.Request{
        Function: "mode01",
        Metadata: map[string]any{"pid": float64(0x0C)},
    }
    raw, _ := codec.Encode(req)
    _ = raw

    // Mode $03: Request emission-related DTCs
    req2 := &kernel.Request{Function: "mode03"}
    raw2, _ := codec.Encode(req2)
    _ = raw2
}
```

## Funções Suportadas

| Função        | Modo OBD-II | Descrição                         |
|----------------|-------------|-------------------------------------|
| `mode01`       | $01         | Solicitar dados atuais do trem de força |
| `mode03`       | $03         | Solicitar DTCs relacionados a emissões |
| `mode0A`       | $0A         | Solicitar DTCs permanentes         |
| `diag_request` | Personalizado | Solicitação de diagnóstico genérica |
| `diag_response`| -           | Quadro de resposta de diagnóstico  |
| `broadcast`    | -           | Transmissão funcional para todas as ECUs |

## Testes

```bash
go test ./... -v
```
