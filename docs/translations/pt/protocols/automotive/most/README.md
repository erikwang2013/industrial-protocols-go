# SDK do Protocolo MOST Serial

O MOST (Media Oriented Systems Transport) é uma tecnologia de rede multimídia de alta velocidade usada principalmente em sistemas de infotainment automotivo. Este pacote fornece um codec MOST sobre um adaptador serial com interface de comandos AT.

## Visão Geral do Protocolo

O MOST usa comunicação serial síncrona sobre camada física de fibra óptica. Esta implementação se conecta por meio de um adaptador serial que expõe uma interface de comandos AT a 115200 baud.

### Comandos AT

| Comando            | Descrição                  |
|--------------------|------------------------------|
| `AT+READ=<addr>,<len>`  | Ler bytes do endereço |
| `AT+WRITE=<addr>,<hex>` | Gravar bytes hex no endereço |
| `AT+STATUS`             | Consultar status do anel/rede |

### Formato de Resposta

- `+OK:<hex_data>` -- resposta de sucesso
- `+ERR:<code>` -- resposta de erro

## Requisitos de Hardware

- Rede MOST de fibra óptica com terminação adequada
- Adaptador MOST-para-serial (ex.: adaptador USB MOST150)
- Porta serial a 115200 baud, 8N1

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

    // Read 4 bytes from address 0x0100
    req := &kernel.Request{
        Function: "read",
        Address:  "0x0100",
        Count:    4,
    }
    raw, _ := codec.Encode(req)
    // raw = "AT+READ=0x0100,4\r\n"
    _ = raw

    // Write data to address 0x0200
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

## Funções Suportadas

| Função | Descrição                 |
|----------|-----------------------------|
| `read`   | Ler de um endereço MOST    |
| `write`  | Gravar dados em um endereço MOST |
| `status` | Consultar status do anel/rede   |

## Testes

```bash
go test ./... -v
```
